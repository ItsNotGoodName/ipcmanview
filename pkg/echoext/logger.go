package echoext

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type (
	// LoggerConfig defines the config for Logger middleware.
	LoggerConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper

		// Tags to construct the logger format.
		//
		// - time_unix
		// - time_unix_milli
		// - time_unix_micro
		// - time_unix_nano
		// - time_rfc3339
		// - time_rfc3339_nano
		// - time_custom
		// - id (Request ID)
		// - remote_ip
		// - uri
		// - host
		// - method
		// - path
		// - route
		// - protocol
		// - referer
		// - user_agent
		// - status
		// - error
		// - latency (In nanoseconds)
		// - latency_human (Human readable)
		// - bytes_in (Bytes received)
		// - bytes_out (Bytes sent)
		// - header:<NAME>
		// - query:<NAME>
		// - form:<NAME>
		// - custom (see CustomTagFunc field)
		//
		// Example []string{"remote_ip", "status"}
		//
		// Optional. Default value DefaultLoggerConfig.Format.
		Format []string

		// Optional. Default value DefaultLoggerConfig.CustomTimeFormat.
		CustomTimeFormat string

		Logger *zerolog.Logger
	}
)

var (
	// DefaultLoggerConfig is the default Logger middleware config.
	DefaultLoggerConfig = LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: []string{
			"time",
			"time_rfc3339_nano",
			"id",
			"remote_ip",
			"host",
			"method",
			"uri",
			"user_agent",
			"status",
			"error",
			"latency",
			"latency_human",
			"bytes_in",
			"bytes_out",
		},
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Logger:           &log.Logger,
	}
)

// Logger returns a middleware that logs HTTP requests.
func Logger() echo.MiddlewareFunc {
	return LoggerWithConfig(DefaultLoggerConfig)
}

// LoggerWithConfig returns a Logger middleware with config.
// See: `Logger()`.
func LoggerWithConfig(config LoggerConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultLoggerConfig.Skipper
	}
	if config.Format == nil {
		config.Format = DefaultLoggerConfig.Format
	}
	if config.Logger == nil {
		config.Logger = DefaultLoggerConfig.Logger
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			stop := time.Now()

			log := func() *zerolog.Event {
				if err == nil {
					return config.Logger.Info()
				}

				var e *echo.HTTPError
				if errors.As(err, &e) {
					if e.Internal == nil {
						return config.Logger.Error()
					}
					return config.Logger.Err(e.Internal)
				}

				return config.Logger.Err(err)
			}()

			fn := func(log *zerolog.Event, tag string) (*zerolog.Event, error) {
				switch tag {
				// case "custom":
				// 	if config.CustomTagFunc == nil {
				// 		return 0, nil
				// 	}
				// 	return config.CustomTagFunc(c, buf)
				case "time_unix":
					return log.Str("time_unix", strconv.FormatInt(time.Now().Unix(), 10)), nil
				case "time_unix_milli":
					// go 1.17 or later, it supports time#UnixMilli()
					return log.Str("time_unix_milli", strconv.FormatInt(time.Now().UnixNano()/1000000, 10)), nil
				case "time_unix_micro":
					// go 1.17 or later, it supports time#UnixMicro()
					return log.Str("time_unix_micro", strconv.FormatInt(time.Now().UnixNano()/1000, 10)), nil
				case "time_unix_nano":
					return log.Str("time_unix_nano", strconv.FormatInt(time.Now().UnixNano(), 10)), nil
				case "time_rfc3339":
					return log.Str("time_rfc3339", time.Now().Format(time.RFC3339)), nil
				case "time_rfc3339_nano":
					return log.Str("time_rfc3339_nano", time.Now().Format(time.RFC3339Nano)), nil
				case "time_custom":
					return log.Str("time_custom", time.Now().Format(config.CustomTimeFormat)), nil
				case "id":
					id := req.Header.Get(echo.HeaderXRequestID)
					if id == "" {
						id = res.Header().Get(echo.HeaderXRequestID)
					}
					return log.Str("id", id), nil
				case "remote_ip":
					return log.Str("remote_ip", c.RealIP()), nil
				case "host":
					return log.Str("host", req.Host), nil
				case "uri":
					return log.Str("uri", req.RequestURI), nil
				case "method":
					return log.Str("method", req.Method), nil
				case "path":
					p := req.URL.Path
					if p == "" {
						p = "/"
					}
					return log.Str("path", p), nil
				case "route":
					return log.Str("route", c.Path()), nil
				case "protocol":
					return log.Str("protocol", req.Proto), nil
				case "referer":
					return log.Str("referer", req.Referer()), nil
				case "user_agent":
					return log.Str("user_agent", req.UserAgent()), nil
				case "status":
					// n := res.Status
					// s := config.colorer.Green(n)
					// switch {
					// case n >= 500:
					// 	s = config.colorer.Red(n)
					// case n >= 400:
					// 	s = config.colorer.Yellow(n)
					// case n >= 300:
					// 	s = config.colorer.Cyan(n)
					// }
					return log.Str("status", strconv.Itoa(res.Status)), nil
				case "error":
					if err != nil {
						// Error may contain invalid JSON e.g. `"`
						b, _ := json.Marshal(err.Error())
						b = b[1 : len(b)-1]
						return log.RawJSON("error", b), nil
					}
				case "latency":
					l := stop.Sub(start)
					return log.Str("latency", strconv.FormatInt(int64(l), 10)), nil
				case "latency_human":
					return log.Str("latency_human", stop.Sub(start).String()), nil
				case "bytes_in":
					cl := req.Header.Get(echo.HeaderContentLength)
					if cl == "" {
						cl = "0"
					}
					return log.Str("bytes_in", cl), nil
				case "bytes_out":
					return log.Str("bytes_out", strconv.FormatInt(res.Size, 10)), nil
				default:
					switch {
					case strings.HasPrefix(tag, "header:"):
						return log.Str("header:", c.Request().Header.Get(tag[7:])), nil
					case strings.HasPrefix(tag, "query:"):
						return log.Str("query:", c.QueryParam(tag[6:])), nil
					case strings.HasPrefix(tag, "form:"):
						return log.Str("form:", c.FormValue(tag[5:])), nil
					case strings.HasPrefix(tag, "cookie:"):
						cookie, err := c.Cookie(tag[7:])
						if err == nil {
							return log.Str("cookie:", cookie.Value), nil
						}
					}
				}
				return log, nil
			}

			for _, tag := range config.Format {
				log, err = fn(log, tag)
				if err != nil {
					return
				}
			}

			log.Msg(req.RequestURI)
			return
		}
	}
}
