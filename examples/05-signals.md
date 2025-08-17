# Signals

Signals allow to execute callbacks whenever an event happens.
We have already used a signal `child-exited` in previous examples:

```go
package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/shelepuginivan/gotk3-vte/vte"
)

func main() {
	// ...

	term.Connect("child-exited", gtk.MainQuit)

	// ...
}
```

This signal runs callback when child spawned with
[`Terminal.Spawn`](https://pkg.go.dev/github.com/shelepuginivan/gotk3-vte/vte#Terminal.Spawn).
In our case, it finishes the GTK main loop.

However, for more complex signals it is recommended to use custom methods, such as
[`Terminal.ConnectTermPropChanged`](https://pkg.go.dev/github.com/shelepuginivan/gotk3-vte/vte#Terminal.ConnectTermPropChanged):

```go
package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/shelepuginivan/gotk3-vte/vte"
)

func main() {
	// ...

	term.ConnectTermPropChanged(func(t *vte.Terminal, prop vte.TermProp) {
		value, err := t.GetTermProp(prop)
		if err != nil {
			return
		}

		fmt.Printf("TermProp '%s' changed, new value: '%v'\n", prop, value)
	})

	// ...
}
```

Unlike `Connect`, these functions provide a type-safe interface for callbacks.

If usage of `Connect` is required, remember that the first argument of the
callback is not `*vte.Terminal`, but `*glib.Object`, since gotk3 does not know
about `vte` and cannot cast the underlying GObject to an appropriate type. You
can use `vte.WrapTerminal` to convert `*glib.Object` to `*vte.Terminal`:

```go
term.Connect("button-press-event", func(obj *glib.Object, ev *gdk.Event) {
	t := vte.WrapTerminal(obj)

	// rest of the logic...
})
```

See [gotk3 documentation](https://pkg.go.dev/github.com/gotk3/gotk3/glib#Object.Connect)
for more information about signal handling.
