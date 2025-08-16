package vte_test

import (
	"fmt"
	"math/rand/v2"
	"os"
	"testing"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/shelepuginivan/gotk3-vte/vte"
	"github.com/stretchr/testify/assert"
)

func newTerm(t *testing.T) *vte.Terminal {
	term, err := vte.TerminalNew()
	assert.NoError(t, err)
	return term
}

func TestTerminal_PropertyAllowHyperlink(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, false, term.GetAllowHyperlink())
	term.SetAllowHyperlink(true)
	assert.Equal(t, true, term.GetAllowHyperlink())
}

func TestTerminal_PropertyAudibleBell(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, true, term.GetAudibleBell())
	term.SetAllowHyperlink(false)
	assert.Equal(t, false, term.GetAllowHyperlink())
}

func TestTerminal_PropertyBackspaceBinding(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, vte.ERASE_AUTO, term.GetBackspaceBinding())
	term.SetBackspaceBinding(vte.ERASE_ASCII_BACKSPACE)
	assert.Equal(t, vte.ERASE_ASCII_BACKSPACE, term.GetBackspaceBinding())
}

func TestTerminal_PropertyBoldIsBright(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, false, term.GetBoldIsBright())
	term.SetBoldIsBright(true)
	assert.Equal(t, true, term.GetBoldIsBright())
}

func TestTerminal_PropertyCellHeightScale(t *testing.T) {
	term := newTerm(t)

	assert.InDelta(t, 1.000000, term.GetCellHeightScale(), 0.00001)
	term.SetCellHeightScale(2.0)
	assert.InDelta(t, 2.000000, term.GetCellHeightScale(), 0.00001)
}

func TestTerminal_PropertyCellWidthScale(t *testing.T) {
	term := newTerm(t)

	assert.InDelta(t, 1.000000, term.GetCellWidthScale(), 0.00001)
	term.SetCellWidthScale(2.0)
	assert.InDelta(t, 2.000000, term.GetCellWidthScale(), 0.00001)
}

func TestTerminal_PropertyCJKAmbiguousWidth(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, vte.CJK_AMBIGUOUS_WIDTH_NARROW, term.GetCJKAmbiguousWidth())
	term.SetCJKAmbiguousWidth(vte.CJK_AMBIGUOUS_WIDTH_WIDE)
	assert.Equal(t, vte.CJK_AMBIGUOUS_WIDTH_WIDE, term.GetCJKAmbiguousWidth())
}

func TestTerminal_PropertyContextMenu(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, nil, term.GetContextMenu())

	menu, err := gtk.MenuNew()
	assert.NoError(t, err)

	term.SetContextMenu(menu)
	assert.Equal(t, menu.Native(), term.GetContextMenu().ToWidget().Native())

	term.SetContextMenu(nil)
	assert.Equal(t, nil, term.GetContextMenu())

	button, err := gtk.ButtonNew()
	assert.NoError(t, err)

	term.SetContextMenu(button)
	assert.Equal(t, nil, term.GetContextMenu())
}

func TestTerminal_PropertyContextMenuModel(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, (*glib.MenuModel)(nil), term.GetContextMenuModel())

	b, err := gtk.MenuButtonNew()
	assert.NoError(t, err)

	model := b.GetMenuModel()

	term.SetContextMenuModel(model)
	assert.Equal(t, model.Native(), term.GetContextMenuModel().Native())

	term.SetContextMenuModel(nil)
	assert.Equal(t, (*glib.MenuModel)(nil), term.GetContextMenuModel())
}

func TestTerminal_PropertyCursorBlinkMode(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, vte.CURSOR_BLINK_SYSTEM, term.GetCursorBlinkMode())
	term.SetCursorBlinkMode(vte.CURSOR_BLINK_OFF)
	assert.Equal(t, vte.CURSOR_BLINK_OFF, term.GetCursorBlinkMode())
}

func TestTerminal_PropertyCursorShape(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, vte.CURSOR_SHAPE_BLOCK, term.GetCursorShape())
	term.SetCursorShape(vte.CURSOR_SHAPE_IBEAM)
	assert.Equal(t, vte.CURSOR_SHAPE_IBEAM, term.GetCursorShape())
}

