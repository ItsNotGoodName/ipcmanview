package dahua

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/types"
	"github.com/ItsNotGoodName/ipcmanview/pkg/dahuacgi"
	"github.com/ItsNotGoodName/ipcmanview/pkg/pagination"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/oklog/ulid/v2"
	"github.com/thejerf/suture/v4"
)

type Event struct {
	ID         string
	Device_ID  int64
	Code       string
	Action     string
	Index      int64
	Data       types.JSON
	Created_At types.Time
}

func DeleteEvents(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM dahua_events
	`)
	return err
}

type ListEventsFilter struct {
	DeviceIDs []int64
	Codes     []string
	Actions   []string
}

func (arg ListEventsFilter) where() sq.Eq {
	where := sq.Eq{}
	if len(arg.DeviceIDs) != 0 {
		where["dahua_events.device_id"] = arg.DeviceIDs
	}
	if len(arg.Codes) != 0 {
		where["dahua_events.code"] = arg.Codes
	}
	if len(arg.Actions) != 0 {
		where["dahua_events.action"] = arg.Actions
	}
	return where
}

type ListEventsParams struct {
	pagination.Page
	Ascending bool
	Filter    ListEventsFilter
}

type ListEventsResult struct {
	pagination.PageResult
	Items []ListEventsItem
}

type ListEventsItem struct {
	Event
	Device_UUID string
}

func ListEvents(ctx context.Context, db *sqlx.DB, arg ListEventsParams) (ListEventsResult, error) {
	where := arg.Filter.where()

	order := "dahua_events.id"
	if arg.Ascending {
		order += " ASC"
	} else {
		order += " DESC"
	}
	sb := sq.
		Select(
			"dahua_events.*",
			"dahua_devices.uuid AS device_uuid",
		).
		From("dahua_events").
		LeftJoin("dahua_devices ON dahua_devices.id = dahua_events.device_id").
		Where(where).
		OrderBy(order).
		Offset(uint64(arg.Offset())).
		Limit(uint64(arg.Limit()))

	var items []ListEventsItem
	q, a, err := sb.ToSql()
	if err != nil {
		return ListEventsResult{}, err
	}
	if err := db.SelectContext(ctx, &items, q, a...); err != nil {
		return ListEventsResult{}, err
	}

	sb = sq.
		Select("COUNT(*) AS count").
		From("dahua_events").
		Where(where)

	q, a, err = sb.ToSql()
	if err != nil {
		return ListEventsResult{}, err
	}
	var count int
	if err := db.GetContext(ctx, &count, q, a...); err != nil {
		return ListEventsResult{}, err
	}

	return ListEventsResult{
		PageResult: arg.Result(int(count)),
		Items:      items,
	}, nil
}

type EventRule struct {
	types.Key
	Code       string
	Allow_DB   bool
	Allow_Live bool
	Allow_MQTT bool
	Can_Delete bool
}

func NormalizeEventRules(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, `
		INSERT INTO dahua_event_rules(code, can_delete, uuid) VALUES ('All', false, ?) ON CONFLICT DO NOTHING
	`, uuid.NewString())
	return err
}

type CreateEventRuleArgs struct {
	UUID      string
	Code      string
	AllowDB   bool
	AllowLive bool
	AllowMQTT bool
}

func CreateEventRule(ctx context.Context, db *sqlx.DB, args CreateEventRuleArgs) (EventRule, error) {
	var row EventRule
	err := sqlx.GetContext(ctx, db, &row, `
		INSERT INTO dahua_event_rules (
			uuid, code, allow_db, allow_live, allow_mqtt
		) 
		VALUES (?, ?, ?, ?, ?) 
		RETURNING *;
	`,
		args.UUID,
		args.Code,
		args.AllowDB,
		args.AllowLive,
		args.AllowMQTT,
	)
	return row, err
}

type UpdateEventRuleArgs struct {
	UUID      string
	Code      string
	AllowDB   bool
	AllowLive bool
	AllowMQTT bool
}

func UpdateEventRule(ctx context.Context, db *sqlx.DB, args UpdateEventRuleArgs) error {
	var row EventRule
	err := db.GetContext(ctx, &row, `
		SELECT * FROM dahua_event_rules WHERE uuid = ?
	`, args.UUID)
	if err != nil {
		return err
	}

	row.Allow_DB = args.AllowDB
	row.Allow_Live = args.AllowLive
	row.Allow_MQTT = args.AllowMQTT
	if row.Can_Delete {
		row.Code = args.Code
	}

	_, err = db.ExecContext(ctx, `
		UPDATE dahua_event_rules SET 
			code = ?,
			allow_db = ?,
			allow_live = ?,
			allow_mqtt = ? 
		WHERE uuid = ?
		RETURNING *
	`,
		row.Code,
		row.Allow_DB,
		row.Allow_Live,
		row.Allow_MQTT,
		row.UUID,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEventRule(ctx context.Context, db *sqlx.DB, uuid string) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM dahua_event_rules WHERE uuid = ? AND can_delete IS TRUE
	`, uuid)
	return err
}

