package goome

/*
#include <stdlib.h>
#include "native/goome.h"
#include "helpers.h"
*/
import "C"

import (
	"runtime/cgo"
	"sync"
	"unsafe"
)

// Skin is a nine-slice bitmap that replaces how a widget part is drawn. The zero Skin means the normal look.
type Skin struct {
	bm unsafe.Pointer `c:"bm"`
	// L is the left border size in bitmap pixels.
	L float64 `c:"l"`
	// T is the top border size in bitmap pixels.
	T float64 `c:"t"`
	// R is the right border size in bitmap pixels.
	R float64 `c:"r"`
	// B is the bottom border size in bitmap pixels.
	B float64 `c:"b"`
}

// NewSkin makes a nine-slice skin from b with l, t, r and bottom as the border sizes in bitmap pixels. Keep the bitmap alive while a theme uses the skin.
func NewSkin(b *Bitmap, l, t, r, bottom float64) Skin {
	if b == nil {
		return Skin{}
	}
	pin(unsafe.Pointer(b.ptr()), b)
	return Skin{bm: unsafe.Pointer(b.b), L: l, T: t, R: r, B: bottom}
}

// Empty reports whether the skin has no bitmap, so the normal look is used.
func (s Skin) Empty() bool { return s.bm == nil }

// Theme holds the colors, sizes and fonts of the UI kit. Colors are 0xAARRGGBB. A 0 color means work it out from the other colors, and a 0 height means work it out from the font.
type Theme struct {
	// Bg is the color behind everything, app mode fills with it.
	Bg uint32 `c:"bg"`
	// Panel is the window background color.
	Panel uint32 `c:"panel"`
	// Border is the outline color.
	Border uint32 `c:"border"`
	// Text is the normal text color.
	Text uint32 `c:"text"`
	// TextDim is the quieter text color.
	TextDim uint32 `c:"text_dim"`
	// Accent is the color of checked boxes, slider fill and the active tab.
	Accent uint32 `c:"accent"`
	// AccentText is the text color on top of Accent.
	AccentText uint32 `c:"accent_text"`
	// Hover is the widget background while hovered.
	Hover uint32 `c:"hover"`
	// Active is the widget background while pressed.
	Active uint32 `c:"active"`
	// InputBg is the color inside inputs and sliders.
	InputBg uint32 `c:"input_bg"`
	// Shadow is the window shadow color.
	Shadow uint32 `c:"shadow"`
	// Radius is the corner radius.
	Radius float64 `c:"radius"`
	// Pad is the space inside windows.
	Pad float64 `c:"pad"`
	// Gap is the space between widgets.
	Gap float64 `c:"gap"`
	// BorderWidth is the outline width.
	BorderWidth float64 `c:"border_width"`
	// FontSize is the text size.
	FontSize float64 `c:"font_size"`
	// ShadowBlur is how soft shadows are.
	ShadowBlur float64 `c:"shadow_blur"`
	// AnimMs is the animation length in milliseconds, 0 means no animation.
	AnimMs     float64                   `c:"anim_ms"`
	font       unsafe.Pointer            `c:"font"`
	skinButton Skin                      `c:"skin_button"`
	skinPanel  Skin                      `c:"skin_panel"`
	skinInput  Skin                      `c:"skin_input"`
	icons      [IconCount]unsafe.Pointer `c:"icon"`
	// IconTint set to 1 draws icons as they are, with no tint. 0 tints them with the color the glyph would have had.
	IconTint       uint32 `c:"icon_tint"`
	skinCheckbox   Skin   `c:"skin_checkbox"`
	skinRadio      Skin   `c:"skin_radio"`
	skinToggle     Skin   `c:"skin_toggle"`
	skinSlider     Skin   `c:"skin_slider"`
	skinSliderFill Skin   `c:"skin_slider_fill"`
	skinScrollbar  Skin   `c:"skin_scrollbar"`
	skinTab        Skin   `c:"skin_tab"`
	skinTitlebar   Skin   `c:"skin_titlebar"`
	skinTooltip    Skin   `c:"skin_tooltip"`
	skinMenu       Skin   `c:"skin_menu"`
	skinPopup      Skin   `c:"skin_popup"`
	// Focus is the focus ring and input glow color, 0 uses Accent.
	Focus uint32 `c:"focus"`
	// Select is the selected text color, 0 uses a faded Accent.
	Select uint32 `c:"select"`
	// Caret is the text caret color, 0 uses Accent.
	Caret uint32 `c:"caret"`
	// Link is the color of links, 0 uses Accent.
	Link uint32 `c:"link"`
	// Badge is the badge color when you pass 0, 0 uses Accent.
	Badge uint32 `c:"badge"`
	// Toast is the toast background color, 0 mixes Panel with Hover.
	Toast uint32 `c:"toast"`
	// Tooltip is the tooltip background color, 0 mixes Bg with Hover.
	Tooltip uint32 `c:"tooltip"`
	// Danger is the color of the red close buttons, 0 uses the built-in red.
	Danger uint32 `c:"danger"`
	// Knob is the color of slider and toggle knobs, 0 is white.
	Knob uint32 `c:"knob"`
	// Dim is the color behind modals, 0 uses 0x73000000.
	Dim uint32 `c:"dim"`
	// CheckSize is the size of the checkbox and radio box. It can't be 0 or less, that goes back to the default.
	CheckSize float64 `c:"check_size"`
	// CheckRadius is the corner radius of checkboxes.
	CheckRadius float64 `c:"check_radius"`
	// ToggleW is the width of the toggle track. It can't be 0 or less, that goes back to the default.
	ToggleW float64 `c:"toggle_w"`
	// ToggleH is the height of the toggle track. It can't be 0 or less, that goes back to the default.
	ToggleH float64 `c:"toggle_h"`
	// SliderH is the slider height. It can't be 0 or less, that goes back to the default.
	SliderH float64 `c:"slider_h"`
	// SliderTrack is the thickness of the slider track. It can't be 0 or less, that goes back to the default.
	SliderTrack float64 `c:"slider_track"`
	// KnobSize is the slider knob radius. It can't be 0 or less, that goes back to the default.
	KnobSize float64 `c:"knob_size"`
	// ScrollW is the scrollbar width. It can't be 0 or less, that goes back to the default.
	ScrollW float64 `c:"scroll_w"`
	// ScrollMin is the shortest a scrollbar thumb gets. It can't be 0 or less, that goes back to the default.
	ScrollMin float64 `c:"scroll_min"`
	// ItemH is the height of every widget row, 0 means line height plus Pad.
	ItemH float64 `c:"item_h"`
	// MenuItemH is the height of menu rows, 0 comes from the font.
	MenuItemH float64 `c:"menu_item_h"`
	// TabH is the height of tabs, 0 comes from the font.
	TabH float64 `c:"tab_h"`
	// WinTitleH is the height of the window title bar, 0 comes from the font.
	WinTitleH float64 `c:"win_title_h"`
	// AppTitleH is the height of the app mode title bar, 0 comes from the font.
	AppTitleH float64 `c:"app_title_h"`
	// MenubarH is the height of the app menubar, 0 comes from the font.
	MenubarH float64 `c:"menubar_h"`
	// StatusbarH is the height of the app statusbar, 0 comes from the font.
	StatusbarH float64 `c:"statusbar_h"`
	// FocusW is the width of the focus ring.
	FocusW float64 `c:"focus_w"`
	// FocusGap is the space between a widget and its focus ring.
	FocusGap float64 `c:"focus_gap"`
	// SepW is the thickness of separators.
	SepW float64 `c:"sep_w"`
	// Grip is the size of the window resize grip. It can't be 0 or less, that goes back to the default.
	Grip float64 `c:"grip"`
	// Indent is the tree indent.
	Indent float64 `c:"indent"`
	// PlotW is the width of plot lines.
	PlotW float64 `c:"plot_w"`
	// Spinner is the spinner size when you pass 0, 0 means one line height.
	Spinner float64 `c:"spinner"`
	// CaretW is the caret width, 0 is thin. Go wide for a block caret.
	CaretW float64 `c:"caret_w"`
	// WinMinW is the smallest width a window resizes to.
	WinMinW float64 `c:"win_min_w"`
	// WinMinH is the smallest height a window resizes to.
	WinMinH float64 `c:"win_min_h"`
	// PopGap is the space between a combo or popup and the thing that opened it.
	PopGap float64 `c:"pop_gap"`
	// MenuGap is the offset of submenus.
	MenuGap float64 `c:"menu_gap"`
	// TipPadX is the horizontal padding inside tooltips.
	TipPadX float64 `c:"tip_pad_x"`
	// TipPadY is the vertical padding inside tooltips.
	TipPadY float64 `c:"tip_pad_y"`
	// TipX is how far right of the mouse a tooltip shows.
	TipX float64 `c:"tip_x"`
	// TipY is how far below the mouse a tooltip shows.
	TipY float64 `c:"tip_y"`
	// TipRadius is the corner radius of tooltips.
	TipRadius float64 `c:"tip_radius"`
	// ShadowX is the horizontal shadow offset.
	ShadowX float64 `c:"shadow_x"`
	// ShadowY is the vertical shadow offset.
	ShadowY float64 `c:"shadow_y"`
	// Align is the text alignment for buttons and labels, one of the Align constants. Auto lets each widget do its usual thing.
	Align float64 `c:"align"`
	// HeadingSize is the size of headings, 0 means 1.3 times FontSize.
	HeadingSize float64 `c:"heading_size"`
	// SmallSize is the text size of badges and the statusbar, 0 means a bit under FontSize.
	SmallSize float64 `c:"small_size"`
	// Ease is the easing curve for animations, one of the Ease constants. Use the custom one together with UIConfig.SetEase.
	Ease float64 `c:"ease"`
	// HoverMs is how long hover fades take, below 0 follows AnimMs.
	HoverMs float64 `c:"hover_ms"`
	// PressMs is how long press fades take, below 0 follows AnimMs.
	PressMs float64 `c:"press_ms"`
	// PopupMs is how long popups and menus take to open, below 0 follows AnimMs.
	PopupMs float64 `c:"popup_ms"`
	// CollapseMs is how long anything that folds open or shut animates, below 0 follows AnimMs.
	CollapseMs float64 `c:"collapse_ms"`
	// ToastFadeMs is how long toasts take to fade out, below 0 follows AnimMs.
	ToastFadeMs float64        `c:"toast_fade_ms"`
	fontHeading unsafe.Pointer `c:"font_heading"`
	fontMono    unsafe.Pointer `c:"font_mono"`
	fontSmall   unsafe.Pointer `c:"font_small"`
}