func TestTerminal_PropertyDeleteBinding(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, vte.ERASE_AUTO, term.GetDeleteBinding())
	term.SetDeleteBinding(vte.ERASE_ASCII_BACKSPACE)
	assert.Equal(t, vte.ERASE_ASCII_BACKSPACE, term.GetDeleteBinding())
}

func TestTerminal_PropertyEnableA11y(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, true, term.GetEnableA11y())
	term.SetEnableA11y(false)
	assert.Equal(t, false, term.GetEnableA11y())
}

func TestTerminal_PropertyEnableBidi(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, true, term.GetEnableBidi())
	term.SetEnableBidi(false)
	assert.Equal(t, false, term.GetEnableBidi())
}

func TestTerminal_PropertyEnableFallbackScrolling(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, true, term.GetEnableFallbackScrolling())
	term.SetEnableFallbackScrolling(false)
	assert.Equal(t, false, term.GetEnableFallbackScrolling())
}

func TestTerminal_PropertyEnableLegacyOSC777(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, false, term.GetEnableLegacyOSC777())
	term.SetEnableLegacyOSC777(true)
	assert.Equal(t, true, term.GetEnableLegacyOSC777())
}

func TestTerminal_PropertyEnableShaping(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, true, term.GetEnableShaping())
	term.SetEnableShaping(false)
	assert.Equal(t, false, term.GetEnableShaping())
}

func TestTerminal_PropertyEnableSixel(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, false, term.GetEnableSixel())
	term.SetEnableSixel(true)
	assert.Equal(t, false, term.GetEnableSixel())
}

func TestTerminal_PropertyFontDesc(t *testing.T) {
	term := newTerm(t)

	expected := pango.FontDescriptionNew()
	expected.SetFamily("sans-serif")
	expected.SetSize(12 * 1024)

	term.SetFont(expected)
	actual := term.GetFont()

	assert.Equal(t, expected.GetFamily(), actual.GetFamily())
	assert.Equal(t, expected.GetSize(), actual.GetSize())
	assert.Equal(t, expected.GetGravity(), actual.GetGravity())
	assert.Equal(t, expected.GetSizeIsAbsolute(), actual.GetSizeIsAbsolute())
	assert.Equal(t, expected.GetStretch(), actual.GetStretch())
	assert.Equal(t, expected.GetStyle(), actual.GetStyle())
	assert.Equal(t, expected.GetFamily(), actual.GetFamily())
	assert.Equal(t, expected.GetWeight(), actual.GetWeight())

	term.SetFontFromString("monospace Bold 16")
	actual = term.GetFont()

	assert.Equal(t, "monospace", actual.GetFamily())
	assert.Equal(t, 16*1024, actual.GetSize())
	assert.Equal(t, pango.WEIGHT_BOLD, actual.GetWeight())
}

func TestTerminal_PropertyFontOptions(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, (*cairo.FontOptions)(nil), term.GetFontOptions())

	expected := cairo.CreateFontOptions()
	expected.SetAntialias(cairo.ANTIALIAS_BEST)
	expected.SetHintStyle(cairo.HINT_STYLE_MEDIUM)

	term.SetFontOptions(expected)
	actual := term.GetFontOptions()

	assert.Equal(t, expected.GetAntialias(), actual.GetAntialias())
	assert.Equal(t, expected.GetHintMetrics(), actual.GetHintMetrics())
	assert.Equal(t, expected.GetHintStyle(), actual.GetHintStyle())
	assert.Equal(t, expected.GetSubpixelOrder(), actual.GetSubpixelOrder())
	assert.Equal(t, expected.GetVariations(), actual.GetVariations())
}

func TestTerminal_PropertyFontScale(t *testing.T) {
	term := newTerm(t)

	assert.InDelta(t, 1.000000, term.GetFontScale(), 0.00001)
	term.SetFontScale(2.0)
	assert.InDelta(t, 2.000000, term.GetFontScale(), 0.00001)
}

func TestTerminal_PropertyInputEnabled(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, true, term.GetInputEnabled())
	term.SetInputEnabled(false)
	assert.Equal(t, false, term.GetInputEnabled())
}

func TestTerminal_PropertyPointerAutohide(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, false, term.GetPointerAutohide())
	term.SetPointerAutohide(true)
	assert.Equal(t, true, term.GetPointerAutohide())
}

