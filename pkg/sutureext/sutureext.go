package sutureext

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func EventHook() suture.EventHook {
	var prevTerminate suture.EventServiceTerminate
	return func(ei suture.Event) {
		switch e := ei.(type) {
		case suture.EventStopTimeout:
			log.Info().Str("supervisor", e.SupervisorName).Str("service", e.ServiceName).Msg("Service failed to terminate in a timely manner")
		case suture.EventServicePanic:
			log.Warn().Msg("Caught a service panic, which shouldn't happen")
			logJSON(log.Info(), e)
		case suture.EventServiceTerminate:
			if e.ServiceName == prevTerminate.ServiceName && e.Err == prevTerminate.Err {
				log.Debug().Str("supervisor", e.SupervisorName).Str("service", e.ServiceName).Str("err", fmt.Sprint(e.Err)).Msg("Service failed")
			} else {
				log.Info().Str("supervisor", e.SupervisorName).Str("service", e.ServiceName).Str("err", fmt.Sprint(e.Err)).Msg("Service failed")
			}
			prevTerminate = e
			logJSON(log.Debug(), e)
		case suture.EventBackoff:
			log.Debug().Str("supervisor", e.SupervisorName).Msg("Exiting backoff state")
		case suture.EventResume:
			log.Debug().Str("supervisor", e.SupervisorName).Msg("Too many service failures - entering the backoff state")
		default:
			log.Warn().Int("type", int(e.Type())).Msg("Unknown suture supervisor event type")
			logJSON(log.Info(), e)
		}
	}
}

func logJSON(event *zerolog.Event, v any) {
	b, _ := json.Marshal(v)
	event.Msg(string(b))
}

type ServiceFunc struct {
	fn func(ctx context.Context) error
}

func NewServiceFunc(fn func(ctx context.Context) error) ServiceFunc {
	return ServiceFunc{
		fn: fn,
	}
}

func (s ServiceFunc) Serve(ctx context.Context) error {
	return s.fn(ctx)
}