func NewEventWorker(conn Conn, db *sqlx.DB) EventWorker {
	return EventWorker{
		conn: conn,
		db:   db,
	}
}

type EventWorker struct {
	conn Conn
	db   *sqlx.DB
}

func (w EventWorker) String() string {
	return fmt.Sprintf("dahua.EventWorker(name=%s)", w.conn.Name)
}

func (w EventWorker) Serve(ctx context.Context) error {
	slog.Info("Started service", "service", w.String())
	defer slog.Info("Stopped service", "service", w.String())

	c := dahuacgi.NewClient(w.conn.IP, w.conn.Username, w.conn.Password)

	manager, err := dahuacgi.EventManagerGet(ctx, c, 0)
	if err != nil {
		var httpErr dahuacgi.HTTPError
		if errors.As(err, &httpErr) && slices.Contains([]int{
			http.StatusUnauthorized,
			http.StatusForbidden,
			http.StatusNotFound,
		}, httpErr.StatusCode) {
			slog.Error("Failed to get EventManager", slog.Any("error", err), slog.String("service", w.String()))
			return errors.Join(suture.ErrDoNotRestart, err)
		}

		return err
	}
	defer manager.Close()

	for reader := manager.Reader(); ; {
		if err := reader.Poll(); err != nil {
			return err
		}

		event, err := reader.ReadEvent()
		if err != nil {
			return err
		}

		if err := HandleEvent(ctx, w.db, w.conn.Key, event); err != nil {
			return err
		}
	}
}

func HandleEvent(ctx context.Context, db *sqlx.DB, deviceKey types.Key, event dahuacgi.Event) error {
	var eventRule struct {
		Allow_DB   bool
		Allow_Live bool
		Allow_MQTT bool
		Code       string
	}
	err := db.GetContext(ctx, &eventRule, `
		SELECT
			allow_db,
			allow_live,
			allow_mqtt,
			code
		FROM
			dahua_event_device_rules
		WHERE
			device_id = ?
			AND (
				dahua_event_device_rules.code = ?
				OR dahua_event_device_rules.code = 'All'
			)
		UNION ALL
		SELECT
			allow_db,
			allow_live,
			allow_mqtt,
			code
		FROM
			dahua_event_rules
		WHERE
			dahua_event_rules.code = ?
			OR dahua_event_rules.code = 'All'
		ORDER BY
			code DESC;
	`, deviceKey.ID, event.Code, event.Code)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	busEvent := bus.EventCreated{
		EventID:   ulid.Make().String(),
		DeviceKey: deviceKey,
		AllowDB:   eventRule.Allow_DB,
		AllowMQTT: eventRule.Allow_MQTT,
		AllowLive: eventRule.Allow_Live,
		Event:     event,
		CreatedAt: time.Now(),
	}
	if busEvent.AllowDB {
		v, err := json.MarshalIndent(busEvent.Event.Data, "", "  ")
		if err != nil {
			return err
		}
		data := types.NewJSON(v)
		createdAt := types.NewTime(busEvent.CreatedAt)
		_, err = db.ExecContext(ctx, `
			INSERT INTO dahua_events (
				id,
				device_id,
				code,
			  action,
				'index',
				data,
				created_at
			) 
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
			busEvent.EventID,
			deviceKey.ID,
			busEvent.Event.Code,
			busEvent.Event.Action,
			busEvent.Event.Index,
			data,
			createdAt,
		)
		if err != nil {
			return err
		}
	}

	bus.Publish(busEvent)

	return nil
}