func TestTerminal_PropertyPty(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, (*vte.Pty)(nil), term.GetPty())

	pty := newPty(t)
	term.SetPty(pty)

	assert.Equal(t, pty.Native(), term.GetPty().Native())

	term.SetPty(nil)
	assert.Equal(t, pty.Native(), term.GetPty().Native())
}

func TestTerminal_PropertyScrollOnInsert(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, false, term.GetScrollOnInsert())
	term.SetScrollOnInsert(true)
	assert.Equal(t, true, term.GetScrollOnInsert())
}

func TestTerminal_PropertyScrollOnKeystroke(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, true, term.GetScrollOnKeystroke())
	term.SetScrollOnKeystroke(false)
	assert.Equal(t, false, term.GetScrollOnKeystroke())
}

func TestTerminal_PropertyScrollOnOutput(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, false, term.GetScrollOnOutput())
	term.SetScrollOnOutput(true)
	assert.Equal(t, true, term.GetScrollOnOutput())
}

func TestTerminal_PropertyScrollUnitIsPixels(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, vte.SCROLL_UNIT_LINES, term.GetScrollUnit())
	term.SetScrollUnit(vte.SCROLL_UNIT_PIXELS)
	assert.Equal(t, vte.SCROLL_UNIT_PIXELS, term.GetScrollUnit())
}

func TestTerminal_PropertyScrollbackLines(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, uint(512), term.GetScrollbackLines())
	term.SetScrollbackLines(1024)
	assert.Equal(t, uint(1024), term.GetScrollbackLines())
}

func TestTerminal_PropertyTextBlinkMode(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, vte.TEXT_BLINK_ALWAYS, term.GetTextBlinkMode())
	term.SetTextBlinkMode(vte.TEXT_BLINK_FOCUSED)
	assert.Equal(t, vte.TEXT_BLINK_FOCUSED, term.GetTextBlinkMode())
}

func TestTerminal_PropertyXAlign(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, vte.ALIGN_START, term.GetXAlign())
	term.SetXAlign(vte.ALIGN_CENTER)
	assert.Equal(t, vte.ALIGN_CENTER, term.GetXAlign())
}

func TestTerminal_PropertyXFill(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, true, term.GetXFill())
	term.SetXFill(false)
	assert.Equal(t, false, term.GetXFill())
}

func TestTerminal_PropertyYAlign(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, vte.ALIGN_START, term.GetYAlign())
	term.SetYAlign(vte.ALIGN_END)
	assert.Equal(t, vte.ALIGN_END, term.GetYAlign())
}

func TestTerminal_PropertyYFill(t *testing.T) {
	term := newTerm(t)

	assert.Equal(t, true, term.GetYFill())
	term.SetYFill(false)
	assert.Equal(t, false, term.GetYFill())
}

func TestTerminal_Spawn(t *testing.T) {
	gtk.Init(nil)

	cmd := vte.CommandNew(
		[]string{"/usr/bin/false"},
		vte.CommandWithOnSpawn(func(pid int, err error) {
			assert.NoError(t, err)
			assert.Greater(t, pid, os.Getpid())
			gtk.MainQuit()
		}),
	)

	newTerm(t).Spawn(cmd)

	// This will block. Unless gtk.MainQuit is called in the OnSpawn callback, the
	// test will timeout after 10 minutes.
	gtk.Main()
}

