// Package vte provides Go bindings for Vte.
package vte

// #cgo pkg-config: gtk+-3.0 vte-2.91
//
// #include <cairo.h>
// #include <gtk/gtk.h>
// #include <pango/pango.h>
// #include <vte/vte.h>
//
// #include "exec.go.h"
// #include "glib.go.h"
// #include "gtk.go.h"
// #include "vte.go.h"
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

// Terminal is a wrapper around VteTerminal.
type Terminal struct {
	gtk.Widget

	ptr *C.VteTerminal
}

// TerminalNew creates a new instance of [*Terminal].
func TerminalNew() (*Terminal, error) {
	ptr := C.vte_terminal_new()
	if ptr == nil {
		return nil, fmt.Errorf("vte_terminal_new returned nil pointer")
	}

	cGObject := glib.ToGObject(unsafe.Pointer(ptr))

	gObject := glib.Object{
		GObject: cGObject,
	}

	initiallyUnowned := glib.InitiallyUnowned{
		Object: &gObject,
	}

	widget := gtk.Widget{
		InitiallyUnowned: initiallyUnowned,
	}

	return &Terminal{
		Widget: widget,

		ptr: C.toTerminal(unsafe.Pointer(ptr)),
	}, nil
}

// CopyClipboardFormat copies selected text to clipboard in the specified
// format.
func (t *Terminal) CopyClipboardFormat(format Format) {
	C.vte_terminal_copy_clipboard_format(t.ptr, C.VteFormat(format))
}

// CopyPrimary copies selected text in the primary selection.
func (t *Terminal) CopyPrimary() {
	C.vte_terminal_copy_primary(t.ptr)
}

// PasteClipboard pastes contents of clipboard to the terminal.
func (t *Terminal) PasteClipboard() {
	C.vte_terminal_paste_clipboard(t.ptr)
}

// PastePrimary pastes contents of the primary selection to the terminal.
func (t *Terminal) PastePrimary() {
	C.vte_terminal_paste_primary(t.ptr)
}

// PasteText pastes text to the terminal.
func (t *Terminal) PasteText(text string) {
	s := C.CString(text)
	C.vte_terminal_paste_text(t.ptr, s)
	C.free(unsafe.Pointer(s))
}

// GetPty returns [Pty] associated with the terminal.
func (t *Terminal) GetPty() *Pty {
	pty := C.vte_terminal_get_pty(t.ptr)
	return &Pty{pty}
}

// SetPty sets [Pty] to use in terminal. Use nil to unset the pty.
func (t *Terminal) SetPty(pty *Pty) {
	if pty == nil {
		pty = &Pty{}
	}

	C.vte_terminal_set_pty(t.ptr, pty.ptr)
}

// GetAllowHyperlink reports whether or not hyperlinks (OSC 8 escape sequence)
// are allowed.
func (t *Terminal) GetAllowHyperlink() bool {
	return goBool(C.vte_terminal_get_allow_hyperlink(t.ptr))
}

// SetAllowHyperlink controls whether or not hyperlinks (OSC 8 escape sequence)
// are allowed.
func (t *Terminal) SetAllowHyperlink(v bool) {
	C.vte_terminal_set_allow_hyperlink(t.ptr, gboolean(v))
}

// GetAudibleBell reports whether or not the terminal will beep when the
// child outputs the "bl" sequence.
func (t *Terminal) GetAudibleBell() bool {
	return goBool(C.vte_terminal_get_audible_bell(t.ptr))
}

// SetAudibleBell controls whether or not the terminal will beep when the
// child outputs the "bl" sequence.
func (t *Terminal) SetAudibleBell(v bool) {
	C.vte_terminal_set_audible_bell(t.ptr, gboolean(v))
}

// GetBackspaceBinding reports what string or control sequence the terminal
// sends to its child when the user presses the backspace key.
func (t *Terminal) GetBackspaceBinding() EraseBinding {
	value, _ := t.GetProperty("backspace-binding")
	binding, ok := value.(int)
	if !ok {
		return ERASE_AUTO
	}
	return EraseBinding(binding)
}

