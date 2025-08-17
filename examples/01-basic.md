# Basic example

```go
package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
	"github.com/shelepuginivan/gotk3-vte/vte"
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal(err)
	}

	term, err := vte.TerminalNew()
	if err != nil {
		log.Fatal(err)
	}

	cmd := vte.CommandNew([]string{"/usr/bin/bash"})
	term.Spawn(cmd)
	term.Connect("child-exited", gtk.MainQuit)

	win.Add(term)
	win.ShowAll()
	win.Connect("destroy", gtk.MainQuit)

	gtk.Main()
}
```

After running the above code, a simple terminal window should appear:

![Basic terminal window](./img/01-basic.webp)
