package vte_test

import (
	"os"
	"strings"
	"testing"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/shelepuginivan/gotk3-vte/vte"
	"github.com/stretchr/testify/assert"
)

func newPty(t *testing.T) *vte.Pty {
	cancellable, err := glib.CancellableNew()
	assert.NoError(t, err)

	pty, err := vte.PtyNewSync(vte.PTY_DEFAULT, cancellable)
	assert.NoError(t, err)

	return pty
}

func TestPtyNewSync(t *testing.T) {
	pty := newPty(t)
	assert.NotNil(t, pty)
	assert.NotEqual(t, uintptr(unsafe.Pointer(nil)), pty.Native())

	t.Run("GCancellable as nil pointer", func(t *testing.T) {
		pty, err := vte.PtyNewSync(vte.PTY_DEFAULT, nil)
		assert.Nil(t, pty)
		assert.Error(t, err)
	})
}

func TestPtyNewForeignSync(t *testing.T) {
	cancellable, err := glib.CancellableNew()
	assert.NoError(t, err)

	ptmx, err := os.Open("/dev/ptmx")
	assert.NoError(t, err)
	defer ptmx.Close()

	pty, err := vte.PtyNewForeignSync(ptmx, cancellable)
	assert.NoError(t, err)

	assert.Equal(t, ptmx.Fd(), pty.GetFd())

	t.Run("GCancellable as nil pointer", func(t *testing.T) {
		ptmx, err := os.Open("/dev/ptmx")
		assert.NoError(t, err)
		defer ptmx.Close()

		pty, err := vte.PtyNewForeignSync(ptmx, nil)
		assert.Nil(t, pty)
		assert.Error(t, err)
	})

	t.Run("Invalid ptmx file descriptor", func(t *testing.T) {
		cancellable, err := glib.CancellableNew()
		assert.NoError(t, err)

		file, err := os.CreateTemp(t.TempDir(), "*")
		assert.NoError(t, err)
		defer file.Close()

		pty, err := vte.PtyNewForeignSync(file, cancellable)
		assert.Nil(t, pty)
		assert.Error(t, err)
	})
}

func TestPty_GetFd(t *testing.T) {
	pty := newPty(t)

	pty.GetFd()

	file := os.NewFile(pty.GetFd(), "")
	info, err := file.Stat()
	assert.NoError(t, err)

	// Check that file descriptor points to a device.
	assert.True(t, strings.HasPrefix(info.Mode().String(), "Dc"))
}

func TestPty_Size(t *testing.T) {
	pty := newPty(t)

	size, err := pty.GetSize()
	assert.NoError(t, err)
	assert.Equal(t, 0, size.Rows)
	assert.Equal(t, 0, size.Columns)

	err = pty.SetSize(&vte.PtySize{
		Rows:    20,
		Columns: 80,
	})
	assert.NoError(t, err)

	size, err = pty.GetSize()
	assert.NoError(t, err)
	assert.Equal(t, 20, size.Rows)
	assert.Equal(t, 80, size.Columns)
}

func TestPty_SetUTF8(t *testing.T) {
	pty := newPty(t)

	assert.NoError(t, pty.SetUTF8(true))
	assert.NoError(t, pty.SetUTF8(false))
}

func TestPty_Spawn(t *testing.T) {
	gtk.Init(nil)

	cmd := vte.CommandNew(
		[]string{"/usr/bin/true"},
		vte.CommandWithOnSpawn(func(pid int, err error) {
			assert.NoError(t, err)
			assert.Greater(t, pid, os.Getpid())
			gtk.MainQuit()
		}),
	)

	newPty(t).Spawn(cmd)

	// This will block. Unless gtk.MainQuit is called in the OnSpawn callback, the
	// test will timeout after 10 minutes.
	gtk.Main()
}
