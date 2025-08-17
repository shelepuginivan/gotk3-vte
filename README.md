# gotk3-vte

[![Go Reference](https://pkg.go.dev/badge/github.com/shelepuginivan/gotk3-vte.svg)](https://pkg.go.dev/github.com/shelepuginivan/gotk3-vte/vte)
[![Go Report Card](https://goreportcard.com/badge/github.com/shelepuginivan/gotk3-vte)](https://goreportcard.com/report/github.com/shelepuginivan/gotk3-vte)
[![License: MIT](https://img.shields.io/badge/License-MIT-00cc00.svg)](https://github.com/shelepuginivan/gotk3-vte/blob/main/LICENSE)

Package `gotk3-vte` provides [gotk3](https://github.com/gotk3/gotk3)-compatible
bindings for [VTE](https://gitlab.gnome.org/GNOME/vte).

Most of the library API is implemented, except for deprecated features. Some
parts of the API is modified to be more idiomatic.

## Installation

```sh
go get -u github.com/shelepuginivan/gotk3-vte
```

## Example usage

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

## Documentation

See [examples](./examples) and
[API reference](https://pkg.go.dev/github.com/shelepuginivan/gotk3-vte)

## License

[MIT](./LICENSE).