func TestTerminal_SetColors(t *testing.T) {
	term := newTerm(t)

	assert.NoError(t, term.SetColors(
		nil,
		nil,
		nil,
	))

	assert.NoError(t, term.SetColors(
		gdk.NewRGBA(0, 0, 0, 1),
		gdk.NewRGBA(1, 1, 1, 1),
		nil,
	))

	assert.NoError(t, term.SetColors(
		gdk.NewRGBA(0, 0, 0, 1),
		gdk.NewRGBA(1, 1, 1, 1),
		nil,
	))

	assert.NoError(t, term.SetColors(
		gdk.NewRGBA(0.157, 0.157, 0.157, 1),
		gdk.NewRGBA(0.659, 0.600, 0.518, 1),
		[]*gdk.RGBA{
			gdk.NewRGBA(0.157, 0.157, 0.157, 1),
			gdk.NewRGBA(0.800, 0.141, 0.114, 1),
			gdk.NewRGBA(0.596, 0.592, 0.102, 1),
			gdk.NewRGBA(0.843, 0.600, 0.129, 1),
			gdk.NewRGBA(0.271, 0.522, 0.533, 1),
			gdk.NewRGBA(0.694, 0.384, 0.525, 1),
			gdk.NewRGBA(0.408, 0.616, 0.416, 1),
			gdk.NewRGBA(0.659, 0.600, 0.518, 1),
		},
	))

	assert.NoError(t, term.SetColors(
		nil,
		nil,
		[]*gdk.RGBA{
			gdk.NewRGBA(0.157, 0.157, 0.157, 1),
			gdk.NewRGBA(0.800, 0.141, 0.114, 1),
			gdk.NewRGBA(0.596, 0.592, 0.102, 1),
			gdk.NewRGBA(0.843, 0.600, 0.129, 1),
			gdk.NewRGBA(0.271, 0.522, 0.533, 1),
			gdk.NewRGBA(0.694, 0.384, 0.525, 1),
			gdk.NewRGBA(0.408, 0.616, 0.416, 1),
			gdk.NewRGBA(0.659, 0.600, 0.518, 1),
		},
	))

	assert.Error(t, term.SetColors(
		gdk.NewRGBA(0.157, 0.157, 0.157, 1),
		gdk.NewRGBA(0.659, 0.600, 0.518, 1),
		[]*gdk.RGBA{
			gdk.NewRGBA(0.157, 0.157, 0.157, 1),
			gdk.NewRGBA(0.800, 0.141, 0.114, 1),
			gdk.NewRGBA(0.843, 0.600, 0.129, 1),
		},
	))
}

func TestTerminal_SearchRegex(t *testing.T) {
	term := newTerm(t)

	reg, err := vte.RegexNew("bar")
	assert.NoError(t, err)

	assert.Equal(t, (*vte.Regex)(nil), term.SearchGetRegex())
	term.SearchSetRegex(reg, 0)
	assert.Equal(t, reg.Native(), term.SearchGetRegex().Native())

	invalid, err := vte.RegexNew(".", vte.RegexWithPurpose(vte.REGEX_PURPOSE_MATCH))
	assert.NoError(t, err)
	term.SearchSetRegex(invalid, 0)
	assert.Equal(t, reg.Native(), term.SearchGetRegex().Native())
}

func TestTerminal_Match(t *testing.T) {
	term := newTerm(t)

	reg0, err := vte.RegexNew(".", vte.RegexWithPurpose(vte.REGEX_PURPOSE_MATCH))
	assert.NoError(t, err)

	reg1, err := vte.RegexNew(".", vte.RegexWithPurpose(vte.REGEX_PURPOSE_MATCH))
	assert.NoError(t, err)

	reg2, err := vte.RegexNew(".", vte.RegexWithPurpose(vte.REGEX_PURPOSE_MATCH))
	assert.NoError(t, err)

	h, err := term.MatchAddRegex(nil, 0)
	assert.Equal(t, vte.MatchHandle(-1), h)
	assert.Error(t, err)

	h, err = term.MatchAddRegex(reg0, 0)
	assert.Equal(t, vte.MatchHandle(0), h)
	assert.NoError(t, err)

	h, err = term.MatchAddRegex(reg1, 0)
	assert.Equal(t, vte.MatchHandle(1), h)
	assert.NoError(t, err)

	h, err = term.MatchAddRegex(reg2, 0)
	assert.Equal(t, vte.MatchHandle(2), h)
	assert.NoError(t, err)

	invalid, err := vte.RegexNew(".", vte.RegexWithPurpose(vte.REGEX_PURPOSE_SEARCH))
	assert.NoError(t, err)
	h, err = term.MatchAddRegex(invalid, 0)
	assert.Equal(t, vte.MatchHandle(-1), h)
	assert.Error(t, err)
}

