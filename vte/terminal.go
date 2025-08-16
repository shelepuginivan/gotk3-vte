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

// MatchHandle represents a tag associated with the [Regex].
type MatchHandle int

// Terminal is a wrapper around VteTerminal.
type Terminal struct {
	gtk.Widget
}

// WrapTerminal wraps [github.com/gotk3/gotk3/glib.Object] and casts it into
// [Terminal].
//
// This function is intended to be used in callbacks when calling
// [github.com/gotk3/gotk3/glib.Object.Connect] directly on [Terminal]:
//
//	term, err := vte.TerminalNew()
//	if err != nil {
//	  	 log.Fatal(err)
//	}
//
//	term.Connect("button-press-event", func(obj *glib.Object, ev *gdk.Event) {
//	    term := vte.WrapTerminal(obj)
//
//	    // rest of the handler logic...
//	})
func WrapTerminal(obj *glib.Object) *Terminal {
	if obj == nil {
		return nil
	}

	return &Terminal{
		Widget: gtk.Widget{
			InitiallyUnowned: glib.InitiallyUnowned{
				Object: obj,
			},
		},
	}
}

func (t *Terminal) native() *C.VteTerminal {
	if t == nil || t.GObject == nil {
		return nil
	}
	return C.toTerminal(unsafe.Pointer(t.GObject))
}

// TerminalNew creates a new instance of [*Terminal].
func TerminalNew() (*Terminal, error) {
	ptr := C.vte_terminal_new()
	if ptr == nil {
		return nil, errNilPointer("vte_terminal_new")
	}
	return WrapTerminal(glib.Take(unsafe.Pointer(ptr))), nil
}

// CopyClipboardFormat copies selected text to clipboard in the specified
// format.
func (t *Terminal) CopyClipboardFormat(format Format) {
	C.vte_terminal_copy_clipboard_format(t.native(), C.VteFormat(format))
}

// CopyPrimary copies selected text in the primary selection.
func (t *Terminal) CopyPrimary() {
	C.vte_terminal_copy_primary(t.native())
}

// PasteClipboard pastes contents of clipboard to the terminal.
func (t *Terminal) PasteClipboard() {
	C.vte_terminal_paste_clipboard(t.native())
}

// PastePrimary pastes contents of the primary selection to the terminal.
func (t *Terminal) PastePrimary() {
	C.vte_terminal_paste_primary(t.native())
}

// PasteText pastes text to the terminal.
func (t *Terminal) PasteText(text string) {
	s := C.CString(text)
	C.vte_terminal_paste_text(t.native(), s)
	C.free(unsafe.Pointer(s))
}

// Feed writes text to the standard output of the terminal as if it were
// received from a child process.
func (t *Terminal) Feed(text string) {
	cstr := C.CString(text)
	length := C.intToGssize(C.int(len(text)))
	C.vte_terminal_feed(t.native(), cstr, length)
	C.free(unsafe.Pointer(cstr))
}