// SetBackspaceBinding controls what string or control sequence the terminal
// sends to its child when the user presses the backspace key.
func (t *Terminal) SetBackspaceBinding(binding EraseBinding) {
	C.vte_terminal_set_backspace_binding(t.ptr, C.VteEraseBinding(binding))
}

// GetDeleteBinding reports what string or control sequence the terminal sends
// to its child when the user presses the delete key.
func (t *Terminal) GetDeleteBinding() EraseBinding {
	value, _ := t.GetProperty("delete-binding")
	binding, ok := value.(int)
	if !ok {
		return ERASE_AUTO
	}
	return EraseBinding(binding)
}

// SetDeleteBinding controls what string or control sequence the terminal sends
// to its child when the user presses the delete key.
func (t *Terminal) SetDeleteBinding(binding EraseBinding) {
	C.vte_terminal_set_delete_binding(t.ptr, C.VteEraseBinding(binding))
}

// GetBoldIsBright reports whether the SGR 1 attribute also switches to the
// bright counterpart of the first 8 palette colors, in addition to making them
// bold (legacy behavior) or if SGR 1 only enables bold and leaves the color
// intact.
func (t *Terminal) GetBoldIsBright() bool {
	return goBool(C.vte_terminal_get_bold_is_bright(t.ptr))
}

// SetBoldIsBright controls whether the SGR 1 attribute also switches to the
// bright counterpart of the first 8 palette colors, in addition to making them
// bold (legacy behavior) or if SGR 1 only enables bold and leaves the color
// intact.
func (t *Terminal) SetBoldIsBright(v bool) {
	C.vte_terminal_set_bold_is_bright(t.ptr, gboolean(v))
}

// GetCellHeightScale returns scale factor for the cell height that increases
// line spacing.
func (t *Terminal) GetCellHeightScale() float64 {
	return float64(C.vte_terminal_get_cell_height_scale(t.ptr))
}

// SetCellHeightScale controls scale factor for the cell height that increases
// line spacing. The font's height is not affected.
func (t *Terminal) SetCellHeightScale(scale float64) {
	C.vte_terminal_set_cell_height_scale(t.ptr, C.gdouble(scale))
}

// GetCellWidthScale returns scale factor for the cell width that increases
// letter spacing.
func (t *Terminal) GetCellWidthScale() float64 {
	return float64(C.vte_terminal_get_cell_width_scale(t.ptr))
}

// SetCellWidthScale controls scale factor for the cell width that increases
// letter spacing. The font's width is not affected.
func (t *Terminal) SetCellWidthScale(scale float64) {
	C.vte_terminal_set_cell_width_scale(t.ptr, C.gdouble(scale))
}

// GetCJKAmbiguousWidth reports whether ambiguous-width characters are narrow
// or wide.
func (t *Terminal) GetCJKAmbiguousWidth() CJKAmbiguousWidth {
	return CJKAmbiguousWidth(C.vte_terminal_get_cjk_ambiguous_width(t.ptr))
}

// SetCJKAmbiguousWidth controls whether ambiguous-width characters are narrow
// or wide.
//
// This setting only takes effect the next time the terminal is reset, either
// via escape sequence or with [Terminal.Reset].
func (t *Terminal) SetCJKAmbiguousWidth(width CJKAmbiguousWidth) {
	C.vte_terminal_set_cjk_ambiguous_width(t.ptr, C.int(width))
}

// GetContextMenu returns context menu of the terminal.
func (t *Terminal) GetContextMenu() gtk.IWidget {
	menu, err := t.GetProperty("context-menu")
	if err != nil {
		return nil
	}

	widget, ok := menu.(gtk.IWidget)
	if !ok {
		return nil
	}

	return widget
}

// SetContextMenu sets menu as the context menu in terminal.
// Use nil to unset the current menu.
func (t *Terminal) SetContextMenu(menu gtk.IWidget) {
	var widget *C.GtkWidget

	if menu != nil {
		widget = C.toGtkWidget(unsafe.Pointer(menu.ToWidget().GObject))
	}

	C.vte_terminal_set_context_menu(t.ptr, widget)
}

