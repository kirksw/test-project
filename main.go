// Command hello-fx is a sample Go CLI built on the Uber fx dependency
// injection framework.
//
// It demonstrates the wiring pattern used for one-shot CLI commands with fx:
// the command runs in an OnStart lifecycle hook and its exit code is
// propagated through fx.Shutdowner so that stop hooks still run before exit.
package main

import (
	"github.com/kirksw/test-project/hello-fx/internal/app"
)

func main() {
	app.New().Run()
}
