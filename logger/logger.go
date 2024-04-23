package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

var logger zerolog.Logger = CreateLoggerContext().Logger()
var cfg *LoggerConfig

// LoggerConfig holds configurations for the logger, including context and event modifiers.
type LoggerConfig struct {
	ctxFields   []LoggerContextOption // Context modifiers to add additional contextual information to each log.
	eventFields []LogEventOption      // Event modifiers to customize log events on-the-fly.
}

// WithContextFields adds default fields to the logger context.
func (lc *LoggerConfig) WithContextFields(m LoggerContextOption) {
	lc.ctxFields = append(lc.ctxFields, m)
}

// WithEventFields adds a fields to the log event.
func (lc *LoggerConfig) WithEventFields(m LogEventOption) {
	lc.eventFields = append(lc.eventFields, m)
}

// LoggerOption represents a function that modifies LoggerConfig.
type LoggerOption func(l *LoggerConfig)

// LoggerContextOption represents a function that modifies zerolog.Context for additional contextual logging setup.
type LoggerContextOption func(c zerolog.Context) zerolog.Context

// LogEventOption represents a function that modifies a logging event, allowing dynamic changes to the log output.
type LogEventOption func(ctx context.Context, e *zerolog.Event) *zerolog.Event

// CreateLoggerContext initializes a zerolog.Context with standard fields and applies any provided LoggerContextOptions.
// This function is typically used to set up a base logger from which all other loggers inherit.
//
// Example usage:
//
//	logContext := logger.CreateLoggerContext(
//		func(c zerolog.Context) zerolog.Context {
//			return c.Str("service", "myService")
//		},
//	)
//
// Returns:
//
//	zerolog.Context: A configured context for logging.
func CreateLoggerContext(opts ...LoggerContextOption) zerolog.Context {
	logCtx := zerolog.New(os.Stdout).With().Caller().Timestamp()

	for _, opt := range opts {
		logCtx = opt(logCtx)
	}

	return logCtx
}

// Enrich configures the global logger with specified LoggerOptions which can modify both context and event behaviors.
// This function initializes the logger configuration and applies the options to set up context and event modifiers.
//
// Example usage:
//
//	logger.Enrich(
//		func(l *logger.LoggerConfig) {
//			l.WithContextFields(func(c zerolog.Context) zerolog.Context {
//				return c.Str("version", "1.0")
//			})
//			l.WithEventFields(func(ctx context.Context, e *zerolog.Event) *zerolog.Event {
//				return e.Str("session", getSessionID(ctx))
//			})
//		},
//	)
//
// Returns:
//
//	zerolog.Logger: The configured logger instance.
func Enrich(opts ...LoggerOption) zerolog.Logger {
	cfg = &LoggerConfig{}

	for _, opt := range opts {
		opt(cfg)
	}

	logger = CreateLoggerContext(cfg.ctxFields...).Logger()

	return logger
}

// Info starts a new logging event at the "info" level.
// This function uses a context.Context to extract necessary tracing information.
// It returns a *zerolog.Event that is not sent until the Msg method is called.
//
// Example usage:
//
//	logger.Info(ctx).Msg("This is an info level log message")
//
// Params:
//
//	ctx (context.Context): The context from which to extract tracing information.
//
// Returns:
//
//	*zerolog.Event: A pointer to the log event. Ensure to call Msg to emit the log.
func Info(ctx context.Context) *zerolog.Event {
	e := logger.Info().Ctx(ctx)

	return event(ctx, e)
}

// Warn starts a new logging event at the "warn" level.
// This function uses a context.Context to extract necessary tracing information.
// It returns a *zerolog.Event that is not sent until the Msg method is called.
//
// Example usage:
//
//	logger.Warn(ctx).Msg("This is an warn level log message")
//
// Params:
//
//	ctx (context.Context): The context from which to extract tracing information.
//
// Returns:
//
//	*zerolog.Event: A pointer to the log event. Ensure to call Msg to emit the log.
func Warn(ctx context.Context) *zerolog.Event {
	e := logger.Warn().Ctx(ctx)

	return event(ctx, e)
}

// Err initializes a new logging event at the "error" level with err as field if not nil or with "info" level if err is nil.
// This function requires a context.Context to extract necessary tracing information
// and an error which will be logged. It returns a *zerolog.Event that is not sent
// until the Msg method is called.
//
// Example usage:
//
//	logger.Err(ctx, err).Msg("This is an error level log message")
//
// Params:
//
//	ctx (context.Context): The context from which to extract tracing information.
//	err (error): The error to log.
//
// Returns:
//
//	*zerolog.Event: A pointer to the log event. Ensure to call Msg to emit the log.
func Err(ctx context.Context, err error) *zerolog.Event {
	e := logger.Err(err).Ctx(ctx)

	return event(ctx, e)
}

// Error starts a new logging event at the "error" level.
// This function uses a context.Context to extract necessary tracing information.
// It returns a *zerolog.Event that is not sent until the Msg method is called.
//
// Example usage:
//
//	logger.Error(ctx).Msg("This is an error level log message")
//
// Params:
//
//	ctx (context.Context): The context from which to extract tracing information.
//
// Returns:
//
//	*zerolog.Event: A pointer to the log event. Ensure to call Msg to emit the log.
func Error(ctx context.Context) *zerolog.Event {
	e := logger.Error().Ctx(ctx)

	return event(ctx, e)
}

// Debug starts a new logging event at the "debug" level.
// This function uses a context.Context to extract necessary tracing information.
// It returns a *zerolog.Event that is not sent until the Msg method is called.
//
// Example usage:
//
//	logger.Debug(ctx).Msg("This is an debug level log message")
//
// Params:
//
//	ctx (context.Context): The context from which to extract tracing information.
//
// Returns:
//
//	*zerolog.Event: A pointer to the log event. Ensure to call Msg to emit the log.
func Debug(ctx context.Context) *zerolog.Event {
	e := logger.Debug().Ctx(ctx)

	return event(ctx, e)
}

// Fatal starts a new logging event at the "fatal" level.
// This function uses a context.Context to extract necessary tracing information.
// It returns a *zerolog.Event that is not sent until the Msg method is called.
// The os.Exit(1) function is called by the Msg method, which terminates the program immediately.
//
// Example usage:
//
//	logger.Fatal(ctx).Msg("This is an fatal level log message")
//
// Params:
//
//	ctx (context.Context): The context from which to extract tracing information.
//
// Returns:
//
//	*zerolog.Event: A pointer to the log event. Ensure to call Msg to emit the log.
func Fatal(ctx context.Context) *zerolog.Event {
	e := logger.Fatal().Ctx(ctx)

	return event(ctx, e)
}

func event(ctx context.Context, event *zerolog.Event) *zerolog.Event {
	for _, opt := range cfg.eventFields {
		event = opt(ctx, event)
	}
	return event
}