// FeedChild writes text to the standard input of the terminal as if it were
// entered by the user at the keyboard.
func (t *Terminal) FeedChild(text string) {
	cstr := C.CString(text)
	length := C.intToGssize(C.int(len(text)))
	C.vte_terminal_feed_child(t.native(), cstr, length)
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

// GetTermProp returns termprop by name. It should be called inside the
// callback of [Terminal.ConnectTermPropChanged] or
// [Terminal.ConnectAfterTermPropChanged].
func (t *Terminal) GetTermProp(prop TermProp) (any, error) {
	v, err := glib.ValueAlloc()
	if err != nil {
		return nil, err
	}

	gvalue := C.toGValue(unsafe.Pointer(v.GValue))
	cprop := C.CString(string(prop))
	defer C.free(unsafe.Pointer(cprop))

	if !goBool(C.vte_terminal_get_termprop_value(t.native(), cprop, gvalue)) {
		return nil, errFailed("vte_terminal_get_termprop_value")
	}

	return v.GoValue()
}

// GetPty returns [Pty] associated with the terminal.
func (t *Terminal) GetPty() *Pty {
	pty := C.vte_terminal_get_pty(t.native())
	return wrapPty(glib.Take(unsafe.Pointer(pty)))
}

// SetPty sets [Pty] to use in terminal. If pty is nil, this function is no-op.
func (t *Terminal) SetPty(pty *Pty) {
	if pty == nil {
		return
	}

	C.vte_terminal_set_pty(t.native(), pty.native())
}

// GetAllowHyperlink reports whether or not hyperlinks (OSC 8 escape sequence)
// are allowed.
func (t *Terminal) GetAllowHyperlink() bool {
	return goBool(C.vte_terminal_get_allow_hyperlink(t.native()))
}

// SetAllowHyperlink controls whether or not hyperlinks (OSC 8 escape sequence)
// are allowed.
func (t *Terminal) SetAllowHyperlink(v bool) {
	C.vte_terminal_set_allow_hyperlink(t.native(), gboolean(v))
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
	return goBool(C.vte_terminal_get_audible_bell(t.native()))
}

// SetAudibleBell controls whether or not the terminal will beep when the
// child outputs the "bl" sequence.
func (t *Terminal) SetAudibleBell(v bool) {
	C.vte_terminal_set_audible_bell(t.native(), gboolean(v))
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
	C.vte_terminal_set_backspace_binding(t.native(), C.VteEraseBinding(binding))
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
	C.vte_terminal_set_delete_binding(t.native(), C.VteEraseBinding(binding))
}

// GetBoldIsBright reports whether the SGR 1 attribute also switches to the
// bright counterpart of the first 8 palette colors, in addition to making them
// bold (legacy behavior) or if SGR 1 only enables bold and leaves the color
// intact.
func (t *Terminal) GetBoldIsBright() bool {
	return goBool(C.vte_terminal_get_bold_is_bright(t.native()))
}

// SetBoldIsBright controls whether the SGR 1 attribute also switches to the
// bright counterpart of the first 8 palette colors, in addition to making them
// bold (legacy behavior) or if SGR 1 only enables bold and leaves the color
// intact.
func (t *Terminal) SetBoldIsBright(v bool) {
	C.vte_terminal_set_bold_is_bright(t.native(), gboolean(v))
}

// GetCellHeightScale returns scale factor for the cell height that increases
// line spacing.
func (t *Terminal) GetCellHeightScale() float64 {
	return float64(C.vte_terminal_get_cell_height_scale(t.native()))
}

// SetCellHeightScale controls scale factor for the cell height that increases
// line spacing. The font's height is not affected.
func (t *Terminal) SetCellHeightScale(scale float64) {
	C.vte_terminal_set_cell_height_scale(t.native(), C.gdouble(scale))
}

// GetCellWidthScale returns scale factor for the cell width that increases
// letter spacing.
func (t *Terminal) GetCellWidthScale() float64 {
	return float64(C.vte_terminal_get_cell_width_scale(t.native()))
}

// SetCellWidthScale controls scale factor for the cell width that increases
// letter spacing. The font's width is not affected.
func (t *Terminal) SetCellWidthScale(scale float64) {
	C.vte_terminal_set_cell_width_scale(t.native(), C.gdouble(scale))
}

// GetCJKAmbiguousWidth reports whether ambiguous-width characters are narrow
// or wide.
func (t *Terminal) GetCJKAmbiguousWidth() CJKAmbiguousWidth {
	return CJKAmbiguousWidth(C.vte_terminal_get_cjk_ambiguous_width(t.native()))
}

// SetCJKAmbiguousWidth controls whether ambiguous-width characters are narrow
// or wide.
//
// This setting only takes effect the next time the terminal is reset, either
// via escape sequence or with [Terminal.Reset].
func (t *Terminal) SetCJKAmbiguousWidth(width CJKAmbiguousWidth) {
	C.vte_terminal_set_cjk_ambiguous_width(t.native(), C.int(width))
}

// GetContextMenu returns context menu of the terminal.
func (t *Terminal) GetContextMenu() gtk.IWidget {
	menu, err := t.GetProperty("context-menu")
	if err != nil {
		return nil
	}

	if menu == (*gtk.Menu)(nil) {
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

	C.vte_terminal_set_context_menu(t.native(), widget)
}

// GetContextMenuModel returns context menu model of the terminal.
func (t *Terminal) GetContextMenuModel() *glib.MenuModel {
	model, err := t.GetProperty("context-menu-model")
	if err != nil {
		return nil
	}

	if model == (*glib.Object)(nil) {
		return nil
	}

	widget, ok := model.(*glib.MenuModel)
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

	C.vte_terminal_set_context_menu_model(t.native(), menuModel)
}

// GetCursorBlinkMode reports whether or not the cursor will blink.
func (t *Terminal) GetCursorBlinkMode() CursorBlinkMode {
	return CursorBlinkMode(C.vte_terminal_get_cursor_blink_mode(t.native()))
}

// SetCursorBlinkMode controls whether or not the cursor will blink.
// Using [CURSOR_BLINK_SYSTEM] will use the GtkSettings::gtk-cursor-blink
// setting.
func (t *Terminal) SetCursorBlinkMode(mode CursorBlinkMode) {
	C.vte_terminal_set_cursor_blink_mode(t.native(), C.VteCursorBlinkMode(mode))
}

// GetCursorShape returns the shape of the cursor drawn.
func (t *Terminal) GetCursorShape() CursorShape {
	return CursorShape(C.vte_terminal_get_cursor_shape(t.native()))
}

// SetCursorShape sets the shape of the cursor drawn.
func (t *Terminal) SetCursorShape(shape CursorShape) {
	C.vte_terminal_set_cursor_shape(t.native(), C.VteCursorShape(shape))
}

// GetEnableA11y reports whether or not a11y is enabled for the widget.
func (t *Terminal) GetEnableA11y() bool {
	return goBool(C.vte_terminal_get_enable_a11y(t.native()))
}

// SetEnableA11y controls whether or not a11y is enabled for the widget.
func (t *Terminal) SetEnableA11y(v bool) {
	C.vte_terminal_set_enable_a11y(t.native(), gboolean(v))
}

// GetEnableBidi reports whether or not the terminal will perform
// bidirectional text rendering.
func (t *Terminal) GetEnableBidi() bool {
	return goBool(C.vte_terminal_get_enable_bidi(t.native()))
}

// SetEnableBidi controls whether or not the terminal will perform
// bidirectional text rendering.
func (t *Terminal) SetEnableBidi(v bool) {
	C.vte_terminal_set_enable_bidi(t.native(), gboolean(v))
}

// GetEnableFallbackScrolling reports whether the terminal uses scroll events
// to scroll the history if the event was not otherwise consumed by it.
func (t *Terminal) GetEnableFallbackScrolling() bool {
	return goBool(C.vte_terminal_get_enable_fallback_scrolling(t.native()))
}

// SetEnableFallbackScrolling controls whether the terminal uses scroll events
// to scroll the history if the event was not otherwise consumed by it.
//
// This function is rarely useful, except when the terminal is added to a
// [github.com/gotk3/gotk3/gtk.ScrolledWindow], to perform kinetic scrolling
func (t *Terminal) SetEnableFallbackScrolling(v bool) {
	C.vte_terminal_set_enable_fallback_scrolling(t.native(), gboolean(v))
}

// GetEnableLegacyOSC777 reports whether legacy OSC 777 sequences are
// translated to their corresponding termprops.
func (t *Terminal) GetEnableLegacyOSC777() bool {
	return goBool(C.vte_terminal_get_enable_legacy_osc777(t.native()))
}

// SetEnableLegacyOSC777 controls whether legacy OSC 777 sequences are
// translated to their corresponding termprops.
func (t *Terminal) SetEnableLegacyOSC777(v bool) {
	C.vte_terminal_set_enable_legacy_osc777(t.native(), gboolean(v))
}

// GetEnableShaping reports whether or not the terminal will shape Arabic text.
func (t *Terminal) GetEnableShaping() bool {
	return goBool(C.vte_terminal_get_enable_shaping(t.native()))
}

// SetEnableShaping controls whether or not the terminal will shape Arabic text.
func (t *Terminal) SetEnableShaping(v bool) {
	C.vte_terminal_set_enable_shaping(t.native(), gboolean(v))
}

// GetEnableSixel reports whether [SIXEL] image support is enabled.
//
// [SIXEL]: https://en.wikipedia.org/wiki/Sixel
func (t *Terminal) GetEnableSixel() bool {
	return goBool(C.vte_terminal_get_enable_sixel(t.native()))
}

// SetEnableSixel controls whether [SIXEL] image support is enabled.
//
// [SIXEL]: https://en.wikipedia.org/wiki/Sixel
func (t *Terminal) SetEnableSixel(v bool) {
	C.vte_terminal_set_enable_sixel(t.native(), gboolean(v))
}

// GetFont queries the terminal for information about the font which is used to
// draw text in the terminal.
//
// The actual font takes the font scale into account, this is not reflected in
// the return value, the unscaled font is returned.
func (t *Terminal) GetFont() *pango.FontDescription {
	desc := C.vte_terminal_get_font(t.native())
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
	C.vte_terminal_set_font(t.native(), d)
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
	C.vte_terminal_set_font(t.native(), desc)
	C.free(unsafe.Pointer(desc))
}

// GetFontOptions returns terminal's font options.
func (t *Terminal) GetFontOptions() *cairo.FontOptions {
	options := C.vte_terminal_get_font_options(t.native())
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
	C.vte_terminal_set_font_options(t.native(), opt)
}

// GetFontScale returns font scale of the terminal.
func (t *Terminal) GetFontScale() float64 {
	return float64(C.vte_terminal_get_font_scale(t.native()))
}

// SetFontScale sets font scale of the terminal.
func (t *Terminal) SetFontScale(scale float64) {
	C.vte_terminal_set_font_scale(t.native(), C.gdouble(scale))
}

// GetInputEnabled reports whether the terminal allows user input. When user
// input is disabled, key press and mouse button press and motion events are
// not sent to the terminal's child.
func (t *Terminal) GetInputEnabled() bool {
	return goBool(C.vte_terminal_get_input_enabled(t.native()))
}

// SetInputEnabled controls whether the terminal allows user input. When user
// input is disabled, key press and mouse button press and motion events are
// not sent to the terminal's child.
func (t *Terminal) SetInputEnabled(v bool) {
	C.vte_terminal_set_input_enabled(t.native(), gboolean(v))
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
	return goBool(C.vte_terminal_get_scroll_on_insert(t.native()))
}

// SetScrollOnInsert controls whether or not the terminal will forcibly scroll
// to the bottom of the viewable history when the text is inserted
// (e.g. by a paste).
func (t *Terminal) SetScrollOnInsert(v bool) {
	C.vte_terminal_set_scroll_on_insert(t.native(), gboolean(v))
}

// GetScrollOnKeystroke reports whether or not the terminal will forcibly
// scroll to the bottom of the viewable history when the user presses a key.
// Modifier keys do not trigger this behavior.
func (t *Terminal) GetScrollOnKeystroke() bool {
	return goBool(C.vte_terminal_get_scroll_on_keystroke(t.native()))
}

// SetScrollOnKeystroke controls whether or not the terminal will forcibly
// scroll to the bottom of the viewable history when the user presses a key.
// Modifier keys do not trigger this behavior.
func (t *Terminal) SetScrollOnKeystroke(v bool) {
	C.vte_terminal_set_scroll_on_keystroke(t.native(), gboolean(v))
}

// GetScrollOnOutput reports whether or not the terminal will forcibly scroll
// to the bottom of the viewable history when the new data is received from the
// child.
func (t *Terminal) GetScrollOnOutput() bool {
	return goBool(C.vte_terminal_get_scroll_on_output(t.native()))
}

// SetScrollOnOutput controls whether or not the terminal will forcibly scroll
// to the bottom of the viewable history when the new data is received from the
// child.
func (t *Terminal) SetScrollOnOutput(v bool) {
	C.vte_terminal_set_scroll_on_output(t.native(), gboolean(v))
}

// GetScrollUnit returns scroll measurement unit of terminal's
// [github.com/gotk3/gotk3/gtk.Adjustment].
func (t *Terminal) GetScrollUnit() ScrollUnit {
	isPixels := goBool(C.vte_terminal_get_scroll_unit_is_pixels(t.native()))
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
	C.vte_terminal_set_scroll_unit_is_pixels(t.native(), gboolean(isPixels))
}

// GetScrollbackLines returns the length of the scrollback buffer used by the
// terminal.
func (t *Terminal) GetScrollbackLines() uint {
	return uint(C.vte_terminal_get_scrollback_lines(t.native()))
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
	C.vte_terminal_set_scrollback_lines(t.native(), C.uintToGlong(C.uint(size)))
}

// GetTextBlinkMode reports whether or not the terminal will allow blinking text.
func (t *Terminal) GetTextBlinkMode() TextBlinkMode {
	return TextBlinkMode(C.vte_terminal_get_text_blink_mode(t.native()))
}

// SetTextBlinkMode controls whether or not the terminal will allow blinking text.
func (t *Terminal) SetTextBlinkMode(mode TextBlinkMode) {
	C.vte_terminal_set_text_blink_mode(t.native(), C.VteTextBlinkMode(mode))
}

// GetXAlign returns the horizontal alignment of terminal within its allocation.
func (t *Terminal) GetXAlign() Align {
	return Align(C.vte_terminal_get_xalign(t.native()))
}

// SetXAlign sets the horizontal alignment of terminal within its allocation.
func (t *Terminal) SetXAlign(align Align) {
	C.vte_terminal_set_xalign(t.native(), C.VteAlign(align))
}

// GetYAlign returns the vertical alignment of terminal within its allocation.
func (t *Terminal) GetYAlign() Align {
	return Align(C.vte_terminal_get_yalign(t.native()))
}

// SetYAlign sets the vertical alignment of terminal within its allocation.
func (t *Terminal) SetYAlign(align Align) {
	C.vte_terminal_set_yalign(t.native(), C.VteAlign(align))
}

// GetXFill returns the horizontal fillment of terminal within its allocation.
func (t *Terminal) GetXFill() bool {
	return goBool(C.vte_terminal_get_xfill(t.native()))
}

// SetXFill sets the horizontal fillment of terminal within its allocation.
func (t *Terminal) SetXFill(v bool) {
	C.vte_terminal_set_xfill(t.native(), gboolean(v))
}

// GetYFill returns the vertical fillment of terminal within its allocation.
func (t *Terminal) GetYFill() bool {
	return goBool(C.vte_terminal_get_yfill(t.native()))
}

// SetYFill sets the vertical fillment of terminal within its allocation.
func (t *Terminal) SetYFill(v bool) {
	C.vte_terminal_set_yfill(t.native(), gboolean(v))
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
		t.native(),
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

// WatchChild watches pid. When the process exists, the "child-exited" signal
// will be called with the child's exit status.
//
// This method is only required if [Pty] was set with [Terminal.SetPty], and
// the child process was spawned with [Pty.Spawn]. If [Terminal.Spawn] is used,
// this is handled automatically.
func (t *Terminal) WatchChild(pid int) {
	C.vte_terminal_watch_child(t.native(), C.GPid(pid))
}

// Reset resets as much of the terminal's internal state as possible,
// discarding any unprocessed input data, resetting character attributes,
// cursor state, national character set state, status line, terminal modes
// (insert/delete), selection state, and encoding.
func (t *Terminal) Reset(clearTabstops, clearHistory bool) {
	C.vte_terminal_reset(t.native(), gboolean(clearTabstops), gboolean(clearTabstops))
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

	C.vte_terminal_set_colors(t.native(), fg, bg, p[0], size)
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

	C.vte_terminal_set_color_cursor(t.native(), bg)
	C.vte_terminal_set_color_cursor_foreground(t.native(), fg)
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

	C.vte_terminal_set_color_highlight(t.native(), bg)
	C.vte_terminal_set_color_highlight_foreground(t.native(), fg)
}

// MatchAddRegex adds the regular expression regex to the list of matching
// expressions. When the user moves the mouse cursor over a section of
// displayed text which matches this expression, the text will be highlighted.
//
// The primary use case for this is handling clicks on URLs or other types of
// links. See [Terminal.MatchCheckEvent] for more information.
func (t *Terminal) MatchAddRegex(regex *Regex, flags RegexMatchFlags) (MatchHandle, error) {
	if regex == nil {
		return -1, fmt.Errorf("regex must not be nil")
	}
	return MatchHandle(C.vte_terminal_match_add_regex(t.native(), regex.ptr, C.uint(flags))), nil
}

// MatchRemove removes the regular expression which is associated with handle
// from the list of expressions which the terminal will highlight when the user
// moves the mouse cursor over matching text.
func (t *Terminal) MatchRemove(handle MatchHandle) {
	C.vte_terminal_match_remove(t.native(), C.int(handle))
}

// MatchRemoveAll clears the list of regular expressions the terminal uses to
// highlight text when the user moves the mouse cursor.
func (t *Terminal) MatchRemoveAll() {
	C.vte_terminal_match_remove_all(t.native())
}

// MatchSetCursorName sets which cursor the terminal will use if the pointer is
// over the pattern specified by tag.
func (t *Terminal) MatchSetCursorName(handle MatchHandle, cursor string) {
	cstr := C.CString(cursor)
	C.vte_terminal_match_set_cursor_name(t.native(), C.int(handle), cstr)
	C.free(unsafe.Pointer(cstr))
}

// MatchCheckEvent checks if the text in and around the position of the event
// matches any of the regular expressions previously set using
// [Terminal.MatchAddRegex]. If a match exists, the text string and handle is
// returned, otherwise MatchCheckEvent returns an error.
//
// Typically, event should be received from the "button-press-event" signal.
// MatchCheckEvent allows to determine whether one of the matching strings were
// clicked. A callback may be provided depending on the returned [MatchHandle],
// e.g. if user clicked on a URL, open it in a web browser.
//
// If more than one regular expression has been set with
// [Terminal.MatchAddRegex], then expressions are checked in the order in which
// they were added.
func (t *Terminal) MatchCheckEvent(event *gdk.Event) (string, MatchHandle, error) {
	handle := C.int(0)
	cstr := C.vte_terminal_match_check_event(t.native(), (*C.GdkEvent)(event.GdkEvent), &handle)
	if cstr == nil {
		return "", -1, fmt.Errorf("no matches found")
	}
	defer func() {
		// NOTE: apparently cstr is freed, even though the caller takes ownership of
		// the returned string according to the documentation.
		if unsafe.Pointer(cstr) != C.NULL {
			C.free(unsafe.Pointer(cstr))
		}
	}()
	return goString(cstr), MatchHandle(handle), nil
}

// SearchFindNext searches the next string matching the search regex set with
// [SearchSetRegex].
func (t *Terminal) SearchFindNext() bool {
	return goBool(C.vte_terminal_search_find_next(t.native()))
}

// SearchFindPrev searches the previous string matching the search regex set
// with [SearchSetRegex].
func (t *Terminal) SearchFindPrev() bool {
	return goBool(C.vte_terminal_search_find_previous(t.native()))
}

// SearchGetRegex returns [Regex] used for search, or nil if it is unset.
func (t *Terminal) SearchGetRegex() *Regex {
	ptr := C.vte_terminal_search_get_regex(t.native())
	if ptr == nil {
		return nil
	}
	return wrapRegex(ptr)
}

// SearchSetRegex sets [Regex] for search. Use nil to reset regex.
func (t *Terminal) SearchSetRegex(regex *Regex, flags RegexMatchFlags) {
	var r *C.VteRegex
	if regex != nil {
		r = regex.ptr
	}
	C.vte_terminal_search_set_regex(t.native(), r, C.uint(flags))
}

// SearchSetWrapAround controls whether [Terminal.SearchFindNext] and
// [Terminal.SearchFindPrev] wrap around to the beginning/end of the terminal
// when reaching its end/beginning.
func (t *Terminal) SearchSetWrapAround(v bool) {
	C.vte_terminal_search_set_wrap_around(t.native(), gboolean(v))
}

// ConnectBell calls callback when the a child sends a bell request to the
// terminal.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectBell(callback func(t *Terminal)) glib.SignalHandle {
	return t.Connect("bell", t.signalCb(callback))
}

// ConnectAfterBell is like [Terminal.ConnectBell], but is invoked after the
// default handler.
func (t *Terminal) ConnectAfterBell(callback func(t *Terminal)) glib.SignalHandle {
	return t.ConnectAfter("bell", t.signalCb(callback))
}

// ConnectCellSizeChanged calls callback when the cell size changes, e.g. due
// to a change in font, font-scale or cell-width/height-scale.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectCellSizeChanged(callback func(t *Terminal, width, height uint)) glib.SignalHandle {
	return t.Connect("char-size-changed", t.signalCbUU(callback))
}

// ConnectAfterCellSizeChanged is like [Terminal.ConnectCellSizeChanged], but
// is invoked after the default handler.
func (t *Terminal) ConnectAfterCellSizeChanged(callback func(t *Terminal, width, height uint)) glib.SignalHandle {
	return t.ConnectAfter("char-size-changed", t.signalCbUU(callback))
}

// ConnectChildExited calls callback when the child has exited.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectChildExited(callback func(t *Terminal, status int)) glib.SignalHandle {
	return t.Connect("child-exited", t.signalCbI(callback))
}