var (
	pinMu  sync.Mutex
	pinned = map[unsafe.Pointer]any{}
)

func pin(p unsafe.Pointer, owner any) {
	pinMu.Lock()
	pinned[p] = owner
	pinMu.Unlock()
}

func unpin(p unsafe.Pointer) {
	pinMu.Lock()
	delete(pinned, p)
	pinMu.Unlock()
}

func (t *Theme) c() *C.goome_theme { return (*C.goome_theme)(unsafe.Pointer(t)) }

// Dark loads the default dark theme into t. It leaves the font alone.
func (t *Theme) Dark() { C.goome_ui_theme_dark(t.c()) }

// Light loads the default light theme into t. It leaves the font alone.
func (t *Theme) Light() { C.goome_ui_theme_light(t.c()) }

// Make builds a theme from one accent color, a dark or light base and a corner radius. A radius below 0 means 10. It leaves the font and your pictures alone.
func (t *Theme) Make(accent uint32, dark bool, radius float64) {
	C.goome_ui_theme_make(t.c(), C.uint32_t(accent), b2i(dark), C.double(radius))
}

// Preset loads one of the built-in presets into t. It returns an error and leaves the theme alone if the number isn't a preset.
func (t *Theme) Preset(preset int) error {
	if r := C.goome_ui_theme_preset(t.c(), C.int(preset)); r < 0 {
		return newErr(int(r), "theres no preset %d", preset)
	}
	return nil
}

