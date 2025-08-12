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
	"errors"
	"fmt"
	"unicode/utf8"
	"unsafe"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
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
		return nil, errNilPointer("vte_terminal_new")
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

// terminalFromGObject converts object to [Terminal]. If assertion
// VTE_IS_TERMINAL fails on the underlying GObject, error is returned.
func terminalFromGObject(object *glib.Object) (*Terminal, error) {
	cGObject := object.GObject

	if !goBool(C.isTerminal(unsafe.Pointer(cGObject))) {
		return nil, fmt.Errorf("object is not a Terminal")
	}

	initiallyUnowned := glib.InitiallyUnowned{
		Object: object,
	}

	widget := gtk.Widget{
		InitiallyUnowned: initiallyUnowned,
	}

	return &Terminal{
		Widget: widget,

		ptr: C.toTerminal(unsafe.Pointer(cGObject)),
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

// Feed writes text to the standard output of the terminal as if it were
// received from a child process.
func (t *Terminal) Feed(text string) {
	cstr := C.CString(text)
	length := C.intToGssize(C.int(len(text)))
	C.vte_terminal_feed(t.ptr, cstr, length)
	C.free(unsafe.Pointer(cstr))
}

// FeedChild writes text to the standard input of the terminal as if it were
// entered by the user at the keyboard.
func (t *Terminal) FeedChild(text string) {
	cstr := C.CString(text)
	length := C.intToGssize(C.int(len(text)))
	C.vte_terminal_feed_child(t.ptr, cstr, length)
	C.free(unsafe.Pointer(cstr))
}

// Write is an alias for [Feed].
//
// Error is only returned if p is not a valid UTF-8 string.
func (t *Terminal) Write(p []byte) (int, error) {
	if !utf8.Valid(p) {
		return 0, errors.New("argument must be a valid UTF-8 string")
	}
	t.Feed(string(p))
	return len(p), nil
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

// GetHoveredURI returns the currently hovered hyperlink URI, or an empty
// string ("") if no hyperlinks are hovered.
func (t *Terminal) GetHoveredURI() string {
	v, err := t.GetProperty("hyperlink-hover-uri")
	if err != nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
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
	desc := C.vte_terminal_get_font(t.ptr)
	if desc == nil {
		return nil
	}
	return wrapPangoFontDescription(desc)
}

// SetFont sets the font used for rendering all text displayed by the terminal.
// The terminal will immediately attempt to load the desired font, retrieve its
// metrics, and attempt to resize itself to keep the same number of rows and
// columns. The font scale is applied to the specified font.
func (t *Terminal) SetFont(desc *pango.FontDescription) {
	var d *C.PangoFontDescription
	if desc != nil {
		d = unwrapPangoFontDescription(desc)
	}
	C.vte_terminal_set_font(t.ptr, d)
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
// #‍feature1=value,feature2=value,... The =value part can be omitted if
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
	options := C.vte_terminal_get_font_options(t.ptr)
	if options == nil {
		return nil
	}
	return wrapCairoFontOptions(options)
}

// SetFontOptions sets terminal's font options. Use nil to use the default
// options.
func (t *Terminal) SetFontOptions(options *cairo.FontOptions) {
	var opt *C.cairo_font_options_t
	if options != nil {
		opt = unwrapCairoFontOptions(options)
	}
	C.vte_terminal_set_font_options(t.ptr, opt)
}

// GetInputEnabled reports whether the terminal allows user input. When user
// input is disabled, key press and mouse button press and motion events are
// not sent to the terminal's child.
func (t *Terminal) GetInputEnabled() bool {
	return goBool(C.vte_terminal_get_input_enabled(t.ptr))
}

// SetInputEnabled controls whether the terminal allows user input. When user
// input is disabled, key press and mouse button press and motion events are
// not sent to the terminal's child.
func (t *Terminal) SetInputEnabled(v bool) {
	C.vte_terminal_set_input_enabled(t.ptr, gboolean(v))
}

// GetPointerAutohide returns mouse autohide setting. When autohiding is
// enabled, the mouse cursor will be hidden when the user presses a key and
// shown when the user moves the mouse.
func (t *Terminal) GetPointerAutohide() bool {
	v, err := t.GetProperty("pointer-autohide")
	if err != nil {
		return false
	}
	b, ok := v.(bool)
	if !ok {
		return false
	}
	return b
}

// SetPointerAutohide sets mouse autohide setting. When autohiding is
// enabled, the mouse cursor will be hidden when the user presses a key and
// shown when the user moves the mouse.
func (t *Terminal) SetPointerAutohide(v bool) {
	t.SetProperty("pointer-autohide", v)
}

// GetScrollOnInsert reports whether or not the terminal will forcibly scroll
// to the bottom of the viewable history when the text is inserted
// (e.g. by a paste).
func (t *Terminal) GetScrollOnInsert() bool {
	return goBool(C.vte_terminal_get_scroll_on_insert(t.ptr))
}

// SetScrollOnInsert controls whether or not the terminal will forcibly scroll
// to the bottom of the viewable history when the text is inserted
// (e.g. by a paste).
func (t *Terminal) SetScrollOnInsert(v bool) {
	C.vte_terminal_set_scroll_on_insert(t.ptr, gboolean(v))
}

// GetScrollOnKeystroke reports whether or not the terminal will forcibly
// scroll to the bottom of the viewable history when the user presses a key.
// Modifier keys do not trigger this behavior.
func (t *Terminal) GetScrollOnKeystroke() bool {
	return goBool(C.vte_terminal_get_scroll_on_keystroke(t.ptr))
}

// SetScrollOnKeystroke controls whether or not the terminal will forcibly
// scroll to the bottom of the viewable history when the user presses a key.
// Modifier keys do not trigger this behavior.
func (t *Terminal) SetScrollOnKeystroke(v bool) {
	C.vte_terminal_set_scroll_on_keystroke(t.ptr, gboolean(v))
}

// GetScrollOnOutput reports whether or not the terminal will forcibly scroll
// to the bottom of the viewable history when the new data is received from the
// child.
func (t *Terminal) GetScrollOnOutput() bool {
	return goBool(C.vte_terminal_get_scroll_on_output(t.ptr))
}

// SetScrollOnOutput controls whether or not the terminal will forcibly scroll
// to the bottom of the viewable history when the new data is received from the
// child.
func (t *Terminal) SetScrollOnOutput(v bool) {
	C.vte_terminal_set_scroll_on_output(t.ptr, gboolean(v))
}

// GetScrollUnit returns scroll measurement unit of terminal's
// [github.com/gotk3/gotk3/gtk.Adjustment].
func (t *Terminal) GetScrollUnit() ScrollUnit {
	isPixels := goBool(C.vte_terminal_get_scroll_unit_is_pixels(t.ptr))
	if isPixels {
		return SCROLL_UNIT_PIXELS
	} else {
		return SCROLL_UNIT_LINES
	}
}

// SetScrollUnit controls scroll measurement unit of terminal's
// [github.com/gotk3/gotk3/gtk.Adjustment].
//
// It may be useful to set unit to [SCROLL_UNIT_PIXELS] when the terminal is the
// child of a [github.com/gotk3/gotk3/gtk.ScrolledWindow] to fix some bugs with
// kinetic scrolling.
func (t *Terminal) SetScrollUnit(unit ScrollUnit) {
	isPixels := unit == SCROLL_UNIT_PIXELS
	C.vte_terminal_set_scroll_unit_is_pixels(t.ptr, gboolean(isPixels))
}

// GetScrollbackLines returns the length of the scrollback buffer used by the
// terminal.
func (t *Terminal) GetScrollbackLines() uint {
	return uint(C.vte_terminal_get_scrollback_lines(t.ptr))
}

// SetScrollbackLines sets the length of the scrollback buffer used by the
// terminal.
//
// The size of the scrollback buffer will be set to the larger of this value
// and the number of visible rows the widget can display, so 0 can safely be
// used to disable scrollback.
//
// Note that this setting only affects the normal screen buffer. For terminal
// types which have an alternate screen buffer, no scrollback is allowed on the
// alternate screen buffer.
func (t *Terminal) SetScrollbackLines(size int) {
	C.vte_terminal_set_scrollback_lines(t.ptr, C.uintToGlong(C.uint(size)))
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

// Spawn is a convenience function that wraps creating the [Pty] and
// spawning the child process on it. See [Pty.Spawn] for more information.
func (t *Terminal) Spawn(cmd *Command) {
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

// Reset resets as much of the terminal's internal state as possible,
// discarding any unprocessed input data, resetting character attributes,
// cursor state, national character set state, status line, terminal modes
// (insert/delete), selection state, and encoding.
func (t *Terminal) Reset(clearTabstops, clearHistory bool) {
	C.vte_terminal_reset(t.ptr, gboolean(clearTabstops), gboolean(clearTabstops))
}

// SetColors sets terminal colors.
//
// palette must contain 0, 8, 16, 232, or 256 colors. It specifies the new
// values for the 256 palette colors in the following order:
//
//   - 8 standard colors
//   - their 8 bright counterparts
//   - 6x6x6 color cube
//   - 24 grayscale colors
//
// More information about terminal colors can be found [here].
//
// Special cases are:
//
//   - If foreground is nil and palette length is greater than 0, the new
//     foreground color is taken from palette[7].
//   - If background is nil and palette length is greater than 0, the new
//     background color is taken from palette[0].
//   - Omitted entries will default to a hardcoded value.
//
// [here]: https://en.wikipedia.org/wiki/ANSI_escape_code#8-bit
func (t *Terminal) SetColors(background, foreground *gdk.RGBA, palette []*gdk.RGBA) error {
	l := len(palette)

	if l != 0 && l != 8 && l != 16 && l != 232 && l != 256 {
		return fmt.Errorf("palette must contain 0, 8, 16, 232, or 256 colors")
	}

	var (
		bg *C.GdkRGBA
		fg *C.GdkRGBA

		p    = make([]*C.GdkRGBA, len(palette))
		size = C.uintToGsize(C.uint(len(palette)))
	)

	if background != nil {
		bg = unwrapGdkRGBA(background)
	}

	if foreground != nil {
		fg = unwrapGdkRGBA(foreground)
	}

	for i, color := range palette {
		p[i] = unwrapGdkRGBA(color)
	}

	C.vte_terminal_set_colors(t.ptr, fg, bg, p[0], size)
	return nil
}

// SetCursorColor sets color for text which is under the cursor.
// Use nil to unset a color. If both background and foreground are nil, text
// under the cursor will be drawn with foreground and background colors
// reversed.
func (t *Terminal) SetCursorColor(background, foreground *gdk.RGBA) {
	var (
		bg *C.GdkRGBA
		fg *C.GdkRGBA
	)

	if background != nil {
		bg = unwrapGdkRGBA(background)
	}

	if foreground != nil {
		fg = unwrapGdkRGBA(foreground)
	}

	C.vte_terminal_set_color_cursor(t.ptr, bg)
	C.vte_terminal_set_color_cursor_foreground(t.ptr, fg)
}

// SetHighlightColor sets the color for the text which is highlighted.
// Use nil to unset a color. If both background and foreground, highlighted
// text (which is usually highlighted because it is selected) will be drawn
// with foreground and background colors reversed.
func (t *Terminal) SetHighlightColor(background, foreground *gdk.RGBA) {
	var (
		bg *C.GdkRGBA
		fg *C.GdkRGBA
	)

	if background != nil {
		bg = unwrapGdkRGBA(background)
	}

	if foreground != nil {
		fg = unwrapGdkRGBA(foreground)
	}

	C.vte_terminal_set_color_highlight(t.ptr, bg)
	C.vte_terminal_set_color_highlight_foreground(t.ptr, fg)
}

// SearchFindNext searches the next string matching the search regex set with
// [SearchSetRegex].
func (t *Terminal) SearchFindNext() bool {
	return goBool(C.vte_terminal_search_find_next(t.ptr))
}

// SearchFindPrev searches the previous string matching the search regex set
// with [SearchSetRegex].
func (t *Terminal) SearchFindPrev() bool {
	return goBool(C.vte_terminal_search_find_previous(t.ptr))
}