// ConnectAfterChildExited is like [Terminal.ConnectChildExited], but is
// invoked after the default handler.
func (t *Terminal) ConnectAfterChildExited(callback func(t *Terminal, status int)) glib.SignalHandle {
	return t.ConnectAfter("child-exited", t.signalCbI(callback))
}

// ConnectCommit calls callback when the terminal receives input from the user
// and prepares to send it to the child process..
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectCommit(callback func(t *Terminal, text string)) glib.SignalHandle {
	return t.ConnectAfter("commit", t.signalCbS(callback))
}

// ConnectAfterCommit is like [Terminal.ConnectCommit], but is invoked after
// the default handler.
func (t *Terminal) ConnectAfterCommit(callback func(t *Terminal, text string)) glib.SignalHandle {
	return t.ConnectAfter("commit", t.signalCbS(callback))
}

// ConnectContentsChanged calls callback when the visible appearance of the
// terminal has changed.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectContentsChanged(callback func(t *Terminal)) glib.SignalHandle {
	return t.Connect("contents-changed", t.signalCb(callback))
}

// ConnectAfterContentsChanged is like [Terminal.ConnectContentsChanged], but
// is invoked after the default handler.
func (t *Terminal) ConnectAfterContentsChanged(callback func(t *Terminal)) glib.SignalHandle {
	return t.ConnectAfter("contents-changed", t.signalCb(callback))
}

