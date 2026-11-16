package goome

/*
#include "native/goome.h"
*/
import "C"

const (
	BackendAuto          = C.GOOME_AUTO
	BackendWayland       = C.GOOME_WAYLAND
	BackendGnome         = C.GOOME_GNOME
	BackendPortal        = C.GOOME_PORTAL
	BackendX11           = C.GOOME_X11
	BackendWin32         = C.GOOME_WIN32
	BackendMacOS         = C.GOOME_MACOS
	BackendTest          = C.GOOME_TEST
	BackendOffscreen     = C.GOOME_OFFSCREEN
	ModeOverlay          = C.GOOME_MODE_OVERLAY
	ModeWindow           = C.GOOME_MODE_WINDOW
	ModeOffscreen        = C.GOOME_MODE_OFFSCREEN
	WindowResizable      = C.GOOME_WINDOW_RESIZABLE
	WindowBorderless     = C.GOOME_WINDOW_BORDERLESS
	WindowTransparent    = C.GOOME_WINDOW_TRANSPARENT
	WindowOnTop          = C.GOOME_WINDOW_ON_TOP
	WindowHidden         = C.GOOME_WINDOW_HIDDEN
	WindowCenter         = C.GOOME_WINDOW_CENTER
	OK                   = C.GOOME_OK
	CodeNotSupported     = C.GOOME_ENOTSUP
	CodeDenied           = C.GOOME_EDENIED
	CodeNoDisplay        = C.GOOME_ENODISPLAY
	CodeGone             = C.GOOME_EGONE
	CodeAgain            = C.GOOME_EAGAIN
	CodeNoMem            = C.GOOME_ENOMEM
	CodeArg              = C.GOOME_EARG
	CapCapture           = C.GOOME_CAP_CAPTURE
	CapGlobalPointer     = C.GOOME_CAP_GLOBAL_POINTER
	CapHotkeys           = C.GOOME_CAP_HOTKEYS
	CapHideSelf          = C.GOOME_CAP_HIDE_SELF
	CapKeyboard          = C.GOOME_CAP_KEYBOARD
	CapWindow            = C.GOOME_CAP_WINDOW
	CapClipboard         = C.GOOME_CAP_CLIPBOARD
	CapCursor            = C.GOOME_CAP_CURSOR
	CapGrab              = C.GOOME_CAP_GRAB
	WantMouse            = C.GOOME_WANT_MOUSE
	WantKeyboard         = C.GOOME_WANT_KEYBOARD
	EventNone            = C.GOOME_NONE
	EventDown            = C.GOOME_DOWN
	EventUp              = C.GOOME_UP
	EventMotion          = C.GOOME_MOTION
	EventScroll          = C.GOOME_SCROLL
	EventKey             = C.GOOME_KEY
	EventHotkey          = C.GOOME_HOTKEY
	EventControl         = C.GOOME_CONTROL
	EventMonitors        = C.GOOME_MONITORS
	EventQuit            = C.GOOME_QUIT
	EventClose           = C.GOOME_CLOSE
	EventFocus           = C.GOOME_FOCUS
	EventDrop            = C.GOOME_DROP
	EventText            = C.GOOME_TEXT
	EventPreedit         = C.GOOME_PREEDIT
	CursorArrow          = C.GOOME_CURSOR_ARROW
	CursorIBeam          = C.GOOME_CURSOR_IBEAM
	CursorHand           = C.GOOME_CURSOR_HAND
	CursorMove           = C.GOOME_CURSOR_MOVE
	CursorResizeEW       = C.GOOME_CURSOR_RESIZE_EW
	CursorResizeNS       = C.GOOME_CURSOR_RESIZE_NS
	CursorResizeNWSE     = C.GOOME_CURSOR_RESIZE_NWSE
	CursorResizeNESW     = C.GOOME_CURSOR_RESIZE_NESW
	CursorWait           = C.GOOME_CURSOR_WAIT
	CursorNo             = C.GOOME_CURSOR_NO
	CursorHidden         = C.GOOME_CURSOR_HIDDEN
	EdgeLeft             = C.GOOME_EDGE_LEFT
	EdgeRight            = C.GOOME_EDGE_RIGHT
	EdgeTop              = C.GOOME_EDGE_TOP
	EdgeBottom           = C.GOOME_EDGE_BOTTOM
	StateNormal          = C.GOOME_STATE_NORMAL
	StateMinimized       = C.GOOME_STATE_MINIMIZED
	StateMaximized       = C.GOOME_STATE_MAXIMIZED
	StateFullscreen      = C.GOOME_STATE_FULLSCREEN
	KeyBackspace         = C.GOOME_K_BACKSPACE
	KeyTab               = C.GOOME_K_TAB
	KeyEnter             = C.GOOME_K_ENTER
	KeyEsc               = C.GOOME_K_ESC
	KeySpace             = C.GOOME_K_SPACE
	KeyDelete            = C.GOOME_K_DELETE
	KeyF1                = C.GOOME_K_F1
	KeyF12               = C.GOOME_K_F12
	KeyLeft              = C.GOOME_K_LEFT
	KeyRight             = C.GOOME_K_RIGHT
	KeyUp                = C.GOOME_K_UP
	KeyDown              = C.GOOME_K_DOWN
	KeyHome              = C.GOOME_K_HOME
	KeyEnd               = C.GOOME_K_END
	KeyPgUp              = C.GOOME_K_PGUP
	KeyPgDn              = C.GOOME_K_PGDN
	KeyInsert            = C.GOOME_K_INSERT
	ModShift             = C.GOOME_SHIFT
	ModCtrl              = C.GOOME_CTRL
	ModAlt               = C.GOOME_ALT
	ModSuper             = C.GOOME_SUPER
	FmtBGRAPremul        = C.GOOME_FMT_BGRA_PREMUL
	FmtRGBAPremul        = C.GOOME_FMT_RGBA_PREMUL
	FmtBGRA              = C.GOOME_FMT_BGRA
	FmtRGBA              = C.GOOME_FMT_RGBA
	DirAuto              = C.GOOME_DIR_AUTO
	DirLTR               = C.GOOME_DIR_LTR
	DirRTL               = C.GOOME_DIR_RTL
	IconArrowDown        = C.GOOME_ICON_ARROW_DOWN
	IconArrowRight       = C.GOOME_ICON_ARROW_RIGHT
	IconArrowUp          = C.GOOME_ICON_ARROW_UP
	IconArrowLeft        = C.GOOME_ICON_ARROW_LEFT
	IconCheck            = C.GOOME_ICON_CHECK
	IconRadioDot         = C.GOOME_ICON_RADIO_DOT
	IconClose            = C.GOOME_ICON_CLOSE
	IconMinimize         = C.GOOME_ICON_MINIMIZE
	IconMaximize         = C.GOOME_ICON_MAXIMIZE
	IconRestore          = C.GOOME_ICON_RESTORE
	IconResize           = C.GOOME_ICON_RESIZE
	IconSortUp           = C.GOOME_ICON_SORT_UP
	IconSortDown         = C.GOOME_ICON_SORT_DOWN
	IconPlus             = C.GOOME_ICON_PLUS
	IconMinus            = C.GOOME_ICON_MINUS
	IconSliderKnob       = C.GOOME_ICON_SLIDER_KNOB
	IconToggleKnob       = C.GOOME_ICON_TOGGLE_KNOB
	IconCollapse         = C.GOOME_ICON_COLLAPSE
	IconExpand           = C.GOOME_ICON_EXPAND
	IconDragHandle       = C.GOOME_ICON_DRAG_HANDLE
	IconToast            = C.GOOME_ICON_TOAST
	IconPicker           = C.GOOME_ICON_PICKER
	IconCount            = C.GOOME_ICON_COUNT
	WinMove              = C.GOOME_WIN_MOVE
	WinResize            = C.GOOME_WIN_RESIZE
	WinCollapse          = C.GOOME_WIN_COLLAPSE
	WinNoTitle           = C.GOOME_WIN_NO_TITLE
	WinClose             = C.GOOME_WIN_CLOSE
	InputPassword        = C.GOOME_INPUT_PASSWORD
	InputReadonly        = C.GOOME_INPUT_READONLY
	InputNumber          = C.GOOME_INPUT_NUMBER
	InputNoSelectOnFocus = C.GOOME_INPUT_NO_SELECT_ON_FOCUS
	InputEnterReturns    = C.GOOME_INPUT_ENTER_RETURNS
	AppTitlebar          = C.GOOME_APP_TITLEBAR
	AppMenubar           = C.GOOME_APP_MENUBAR
	AppStatusbar         = C.GOOME_APP_STATUSBAR
	AppNoPad             = C.GOOME_APP_NO_PAD
	EditReadonly         = C.GOOME_EDIT_READONLY
	EditLineNumbers      = C.GOOME_EDIT_LINE_NUMBERS
	EditWrap             = C.GOOME_EDIT_WRAP
	EditMono             = C.GOOME_EDIT_MONO
	TableSortable        = C.GOOME_TABLE_SORTABLE
	TableResizable       = C.GOOME_TABLE_RESIZABLE
	TableStripes         = C.GOOME_TABLE_STRIPES
	TableBorders         = C.GOOME_TABLE_BORDERS
	StyleBg              = C.GOOME_S_BG
	StylePanel           = C.GOOME_S_PANEL
	StyleBorder          = C.GOOME_S_BORDER
	StyleText            = C.GOOME_S_TEXT
	StyleTextDim         = C.GOOME_S_TEXT_DIM
	StyleAccent          = C.GOOME_S_ACCENT
	StyleAccentText      = C.GOOME_S_ACCENT_TEXT
	StyleHover           = C.GOOME_S_HOVER
	StyleActive          = C.GOOME_S_ACTIVE
	StyleInputBg         = C.GOOME_S_INPUT_BG
	StyleShadow          = C.GOOME_S_SHADOW
	StyleRadius          = C.GOOME_S_RADIUS
	StylePad             = C.GOOME_S_PAD
	StyleGap             = C.GOOME_S_GAP
	StyleBorderWidth     = C.GOOME_S_BORDER_WIDTH
	StyleFontSize        = C.GOOME_S_FONT_SIZE
	StyleShadowBlur      = C.GOOME_S_SHADOW_BLUR
	StyleAnimMs          = C.GOOME_S_ANIM_MS
	StyleFocus           = C.GOOME_S_FOCUS
	StyleSelect          = C.GOOME_S_SELECT
	StyleCaret           = C.GOOME_S_CARET
	StyleLink            = C.GOOME_S_LINK
	StyleBadge           = C.GOOME_S_BADGE
	StyleToast           = C.GOOME_S_TOAST
	StyleTooltip         = C.GOOME_S_TOOLTIP
	StyleDanger          = C.GOOME_S_DANGER
	StyleKnob            = C.GOOME_S_KNOB
	StyleDim             = C.GOOME_S_DIM
	StyleCheckSize       = C.GOOME_S_CHECK_SIZE
	StyleCheckRadius     = C.GOOME_S_CHECK_RADIUS
	StyleToggleW         = C.GOOME_S_TOGGLE_W
	StyleToggleH         = C.GOOME_S_TOGGLE_H
	StyleSliderH         = C.GOOME_S_SLIDER_H
	StyleSliderTrack     = C.GOOME_S_SLIDER_TRACK
	StyleKnobSize        = C.GOOME_S_KNOB_SIZE
	StyleScrollW         = C.GOOME_S_SCROLL_W
	StyleScrollMin       = C.GOOME_S_SCROLL_MIN
	StyleItemH           = C.GOOME_S_ITEM_H
	StyleMenuItemH       = C.GOOME_S_MENU_ITEM_H
	StyleTabH            = C.GOOME_S_TAB_H
	StyleWinTitleH       = C.GOOME_S_WIN_TITLE_H
	StyleAppTitleH       = C.GOOME_S_APP_TITLE_H
	StyleMenubarH        = C.GOOME_S_MENUBAR_H
	StyleStatusbarH      = C.GOOME_S_STATUSBAR_H
	StyleFocusW          = C.GOOME_S_FOCUS_W
	StyleFocusGap        = C.GOOME_S_FOCUS_GAP
	StyleSepW            = C.GOOME_S_SEP_W
	StyleGrip            = C.GOOME_S_GRIP
	StyleIndent          = C.GOOME_S_INDENT
	StylePlotW           = C.GOOME_S_PLOT_W
	StyleSpinner         = C.GOOME_S_SPINNER
	StyleCaretW          = C.GOOME_S_CARET_W
	StyleWinMinW         = C.GOOME_S_WIN_MIN_W
	StyleWinMinH         = C.GOOME_S_WIN_MIN_H
	StylePopGap          = C.GOOME_S_POP_GAP
	StyleMenuGap         = C.GOOME_S_MENU_GAP
	StyleTipPadX         = C.GOOME_S_TIP_PAD_X
	StyleTipPadY         = C.GOOME_S_TIP_PAD_Y
	StyleTipX            = C.GOOME_S_TIP_X
	StyleTipY            = C.GOOME_S_TIP_Y
	StyleTipRadius       = C.GOOME_S_TIP_RADIUS
	StyleShadowX         = C.GOOME_S_SHADOW_X
	StyleShadowY         = C.GOOME_S_SHADOW_Y
	StyleAlign           = C.GOOME_S_ALIGN
	StyleHeadingSize     = C.GOOME_S_HEADING_SIZE
	StyleSmallSize       = C.GOOME_S_SMALL_SIZE
	StyleEase            = C.GOOME_S_EASE
	StyleHoverMs         = C.GOOME_S_HOVER_MS
	StylePressMs         = C.GOOME_S_PRESS_MS
	StylePopupMs         = C.GOOME_S_POPUP_MS
	StyleCollapseMs      = C.GOOME_S_COLLAPSE_MS
	StyleToastFadeMs     = C.GOOME_S_TOAST_FADE_MS
	StyleCount           = C.GOOME_S_COUNT
	AlignAuto            = C.GOOME_ALIGN_AUTO
	AlignLeft            = C.GOOME_ALIGN_LEFT
	AlignCenter          = C.GOOME_ALIGN_CENTER
	AlignRight           = C.GOOME_ALIGN_RIGHT
	EaseSmooth           = C.GOOME_EASE_SMOOTH
	EaseLinear           = C.GOOME_EASE_LINEAR
	EaseOut              = C.GOOME_EASE_OUT
	EaseInOut            = C.GOOME_EASE_IN_OUT
	EaseBack             = C.GOOME_EASE_BACK
	EaseCustom           = C.GOOME_EASE_CUSTOM
	EaseIn               = C.GOOME_EASE_IN
	EaseSpring           = C.GOOME_EASE_SPRING
	EaseBounce           = C.GOOME_EASE_BOUNCE
	SkinButton           = C.GOOME_SKIN_BUTTON
	SkinPanel            = C.GOOME_SKIN_PANEL
	SkinInput            = C.GOOME_SKIN_INPUT
	SkinCheckbox         = C.GOOME_SKIN_CHECKBOX
	SkinRadio            = C.GOOME_SKIN_RADIO
	SkinToggle           = C.GOOME_SKIN_TOGGLE
	SkinSliderTrack      = C.GOOME_SKIN_SLIDER_TRACK
	SkinSliderFill       = C.GOOME_SKIN_SLIDER_FILL
	SkinScrollbar        = C.GOOME_SKIN_SCROLLBAR
	SkinTab              = C.GOOME_SKIN_TAB
	SkinTitlebar         = C.GOOME_SKIN_TITLEBAR
	SkinTooltip          = C.GOOME_SKIN_TOOLTIP
	SkinMenu             = C.GOOME_SKIN_MENU
	SkinPopup            = C.GOOME_SKIN_POPUP
	SkinCount            = C.GOOME_SKIN_COUNT
	ItemDragX            = C.GOOME_ITEM_DRAG_X
	ItemDragY            = C.GOOME_ITEM_DRAG_Y
	ItemFocusable        = C.GOOME_ITEM_FOCUSABLE
	ItemNoHit            = C.GOOME_ITEM_NO_HIT
	ItemRepeat           = C.GOOME_ITEM_REPEAT
	PresetDark           = C.GOOME_PRESET_DARK
	PresetLight          = C.GOOME_PRESET_LIGHT
	PresetMidnight       = C.GOOME_PRESET_MIDNIGHT
	PresetPaper          = C.GOOME_PRESET_PAPER
	PresetMint           = C.GOOME_PRESET_MINT
	PresetHighContrast   = C.GOOME_PRESET_HIGH_CONTRAST
	PresetRose           = C.GOOME_PRESET_ROSE
	PresetCount          = C.GOOME_PRESET_COUNT
	ActCopy              = C.GOOME_ACT_COPY
	ActCut               = C.GOOME_ACT_CUT
	ActPaste             = C.GOOME_ACT_PASTE
	ActUndo              = C.GOOME_ACT_UNDO
	ActRedo              = C.GOOME_ACT_REDO
	ActSelectAll         = C.GOOME_ACT_SELECT_ALL
	ActWordLeft          = C.GOOME_ACT_WORD_LEFT
	ActWordRight         = C.GOOME_ACT_WORD_RIGHT
	ActDelWordLeft       = C.GOOME_ACT_DEL_WORD_LEFT
	ActDelWordRight      = C.GOOME_ACT_DEL_WORD_RIGHT
	ActDocStart          = C.GOOME_ACT_DOC_START
	ActDocEnd            = C.GOOME_ACT_DOC_END
	ActCount             = C.GOOME_ACT_COUNT
	CurButton            = C.GOOME_CUR_BUTTON
	CurInput             = C.GOOME_CUR_INPUT
	CurLink              = C.GOOME_CUR_LINK
	CurDragValue         = C.GOOME_CUR_DRAG_VALUE
	CurSplitX            = C.GOOME_CUR_SPLIT_X
	CurSplitY            = C.GOOME_CUR_SPLIT_Y
	CurGrip              = C.GOOME_CUR_GRIP
	CurColumn            = C.GOOME_CUR_COLUMN
	CurDragDrop          = C.GOOME_CUR_DRAG_DROP
	CurCustomX           = C.GOOME_CUR_CUSTOM_X
	CurCustomY           = C.GOOME_CUR_CUSTOM_Y
	CurCustomXy          = C.GOOME_CUR_CUSTOM_XY
	CurCount             = C.GOOME_CUR_COUNT
	RoleNone             = C.GOOME_ROLE_NONE
	RoleWindow           = C.GOOME_ROLE_WINDOW
	RoleButton           = C.GOOME_ROLE_BUTTON
	RoleCheckbox         = C.GOOME_ROLE_CHECKBOX
	RoleRadio            = C.GOOME_ROLE_RADIO
	RoleSlider           = C.GOOME_ROLE_SLIDER
	RoleText             = C.GOOME_ROLE_TEXT
	RoleInput            = C.GOOME_ROLE_INPUT
	RoleCombo            = C.GOOME_ROLE_COMBO
	RoleList             = C.GOOME_ROLE_LIST
	RoleItem             = C.GOOME_ROLE_ITEM
	RoleTab              = C.GOOME_ROLE_TAB
	RoleMenu             = C.GOOME_ROLE_MENU
	RoleMenuItem         = C.GOOME_ROLE_MENUITEM
	RoleTable            = C.GOOME_ROLE_TABLE
	RoleTree             = C.GOOME_ROLE_TREE
	RoleProgress         = C.GOOME_ROLE_PROGRESS
	RoleImage            = C.GOOME_ROLE_IMAGE
	RoleGroup            = C.GOOME_ROLE_GROUP
	RoleDialog           = C.GOOME_ROLE_DIALOG
	A11yFocused          = C.GOOME_A11Y_FOCUSED
	A11yChecked          = C.GOOME_A11Y_CHECKED
	A11yDisabled         = C.GOOME_A11Y_DISABLED
	A11yExpanded         = C.GOOME_A11Y_EXPANDED
	A11ySelected         = C.GOOME_A11Y_SELECTED
)

