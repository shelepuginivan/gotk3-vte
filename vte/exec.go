package vte

// #include <glib.h>
// #include <gtk/gtk.h>
// #include <vte/vte.h>
// #include "glib.go.h"
import "C"
import (
	"errors"
	"sync"
)

var (
	vteAsyncExecLock sync.Mutex
	vteAsyncExecMap  = make(map[uint]*Command)
)

func assignCallID(cmd *Command) uint {
	callID := uint(1)

	vteAsyncExecLock.Lock()
	defer vteAsyncExecLock.Unlock()

	for callID != 0 {
		_, exists := vteAsyncExecMap[callID]
		if !exists {
			vteAsyncExecMap[callID] = cmd
			return callID
		}
		callID++
	}
	return 0
}

//export ptySpawnAsyncCallback
func ptySpawnAsyncCallback(source *C.VtePty, res *C.GAsyncResult, cCallID C.gpointer) {
	callID := uint(C.gpointerToUint(cCallID))
	if callID == 0 {
		return
	}

	vteAsyncExecLock.Lock()

	cmd, exists := vteAsyncExecMap[callID]
	if !exists {
		vteAsyncExecLock.Unlock()
		return
	}

	delete(vteAsyncExecMap, callID)
	vteAsyncExecLock.Unlock()

	if cmd.OnSpawn != nil {
		pty := &Pty{source}
		cmd.OnSpawn(pty.spawnFinish(res))
	}
}

//export terminalSpawnAsyncCallback
func terminalSpawnAsyncCallback(_ *C.VteTerminal, pid C.GPid, gerr *C.GError, cCallID C.gpointer) {
	callID := uint(C.gpointerToUint(cCallID))
	if callID == 0 {
		return
	}

	vteAsyncExecLock.Lock()

	cmd, exists := vteAsyncExecMap[callID]
	if !exists {
		vteAsyncExecLock.Unlock()
		return
	}

	delete(vteAsyncExecMap, callID)
	vteAsyncExecLock.Unlock()

	var err error
	if gerr != nil {
		err = errors.New(goString(gerr.message))
	}
	cmd.OnSpawn(int(pid), err)
}