// ConnectCopyClipboard calls callback when [Terminal.CopyClipboardFormat] is
// called.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectCopyClipboard(callback func(t *Terminal)) glib.SignalHandle {
	return t.Connect("copy-clipboard", t.signalCb(callback))
}

// ConnectAfterCopyClipboard is like [Terminal.ConnectCopyClipboard], but is
// invoked after the default handler.
func (t *Terminal) ConnectAfterCopyClipboard(callback func(t *Terminal)) glib.SignalHandle {
	return t.ConnectAfter("copy-clipboard", t.signalCb(callback))
}

// ConnectCursorMoved calls callback when the cursor moves to a new character
// cell.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectCursorMoved(callback func(t *Terminal)) glib.SignalHandle {
	return t.Connect("cursor-moved", t.signalCb(callback))
}

// ConnectAfterCursorMoved is like [Terminal.ConnectCursorMoved], but is
// invoked after the default handler.
func (t *Terminal) ConnectAfterCursorMoved(callback func(t *Terminal)) glib.SignalHandle {
	return t.ConnectAfter("cursor-moved", t.signalCb(callback))
}

// ConnectDecreaseFontSize calls callback when when user hits the '-' key while
// holding the Control key.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectDecreaseFontSize(callback func(t *Terminal)) glib.SignalHandle {
	return t.Connect("decrease-font-size", t.signalCb(callback))
}

