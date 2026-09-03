// Package cli implements the hello-fx command-line interface.
package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"go.uber.org/fx"

	"github.com/kirksw/test-project/hello-fx/internal/greeter"
)

// Version is the CLI version reported by the version command.
// It is overridden at build time with:
//
//	-ldflags "-X github.com/kirksw/test-project/hello-fx/internal/cli.Version=v1.2.3"
var Version = "dev"

// Exit codes returned by Run.
const (
	ExitSuccess = 0
	ExitFailure = 1
	ExitUsage   = 2
)

// Runner executes CLI commands against the fx-provided service graph.
type Runner struct {
	greeter *greeter.Greeter
}

// Module wires the cli package into the fx application graph.
var Module = fx.Module("cli",
	fx.Provide(New),
)

// New constructs a Runner.
func New(g *greeter.Greeter) *Runner {
	return &Runner{greeter: g}
}

// Run dispatches args to a command, writes output to stdout and stderr, and
// returns the process exit code.
func (r *Runner) Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return ExitUsage
	}

	switch cmd, rest := args[0], args[1:]; cmd {
	case "hello":
		return r.runHello(rest, stdout, stderr)
	case "version":
		fmt.Fprintln(stdout, Version)
		return ExitSuccess
	case "help", "-h", "--help":
		usage(stdout)
		return ExitSuccess
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", cmd)
		usage(stderr)
		return ExitUsage
	}
}

// runHello implements the hello command: hello <name> [--shout].
// Flags may appear before or after the name.
func (r *Runner) runHello(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("hello", flag.ContinueOnError)
	fs.SetOutput(stderr)
	shout := fs.Bool("shout", false, "print the greeting in upper case")

	if err := fs.Parse(partitionFlags(args)); err != nil {
		return ExitUsage
	}

	name := strings.Join(fs.Args(), " ")
	if name == "" {
		fmt.Fprintln(stderr, `hello requires a name, e.g. "hello-fx hello World"`)
		usage(stderr)
		return ExitUsage
	}

	greeting := r.greeter.Greet(name)
	if *shout {
		greeting = strings.ToUpper(greeting)
	}
	fmt.Fprintln(stdout, greeting)
	return ExitSuccess
}

// partitionFlags reorders args into flag arguments followed by positional
// arguments, so that both "--shout World" and "World --shout" parse.
// It assumes all flags are boolean, so no flag consumes a following value.
// A bare "--" terminator is kept: everything after it stays positional.
func partitionFlags(args []string) []string {
	flags := make([]string, 0, len(args))
	rest := make([]string, 0, len(args))
	for i, a := range args {
		switch {
		case a == "--":
			return append(append(flags, "--"), args[i+1:]...)
		case a != "-" && strings.HasPrefix(a, "-"):
			flags = append(flags, a)
		default:
			rest = append(rest, a)
		}
	}
	return append(flags, rest...)
}

// usageTemplate is the CLI usage summary, written verbatim (not through a
// printf-style function).
const usageTemplate = `Usage: hello-fx <command> [arguments]

Commands:
  hello <name> [--shout]   greet someone by name
  version                  print the CLI version
  help                     show this help

Environment:
  HELLO_FX_GREETING   greeting template, must contain one %s verb (default "Hello, %s!")
  HELLO_FX_LOG_LEVEL  log level for the application logger (default "info")
`

// usage prints the CLI usage summary.
func usage(w io.Writer) {
	// Best effort: there is nothing actionable if writing usage text fails.
	_, _ = io.WriteString(w, usageTemplate)
}