var cConsts = map[string]int{
	"GOOME_AUTO":                     BackendAuto,
	"GOOME_WAYLAND":                  BackendWayland,
	"GOOME_GNOME":                    BackendGnome,
	"GOOME_PORTAL":                   BackendPortal,
	"GOOME_X11":                      BackendX11,
	"GOOME_WIN32":                    BackendWin32,
	"GOOME_MACOS":                    BackendMacOS,
	"GOOME_TEST":                     BackendTest,
	"GOOME_OFFSCREEN":                BackendOffscreen,
	"GOOME_MODE_OVERLAY":             ModeOverlay,
	"GOOME_MODE_WINDOW":              ModeWindow,
	"GOOME_MODE_OFFSCREEN":           ModeOffscreen,
	"GOOME_WINDOW_RESIZABLE":         WindowResizable,
	"GOOME_WINDOW_BORDERLESS":        WindowBorderless,
	"GOOME_WINDOW_TRANSPARENT":       WindowTransparent,
	"GOOME_WINDOW_ON_TOP":            WindowOnTop,
	"GOOME_WINDOW_HIDDEN":            WindowHidden,
	"GOOME_WINDOW_CENTER":            WindowCenter,
	"GOOME_OK":                       OK,
	"GOOME_ENOTSUP":                  CodeNotSupported,
	"GOOME_EDENIED":                  CodeDenied,
	"GOOME_ENODISPLAY":               CodeNoDisplay,
	"GOOME_EGONE":                    CodeGone,
	"GOOME_EAGAIN":                   CodeAgain,
	"GOOME_ENOMEM":                   CodeNoMem,
	"GOOME_EARG":                     CodeArg,
	"GOOME_CAP_CAPTURE":              CapCapture,
	"GOOME_CAP_GLOBAL_POINTER":       CapGlobalPointer,
	"GOOME_CAP_HOTKEYS":              CapHotkeys,
	"GOOME_CAP_HIDE_SELF":            CapHideSelf,
	"GOOME_CAP_KEYBOARD":             CapKeyboard,
	"GOOME_CAP_WINDOW":               CapWindow,
	"GOOME_CAP_CLIPBOARD":            CapClipboard,
	"GOOME_CAP_CURSOR":               CapCursor,
	"GOOME_CAP_GRAB":                 CapGrab,
	"GOOME_WANT_MOUSE":               WantMouse,
	"GOOME_WANT_KEYBOARD":            WantKeyboard,
	"GOOME_NONE":                     EventNone,
	"GOOME_DOWN":                     EventDown,
	"GOOME_UP":                       EventUp,
	"GOOME_MOTION":                   EventMotion,
	"GOOME_SCROLL":                   EventScroll,
	"GOOME_KEY":                      EventKey,
	"GOOME_HOTKEY":                   EventHotkey,
	"GOOME_CONTROL":                  EventControl,
	"GOOME_MONITORS":                 EventMonitors,
	"GOOME_QUIT":                     EventQuit,
	"GOOME_CLOSE":                    EventClose,
	"GOOME_FOCUS":                    EventFocus,
	"GOOME_DROP":                     EventDrop,
	"GOOME_TEXT":                     EventText,
	"GOOME_PREEDIT":                  EventPreedit,
	"GOOME_CURSOR_ARROW":             CursorArrow,
	"GOOME_CURSOR_IBEAM":             CursorIBeam,
	"GOOME_CURSOR_HAND":              CursorHand,
	"GOOME_CURSOR_MOVE":              CursorMove,
	"GOOME_CURSOR_RESIZE_EW":         CursorResizeEW,
	"GOOME_CURSOR_RESIZE_NS":         CursorResizeNS,
	"GOOME_CURSOR_RESIZE_NWSE":       CursorResizeNWSE,
	"GOOME_CURSOR_RESIZE_NESW":       CursorResizeNESW,
	"GOOME_CURSOR_WAIT":              CursorWait,
	"GOOME_CURSOR_NO":                CursorNo,
	"GOOME_CURSOR_HIDDEN":            CursorHidden,
	"GOOME_EDGE_LEFT":                EdgeLeft,
	"GOOME_EDGE_RIGHT":               EdgeRight,
	"GOOME_EDGE_TOP":                 EdgeTop,
	"GOOME_EDGE_BOTTOM":              EdgeBottom,
	"GOOME_STATE_NORMAL":             StateNormal,
	"GOOME_STATE_MINIMIZED":          StateMinimized,
	"GOOME_STATE_MAXIMIZED":          StateMaximized,
	"GOOME_STATE_FULLSCREEN":         StateFullscreen,
	"GOOME_K_BACKSPACE":              KeyBackspace,
	"GOOME_K_TAB":                    KeyTab,
	"GOOME_K_ENTER":                  KeyEnter,
	"GOOME_K_ESC":                    KeyEsc,
	"GOOME_K_SPACE":                  KeySpace,
	"GOOME_K_DELETE":                 KeyDelete,
	"GOOME_K_F1":                     KeyF1,
	"GOOME_K_F12":                    KeyF12,
	"GOOME_K_LEFT":                   KeyLeft,
	"GOOME_K_RIGHT":                  KeyRight,
	"GOOME_K_UP":                     KeyUp,
	"GOOME_K_DOWN":                   KeyDown,
	"GOOME_K_HOME":                   KeyHome,
	"GOOME_K_END":                    KeyEnd,
	"GOOME_K_PGUP":                   KeyPgUp,
	"GOOME_K_PGDN":                   KeyPgDn,
	"GOOME_K_INSERT":                 KeyInsert,
	"GOOME_SHIFT":                    ModShift,
	"GOOME_CTRL":                     ModCtrl,
	"GOOME_ALT":                      ModAlt,
	"GOOME_SUPER":                    ModSuper,
	"GOOME_FMT_BGRA_PREMUL":          FmtBGRAPremul,
	"GOOME_FMT_RGBA_PREMUL":          FmtRGBAPremul,
	"GOOME_FMT_BGRA":                 FmtBGRA,
	"GOOME_FMT_RGBA":                 FmtRGBA,
	"GOOME_DIR_AUTO":                 DirAuto,
	"GOOME_DIR_LTR":                  DirLTR,
	"GOOME_DIR_RTL":                  DirRTL,
	"GOOME_ICON_ARROW_DOWN":          IconArrowDown,
	"GOOME_ICON_ARROW_RIGHT":         IconArrowRight,
	"GOOME_ICON_ARROW_UP":            IconArrowUp,
	"GOOME_ICON_ARROW_LEFT":          IconArrowLeft,
	"GOOME_ICON_CHECK":               IconCheck,
	"GOOME_ICON_RADIO_DOT":           IconRadioDot,
	"GOOME_ICON_CLOSE":               IconClose,
	"GOOME_ICON_MINIMIZE":            IconMinimize,
	"GOOME_ICON_MAXIMIZE":            IconMaximize,
	"GOOME_ICON_RESTORE":             IconRestore,
	"GOOME_ICON_RESIZE":              IconResize,
	"GOOME_ICON_SORT_UP":             IconSortUp,
	"GOOME_ICON_SORT_DOWN":           IconSortDown,
	"GOOME_ICON_PLUS":                IconPlus,
	"GOOME_ICON_MINUS":               IconMinus,
	"GOOME_ICON_SLIDER_KNOB":         IconSliderKnob,
	"GOOME_ICON_TOGGLE_KNOB":         IconToggleKnob,
	"GOOME_ICON_COLLAPSE":            IconCollapse,
	"GOOME_ICON_EXPAND":              IconExpand,
	"GOOME_ICON_DRAG_HANDLE":         IconDragHandle,
	"GOOME_ICON_TOAST":               IconToast,
	"GOOME_ICON_PICKER":              IconPicker,
	"GOOME_ICON_COUNT":               IconCount,
	"GOOME_WIN_MOVE":                 WinMove,
	"GOOME_WIN_RESIZE":               WinResize,
	"GOOME_WIN_COLLAPSE":             WinCollapse,
	"GOOME_WIN_NO_TITLE":             WinNoTitle,
	"GOOME_WIN_CLOSE":                WinClose,
	"GOOME_INPUT_PASSWORD":           InputPassword,
	"GOOME_INPUT_READONLY":           InputReadonly,
	"GOOME_INPUT_NUMBER":             InputNumber,
	"GOOME_INPUT_NO_SELECT_ON_FOCUS": InputNoSelectOnFocus,
	"GOOME_INPUT_ENTER_RETURNS":      InputEnterReturns,
	"GOOME_APP_TITLEBAR":             AppTitlebar,
	"GOOME_APP_MENUBAR":              AppMenubar,
	"GOOME_APP_STATUSBAR":            AppStatusbar,
	"GOOME_APP_NO_PAD":               AppNoPad,
	"GOOME_EDIT_READONLY":            EditReadonly,
	"GOOME_EDIT_LINE_NUMBERS":        EditLineNumbers,
	"GOOME_EDIT_WRAP":                EditWrap,
	"GOOME_EDIT_MONO":                EditMono,
	"GOOME_TABLE_SORTABLE":           TableSortable,
	"GOOME_TABLE_RESIZABLE":          TableResizable,
	"GOOME_TABLE_STRIPES":            TableStripes,
	"GOOME_TABLE_BORDERS":            TableBorders,
	"GOOME_S_BG":                     StyleBg,
	"GOOME_S_PANEL":                  StylePanel,
	"GOOME_S_BORDER":                 StyleBorder,
	"GOOME_S_TEXT":                   StyleText,
	"GOOME_S_TEXT_DIM":               StyleTextDim,
	"GOOME_S_ACCENT":                 StyleAccent,
	"GOOME_S_ACCENT_TEXT":            StyleAccentText,
	"GOOME_S_HOVER":                  StyleHover,
	"GOOME_S_ACTIVE":                 StyleActive,
	"GOOME_S_INPUT_BG":               StyleInputBg,
	"GOOME_S_SHADOW":                 StyleShadow,
	"GOOME_S_RADIUS":                 StyleRadius,
	"GOOME_S_PAD":                    StylePad,
	"GOOME_S_GAP":                    StyleGap,
	"GOOME_S_BORDER_WIDTH":           StyleBorderWidth,
	"GOOME_S_FONT_SIZE":              StyleFontSize,
	"GOOME_S_SHADOW_BLUR":            StyleShadowBlur,
	"GOOME_S_ANIM_MS":                StyleAnimMs,
	"GOOME_S_FOCUS":                  StyleFocus,
	"GOOME_S_SELECT":                 StyleSelect,
	"GOOME_S_CARET":                  StyleCaret,
	"GOOME_S_LINK":                   StyleLink,
	"GOOME_S_BADGE":                  StyleBadge,
	"GOOME_S_TOAST":                  StyleToast,
	"GOOME_S_TOOLTIP":                StyleTooltip,
	"GOOME_S_DANGER":                 StyleDanger,
	"GOOME_S_KNOB":                   StyleKnob,
	"GOOME_S_DIM":                    StyleDim,
	"GOOME_S_CHECK_SIZE":             StyleCheckSize,
	"GOOME_S_CHECK_RADIUS":           StyleCheckRadius,
	"GOOME_S_TOGGLE_W":               StyleToggleW,
	"GOOME_S_TOGGLE_H":               StyleToggleH,
	"GOOME_S_SLIDER_H":               StyleSliderH,
	"GOOME_S_SLIDER_TRACK":           StyleSliderTrack,
	"GOOME_S_KNOB_SIZE":              StyleKnobSize,
	"GOOME_S_SCROLL_W":               StyleScrollW,
	"GOOME_S_SCROLL_MIN":             StyleScrollMin,
	"GOOME_S_ITEM_H":                 StyleItemH,
	"GOOME_S_MENU_ITEM_H":            StyleMenuItemH,
	"GOOME_S_TAB_H":                  StyleTabH,
	"GOOME_S_WIN_TITLE_H":            StyleWinTitleH,
	"GOOME_S_APP_TITLE_H":            StyleAppTitleH,
	"GOOME_S_MENUBAR_H":              StyleMenubarH,
	"GOOME_S_STATUSBAR_H":            StyleStatusbarH,
	"GOOME_S_FOCUS_W":                StyleFocusW,
	"GOOME_S_FOCUS_GAP":              StyleFocusGap,
	"GOOME_S_SEP_W":                  StyleSepW,
	"GOOME_S_GRIP":                   StyleGrip,
	"GOOME_S_INDENT":                 StyleIndent,
	"GOOME_S_PLOT_W":                 StylePlotW,
	"GOOME_S_SPINNER":                StyleSpinner,
	"GOOME_S_CARET_W":                StyleCaretW,
	"GOOME_S_WIN_MIN_W":              StyleWinMinW,
	"GOOME_S_WIN_MIN_H":              StyleWinMinH,
	"GOOME_S_POP_GAP":                StylePopGap,
	"GOOME_S_MENU_GAP":               StyleMenuGap,
	"GOOME_S_TIP_PAD_X":              StyleTipPadX,
	"GOOME_S_TIP_PAD_Y":              StyleTipPadY,
	"GOOME_S_TIP_X":                  StyleTipX,
	"GOOME_S_TIP_Y":                  StyleTipY,
	"GOOME_S_TIP_RADIUS":             StyleTipRadius,
	"GOOME_S_SHADOW_X":               StyleShadowX,
	"GOOME_S_SHADOW_Y":               StyleShadowY,
	"GOOME_S_ALIGN":                  StyleAlign,
	"GOOME_S_HEADING_SIZE":           StyleHeadingSize,
	"GOOME_S_SMALL_SIZE":             StyleSmallSize,
	"GOOME_S_EASE":                   StyleEase,
	"GOOME_S_HOVER_MS":               StyleHoverMs,
	"GOOME_S_PRESS_MS":               StylePressMs,
	"GOOME_S_POPUP_MS":               StylePopupMs,
	"GOOME_S_COLLAPSE_MS":            StyleCollapseMs,
	"GOOME_S_TOAST_FADE_MS":          StyleToastFadeMs,
	"GOOME_S_COUNT":                  StyleCount,
	"GOOME_ALIGN_AUTO":               AlignAuto,
	"GOOME_ALIGN_LEFT":               AlignLeft,
	"GOOME_ALIGN_CENTER":             AlignCenter,
	"GOOME_ALIGN_RIGHT":              AlignRight,
	"GOOME_EASE_SMOOTH":              EaseSmooth,
	"GOOME_EASE_LINEAR":              EaseLinear,
	"GOOME_EASE_OUT":                 EaseOut,
	"GOOME_EASE_IN_OUT":              EaseInOut,
	"GOOME_EASE_BACK":                EaseBack,
	"GOOME_EASE_CUSTOM":              EaseCustom,
	"GOOME_EASE_IN":                  EaseIn,
	"GOOME_EASE_SPRING":              EaseSpring,
	"GOOME_EASE_BOUNCE":              EaseBounce,
	"GOOME_SKIN_BUTTON":              SkinButton,
	"GOOME_SKIN_PANEL":               SkinPanel,
	"GOOME_SKIN_INPUT":               SkinInput,
	"GOOME_SKIN_CHECKBOX":            SkinCheckbox,
	"GOOME_SKIN_RADIO":               SkinRadio,
	"GOOME_SKIN_TOGGLE":              SkinToggle,
	"GOOME_SKIN_SLIDER_TRACK":        SkinSliderTrack,
	"GOOME_SKIN_SLIDER_FILL":         SkinSliderFill,
	"GOOME_SKIN_SCROLLBAR":           SkinScrollbar,
	"GOOME_SKIN_TAB":                 SkinTab,
	"GOOME_SKIN_TITLEBAR":            SkinTitlebar,
	"GOOME_SKIN_TOOLTIP":             SkinTooltip,
	"GOOME_SKIN_MENU":                SkinMenu,
	"GOOME_SKIN_POPUP":               SkinPopup,
	"GOOME_SKIN_COUNT":               SkinCount,
	"GOOME_ITEM_DRAG_X":              ItemDragX,
	"GOOME_ITEM_DRAG_Y":              ItemDragY,
	"GOOME_ITEM_FOCUSABLE":           ItemFocusable,
	"GOOME_ITEM_NO_HIT":              ItemNoHit,
	"GOOME_ITEM_REPEAT":              ItemRepeat,
	"GOOME_PRESET_DARK":              PresetDark,
	"GOOME_PRESET_LIGHT":             PresetLight,
	"GOOME_PRESET_MIDNIGHT":          PresetMidnight,
	"GOOME_PRESET_PAPER":             PresetPaper,
	"GOOME_PRESET_MINT":              PresetMint,
	"GOOME_PRESET_HIGH_CONTRAST":     PresetHighContrast,
	"GOOME_PRESET_ROSE":              PresetRose,
	"GOOME_PRESET_COUNT":             PresetCount,
	"GOOME_ACT_COPY":                 ActCopy,
	"GOOME_ACT_CUT":                  ActCut,
	"GOOME_ACT_PASTE":                ActPaste,
	"GOOME_ACT_UNDO":                 ActUndo,
	"GOOME_ACT_REDO":                 ActRedo,
	"GOOME_ACT_SELECT_ALL":           ActSelectAll,
	"GOOME_ACT_WORD_LEFT":            ActWordLeft,
	"GOOME_ACT_WORD_RIGHT":           ActWordRight,
	"GOOME_ACT_DEL_WORD_LEFT":        ActDelWordLeft,
	"GOOME_ACT_DEL_WORD_RIGHT":       ActDelWordRight,
	"GOOME_ACT_DOC_START":            ActDocStart,
	"GOOME_ACT_DOC_END":              ActDocEnd,
	"GOOME_ACT_COUNT":                ActCount,
	"GOOME_CUR_BUTTON":               CurButton,
	"GOOME_CUR_INPUT":                CurInput,
	"GOOME_CUR_LINK":                 CurLink,
	"GOOME_CUR_DRAG_VALUE":           CurDragValue,
	"GOOME_CUR_SPLIT_X":              CurSplitX,
	"GOOME_CUR_SPLIT_Y":              CurSplitY,
	"GOOME_CUR_GRIP":                 CurGrip,
	"GOOME_CUR_COLUMN":               CurColumn,
	"GOOME_CUR_DRAG_DROP":            CurDragDrop,
	"GOOME_CUR_CUSTOM_X":             CurCustomX,
	"GOOME_CUR_CUSTOM_Y":             CurCustomY,
	"GOOME_CUR_CUSTOM_XY":            CurCustomXy,
	"GOOME_CUR_COUNT":                CurCount,
	"GOOME_ROLE_NONE":                RoleNone,
	"GOOME_ROLE_WINDOW":              RoleWindow,
	"GOOME_ROLE_BUTTON":              RoleButton,
	"GOOME_ROLE_CHECKBOX":            RoleCheckbox,
	"GOOME_ROLE_RADIO":               RoleRadio,
	"GOOME_ROLE_SLIDER":              RoleSlider,
	"GOOME_ROLE_TEXT":                RoleText,
	"GOOME_ROLE_INPUT":               RoleInput,
	"GOOME_ROLE_COMBO":               RoleCombo,
	"GOOME_ROLE_LIST":                RoleList,
	"GOOME_ROLE_ITEM":                RoleItem,
	"GOOME_ROLE_TAB":                 RoleTab,
	"GOOME_ROLE_MENU":                RoleMenu,
	"GOOME_ROLE_MENUITEM":            RoleMenuItem,
	"GOOME_ROLE_TABLE":               RoleTable,
	"GOOME_ROLE_TREE":                RoleTree,
	"GOOME_ROLE_PROGRESS":            RoleProgress,
	"GOOME_ROLE_IMAGE":               RoleImage,
	"GOOME_ROLE_GROUP":               RoleGroup,
	"GOOME_ROLE_DIALOG":              RoleDialog,
	"GOOME_A11Y_FOCUSED":             A11yFocused,
	"GOOME_A11Y_CHECKED":             A11yChecked,
	"GOOME_A11Y_DISABLED":            A11yDisabled,
	"GOOME_A11Y_EXPANDED":            A11yExpanded,
	"GOOME_A11Y_SELECTED":            A11ySelected,
}