// ConnectAfterDecreaseFontSize is like [Terminal.ConnectDecreaseFontSize], but
// is invoked after the default handler.
func (t *Terminal) ConnectAfterDecreaseFontSize(callback func(t *Terminal)) glib.SignalHandle {
	return t.ConnectAfter("decrease-font-size", t.signalCb(callback))
}

// ConnectEOF calls callback when the terminal receives an end-of-file from a
// child which is running in the terminal. This signal is frequently (but not
// always) emitted with a child-exited signal.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectEOF(callback func(t *Terminal)) glib.SignalHandle {
	return t.Connect("eof", t.signalCb(callback))
}

// ConnectAfterEOF is like [Terminal.ConnectEOF], but is invoked after the
// default handler.
func (t *Terminal) ConnectAfterEOF(callback func(t *Terminal)) glib.SignalHandle {
	return t.ConnectAfter("eof", t.signalCb(callback))
}

// ConnectHyperlinkHoverURIChanged calls callback when the hovered hyperlink
// changes.
//
// The signal is not re-emitted when the bounding box changes for the same
// hyperlink.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectHyperlinkHoverURIChanged(callback func(
	t *Terminal,
	uri string,
	rect *gdk.Rectangle,
)) glib.SignalHandle {
	return t.Connect("hyperlink-hover-uri-changed", t.signalCbSRect(callback))
}

