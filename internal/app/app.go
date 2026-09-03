// Package app assembles the fx application graph for hello-fx.
package app

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kirksw/test-project/hello-fx/internal/cli"
	"github.com/kirksw/test-project/hello-fx/internal/config"
	"github.com/kirksw/test-project/hello-fx/internal/greeter"
)

// New builds the hello-fx fx application.
func New() *fx.App {
	return fx.New(
		fx.Provide(
			config.New,
			NewLogger,
		),
		greeter.Module,
		cli.Module,
		// Route fx's own lifecycle events through the zap logger, demoted to
		// debug, so a normal run stays quiet and HELLO_FX_LOG_LEVEL=debug
		// brings the fx trace back for troubleshooting.
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			fxLog := &fxevent.ZapLogger{Logger: log}
			fxLog.UseLogLevel(zapcore.DebugLevel)
			return fxLog
		}),
		fx.Invoke(Run),
	)
}

// NewLogger builds the application logger. Logs are written as JSON to
// stderr so that stdout stays clean for command output.
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.LogLevel)); err != nil {
		return nil, fmt.Errorf("app: parse log level %q: %w", cfg.LogLevel, err)
	}
	return zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		os.Stderr,
		level,
	)), nil
}

// Run wires the CLI into the application lifecycle: the command executes
// during startup, and its exit code is forwarded through fx.Shutdowner so
// that stop hooks (e.g. logger sync) still run before exit.
func Run(r *cli.Runner, log *zap.Logger, shutdowner fx.Shutdowner, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			code := r.Run(os.Args[1:], os.Stdout, os.Stderr)
			return shutdowner.Shutdown(fx.ExitCode(code))
		},
		OnStop: func(context.Context) error {
			// Best effort: syncing stderr fails on some platforms (e.g. darwin).
			_ = log.Sync()
			return nil
		},
	})
}
