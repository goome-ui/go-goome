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

type Skin struct {
	bm unsafe.Pointer `c:"bm"`
	L  float64        `c:"l"`
	T  float64        `c:"t"`
	R  float64        `c:"r"`
	B  float64        `c:"b"`
}

func NewSkin(b *Bitmap, l, t, r, bottom float64) Skin {
	if b == nil {
		return Skin{}
	}
	pin(unsafe.Pointer(b.ptr()), b)
	return Skin{bm: unsafe.Pointer(b.b), L: l, T: t, R: r, B: bottom}
}

func (s Skin) Empty() bool { return s.bm == nil }

type Theme struct {
	Bg             uint32                    `c:"bg"`
	Panel          uint32                    `c:"panel"`
	Border         uint32                    `c:"border"`
	Text           uint32                    `c:"text"`
	TextDim        uint32                    `c:"text_dim"`
	Accent         uint32                    `c:"accent"`
	AccentText     uint32                    `c:"accent_text"`
	Hover          uint32                    `c:"hover"`
	Active         uint32                    `c:"active"`
	InputBg        uint32                    `c:"input_bg"`
	Shadow         uint32                    `c:"shadow"`
	Radius         float64                   `c:"radius"`
	Pad            float64                   `c:"pad"`
	Gap            float64                   `c:"gap"`
	BorderWidth    float64                   `c:"border_width"`
	FontSize       float64                   `c:"font_size"`
	ShadowBlur     float64                   `c:"shadow_blur"`
	AnimMs         float64                   `c:"anim_ms"`
	font           unsafe.Pointer            `c:"font"`
	skinButton     Skin                      `c:"skin_button"`
	skinPanel      Skin                      `c:"skin_panel"`
	skinInput      Skin                      `c:"skin_input"`
	icons          [IconCount]unsafe.Pointer `c:"icon"`
	IconTint       uint32                    `c:"icon_tint"`
	skinCheckbox   Skin                      `c:"skin_checkbox"`
	skinRadio      Skin                      `c:"skin_radio"`
	skinToggle     Skin                      `c:"skin_toggle"`
	skinSlider     Skin                      `c:"skin_slider"`
	skinSliderFill Skin                      `c:"skin_slider_fill"`
	skinScrollbar  Skin                      `c:"skin_scrollbar"`
	skinTab        Skin                      `c:"skin_tab"`
	skinTitlebar   Skin                      `c:"skin_titlebar"`
	skinTooltip    Skin                      `c:"skin_tooltip"`
	skinMenu       Skin                      `c:"skin_menu"`
	skinPopup      Skin                      `c:"skin_popup"`
	Focus          uint32                    `c:"focus"`
	Select         uint32                    `c:"select"`
	Caret          uint32                    `c:"caret"`
	Link           uint32                    `c:"link"`
	Badge          uint32                    `c:"badge"`
	Toast          uint32                    `c:"toast"`
	Tooltip        uint32                    `c:"tooltip"`
	Danger         uint32                    `c:"danger"`
	Knob           uint32                    `c:"knob"`
	Dim            uint32                    `c:"dim"`
	CheckSize      float64                   `c:"check_size"`
	CheckRadius    float64                   `c:"check_radius"`
	ToggleW        float64                   `c:"toggle_w"`
	ToggleH        float64                   `c:"toggle_h"`
	SliderH        float64                   `c:"slider_h"`
	SliderTrack    float64                   `c:"slider_track"`
	KnobSize       float64                   `c:"knob_size"`
	ScrollW        float64                   `c:"scroll_w"`
	ScrollMin      float64                   `c:"scroll_min"`
	ItemH          float64                   `c:"item_h"`
	MenuItemH      float64                   `c:"menu_item_h"`
	TabH           float64                   `c:"tab_h"`
	WinTitleH      float64                   `c:"win_title_h"`
	AppTitleH      float64                   `c:"app_title_h"`
	MenubarH       float64                   `c:"menubar_h"`
	StatusbarH     float64                   `c:"statusbar_h"`
	FocusW         float64                   `c:"focus_w"`
	FocusGap       float64                   `c:"focus_gap"`
	SepW           float64                   `c:"sep_w"`
	Grip           float64                   `c:"grip"`
	Indent         float64                   `c:"indent"`
	PlotW          float64                   `c:"plot_w"`
	Spinner        float64                   `c:"spinner"`
	CaretW         float64                   `c:"caret_w"`
	WinMinW        float64                   `c:"win_min_w"`
	WinMinH        float64                   `c:"win_min_h"`
	PopGap         float64                   `c:"pop_gap"`
	MenuGap        float64                   `c:"menu_gap"`
	TipPadX        float64                   `c:"tip_pad_x"`
	TipPadY        float64                   `c:"tip_pad_y"`
	TipX           float64                   `c:"tip_x"`
	TipY           float64                   `c:"tip_y"`
	TipRadius      float64                   `c:"tip_radius"`
	ShadowX        float64                   `c:"shadow_x"`
	ShadowY        float64                   `c:"shadow_y"`
	Align          float64                   `c:"align"`
	HeadingSize    float64                   `c:"heading_size"`
	SmallSize      float64                   `c:"small_size"`
	Ease           float64                   `c:"ease"`
	HoverMs        float64                   `c:"hover_ms"`
	PressMs        float64                   `c:"press_ms"`
	PopupMs        float64                   `c:"popup_ms"`
	CollapseMs     float64                   `c:"collapse_ms"`
	ToastFadeMs    float64                   `c:"toast_fade_ms"`
	fontHeading    unsafe.Pointer            `c:"font_heading"`
	fontMono       unsafe.Pointer            `c:"font_mono"`
	fontSmall      unsafe.Pointer            `c:"font_small"`
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

func (t *Theme) Dark()  { C.goome_ui_theme_dark(t.c()) }
func (t *Theme) Light() { C.goome_ui_theme_light(t.c()) }

func (t *Theme) Make(accent uint32, dark bool, radius float64) {
	C.goome_ui_theme_make(t.c(), C.uint32_t(accent), b2i(dark), C.double(radius))
}

func (t *Theme) Preset(preset int) error {
	if r := C.goome_ui_theme_preset(t.c(), C.int(preset)); r < 0 {
		return newErr(int(r), "theres no preset %d", preset)
	}
	return nil
}

const (
	FontMain = iota
	FontHeading
	FontMono
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

func (t *Theme) SetFont(which int, f *Font) {
	p := unsafe.Pointer(f.raw())
	if p != nil {
		pin(p, f)
	}
	*t.fontSlot(which) = p
}

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

func (t *Theme) SetSkin(kind int, s Skin) error {
	slot := t.skinSlot(kind)
	if slot == nil {
		return newErr(CodeArg, "theres no skin %d", kind)
	}
	*slot = s
	return nil
}

func (t *Theme) Skin(kind int) Skin {
	if slot := t.skinSlot(kind); slot != nil {
		return *slot
	}
	return Skin{}
}

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

type Combo struct {
	Key  int32 `c:"key"`
	Mods int32 `c:"mods"`
}

type UIConfig struct {
	TipDelayMs    int32              `c:"tip_delay_ms"`
	DblClickMs    int32              `c:"dbl_click_ms"`
	DblClickPx    int32              `c:"dbl_click_px"`
	DragPx        int32              `c:"drag_px"`
	RepeatDelayMs int32              `c:"repeat_delay_ms"`
	RepeatMs      int32              `c:"repeat_ms"`
	BlinkMs       int32              `c:"blink_ms"`
	MenuDelayMs   int32              `c:"menu_delay_ms"`
	ToastMs       int32              `c:"toast_ms"`
	ScrollLines   float64            `c:"scroll_lines"`
	RaiseOnClick  int32              `c:"raise_on_click"`
	EscCloses     int32              `c:"esc_closes"`
	TabNav        int32              `c:"tab_nav"`
	FocusOnMouse  int32              `c:"focus_on_mouse"`
	Cursor        [CurCount]int32    `c:"cursor"`
	Keys          [ActCount][4]Combo `c:"keys"`
	ease          unsafe.Pointer     `c:"ease"`
	easeUser      unsafe.Pointer     `c:"ease_user"`
	mask          unsafe.Pointer     `c:"mask"`
}

func (c *UIConfig) c() *C.goome_ui_config { return (*C.goome_ui_config)(unsafe.Pointer(c)) }

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

func (c *UIConfig) SetMask(mask string) {
	e := c.extra()
	old := e.mask
	e.mask = cstr(mask)
	c.mask = unsafe.Pointer(e.mask)
	if old != nil {
		C.free(unsafe.Pointer(old))
	}
}

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