const (
	// FontMain is the font slot for normal text.
	FontMain = iota
	// FontHeading is the font slot for headings.
	FontHeading
	// FontMono is the font slot for the text editor in mono mode.
	FontMono
	// FontSmall is the font slot for badges and the statusbar.
	FontSmall
)

func (t *Theme) fontSlot(which int) *unsafe.Pointer {
	switch which {
	case FontHeading:
		return &t.fontHeading
	case FontMono:
		return &t.fontMono
	case FontSmall:
		return &t.fontSmall
	}
	return &t.font
}

// SetFont sets the font for a slot, one of FontMain, FontHeading, FontMono or FontSmall. A nil font means the built-in one. The theme keeps the font alive.
func (t *Theme) SetFont(which int, f *Font) {
	p := unsafe.Pointer(f.raw())
	if p != nil {
		pin(p, f)
	}
	*t.fontSlot(which) = p
}

// HasFont reports whether the font slot holds f.
func (t *Theme) HasFont(which int, f *Font) bool {
	return *t.fontSlot(which) == unsafe.Pointer(f.raw())
}

func (t *Theme) skinSlot(kind int) *Skin {
	switch kind {
	case SkinButton:
		return &t.skinButton
	case SkinPanel:
		return &t.skinPanel
	case SkinInput:
		return &t.skinInput
	case SkinCheckbox:
		return &t.skinCheckbox
	case SkinRadio:
		return &t.skinRadio
	case SkinToggle:
		return &t.skinToggle
	case SkinSliderTrack:
		return &t.skinSlider
	case SkinSliderFill:
		return &t.skinSliderFill
	case SkinScrollbar:
		return &t.skinScrollbar
	case SkinTab:
		return &t.skinTab
	case SkinTitlebar:
		return &t.skinTitlebar
	case SkinTooltip:
		return &t.skinTooltip
	case SkinMenu:
		return &t.skinMenu
	case SkinPopup:
		return &t.skinPopup
	}
	return nil
}