func TestTerminal_SignalBell(t *testing.T) {
	gtk.Init(nil)

	var callStack []int

	term := newTerm(t)

	cmd := vte.CommandNew([]string{"/usr/bin/env", "echo", "-ne", `\007`})

	term.ConnectBell(func(_ *vte.Terminal) {
		callStack = append(callStack, 0)
	})

	term.ConnectAfterBell(func(_ *vte.Terminal) {
		callStack = append(callStack, 1)
		gtk.MainQuit()
	})

	term.Spawn(cmd)

	gtk.Main()

	assert.Len(t, callStack, 2)
	assert.IsIncreasing(t, callStack)
}

func TestTerminal_SignalChildExited(t *testing.T) {
	gtk.Init(nil)

	var callStack []int

	cmd := vte.CommandNew([]string{"/usr/bin/env", "true"})

	term := newTerm(t)

	term.ConnectChildExited(func(t *vte.Terminal, status int) {
		callStack = append(callStack, 0)
	})

	term.ConnectAfterChildExited(func(t *vte.Terminal, status int) {
		callStack = append(callStack, 1)
		gtk.MainQuit()
	})

	term.Spawn(cmd)

	gtk.Main()

	assert.Len(t, callStack, 2)
	assert.IsIncreasing(t, callStack)
}

func TestTerminal_SignalContentsChanged(t *testing.T) {
	gtk.Init(nil)

	var callStack []int

	term := newTerm(t)

	term.ConnectContentsChanged(func(t *vte.Terminal) {
		callStack = append(callStack, 0)
	})

	term.ConnectAfterContentsChanged(func(t *vte.Terminal) {
		callStack = append(callStack, 1)
		gtk.MainQuit()
	})

	term.Feed("some")

	gtk.Main()

	assert.Len(t, callStack, 2)
	assert.IsIncreasing(t, callStack)
}

func TestTerminal_SignalEOF(t *testing.T) {
	gtk.Init(nil)

	var callStack []int

	cmd := vte.CommandNew([]string{"/usr/bin/env", "true"})

	term := newTerm(t)

	term.ConnectEOF(func(t *vte.Terminal) {
		callStack = append(callStack, 0)
	})

	term.ConnectAfterEOF(func(t *vte.Terminal) {
		callStack = append(callStack, 1)
		gtk.MainQuit()
	})

	term.Spawn(cmd)

	gtk.Main()

	assert.Len(t, callStack, 2)
	assert.IsIncreasing(t, callStack)
}

func TestTerminal_SignalTermPropChanged(t *testing.T) {
	gtk.Init(nil)

	var callStack []int

	term := newTerm(t)

	cmd := vte.CommandNew([]string{"/usr/bin/env", "echo", "-ne", `\e]0;something\a`})

	term.ConnectTermPropChanged(func(_ *vte.Terminal, prop vte.TermProp) {
		assert.Equal(t, vte.TERMPROP_XTERM_TITLE, prop)
		callStack = append(callStack, 0)
	})

	term.ConnectAfterTermPropChanged(func(_ *vte.Terminal, prop vte.TermProp) {
		assert.Equal(t, vte.TERMPROP_XTERM_TITLE, prop)
		callStack = append(callStack, 1)
		gtk.MainQuit()
	})

	term.Spawn(cmd)

	gtk.Main()

	assert.Len(t, callStack, 2)
	assert.IsIncreasing(t, callStack)
}

func TestTerminal_SignalResizeWindow(t *testing.T) {
	gtk.Init(nil)

	var callStack []int

	term := newTerm(t)

	expectedWidth := rand.UintN(80) + 20
	expectedHeight := rand.UintN(80) + 20

	cmd := vte.CommandNew([]string{"/usr/bin/env", "echo", "-ne", fmt.Sprintf(`\e[8;%d;%dt`, expectedHeight, expectedWidth)})

	term.ConnectResizeWindow(func(_ *vte.Terminal, w, h uint) {
		assert.Equal(t, expectedWidth, w)
		assert.Equal(t, expectedHeight, h)
		callStack = append(callStack, 0)
	})

	term.ConnectAfterResizeWindow(func(_ *vte.Terminal, w, h uint) {
		assert.Equal(t, expectedWidth, w)
		assert.Equal(t, expectedHeight, h)
		callStack = append(callStack, 1)
		gtk.MainQuit()
	})

	term.Spawn(cmd)

	gtk.Main()

	assert.Len(t, callStack, 2)
	assert.IsIncreasing(t, callStack)
}