// ConnectAfterHyperlinkHoverURIChanged is like
// [Terminal.ConnectHyperlinkHoverURIChanged], but is invoked after the default
// handler.
func (t *Terminal) ConnectAfterHyperlinkHoverURIChanged(callback func(
	t *Terminal,
	uri string,
	rect *gdk.Rectangle,
)) glib.SignalHandle {
	return t.ConnectAfter("hyperlink-hover-uri-changed", t.signalCbSRect(callback))
}

// ConnectIncreaseFontSize calls callback when user hits the '+' key while
// holding the Control key.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectIncreaseFontSize(callback func(t *Terminal)) glib.SignalHandle {
	return t.Connect("increase-font-size", t.signalCb(callback))
}

// ConnectAfterIncreaseFontSize is like [Terminal.ConnectIncreaseFontSize], but
// is invoked after the default handler.
func (t *Terminal) ConnectAfterIncreaseFontSize(callback func(t *Terminal)) glib.SignalHandle {
	return t.ConnectAfter("increase-font-size", t.signalCb(callback))
}

// ConnectPasteClipboard calls callback when [Terminal.PasteClipboard] is called.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectPasteClipboard(callback func(t *Terminal)) glib.SignalHandle {
	return t.Connect("paste-clipboard", t.signalCb(callback))
}