// SetSkin sets the skin for a widget part, one of the Skin constants. It returns an error for a kind that doesn't exist.
func (t *Theme) SetSkin(kind int, s Skin) error {
	slot := t.skinSlot(kind)
	if slot == nil {
		return newErr(CodeArg, "theres no skin %d", kind)
	}
	*slot = s
	return nil
}

// Skin returns the skin set for a widget part, or an empty one if there is none.
func (t *Theme) Skin(kind int) Skin {
	if slot := t.skinSlot(kind); slot != nil {
		return *slot
	}
	return Skin{}
}

// SetIcon replaces one icon with a bitmap, or puts the built-in one back when b is nil. It returns an error for an icon number that doesn't exist.
func (t *Theme) SetIcon(icon int, b *Bitmap) error {
	if icon < 0 || icon >= IconCount {
		return newErr(CodeArg, "theres no icon %d", icon)
	}
	if b == nil {
		t.icons[icon] = nil
		return nil
	}
	pin(unsafe.Pointer(b.ptr()), b)
	t.icons[icon] = unsafe.Pointer(b.b)
	return nil
}

// HasIcon reports whether that icon has a bitmap set.
func (t *Theme) HasIcon(icon int) bool { return icon >= 0 && icon < IconCount && t.icons[icon] != nil }

func (t *Theme) forget(font, bm unsafe.Pointer) {
	for _, w := range []int{FontMain, FontHeading, FontMono, FontSmall} {
		if s := t.fontSlot(w); font != nil && *s == font {
			*s = nil
		}
	}
	if bm == nil {
		return
	}
	for k := 0; k < SkinCount; k++ {
		if s := t.skinSlot(k); s.bm == bm {
			*s = Skin{}
		}
	}
	for i := range t.icons {
		if t.icons[i] == bm {
			t.icons[i] = nil
		}
	}
}

// Combo is one key and modifier pair, like ctrl+c.
type Combo struct {
	// Key is the key, a lowercase ASCII letter or one of the Key constants.
	Key int32 `c:"key"`
	// Mods is the modifier flags ored together.
	Mods int32 `c:"mods"`
}

