package app

import (
	"testing"
)

// TestGraphWires verifies that the fx application graph is valid: all
// constructors resolve, dependencies are satisfied, and the CLI invoke
// registers without error. The command itself never runs because the
// application is not started.
func TestGraphWires(t *testing.T) {
	if err := New().Err(); err != nil {
		t.Fatalf("fx graph error: %v", err)
	}
}
