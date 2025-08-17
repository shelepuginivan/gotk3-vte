# PTY

`vte.Terminal` uses a struct `vte.Pty` to execute commands under the hood. You
can use it directly, although most applications probably don't need to do so:

```go
func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal(err)
	}

	cancellable, err := glib.CancellableNew()
	if err != nil {
		log.Fatal(err)
	}

	pty, err := vte.PtyNewSync(vte.PTY_DEFAULT, cancellable)
	if err != nil {
		log.Fatal(err)
	}

	term, err := vte.TerminalNew()
	if err != nil {
		log.Fatal(err)
	}

	pty.Spawn(vte.CommandNew([]string{"/usr/bin/bash"}))
	term.SetPty(pty)

	win.Add(term)
	win.ShowAll()
	win.Connect("destroy", gtk.MainQuit)

	gtk.Main()
}
```

By default `vte.Terminal` will not handle the `child-exited` signal when
spawning commands with `vte.Pty`. You can call
[`vte.Terminal.WatchChild`](https://pkg.go.dev/github.com/shelepuginivan/gotk3-vte/vte#Terminal.WatchChild)
inside the `OnSpawn` callback, so that the terminal can handle `child-exited`
signal:

```go
pty.Spawn(vte.CommandNew(
	[]string{"/usr/bin/bash"},
	vte.CommandWithOnSpawn(func(pid int, err error) {
		if err != nil {
			log.Fatal(err)
		}

		term.WatchChild(pid)    // Tell the terminal which child to watch.
	}),
))

// ...

term.Connect("child-exited", gtk.MainQuit)
```

## Using external PTY master device

You can utilize an external PTY master device using
[`vte.PtyNewForeignSync`](https://pkg.go.dev/github.com/shelepuginivan/gotk3-vte/vte#PtyNewForeignSync):

```go
ptmx, err := os.Open("/dev/ptmx")
if err != nil {
	log.Fatal(err)
}

pty, err := vte.PtyNewForeignSync(ptmx, cancellable)
if err != nil {
	log.Fatal(err)
}
```

Note that the newly created instance of `vte.Pty` takes ownership of the PTY
master device and will close it on finalize. It must not be closed or have its
flags changed.

Instead of accessing `/dev/ptmx` directly, it is possible to use an external
package that provides PTY capabilities, e.g. https://github.com/creack/pty:

```go
package main

import (
	"log"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/shelepuginivan/gotk3-vte/vte"
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal(err)
	}

	cancellable, err := glib.CancellableNew()
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("/usr/bin/zsh")

	ptmx, err := pty.Start(cmd)
	if err != nil {
		log.Fatal(err)
	}

	p, err := vte.PtyNewForeignSync(ptmx, cancellable)
	if err != nil {
		log.Fatal(err)
	}

	term, err := vte.TerminalNew()
	if err != nil {
		log.Fatal(err)
	}

	term.SetPty(p)

	win.Add(term)
	win.ShowAll()
	win.Connect("destroy", gtk.MainQuit)

	gtk.Main()
}
```

See [`pty(7)`](https://man.archlinux.org/man/pty.7) for more information about
pseudoterminal interfaces.