// ConnectAfterPasteClipboard is like [Terminal.ConnectPasteClipboard], but is
// invoked after the default handler.
func (t *Terminal) ConnectAfterPasteClipboard(callback func(t *Terminal)) glib.SignalHandle {
	return t.ConnectAfter("paste-clipboard", t.signalCb(callback))
}

// ConnectResizeWindow calls callback when child requests a change in the
// terminal's size.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectResizeWindow(callback func(
	t *Terminal,
	width,
	height uint,
)) glib.SignalHandle {
	return t.Connect("resize-window", t.signalCbUU(callback))
}

// ConnectAfterResizeWindow is like [Terminal.ConnectResizeWindow], but is
// invoked after the default handler.
func (t *Terminal) ConnectAfterResizeWindow(callback func(
	t *Terminal,
	width,
	height uint,
)) glib.SignalHandle {
	return t.ConnectAfter("resize-window", t.signalCbUU(callback))
}

// ConnectSelectionChanged calls callback when the contents of terminal's
// selection changes.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectSelectionChanged(callback func(t *Terminal)) glib.SignalHandle {
	return t.Connect("selection-changed", t.signalCb(callback))
}

// ConnectAfterSelectionChanged is like [Terminal.ConnectSelectionChanged], but
// is invoked after the default handler.
func (t *Terminal) ConnectAfterSelectionChanged(callback func(t *Terminal)) glib.SignalHandle {
	return t.ConnectAfter("selection-changed", t.signalCb(callback))
}