// UIConfig holds the behavior settings of a UI, the stuff that isn't a look. There is one per UI and it can't be pushed.
type UIConfig struct {
	// TipDelayMs is the hover time in milliseconds before a tooltip shows.
	TipDelayMs int32 `c:"tip_delay_ms"`
	// DblClickMs is the longest time in milliseconds between the clicks of a double click.
	DblClickMs int32 `c:"dbl_click_ms"`
	// DblClickPx is how far the mouse can move between the clicks of a double click.
	DblClickPx int32 `c:"dbl_click_px"`
	// DragPx is how far the mouse moves before a press turns into a drag.
	DragPx int32 `c:"drag_px"`
	// RepeatDelayMs is how long a key or repeating button is held, in milliseconds, before it starts repeating.
	RepeatDelayMs int32 `c:"repeat_delay_ms"`
	// RepeatMs is the time in milliseconds between repeats once repeating has started.
	RepeatMs int32 `c:"repeat_ms"`
	// BlinkMs is the caret blink time in milliseconds, 0 means no blink.
	BlinkMs int32 `c:"blink_ms"`
	// MenuDelayMs is the hover time in milliseconds before a submenu opens.
	MenuDelayMs int32 `c:"menu_delay_ms"`
	// ToastMs is how long a toast stays in milliseconds when it was given 0.
	ToastMs int32 `c:"toast_ms"`
	// ScrollLines is how many lines one wheel step scrolls.
	ScrollLines float64 `c:"scroll_lines"`
	// RaiseOnClick set to 1 brings a window to the front when you click it.
	RaiseOnClick int32 `c:"raise_on_click"`
	// EscCloses set to 1 lets esc close the top popup or modal.
	EscCloses int32 `c:"esc_closes"`
	// TabNav set to 1 lets tab and shift+tab move focus.
	TabNav int32 `c:"tab_nav"`
	// FocusOnMouse set to 1 makes clicking a widget focus it and show the focus ring.
	FocusOnMouse int32 `c:"focus_on_mouse"`
	// Cursor is the cursor shape each kind of widget asks for, indexed by the cursor kind constants.
	Cursor [CurCount]int32 `c:"cursor"`
	// Keys is the key map, up to 4 combos for each action. Use the UI's bind call to change it.
	Keys     [ActCount][4]Combo `c:"keys"`
	ease     unsafe.Pointer     `c:"ease"`
	easeUser unsafe.Pointer     `c:"ease_user"`
	mask     unsafe.Pointer     `c:"mask"`
}

func (c *UIConfig) c() *C.goome_ui_config { return (*C.goome_ui_config)(unsafe.Pointer(c)) }

// Default resets every setting to its default.
func (c *UIConfig) Default() {
	C.goome_ui_config_default(c.c())
}

type uiExtra struct {
	ease cgo.Handle
	mask *C.char
}

var (
	extrasMu sync.Mutex
	extras   = map[*UIConfig]*uiExtra{}
)

func (c *UIConfig) extra() *uiExtra {
	extrasMu.Lock()
	defer extrasMu.Unlock()
	e := extras[c]
	if e == nil {
		e = &uiExtra{}
		extras[c] = e
	}
	return e
}

// SetEase sets the easing function for the custom ease. t goes from 0 to 1. Pass nil to remove it.
func (c *UIConfig) SetEase(f func(t float64) float64) {
	e := c.extra()
	if e.ease != 0 {
		e.ease.Delete()
		e.ease = 0
	}
	if f == nil {
		C.ggo_set_ease(c.c(), 0, 0)
		return
	}
	e.ease = cgo.NewHandle(f)
	C.ggo_set_ease(c.c(), 1, C.uintptr_t(e.ease))
}

// SetMask sets what a password input draws per letter. An empty string goes back to "*".
func (c *UIConfig) SetMask(mask string) {
	e := c.extra()
	old := e.mask
	e.mask = cstr(mask)
	c.mask = unsafe.Pointer(e.mask)
	if old != nil {
		C.free(unsafe.Pointer(old))
	}
}

// Mask returns what a password input draws per letter, or an empty string for the default "*".
func (c *UIConfig) Mask() string {
	if c.mask == nil {
		return ""
	}
	return C.GoString((*C.char)(c.mask))
}

func (c *UIConfig) release() {
	extrasMu.Lock()
	e := extras[c]
	delete(extras, c)
	extrasMu.Unlock()
	if e == nil {
		return
	}
	if e.ease != 0 {
		e.ease.Delete()
	}
	if e.mask != nil {
		C.free(unsafe.Pointer(e.mask))
	}
}
