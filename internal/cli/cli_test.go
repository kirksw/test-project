package cli

import (
	"bytes"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/kirksw/test-project/hello-fx/internal/config"
	"github.com/kirksw/test-project/hello-fx/internal/greeter"
)

func newRunner(greeting string) *Runner {
	return New(greeter.New(&config.Config{Greeting: greeting}, zap.NewNop()))
}

// run executes the CLI with args and returns stdout, stderr, and exit code.
func run(r *Runner, args ...string) (string, string, int) {
	var stdout, stderr bytes.Buffer
	code := r.Run(args, &stdout, &stderr)
	return stdout.String(), stderr.String(), code
}

func TestRunHello(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantOut string
		wantErr string
		wantCod int
	}{
		{
			name:    "plain",
			args:    []string{"hello", "World"},
			wantOut: "Hello, World!\n",
			wantCod: ExitSuccess,
		},
		{
			name:    "multi-word name",
			args:    []string{"hello", "Grace", "Hopper"},
			wantOut: "Hello, Grace Hopper!\n",
			wantCod: ExitSuccess,
		},
		{
			name:    "shout",
			args:    []string{"hello", "--shout", "World"},
			wantOut: "HELLO, WORLD!\n",
			wantCod: ExitSuccess,
		},
		{
			name:    "flag after positional args",
			args:    []string{"hello", "Ada", "Lovelace", "--shout"},
			wantOut: "HELLO, ADA LOVELACE!\n",
			wantCod: ExitSuccess,
		},
		{
			name:    "terminator keeps flag-like name positional",
			args:    []string{"hello", "--", "--shout"},
			wantOut: "Hello, --shout!\n",
			wantCod: ExitSuccess,
		},
		{
			name:    "missing name",
			args:    []string{"hello"},
			wantErr: "hello requires a name",
			wantCod: ExitUsage,
		},
		{
			name:    "unknown flag",
			args:    []string{"hello", "--nope", "World"},
			wantErr: "flag provided but not defined",
			wantCod: ExitUsage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, code := run(newRunner("Hello, %s!"), tt.args...)
			if code != tt.wantCod {
				t.Errorf("exit code = %d, want %d", code, tt.wantCod)
			}
			if stdout != tt.wantOut {
				t.Errorf("stdout = %q, want %q", stdout, tt.wantOut)
			}
			if tt.wantErr != "" && !strings.Contains(stderr, tt.wantErr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr, tt.wantErr)
			}
		})
	}
}

func TestRunCustomGreeting(t *testing.T) {
	stdout, _, code := run(newRunner("Hi %s, welcome!"), "hello", "Ada")
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want %d", code, ExitSuccess)
	}
	if stdout != "Hi Ada, welcome!\n" {
		t.Errorf("stdout = %q, want %q", stdout, "Hi Ada, welcome!\n")
	}
}

func TestRunVersion(t *testing.T) {
	stdout, stderr, code := run(newRunner("Hello, %s!"), "version")
	if code != ExitSuccess {
		t.Errorf("exit code = %d, want %d", code, ExitSuccess)
	}
	if stdout != Version+"\n" {
		t.Errorf("stdout = %q, want %q", stdout, Version+"\n")
	}
	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}
}

func TestRunHelp(t *testing.T) {
	for _, cmd := range []string{"help", "-h", "--help"} {
		t.Run(cmd, func(t *testing.T) {
			stdout, _, code := run(newRunner("Hello, %s!"), cmd)
			if code != ExitSuccess {
				t.Errorf("exit code = %d, want %d", code, ExitSuccess)
			}
			if !strings.Contains(stdout, "Usage: hello-fx") {
				t.Errorf("stdout = %q, want usage text", stdout)
			}
		})
	}
}

func TestRunNoArgs(t *testing.T) {
	stdout, stderr, code := run(newRunner("Hello, %s!"))
	if code != ExitUsage {
		t.Errorf("exit code = %d, want %d", code, ExitUsage)
	}
	if stdout != "" {
		t.Errorf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "Usage: hello-fx") {
		t.Errorf("stderr = %q, want usage text", stderr)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	stdout, stderr, code := run(newRunner("Hello, %s!"), "frobnicate")
	if code != ExitUsage {
		t.Errorf("exit code = %d, want %d", code, ExitUsage)
	}
	if stdout != "" {
		t.Errorf("stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, `unknown command "frobnicate"`) {
		t.Errorf("stderr = %q, want unknown-command message", stderr)
	}
}
