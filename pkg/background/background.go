package background

import "context"

type Background interface {
	Background(ctx context.Context, doneC chan<- struct{})
}

func Run(ctx context.Context, backgrounds ...Background) <-chan struct{} {
	doneFanIn := make(chan struct{}, len(backgrounds))
	running := 0

	// Start backgrounds
	for _, background := range backgrounds {
		go background.Background(ctx, doneFanIn)
		running++
	}

	done := make(chan struct{})

	go func() {
		// Wait for context
		<-ctx.Done()

		// Wait for backgrounds
		for i := 0; i < running; i++ {
			<-doneFanIn
		}

		close(done)
	}()

	return done
}

type BackgroundFunction func(ctx context.Context)

type Function struct {
	fn       BackgroundFunction
	blocking Blocking
}

type Blocking int

const (
	BlockingForever Blocking = iota // BlockingForever blocks forever.
	BlockingContext                 // BlockingContext blocks until context is done.
)

// NewFunction create an adhoc function that implements the Background interface.
func NewFunction(blocking Blocking, fn BackgroundFunction) Function {
	return Function{
		fn:       fn,
		blocking: blocking,
	}
}

func (f Function) Background(ctx context.Context, doneC chan<- struct{}) {
	switch f.blocking {
	case BlockingForever:
		doneC <- struct{}{}
		f.fn(ctx)
	case BlockingContext:
		f.fn(ctx)
		doneC <- struct{}{}
	default:
		panic("unhandled case")
	}
}