// GetContextMenuModel returns context menu model of the terminal.
func (t *Terminal) GetContextMenuModel() *glib.MenuModel {
	menu, err := t.GetProperty("context-menu-model")
	if err != nil {
		return nil
	}

	widget, ok := menu.(*glib.MenuModel)
	if !ok {
		return nil
	}

	return widget
}

// SetContextMenuModel sets model as the context menu model in terminal.
// Use nil to unset the current menu model.
//
// Takes precedence over [Terminal.SetContextMenu].
func (t *Terminal) SetContextMenuModel(model *glib.MenuModel) {
	var menuModel *C.GMenuModel

	if model != nil {
		menuModel = C.toMenuModel(unsafe.Pointer(model.GObject))
	}

	C.vte_terminal_set_context_menu_model(t.ptr, menuModel)
}

// GetCursorBlinkMode reports whether or not the cursor will blink.
func (t *Terminal) GetCursorBlinkMode() CursorBlinkMode {
	return CursorBlinkMode(C.vte_terminal_get_cursor_blink_mode(t.ptr))
}

// SetCursorBlinkMode controls whether or not the cursor will blink.
// Using [CURSOR_BLINK_SYSTEM] will use the GtkSettings::gtk-cursor-blink
// setting.
func (t *Terminal) SetCursorBlinkMode(mode CursorBlinkMode) {
	C.vte_terminal_set_cursor_blink_mode(t.ptr, C.VteCursorBlinkMode(mode))
}

// GetCursorShape returns the shape of the cursor drawn.
func (t *Terminal) GetCursorShape() CursorShape {
	return CursorShape(C.vte_terminal_get_cursor_shape(t.ptr))
}

// SetCursorShape sets the shape of the cursor drawn.
func (t *Terminal) SetCursorShape(shape CursorShape) {
	C.vte_terminal_set_cursor_shape(t.ptr, C.VteCursorShape(shape))
}

// GetEnableA11y reports whether or not a11y is enabled for the widget.
func (t *Terminal) GetEnableA11y() bool {
	return goBool(C.vte_terminal_get_enable_a11y(t.ptr))
}

// SetEnableA11y controls whether or not a11y is enabled for the widget.
func (t *Terminal) SetEnableA11y(v bool) {
	C.vte_terminal_set_enable_a11y(t.ptr, gboolean(v))
}

// GetEnableBidi reports whether or not the terminal will perform
// bidirectional text rendering.
func (t *Terminal) GetEnableBidi() bool {
	return goBool(C.vte_terminal_get_enable_bidi(t.ptr))
}

// SetEnableBidi controls whether or not the terminal will perform
// bidirectional text rendering.
func (t *Terminal) SetEnableBidi(v bool) {
	C.vte_terminal_set_enable_bidi(t.ptr, gboolean(v))
}

// GetEnableFallbackScrolling reports whether the terminal uses scroll events
// to scroll the history if the event was not otherwise consumed by it.
func (t *Terminal) GetEnableFallbackScrolling() bool {
	return goBool(C.vte_terminal_get_enable_fallback_scrolling(t.ptr))
}

// SetEnableFallbackScrolling controls whether the terminal uses scroll events
// to scroll the history if the event was not otherwise consumed by it.
//
// This function is rarely useful, except when the terminal is added to a
// [github.com/gotk3/gotk3/gtk.ScrolledWindow], to perform kinetic scrolling
func (t *Terminal) SetEnableFallbackScrolling(v bool) {
	C.vte_terminal_set_enable_fallback_scrolling(t.ptr, gboolean(v))
}

// GetEnableLegacyOSC777 reports whether legacy OSC 777 sequences are
// translated to their corresponding termprops.
func (t *Terminal) GetEnableLegacyOSC777() bool {
	return goBool(C.vte_terminal_get_enable_legacy_osc777(t.ptr))
}

// SetEnableLegacyOSC777 controls whether legacy OSC 777 sequences are
// translated to their corresponding termprops.
func (t *Terminal) SetEnableLegacyOSC777(v bool) {
	C.vte_terminal_set_enable_legacy_osc777(t.ptr, gboolean(v))
}

// GetEnableShaping reports whether or not the terminal will shape Arabic text.
func (t *Terminal) GetEnableShaping() bool {
	return goBool(C.vte_terminal_get_enable_shaping(t.ptr))
}