// ConnectSetupContextMenu calls callback when terminal shows a context menu.
//
// Context menu may be set with [Terminal.SetContextMenu] or
// [Terminal.SetContextMenuModel].
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectSetupContextMenu(callback func(t *Terminal)) glib.SignalHandle {
	return t.Connect("setup-context-menu", t.signalCb(callback))
}

// ConnectAfterSetupContextMenu is like [Terminal.ConnectSetupContextMenu], but
// is invoked after the default handler.
func (t *Terminal) ConnectAfterSetupContextMenu(callback func(t *Terminal)) glib.SignalHandle {
	return t.ConnectAfter("setup-context-menu", t.signalCb(callback))
}

// ConnectTermPropChanged calls callback when a termprop has changed or been
// reset.
//
// The handler may use [Terminal.GetTermProp] to retrieve the value of any
// termprop (not just prop) It must not call any other API on terminal,
// including API of its parent classes.
//
// See [github.com/gotk3/gotk3/glib.Object.Connect] for more information about
// signal handling.
func (t *Terminal) ConnectTermPropChanged(callback func(t *Terminal, prop TermProp)) glib.SignalHandle {
	return t.Connect("termprop-changed", t.signalCbProp(callback))
}

// ConnectAfterTermPropChanged is like [Terminal.ConnectTermPropChanged], but
// is invoked after the default handler.
func (t *Terminal) ConnectAfterTermPropChanged(callback func(t *Terminal, prop TermProp)) glib.SignalHandle {
	return t.ConnectAfter("termprop-changed", t.signalCbProp(callback))
}

func (t *Terminal) signalCb(cb func(*Terminal)) any {
	return func(o *glib.Object, s string) {
		term := WrapTerminal(o)
		if term != nil {
			cb(term)
		}
	}
}

func (t *Terminal) signalCbI(cb func(*Terminal, int)) any {
	return func(o *glib.Object, i int) {
		term := WrapTerminal(o)
		if term != nil {
			cb(term, i)
		}
	}
}

func (t *Terminal) signalCbS(cb func(*Terminal, string)) any {
	return func(o *glib.Object, s string) {
		term := WrapTerminal(o)
		if term != nil {
			cb(term, s)
		}
	}
}

func (t *Terminal) signalCbProp(cb func(*Terminal, TermProp)) any {
	return func(o *glib.Object, s string) {
		term := WrapTerminal(o)
		if term != nil {
			cb(term, TermProp(s))
		}
	}
}

func (t *Terminal) signalCbSRect(cb func(*Terminal, string, *gdk.Rectangle)) any {
	return func(o *glib.Object, s string, r uintptr) {
		term := WrapTerminal(o)
		if term != nil {
			cb(term, s, gdk.WrapRectangle(r))
		}
	}
}

func (t *Terminal) signalCbUU(cb func(*Terminal, uint, uint)) any {
	return func(o *glib.Object, w, h uint) {
		term := WrapTerminal(o)
		if term != nil {
			cb(term, w, h)
		}
	}
}
