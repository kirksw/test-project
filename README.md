# hello-fx

A sample Go CLI built with [Uber fx](https://github.com/uber-go/fx), wired for GitHub Actions CI with linting and testing as the pull-request gate.

## Layout

```
main.go                    entry point: build the fx app and run it
internal/app/              fx graph assembly, logger, CLI lifecycle wiring
internal/cli/              command parsing and dispatch (exit codes)
internal/config/           env-driven configuration with validation
internal/greeter/          greeting service (business logic)
```

The fx pattern used here suits one-shot CLI commands:

- the command runs in an `OnStart` lifecycle hook,
- its exit code is propagated through `fx.Shutdowner`,
- stop hooks (for example logger sync) still run before exit.

## Usage

```
hello-fx hello World            # Hello, World!
hello-fx hello Ada --shout      # HELLO, ADA!
hello-fx hello -- Grace         # name may start with "-"
hello-fx version                # dev (override with -ldflags)
hello-fx help
```

Environment:

| Variable             | Default      | Description                                    |
| -------------------- | ------------ | ---------------------------------------------- |
| `HELLO_FX_GREETING`  | `Hello, %s!` | Greeting template; must contain one `%s` verb. |
| `HELLO_FX_LOG_LEVEL` | `info`       | zap level; `debug` also enables the fx trace.  |

Exit codes: `0` success, `1` runtime or configuration failure, `2` usage error.

## Development

```
make verify   # lint + test + build (what CI runs)
make lint     # golangci-lint
make test     # go test -race with coverage profile
make build    # go build
```

Build with a version stamp:

```
go build -ldflags "-X github.com/kirksw/test-project/hello-fx/internal/cli.Version=v0.1.0" -o hello-fx .
```

## CI (PR gate)

`.github/workflows/ci.yml` runs on every pull request and on pushes to `main`.
It has two jobs:

- **Lint**: `golangci-lint-action` with the config in `.golangci.yml`.
- **Test**: `go vet`, `go test -race` with coverage, and `go build`.

To make it a hard PR gate, mark both checks (**CI / Lint**, **CI / Test**) as required in the repository's branch protection settings for `main`.
auto-merge test change
