package vte

// TermProp is an enumeration type that represents terminal property names.
// TermProp values can be accessed with [Terminal.GetTermProp].
type TermProp string

const (
	// Current directory URI as set by OSC 7.
	TERMPROP_CURRENT_DIRECTORY_URI TermProp = "vte.cwd"

	// Current file URI as set by OSC 6.
	TERMPROP_CURRENT_FILE_URI TermProp = "vte.cwf"

	// Xterm window title as set by OSC 0 and OSC 2.
	TERMPROP_XTERM_TITLE TermProp = "xterm.title"

	// Name of the container.
	TERMPROP_CONTAINER_NAME TermProp = "vte.container.name"

	// Runtime of the container.
	TERMPROP_CONTAINER_RUNTIME TermProp = "vte.container.runtime"

	// User ID of the container.
	TERMPROP_CONTAINER_UID TermProp = "vte.container.uid"

	// Signals that the shell is going to prompt.
	TERMPROP_SHELL_PRECMD TermProp = "vte.shell.precmd"

	// Shell is preparing to execute the command entered at the prompt.
	TERMPROP_SHELL_PREEXEC TermProp = "vte.shell.preexec"

	// Signals that the shell has executed the commands entered at the prompt and
	// these commands have returned. The termprop value is the exit code.
	TERMPROP_SHELL_POSTEXEC TermProp = "vte.shell.postexec"
)
