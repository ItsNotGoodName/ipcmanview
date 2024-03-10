package squeuel

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/repo"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	ErrDuplicateTaskID = errors.New("duplicate task id")
	ErrSkipRetry       = errors.New("skip retry")
)

type HandleFunc func(ctx context.Context, task *Task) error

type Option func(c *Task)

func TaskID(taskID string) Option {
	return func(c *Task) {
		c.TaskID.String = c.Queue + "|" + taskID
		c.TaskID.Valid = true
	}
}

func MaxRetry(maxRetry int) Option {
	return func(c *Task) {
		c.MaxRetry = maxRetry
	}
}

func NewTask(queue string, payload []byte, options ...Option) *Task {
	task := &Task{
		ID:       uuid.NewString(),
		Queue:    queue,
		Payload:  payload,
		TaskID:   sql.NullString{},
		MaxRetry: 3,
	}

	for _, option := range options {
		option(task)
	}

	return task
}

type Task struct {
	ID       string
	Queue    string
	Payload  []byte
	TaskID   sql.NullString
	MaxRetry int
}

func NewTaskBuilder[T any](queue string) TaskBuilder[T] {
	return TaskBuilder[T]{
		Queue: queue,
	}
}

type TaskBuilder[T any] struct {
	Queue   string
	payload T
}

func (e TaskBuilder[T]) New(payload T, options ...Option) (*Task, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return NewTask(e.Queue, b, options...), nil
}

func (e TaskBuilder[T]) Payload(task *Task) (T, error) {
	var payload T
	if task.Queue != e.Queue {
		return payload, fmt.Errorf("wrong queue: %s != %s", task.Queue, e.Queue)
	}
	err := json.Unmarshal(task.Payload, &payload)
	return payload, err
}

func dequeueTask(ctx context.Context, db sqlite.DB, queue string, duration time.Duration) (*Task, error) {
	now := time.Now()
	timeout := now.Add(duration)

	v, err := db.C().SqueuelDequeue(ctx, repo.SqueuelDequeueParams{
		Timeout: types.NewTime(timeout),
		Queue:   queue,
		Now:     types.NewTime(timeout),
	})
	if err != nil {
		if core.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &Task{
		ID:       v.ID,
		Queue:    queue,
		Payload:  v.Payload,
		TaskID:   v.TaskID,
		MaxRetry: int(v.MaxReceived),
	}, nil
}

func extendTask(ctx context.Context, db sqlite.DB, task *Task, duration time.Duration) error {
	timeout := time.Now().Add(duration)

	return db.C().SqueuelExtend(ctx, repo.SqueuelExtendParams{
		Timeout: types.NewTime(timeout),
		Queue:   task.Queue,
		ID:      task.ID,
	})
}

func deleteTask(ctx context.Context, db sqlite.DB, task *Task) error {
	return db.C().SqueuelDelete(ctx, repo.SqueuelDeleteParams{
		Queue: task.Queue,
		ID:    task.ID,
	})
}

func deleteExpiredTasks(ctx context.Context, db sqlite.DBTx) error {
	return db.SqueuelDeleteExpired(ctx, types.NewTime(time.Now()))
}

func EnqueueTask(ctx context.Context, db sqlite.DB, hub *bus.Hub, task *Task) (string, error) {
	id, err := enqueueTask(ctx, db.C(), hub, task)
	if err != nil {
		return "", err
	}

	hub.SqueuelEnqueued(bus.SqueuelEnqueued{
		Queue: task.Queue,
	})

	return id, nil
}

func EnqueueTaskTx(ctx context.Context, tx sqlite.Tx, hub *bus.Hub, task *Task) (string, error) {
	id, err := enqueueTask(ctx, tx.C(), hub, task)
	if err != nil {
		return "", err
	}

	tx.CommitHook(func() {
		hub.SqueuelEnqueued(bus.SqueuelEnqueued{
			Queue: task.Queue,
		})
	})

	return id, nil
}

func enqueueTask(ctx context.Context, db sqlite.DBTx, hub *bus.Hub, task *Task) (string, error) {
	now := types.NewTime(time.Now())

	err := deleteExpiredTasks(ctx, db)
	if err != nil {
		return "", err
	}

	err = db.SqueuelEnqueue(ctx, repo.SqueuelEnqueueParams{
		ID:          task.ID,
		TaskID:      task.TaskID,
		Queue:       task.Queue,
		Payload:     task.Payload,
		Timeout:     now,
		Received:    0,
		MaxReceived: int64(task.MaxRetry),
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		if _, ok := sqlite.AsConstraintError(err, sqlite.CONSTRAINT_UNIQUE); ok {
			return "", fmt.Errorf("%w: %s", ErrDuplicateTaskID, task.TaskID.String)
		}
		return "", err
	}

	return task.ID, nil
}

func Do(ctx context.Context, db sqlite.DB, queue string, fn HandleFunc) (bool, error) {
	duration := 30 * time.Second

	// Dequeue task
	task, err := dequeueTask(ctx, db, queue, duration)
	if err != nil {
		return false, err
	}
	if task == nil {
		return false, nil
	}

	// Keep task alive
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	t := time.NewTicker(duration / 2)
	defer t.Stop()

	go func() {
		// TODO: cancel context when there is some sort of signal to cancel task
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := extendTask(ctx, db, task, duration); err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					log.Err(err).Str("package", "squeuel").Msg("Failed to extend task")
				}
			}
		}
	}()

	log.Info().Str("package", "squeuel").Str("id", task.ID).Str("queue", task.Queue).RawJSON("payload", task.Payload).Msg("Starting task")

	// Execute task
	if err := fn(ctx, task); err != nil {
		if errors.Is(err, ErrSkipRetry) {
			// Delete task
			if err := deleteTask(ctx, db, task); err != nil {
				return false, err
			}
		}

		return false, err
	}

	// Delete task
	if err := deleteTask(ctx, db, task); err != nil {
		return false, err
	}

	return true, nil
}
