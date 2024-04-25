package logger

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var suts = map[string]struct {
	arrange func() *bytes.Buffer
	act     func(ctx context.Context)
	assert  func(t *testing.T, b *bytes.Buffer)
}{
	"Info when Msg is invoked should write to buffer": {
		arrange: func() *bytes.Buffer {
			buff := &bytes.Buffer{}
			logger = zerolog.New(buff)
			return buff
		},
		act: func(ctx context.Context) {
			Info(ctx).Msg("info message")
		},
		assert: func(t *testing.T, b *bytes.Buffer) {
			msg := b.String()
			assert.Contains(t, msg, "\"message\":\"info message\"")
			assert.Contains(t, msg, "\"level\":\"info\"")
		},
	},
	"Warn when Msg is invoked should write to buffer": {
		arrange: func() *bytes.Buffer {
			buff := &bytes.Buffer{}
			logger = zerolog.New(buff)
			return buff
		},
		act: func(ctx context.Context) {
			Warn(ctx).Msg("warn message")
		},
		assert: func(t *testing.T, b *bytes.Buffer) {
			msg := b.String()
			assert.Contains(t, msg, "\"message\":\"warn message\"")
			assert.Contains(t, msg, "\"level\":\"warn\"")
		},
	},
	"Err when Msg is invoked should write to buffer with error field": {
		arrange: func() *bytes.Buffer {
			buff := &bytes.Buffer{}
			logger = zerolog.New(buff)
			return buff
		},
		act: func(ctx context.Context) {
			Err(ctx, errors.New("some error")).Msg("err message")
		},
		assert: func(t *testing.T, b *bytes.Buffer) {
			msg := b.String()
			assert.Contains(t, msg, "\"message\":\"err message\"")
			assert.Contains(t, msg, "\"error\":\"some error\"")
			assert.Contains(t, msg, "\"level\":\"error\"")
		},
	},
	"Error when Msg is invoked should write to buffer": {
		arrange: func() *bytes.Buffer {
			buff := &bytes.Buffer{}
			logger = zerolog.New(buff)
			return buff
		},
		act: func(ctx context.Context) {
			Error(ctx).Msg("error message")
		},
		assert: func(t *testing.T, b *bytes.Buffer) {
			msg := b.String()
			assert.Contains(t, msg, "\"message\":\"error message\"")
			assert.Contains(t, msg, "\"level\":\"error\"")
		},
	},
	"Debug when Msg is invoked should write to buffer": {
		arrange: func() *bytes.Buffer {
			buff := &bytes.Buffer{}
			logger = zerolog.New(buff)
			return buff
		},
		act: func(ctx context.Context) {
			Debug(ctx).Msg("debug message")
		},
		assert: func(t *testing.T, b *bytes.Buffer) {
			msg := b.String()
			assert.Contains(t, msg, "\"message\":\"debug message\"")
			assert.Contains(t, msg, "\"level\":\"debug\"")
		},
	},
	"Configure when adding contextual fields should have fields into log message": {
		arrange: func() *bytes.Buffer {
			buff := &bytes.Buffer{}
			logger = Configure(func(cfg *LoggerConfig) {
				cfg.WithWriter(buff)
				cfg.WithContextFields(func(c zerolog.Context) zerolog.Context {
					return c.Str("context", "value")
				})
			})
			return buff
		},
		act: func(ctx context.Context) {
			Info(ctx).Msg("contextual log")
		},
		assert: func(t *testing.T, b *bytes.Buffer) {
			assert.Contains(t, b.String(), "\"context\":\"value\"")
		},
	},
	"Configure when adding event fields should have fields into log message": {
		arrange: func() *bytes.Buffer {
			buff := &bytes.Buffer{}
			logger = Configure(func(cfg *LoggerConfig) {
				cfg.WithWriter(buff)
				cfg.WithEventFields(func(ctx context.Context, e *zerolog.Event) *zerolog.Event {
					return e.Str("trace_id", "123456")
				})
			})
			return buff
		},
		act: func(ctx context.Context) {
			Info(ctx).Msg("contextual log")
		},
		assert: func(t *testing.T, b *bytes.Buffer) {
			assert.Contains(t, b.String(), "\"trace_id\":\"123456\"")
		},
	},
}

func TestLogLevelFuncs(t *testing.T) {
	for name, sut := range suts {
		t.Run(name, func(t *testing.T) {
			b := sut.arrange()

			sut.act(context.TODO())

			sut.assert(t, b)
		})
	}
}