// SetEnableShaping controls whether or not the terminal will shape Arabic text.
func (t *Terminal) SetEnableShaping(v bool) {
	C.vte_terminal_set_enable_shaping(t.ptr, gboolean(v))
}

// GetEnableSixel reports whether [SIXEL] image support is enabled.
//
// [SIXEL]: https://en.wikipedia.org/wiki/Sixel
func (t *Terminal) GetEnableSixel() bool {
	return goBool(C.vte_terminal_get_enable_sixel(t.ptr))
}

// SetEnableSixel controls whether [SIXEL] image support is enabled.
//
// [SIXEL]: https://en.wikipedia.org/wiki/Sixel
func (t *Terminal) SetEnableSixel(v bool) {
	C.vte_terminal_set_enable_sixel(t.ptr, gboolean(v))
}

// GetFont queries the terminal for information about the font which is used to
// draw text in the terminal.
//
// The actual font takes the font scale into account, this is not reflected in
// the return value, the unscaled font is returned.
func (t *Terminal) GetFont() *pango.FontDescription {
	return wrapPangoFontDescription(C.vte_terminal_get_font(t.ptr))
}

// SetFont sets the font used for rendering all text displayed by the terminal.
// The terminal will immediately attempt to load the desired font, retrieve its
// metrics, and attempt to resize itself to keep the same number of rows and
// columns. The font scale is applied to the specified font.
func (t *Terminal) SetFont(desc *pango.FontDescription) {
	C.vte_terminal_set_font(t.ptr, unwrapPangoFontDescription(desc))
}

// SetFontFromString is a convenience method that sets font used for rendering
// all text displayed by the terminal from font description string.
//
// # Description string format
//
// The string must have the form:
//
//	[FAMILY-LIST] [STYLE-OPTIONS] [SIZE] [VARIATIONS] [FEATURES]
//
// FAMILY-LIST is a comma-separated list of families, possibly terminated by a
// comma.
//
// STYLE-OPTIONS is a whitespace-separated list of words. Each word describes
// one of style, variant, weight, stretch, or gravity.
//
//   - Styles: "Normal", "Roman", "Oblique", "Italic".
//
//   - Variants: "Small-Caps", "All-Small-Caps", "Petite-Caps",
//     "All-Petite-Caps", "Unicase", "Title-Caps".
//
//   - Weights: "Thin", "Ultra-Light", "Extra-Light", "Light", "Semi-Light",
//     "Demi-Light", "Book", "Regular", "Medium", "Semi-Bold", "Demi-Bold",
//     "Bold", "Ultra-Bold", "Extra-Bold", "Heavy", "Black", "Ultra-Black",
//     "Extra-Black".
//
//   - Stretch values: "Ultra-Condensed", "Extra-Condensed", "Condensed",
//     "Semi-Condensed", "Semi-Expanded", "Expanded", "Extra-Expanded",
//     "Ultra-Expanded".
//
//   - Gravity values: "Not-Rotated", "South", "Upside-Down", "North",
//     "Rotated-Left", "East", "Rotated-Right", "West".
//
//   - Color values: "With-Color", "Without-Color".
//
// VARIATIONS is a comma-separated list of font variations of the form
// @‍axis1=value,axis2=value,...
//
// FEATURES is a comma-separated list of font features of the form
// #‍feature1=value,feature2=value,... The =value part can be ommitted if
// the value is 1.
//
// Any one of the options may be absent. If FAMILY-LIST is absent, then the
// family_name field of the resulting font description will be initialized to
// NULL. If STYLE-OPTIONS is missing, then all style options will be set to the
// default values. If SIZE is missing, the size in the resulting font
// description will be set to 0.
//
// A typical example:
//
//	Cantarell Italic Light 15 @‍wght=200 #‍tnum=1.
func (t *Terminal) SetFontFromString(s string) {
	cstr := C.CString(s)
	desc := C.pango_font_description_from_string(cstr)
	C.free(unsafe.Pointer(cstr))
	C.vte_terminal_set_font(t.ptr, desc)
	C.free(unsafe.Pointer(desc))
}

