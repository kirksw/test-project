// Package greeter implements the greeting service used by the CLI.
package greeter

import (
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/kirksw/test-project/hello-fx/internal/config"
)

// Greeter builds greetings from a configurable template.
type Greeter struct {
	log      *zap.Logger
	greeting string
}

// Module wires the greeter package into the fx application graph.
var Module = fx.Module("greeter",
	fx.Provide(New),
)

// New constructs a Greeter from configuration and a logger.
func New(cfg *config.Config, log *zap.Logger) *Greeter {
	return &Greeter{
		log:      log,
		greeting: cfg.Greeting,
	}
}

// Greet returns the greeting for name.
func (g *Greeter) Greet(name string) string {
	g.log.Info("building greeting", zap.String("name", name))
	return strings.Replace(g.greeting, "%s", name, 1)
}
