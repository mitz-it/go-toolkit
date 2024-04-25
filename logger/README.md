<!-- markdownlint-disable MD010 -->
# Go Toolkit - Logger

A default set of configurations on top of [zerolog](https://github.com/rs/zerolog).
We use these abstractions to force Go context propagation to the logger instance to create external providers for adding tracing information into the log events.

## Installation

```bash
go get -u github.com/mitz-it/go-toolkit/logger
```

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/mitz-it/go-toolkit/logger"
	"github.com/rs/zerolog"
)


func main() {
	// logger output default is os.Stdout
	logger.Configure(
		func(cfg *logger.LoggerConfig) {
			cfg.WithContextFields(func(c zerolog.Context) zerolog.Context {
				return c.Str("version", "1.0.0")
			})
			cfg.WithEventFields(func(ctx context.Context, e *zerolog.Event) *zerolog.Event {
				return e.Str("trace_id", getTraceID(ctx))
			}),
		}
	)

	ctx := context.Background()

	logger.Info(ctx).Msg("info level log message")

	logger.Warn(ctx).Msg("warn level log message")

	err := fmt.Errorf("some error")
	logger.Err(ctx, err).Msg("error level log message")

	logger.Error(ctx).Msg("error level log message")

	logger.Debug(ctx).Msg("debug level log message")

	logger.Fatal(ctx).Msg("fatal level log message")
}
```
