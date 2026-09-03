package greeter

import (
	"testing"

	"go.uber.org/zap"

	"github.com/kirksw/test-project/hello-fx/internal/config"
)

func TestGreet(t *testing.T) {
	tests := []struct {
		name      string
		greeting  string
		greetName string
		want      string
	}{
		{
			name:      "default greeting",
			greeting:  "Hello, %s!",
			greetName: "World",
			want:      "Hello, World!",
		},
		{
			name:      "custom greeting",
			greeting:  "Hi %s, welcome!",
			greetName: "Ada",
			want:      "Hi Ada, welcome!",
		},
		{
			name:      "verb inside word",
			greeting:  "Hey%shey",
			greetName: "There",
			want:      "HeyTherehey",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New(&config.Config{Greeting: tt.greeting}, zap.NewNop())
			if got := g.Greet(tt.greetName); got != tt.want {
				t.Errorf("Greet(%q) = %q, want %q", tt.greetName, got, tt.want)
			}
		})
	}
}