// GetFontOptions returns terminal's font options.
func (t *Terminal) GetFontOptions() *cairo.FontOptions {
	return wrapCairoFontOptions(C.vte_terminal_get_font_options(t.ptr))
}

// SetFontOptions sets terminal's font options. Use nil to use the default
// options.
func (t *Terminal) SetFontOptions(options *cairo.FontOptions) {
	C.vte_terminal_set_font_options(t.ptr, unwrapCairoFontOptions(options))
}

// GetTextBlinkMode reports whether or not the terminal will allow blinking text.
func (t *Terminal) GetTextBlinkMode() TextBlinkMode {
	return TextBlinkMode(C.vte_terminal_get_text_blink_mode(t.ptr))
}

// SetTextBlinkMode controls whether or not the terminal will allow blinking text.
func (t *Terminal) SetTextBlinkMode(mode TextBlinkMode) {
	C.vte_terminal_set_text_blink_mode(t.ptr, C.VteTextBlinkMode(mode))
}

// GetXAlign returns the horizontal alignment of terminal within its allocation.
func (t *Terminal) GetXAlign() Align {
	return Align(C.vte_terminal_get_xalign(t.ptr))
}

// SetXAlign sets the horizontal alignment of terminal within its allocation.
func (t *Terminal) SetXAlign(align Align) {
	C.vte_terminal_set_xalign(t.ptr, C.VteAlign(align))
}

// GetYAlign returns the vertical alignment of terminal within its allocation.
func (t *Terminal) GetYAlign() Align {
	return Align(C.vte_terminal_get_yalign(t.ptr))
}

// SetYAlign sets the vertical alignment of terminal within its allocation.
func (t *Terminal) SetYAlign(align Align) {
	C.vte_terminal_set_yalign(t.ptr, C.VteAlign(align))
}

// GetXFill returns the horizontal fillment of terminal within its allocation.
func (t *Terminal) GetXFill() bool {
	return goBool(C.vte_terminal_get_xfill(t.ptr))
}

// SetXFill sets the horizontal fillment of terminal within its allocation.
func (t *Terminal) SetXFill(v bool) {
	C.vte_terminal_set_xfill(t.ptr, gboolean(v))
}

// GetYFill returns the vertical fillment of terminal within its allocation.
func (t *Terminal) GetYFill() bool {
	return goBool(C.vte_terminal_get_yfill(t.ptr))
}

// SetYFill sets the vertical fillment of terminal within its allocation.
func (t *Terminal) SetYFill(v bool) {
	C.vte_terminal_set_yfill(t.ptr, gboolean(v))
}

// SpawnAsync is a convenience function that wraps creating the [Pty] and
// spawning the child process on it. See [Pty.SpawnAsync] for more information.
func (t *Terminal) SpawnAsync(cmd *Command) {
	var ccallID C.gpointer
	if cmd.OnSpawn != nil {
		callID := assignCallID(cmd)
		ccallID = C.uintToGpointer(C.uint(callID))
	}

	var (
		ptyFlags              = C.VtePtyFlags(cmd.PtyFlags)
		workdir               = C.CString(cmd.Dir)
		argv                  = cStringArr(cmd.Args)
		envv                  = cStringArr(cmd.Env)
		spawnFlags            = C.GSpawnFlags(cmd.SpawnFlags)
		childSetup            C.GSpawnChildSetupFunc
		childSetupData        C.gpointer
		cTimeout              = C.int(cmd.Timeout.Milliseconds())
		cCancellable          = C.toCancellable(unsafe.Pointer(cmd.Cancellable.GObject))
		childSetupDataDestroy C.GDestroyNotify
		callback              = C.VteTerminalSpawnAsyncCallback(C.terminalSpawnAsyncCallback)
		userData              = ccallID
	)

	defer C.free(unsafe.Pointer(workdir))
	defer cStringArrFree(argv)
	defer cStringArrFree(envv)

	C.vte_terminal_spawn_async(
		t.ptr,
		ptyFlags,
		workdir,
		&argv[0],
		&envv[0],
		spawnFlags,
		childSetup,
		childSetupData,
		childSetupDataDestroy,
		cTimeout,
		cCancellable,
		callback,
		userData,
	)
}
