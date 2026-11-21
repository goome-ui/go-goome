package goome

/*
#include "native/goome.h"
*/
import "C"

const (
	// BackendAuto lets goome pick the backend itself.
	BackendAuto = C.GOOME_AUTO
	// BackendWayland is the Wayland backend.
	BackendWayland = C.GOOME_WAYLAND
	// BackendGnome is the GNOME Wayland backend.
	BackendGnome = C.GOOME_GNOME
	// BackendPortal is the backend that grabs the screen through the desktop portal.
	BackendPortal = C.GOOME_PORTAL
	// BackendX11 is the X11 backend.
	BackendX11 = C.GOOME_X11
	// BackendWin32 is the Windows backend.
	BackendWin32 = C.GOOME_WIN32
	// BackendMacOS is the macOS backend.
	BackendMacOS = C.GOOME_MACOS
	// BackendTest is the headless test backend.
	BackendTest = C.GOOME_TEST
	// BackendOffscreen is the offscreen backend, with no OS window.
	BackendOffscreen = C.GOOME_OFFSCREEN
	// ModeOverlay is the mode that draws a see-through layer over every monitor, where only hit rects take clicks.
	ModeOverlay = C.GOOME_MODE_OVERLAY
	// ModeWindow is the mode that opens one normal window.
	ModeWindow = C.GOOME_MODE_WINDOW
	// ModeOffscreen is the mode with no OS window, just a pixel buffer you read.
	ModeOffscreen = C.GOOME_MODE_OFFSCREEN
	// WindowResizable lets the user resize the window.
	WindowResizable = C.GOOME_WINDOW_RESIZABLE
	// WindowBorderless removes the OS title bar and border so you draw your own.
	WindowBorderless = C.GOOME_WINDOW_BORDERLESS
	// WindowTransparent turns the canvas alpha into real per pixel window alpha, so you can make shaped windows.
	WindowTransparent = C.GOOME_WINDOW_TRANSPARENT
	// WindowOnTop keeps the window above other windows.
	WindowOnTop = C.GOOME_WINDOW_ON_TOP
	// WindowHidden starts the window hidden until you show it.
	WindowHidden = C.GOOME_WINDOW_HIDDEN
	// WindowCenter centers the window on the screen when it opens.
	WindowCenter = C.GOOME_WINDOW_CENTER
	// WindowNoTitlebar hides the title bar.
	WindowNoTitlebar = C.GOOME_WINDOW_NO_TITLEBAR
	// WindowNoConsole hides the console window on Windows. It is off by default.
	WindowNoConsole = C.GOOME_WINDOW_NO_CONSOLE
	// OK is the success return code, anything below it is an error.
	OK = C.GOOME_OK
	// CodeNotSupported is the error code for something this backend cannot do.
	CodeNotSupported = C.GOOME_ENOTSUP
	// CodeDenied is the error code for a call that was not allowed, either by permissions or because the user said no.
	CodeDenied = C.GOOME_EDENIED
	// CodeNoDisplay is the error code for having no display to draw on.
	CodeNoDisplay = C.GOOME_ENODISPLAY
	// CodeGone is the error code for a monitor that is gone.
	CodeGone = C.GOOME_EGONE
	// CodeAgain is the error code for nothing being ready yet, so try again in a bit.
	CodeAgain = C.GOOME_EAGAIN
	// CodeNoMem is the error code for running out of memory.
	CodeNoMem = C.GOOME_ENOMEM
	// CodeArg is the error code for a bad argument.
	CodeArg = C.GOOME_EARG
	// CapCapture means Capture works on this backend.
	CapCapture = C.GOOME_CAP_CAPTURE
	// CapGlobalPointer means Pointer reports the mouse everywhere, not only over hit rects.
	CapGlobalPointer = C.GOOME_CAP_GLOBAL_POINTER
	// CapHotkeys means global hotkeys work.
	CapHotkeys = C.GOOME_CAP_HOTKEYS
	// CapHideSelf means HideSelf really keeps the overlay out of screen grabs.
	CapHideSelf = C.GOOME_CAP_HIDE_SELF
	// CapKeyboard means the keyboard can be taken from other apps.
	CapKeyboard = C.GOOME_CAP_KEYBOARD
	// CapWindow means window mode works here.
	CapWindow = C.GOOME_CAP_WINDOW
	// CapClipboard means the clipboard functions reach the real system clipboard.
	CapClipboard = C.GOOME_CAP_CLIPBOARD
	// CapCursor means the cursor shape calls change the real cursor.
	CapCursor = C.GOOME_CAP_CURSOR
	// CapGrab means Grab gets a real pointer grab from the OS.
	CapGrab = C.GOOME_CAP_GRAB
	// WantMouse is set when the pointer is over a hit rect, a press or drag is going, or a popup or modal is open, so the game should leave the mouse alone.
	WantMouse = C.GOOME_WANT_MOUSE
	// WantKeyboard is set when a widget has keyboard focus or a modal is open, so the game should leave the keyboard alone.
	WantKeyboard = C.GOOME_WANT_KEYBOARD
	// EventNone is the event type for nothing happening.
	EventNone = C.GOOME_NONE
	// EventDown is the event for a mouse button press.
	EventDown = C.GOOME_DOWN
	// EventUp is the event for a mouse button release.
	EventUp = C.GOOME_UP
	// EventMotion is the event for the pointer moving.
	EventMotion = C.GOOME_MOTION
	// EventScroll is the event for the mouse wheel.
	EventScroll = C.GOOME_SCROLL
	// EventKey is the event for a key press.
	EventKey = C.GOOME_KEY
	// EventHotkey is the event for a registered hotkey, with its ID in ID.
	EventHotkey = C.GOOME_HOTKEY
	// EventControl is the event for a line sent through the control fifo, in Msg.
	EventControl = C.GOOME_CONTROL
	// EventMonitors is the event for the monitor list changing, so ask for it again.
	EventMonitors = C.GOOME_MONITORS
	// EventQuit is the event for time to go, such as the session ending.
	EventQuit = C.GOOME_QUIT
	// EventClose is the event for the window close button, alt-f4 or cmd-w, and the app decides what to do.
	EventClose = C.GOOME_CLOSE
	// EventFocus is the event for the window gaining or losing focus, with ID 1 for gained and 0 for lost.
	EventFocus = C.GOOME_FOCUS
	// EventDrop is the event for a file dropped on the window, with the path in Msg.
	EventDrop = C.GOOME_DROP
	// EventText is the event for an IME, dead key or compose finishing some text, with the text in Msg.
	EventText = C.GOOME_TEXT
	// EventPreedit is the event for the text being composed changing, with ID as the caret byte offset and an empty Msg when composing ends.
	EventPreedit = C.GOOME_PREEDIT
	// CursorArrow is the normal arrow cursor.
	CursorArrow = C.GOOME_CURSOR_ARROW
	// CursorIBeam is the text insertion cursor.
	CursorIBeam = C.GOOME_CURSOR_IBEAM
	// CursorHand is the pointing hand cursor.
	CursorHand = C.GOOME_CURSOR_HAND
	// CursorMove is the move cursor.
	CursorMove = C.GOOME_CURSOR_MOVE
	// CursorResizeEW is the left and right resize cursor.
	CursorResizeEW = C.GOOME_CURSOR_RESIZE_EW
	// CursorResizeNS is the up and down resize cursor.
	CursorResizeNS = C.GOOME_CURSOR_RESIZE_NS
	// CursorResizeNWSE is the diagonal resize cursor running from top left to bottom right.
	CursorResizeNWSE = C.GOOME_CURSOR_RESIZE_NWSE
	// CursorResizeNESW is the diagonal resize cursor running from top right to bottom left.
	CursorResizeNESW = C.GOOME_CURSOR_RESIZE_NESW
	// CursorWait is the busy cursor.
	CursorWait = C.GOOME_CURSOR_WAIT
	// CursorNo is the not allowed cursor.
	CursorNo = C.GOOME_CURSOR_NO
	// CursorHidden hides the cursor.
	CursorHidden = C.GOOME_CURSOR_HIDDEN
	// EdgeLeft is the left edge bit for WindowDrag resize requests.
	EdgeLeft = C.GOOME_EDGE_LEFT
	// EdgeRight is the right edge bit for WindowDrag resize requests.
	EdgeRight = C.GOOME_EDGE_RIGHT
	// EdgeTop is the top edge bit for WindowDrag resize requests.
	EdgeTop = C.GOOME_EDGE_TOP
	// EdgeBottom is the bottom edge bit for WindowDrag resize requests.
	EdgeBottom = C.GOOME_EDGE_BOTTOM
	// StateNormal is the normal window state.
	StateNormal = C.GOOME_STATE_NORMAL
	// StateMinimized is the minimized window state.
	StateMinimized = C.GOOME_STATE_MINIMIZED
	// StateMaximized is the maximized window state.
	StateMaximized = C.GOOME_STATE_MAXIMIZED
	// StateFullscreen is the fullscreen window state.
	StateFullscreen = C.GOOME_STATE_FULLSCREEN
	// KeyBackspace is the backspace key.
	KeyBackspace = C.GOOME_K_BACKSPACE
	// KeyTab is the tab key.
	KeyTab = C.GOOME_K_TAB
	// KeyEnter is the enter key.
	KeyEnter = C.GOOME_K_ENTER
	// KeyEsc is the escape key.
	KeyEsc = C.GOOME_K_ESC
	// KeySpace is the space key.
	KeySpace = C.GOOME_K_SPACE
	// KeyDelete is the delete key.
	KeyDelete = C.GOOME_K_DELETE
	// KeyF1 is the F1 key, and the other function keys follow it up to KeyF12.
	KeyF1 = C.GOOME_K_F1
	// KeyF12 is the F12 key.
	KeyF12 = C.GOOME_K_F12
	// KeyLeft is the left arrow key.
	KeyLeft = C.GOOME_K_LEFT
	// KeyRight is the right arrow key.
	KeyRight = C.GOOME_K_RIGHT
	// KeyUp is the up arrow key.
	KeyUp = C.GOOME_K_UP
	// KeyDown is the down arrow key.
	KeyDown = C.GOOME_K_DOWN
	// KeyHome is the home key.
	KeyHome = C.GOOME_K_HOME
	// KeyEnd is the end key.
	KeyEnd = C.GOOME_K_END
	// KeyPgUp is the page up key.
	KeyPgUp = C.GOOME_K_PGUP
	// KeyPgDn is the page down key.
	KeyPgDn = C.GOOME_K_PGDN
	// KeyInsert is the insert key.
	KeyInsert = C.GOOME_K_INSERT
	// ModShift is the shift modifier bit.
	ModShift = C.GOOME_SHIFT
	// ModCtrl is the ctrl modifier bit.
	ModCtrl = C.GOOME_CTRL
	// ModAlt is the alt modifier bit.
	ModAlt = C.GOOME_ALT
	// ModSuper is the super modifier bit, the windows or command key.
	ModSuper = C.GOOME_SUPER
	// FmtBGRAPremul is B G R A byte order in memory with premultiplied alpha, which is what goome uses inside.
	FmtBGRAPremul = C.GOOME_FMT_BGRA_PREMUL
	// FmtRGBAPremul is R G B A byte order in memory with premultiplied alpha.
	FmtRGBAPremul = C.GOOME_FMT_RGBA_PREMUL
	// FmtBGRA is B G R A byte order in memory with straight alpha.
	FmtBGRA = C.GOOME_FMT_BGRA
	// FmtRGBA is R G B A byte order in memory with straight alpha.
	FmtRGBA = C.GOOME_FMT_RGBA
	// DirAuto sets each line to run the way its first strong letter goes.
	DirAuto = C.GOOME_DIR_AUTO
	// DirLTR sets every line to run left to right.
	DirLTR = C.GOOME_DIR_LTR
	// DirRTL sets every line to run right to left.
	DirRTL = C.GOOME_DIR_RTL
	// IconArrowDown is the down arrow icon, used on combos, menubar menus and open headers and tree nodes.
	IconArrowDown = C.GOOME_ICON_ARROW_DOWN
	// IconArrowRight is the right arrow icon, used on closed headers, tree nodes and submenus.
	IconArrowRight = C.GOOME_ICON_ARROW_RIGHT
	// IconArrowUp is the up arrow icon, used on open combos and menubar menus.
	IconArrowUp = C.GOOME_ICON_ARROW_UP
	// IconArrowLeft is the left arrow icon, used on the left side of drag numbers.
	IconArrowLeft = C.GOOME_ICON_ARROW_LEFT
	// IconCheck is the check mark icon, used on checkboxes and checked menu items.
	IconCheck = C.GOOME_ICON_CHECK
	// IconRadioDot is the dot icon inside a selected radio button.
	IconRadioDot = C.GOOME_ICON_RADIO_DOT
	// IconClose is the close X icon, used on windows, modals and the app title bar.
	IconClose = C.GOOME_ICON_CLOSE
	// IconMinimize is the minimize button icon in the app title bar.
	IconMinimize = C.GOOME_ICON_MINIMIZE
	// IconMaximize is the maximize button icon in the app title bar.
	IconMaximize = C.GOOME_ICON_MAXIMIZE
	// IconRestore is the restore button icon in the app title bar.
	IconRestore = C.GOOME_ICON_RESTORE
	// IconResize is the window resize grip icon.
	IconResize = C.GOOME_ICON_RESIZE
	// IconSortUp is the icon for a table column sorted ascending.
	IconSortUp = C.GOOME_ICON_SORT_UP
	// IconSortDown is the icon for a table column sorted descending.
	IconSortDown = C.GOOME_ICON_SORT_DOWN
	// IconPlus is the plus icon on number widget buttons.
	IconPlus = C.GOOME_ICON_PLUS
	// IconMinus is the minus icon on number widget buttons.
	IconMinus = C.GOOME_ICON_MINUS
	// IconSliderKnob is the slider knob icon, also used for the bars in the color picker.
	IconSliderKnob = C.GOOME_ICON_SLIDER_KNOB
	// IconToggleKnob is the toggle knob icon.
	IconToggleKnob = C.GOOME_ICON_TOGGLE_KNOB
	// IconCollapse is the icon for the window collapse button.
	IconCollapse = C.GOOME_ICON_COLLAPSE
	// IconExpand is the icon for the window expand button.
	IconExpand = C.GOOME_ICON_EXPAND
	// IconDragHandle is the split pane handle icon.
	IconDragHandle = C.GOOME_ICON_DRAG_HANDLE
	// IconToast is the mark icon drawn on toasts.
	IconToast = C.GOOME_ICON_TOAST
	// IconPicker is the color picker ring icon.
	IconPicker = C.GOOME_ICON_PICKER
	// IconCount is the number of icon slots, not an icon itself.
	IconCount = C.GOOME_ICON_COUNT
	// WinMove lets the user drag the window by its title bar.
	WinMove = C.GOOME_WIN_MOVE
	// WinResize adds a resize grip in the bottom right corner of the window.
	WinResize = C.GOOME_WIN_RESIZE
	// WinCollapse adds an arrow that folds the window to just its title.
	WinCollapse = C.GOOME_WIN_COLLAPSE
	// WinNoTitle hides the window's title bar.
	WinNoTitle = C.GOOME_WIN_NO_TITLE
	// WinClose adds a close X to the window that sets its open flag to false.
	WinClose = C.GOOME_WIN_CLOSE
	// InputPassword draws the mask character per letter and makes copy and cut do nothing.
	InputPassword = C.GOOME_INPUT_PASSWORD
	// InputReadonly blocks typing and pasting while still allowing select and copy.
	InputReadonly = C.GOOME_INPUT_READONLY
	// InputNumber limits the input to digits, plus, minus and dot.
	InputNumber = C.GOOME_INPUT_NUMBER
	// InputNoSelectOnFocus puts the caret at the end when tabbing in instead of selecting everything.
	InputNoSelectOnFocus = C.GOOME_INPUT_NO_SELECT_ON_FOCUS
	// InputEnterReturns makes the input return 2 when enter is pressed, handy for login and search boxes.
	InputEnterReturns = C.GOOME_INPUT_ENTER_RETURNS
	// AppTitlebar draws the app's own title bar with move, maximize and close handling.
	AppTitlebar = C.GOOME_APP_TITLEBAR
	// AppMenubar reserves room for a menu bar at the top of the app.
	AppMenubar = C.GOOME_APP_MENUBAR
	// AppStatusbar adds a status bar at the bottom of the app.
	AppStatusbar = C.GOOME_APP_STATUSBAR
	// AppNoPad removes the padding around the app content.
	AppNoPad = C.GOOME_APP_NO_PAD
	// EditReadonly lets the text editor select and copy but not type.
	EditReadonly = C.GOOME_EDIT_READONLY
	// EditLineNumbers adds a gutter with line numbers to the text editor.
	EditLineNumbers = C.GOOME_EDIT_LINE_NUMBERS
	// EditWrap wraps long lines in the text editor instead of scrolling sideways.
	EditWrap = C.GOOME_EDIT_WRAP
	// EditMono makes the text editor use a monospace font.
	EditMono = C.GOOME_EDIT_MONO
	// TableSortable lets the user click a header to sort by it, and again to flip it.
	TableSortable = C.GOOME_TABLE_SORTABLE
	// TableResizable lets the user drag between headers to resize columns.
	TableResizable = C.GOOME_TABLE_RESIZABLE
	// TableStripes makes every other table row a bit darker.
	TableStripes = C.GOOME_TABLE_STRIPES
	// TableBorders draws lines between table cells.
	TableBorders = C.GOOME_TABLE_BORDERS
	// StyleBg is the theme color for what sits behind everything, which app mode fills the window with.
	StyleBg = C.GOOME_S_BG
	// StylePanel is the theme color for the window background color.
	StylePanel = C.GOOME_S_PANEL
	// StyleBorder is the theme color for window and widget outlines.
	StyleBorder = C.GOOME_S_BORDER
	// StyleText is the theme color for the normal text color.
	StyleText = C.GOOME_S_TEXT
	// StyleTextDim is the theme color for the quieter text color.
	StyleTextDim = C.GOOME_S_TEXT_DIM
	// StyleAccent is the theme color for checked boxes, slider fill and the active tab.
	StyleAccent = C.GOOME_S_ACCENT
	// StyleAccentText is the theme color for text drawn on top of the accent.
	StyleAccentText = C.GOOME_S_ACCENT_TEXT
	// StyleHover is the theme color for the widget background color while hovered.
	StyleHover = C.GOOME_S_HOVER
	// StyleActive is the theme color for the widget background color while pressed.
	StyleActive = C.GOOME_S_ACTIVE
	// StyleInputBg is the theme color for the inside inputs and sliders.
	StyleInputBg = C.GOOME_S_INPUT_BG
	// StyleShadow is the theme color for the window shadow color.
	StyleShadow = C.GOOME_S_SHADOW
	// StyleRadius is the theme number for the corner radius.
	StyleRadius = C.GOOME_S_RADIUS
	// StylePad is the theme number for the space inside windows.
	StylePad = C.GOOME_S_PAD
	// StyleGap is the theme number for the space between widgets.
	StyleGap = C.GOOME_S_GAP
	// StyleBorderWidth is the theme number for the outline width.
	StyleBorderWidth = C.GOOME_S_BORDER_WIDTH
	// StyleFontSize is the theme number for the text size.
	StyleFontSize = C.GOOME_S_FONT_SIZE
	// StyleShadowBlur is the theme number for how soft shadows are.
	StyleShadowBlur = C.GOOME_S_SHADOW_BLUR
	// StyleAnimMs is the theme number for the animation length in milliseconds, with 0 meaning no animation.
	StyleAnimMs = C.GOOME_S_ANIM_MS
	// StyleFocus is the theme color for the focus ring and input glow color, with 0 meaning the accent.
	StyleFocus = C.GOOME_S_FOCUS
	// StyleSelect is the theme color for the selected text color, with 0 meaning a faded accent.
	StyleSelect = C.GOOME_S_SELECT
	// StyleCaret is the theme color for the text caret color, with 0 meaning the accent.
	StyleCaret = C.GOOME_S_CARET
	// StyleLink is the theme color for the link color, with 0 meaning the accent.
	StyleLink = C.GOOME_S_LINK
	// StyleBadge is the theme color for the default badge color, with 0 meaning the accent.
	StyleBadge = C.GOOME_S_BADGE
	// StyleToast is the theme color for the toast background color, with 0 meaning the panel mixed with the hover color.
	StyleToast = C.GOOME_S_TOAST
	// StyleTooltip is the theme color for the tooltip background color, with 0 meaning the background mixed with the hover color.
	StyleTooltip = C.GOOME_S_TOOLTIP
	// StyleDanger is the theme color for the red close buttons, with 0 meaning the built in red.
	StyleDanger = C.GOOME_S_DANGER
	// StyleKnob is the theme color for the slider and toggle knob color, with 0 meaning white.
	StyleKnob = C.GOOME_S_KNOB
	// StyleDim is the theme color for the drawn behind modals, with 0 meaning 0x73000000.
	StyleDim = C.GOOME_S_DIM
	// StyleCheckSize is the theme number for the checkbox and radio box size.
	StyleCheckSize = C.GOOME_S_CHECK_SIZE
	// StyleCheckRadius is the theme number for the checkbox corner radius.
	StyleCheckRadius = C.GOOME_S_CHECK_RADIUS
	// StyleToggleW is the theme number for the toggle track width.
	StyleToggleW = C.GOOME_S_TOGGLE_W
	// StyleToggleH is the theme number for the toggle track height.
	StyleToggleH = C.GOOME_S_TOGGLE_H
	// StyleSliderH is the theme number for the slider height.
	StyleSliderH = C.GOOME_S_SLIDER_H
	// StyleSliderTrack is the theme number for the slider track thickness.
	StyleSliderTrack = C.GOOME_S_SLIDER_TRACK
	// StyleKnobSize is the theme number for the slider knob radius.
	StyleKnobSize = C.GOOME_S_KNOB_SIZE
	// StyleScrollW is the theme number for the scrollbar width.
	StyleScrollW = C.GOOME_S_SCROLL_W
	// StyleScrollMin is the theme number for the shortest a scrollbar thumb gets.
	StyleScrollMin = C.GOOME_S_SCROLL_MIN
	// StyleItemH is the theme number for the height of every widget row, with 0 meaning line height plus pad.
	StyleItemH = C.GOOME_S_ITEM_H
	// StyleMenuItemH is the theme number for the menu row height, with 0 meaning it comes from the font.
	StyleMenuItemH = C.GOOME_S_MENU_ITEM_H
	// StyleTabH is the theme number for the tab height, with 0 meaning it comes from the font.
	StyleTabH = C.GOOME_S_TAB_H
	// StyleWinTitleH is the theme number for the window title bar height, with 0 meaning it comes from the font.
	StyleWinTitleH = C.GOOME_S_WIN_TITLE_H
	// StyleAppTitleH is the theme number for the app title bar height, with 0 meaning it comes from the font.
	StyleAppTitleH = C.GOOME_S_APP_TITLE_H
	// StyleMenubarH is the theme number for the app menu bar height, with 0 meaning it comes from the font.
	StyleMenubarH = C.GOOME_S_MENUBAR_H
	// StyleStatusbarH is the theme number for the app status bar height, with 0 meaning it comes from the font.
	StyleStatusbarH = C.GOOME_S_STATUSBAR_H
	// StyleFocusW is the theme number for the focus ring width.
	StyleFocusW = C.GOOME_S_FOCUS_W
	// StyleFocusGap is the theme number for the space between a widget and its focus ring.
	StyleFocusGap = C.GOOME_S_FOCUS_GAP
	// StyleSepW is the theme number for the separator thickness.
	StyleSepW = C.GOOME_S_SEP_W
	// StyleGrip is the theme number for the window resize grip size.
	StyleGrip = C.GOOME_S_GRIP
	// StyleIndent is the theme number for the tree indent.
	StyleIndent = C.GOOME_S_INDENT
	// StylePlotW is the theme number for the plot line width.
	StylePlotW = C.GOOME_S_PLOT_W
	// StyleSpinner is the theme number for the default spinner size, with 0 meaning line height.
	StyleSpinner = C.GOOME_S_SPINNER
	// StyleCaretW is the theme number for the caret width, with 0 meaning thin and a wide value giving a block caret.
	StyleCaretW = C.GOOME_S_CARET_W
	// StyleWinMinW is the theme number for the smallest width a window resizes to.
	StyleWinMinW = C.GOOME_S_WIN_MIN_W
	// StyleWinMinH is the theme number for the smallest height a window resizes to.
	StyleWinMinH = C.GOOME_S_WIN_MIN_H
	// StylePopGap is the theme number for the space between a combo or popup and what opened it.
	StylePopGap = C.GOOME_S_POP_GAP
	// StyleMenuGap is the theme number for the offset of submenus.
	StyleMenuGap = C.GOOME_S_MENU_GAP
	// StyleTipPadX is the theme number for the horizontal padding inside tooltips.
	StyleTipPadX = C.GOOME_S_TIP_PAD_X
	// StyleTipPadY is the theme number for the vertical padding inside tooltips.
	StyleTipPadY = C.GOOME_S_TIP_PAD_Y
	// StyleTipX is the theme number for the tooltip horizontal offset from the mouse.
	StyleTipX = C.GOOME_S_TIP_X
	// StyleTipY is the theme number for the tooltip vertical offset from the mouse.
	StyleTipY = C.GOOME_S_TIP_Y
	// StyleTipRadius is the theme number for the tooltip corner radius.
	StyleTipRadius = C.GOOME_S_TIP_RADIUS
	// StyleShadowX is the theme number for the shadow horizontal offset.
	StyleShadowX = C.GOOME_S_SHADOW_X
	// StyleShadowY is the theme number for the shadow vertical offset.
	StyleShadowY = C.GOOME_S_SHADOW_Y
	// StyleAlign is the theme setting for button and label text alignment, one of the Align constants.
	StyleAlign = C.GOOME_S_ALIGN
	// StyleHeadingSize is the theme number for the heading text size, with 0 meaning 1.3 times the font size.
	StyleHeadingSize = C.GOOME_S_HEADING_SIZE
	// StyleSmallSize is the theme number for the size of badges and statusbar text, with 0 meaning a bit under the font size.
	StyleSmallSize = C.GOOME_S_SMALL_SIZE
	// StyleEase is the theme setting for the easing curve widgets use, one of the Ease constants.
	StyleEase = C.GOOME_S_EASE
	// StyleHoverMs is the theme number for how long hover fades take in milliseconds, with below 0 following AnimMs.
	StyleHoverMs = C.GOOME_S_HOVER_MS
	// StylePressMs is the theme number for how long press fades take in milliseconds, with below 0 following AnimMs.
	StylePressMs = C.GOOME_S_PRESS_MS
	// StylePopupMs is the theme number for how long popups and menus take to open in milliseconds, with below 0 following AnimMs.
	StylePopupMs = C.GOOME_S_POPUP_MS
	// StyleCollapseMs is the theme number for how long anything that folds open or shut takes in milliseconds, with below 0 following AnimMs.
	StyleCollapseMs = C.GOOME_S_COLLAPSE_MS
	// StyleToastFadeMs is the theme number for how long toasts take to fade out in milliseconds, with below 0 following AnimMs.
	StyleToastFadeMs = C.GOOME_S_TOAST_FADE_MS
	// StyleCount is the number of style slots, not a style itself.
	StyleCount = C.GOOME_S_COUNT
	// AlignAuto lets each widget pick its own text alignment.
	AlignAuto = C.GOOME_ALIGN_AUTO
	// AlignLeft aligns button and label text to the left.
	AlignLeft = C.GOOME_ALIGN_LEFT
	// AlignCenter centers button and label text.
	AlignCenter = C.GOOME_ALIGN_CENTER
	// AlignRight aligns button and label text to the right.
	AlignRight = C.GOOME_ALIGN_RIGHT
	// EaseSmooth is the smoothstep easing curve the widgets use by default.
	EaseSmooth = C.GOOME_EASE_SMOOTH
	// EaseLinear is an easing curve with the same speed all the way.
	EaseLinear = C.GOOME_EASE_LINEAR
	// EaseOut is an easing curve that starts fast and slows down at the end.
	EaseOut = C.GOOME_EASE_OUT
	// EaseInOut is an easing curve that is slow, then fast, then slow.
	EaseInOut = C.GOOME_EASE_IN_OUT
	// EaseBack is an easing curve that goes a little past the target and comes back.
	EaseBack = C.GOOME_EASE_BACK
	// EaseCustom is an easing curve that calls the function set in the UI config.
	EaseCustom = C.GOOME_EASE_CUSTOM
	// EaseIn is an easing curve that starts slow and ends fast.
	EaseIn = C.GOOME_EASE_IN
	// EaseSpring is an easing curve like a damped spring, about 16 percent over before it settles.
	EaseSpring = C.GOOME_EASE_SPRING
	// EaseBounce is an easing curve that drops and bounces like a ball.
	EaseBounce = C.GOOME_EASE_BOUNCE
	// SkinButton is the skin kind for buttons.
	SkinButton = C.GOOME_SKIN_BUTTON
	// SkinPanel is the skin kind for windows.
	SkinPanel = C.GOOME_SKIN_PANEL
	// SkinInput is the skin kind for text inputs, the text editor and drag numbers.
	SkinInput = C.GOOME_SKIN_INPUT
	// SkinCheckbox is the skin kind for the checkbox box.
	SkinCheckbox = C.GOOME_SKIN_CHECKBOX
	// SkinRadio is the skin kind for the radio circle.
	SkinRadio = C.GOOME_SKIN_RADIO
	// SkinToggle is the skin kind for the track behind the toggle knob.
	SkinToggle = C.GOOME_SKIN_TOGGLE
	// SkinSliderTrack is the skin kind for the slider track.
	SkinSliderTrack = C.GOOME_SKIN_SLIDER_TRACK
	// SkinSliderFill is the skin kind for the filled part of a slider.
	SkinSliderFill = C.GOOME_SKIN_SLIDER_FILL
	// SkinScrollbar is the skin kind for the scrollbar thumb.
	SkinScrollbar = C.GOOME_SKIN_SCROLLBAR
	// SkinTab is the skin kind for tabs.
	SkinTab = C.GOOME_SKIN_TAB
	// SkinTitlebar is the skin kind for window and app title bars.
	SkinTitlebar = C.GOOME_SKIN_TITLEBAR
	// SkinTooltip is the skin kind for tooltips.
	SkinTooltip = C.GOOME_SKIN_TOOLTIP
	// SkinMenu is the skin kind for menus and the open combo list.
	SkinMenu = C.GOOME_SKIN_MENU
	// SkinPopup is the skin kind for popups, modals and the color picker.
	SkinPopup = C.GOOME_SKIN_POPUP
	// SkinCount is the number of skin kinds, not a kind itself.
	SkinCount = C.GOOME_SKIN_COUNT
	// ItemDragX makes a custom item draggable, sideways only when used alone.
	ItemDragX = C.GOOME_ITEM_DRAG_X
	// ItemDragY makes a custom item draggable, up and down only when used alone. With ItemDragX it drags freely.
	ItemDragY = C.GOOME_ITEM_DRAG_Y
	// ItemFocusable gives a custom item tab focus and raw key input.
	ItemFocusable = C.GOOME_ITEM_FOCUSABLE
	// ItemNoHit makes a custom item never hovered or pressed, for backgrounds you put other widgets on.
	ItemNoHit = C.GOOME_ITEM_NO_HIT
	// ItemRepeat makes a custom item's clicked fire on the press and keep firing while held.
	ItemRepeat = C.GOOME_ITEM_REPEAT
	// PresetDark is the dark theme preset.
	PresetDark = C.GOOME_PRESET_DARK
	// PresetLight is the light theme preset.
	PresetLight = C.GOOME_PRESET_LIGHT
	// PresetMidnight is the midnight theme preset.
	PresetMidnight = C.GOOME_PRESET_MIDNIGHT
	// PresetPaper is the paper theme preset.
	PresetPaper = C.GOOME_PRESET_PAPER
	// PresetMint is the mint theme preset.
	PresetMint = C.GOOME_PRESET_MINT
	// PresetHighContrast is the high contrast theme preset.
	PresetHighContrast = C.GOOME_PRESET_HIGH_CONTRAST
	// PresetRose is the rose theme preset.
	PresetRose = C.GOOME_PRESET_ROSE
	// PresetCount is the number of theme presets, not a preset itself.
	PresetCount = C.GOOME_PRESET_COUNT
	// ActCopy is the copy action in the UI key map.
	ActCopy = C.GOOME_ACT_COPY
	// ActCut is the cut action in the UI key map.
	ActCut = C.GOOME_ACT_CUT
	// ActPaste is the paste action in the UI key map.
	ActPaste = C.GOOME_ACT_PASTE
	// ActUndo is the undo action in the UI key map.
	ActUndo = C.GOOME_ACT_UNDO
	// ActRedo is the redo action in the UI key map.
	ActRedo = C.GOOME_ACT_REDO
	// ActSelectAll is the select all action in the UI key map.
	ActSelectAll = C.GOOME_ACT_SELECT_ALL
	// ActWordLeft is the move one word left action in the UI key map.
	ActWordLeft = C.GOOME_ACT_WORD_LEFT
	// ActWordRight is the move one word right action in the UI key map.
	ActWordRight = C.GOOME_ACT_WORD_RIGHT
	// ActDelWordLeft is the delete word to the left action in the UI key map.
	ActDelWordLeft = C.GOOME_ACT_DEL_WORD_LEFT
	// ActDelWordRight is the delete word to the right action in the UI key map.
	ActDelWordRight = C.GOOME_ACT_DEL_WORD_RIGHT
	// ActDocStart is the jump to the start of the text action in the UI key map.
	ActDocStart = C.GOOME_ACT_DOC_START
	// ActDocEnd is the jump to the end of the text action in the UI key map.
	ActDocEnd = C.GOOME_ACT_DOC_END
	// ActCount is the number of key map actions, not an action itself.
	ActCount = C.GOOME_ACT_COUNT
	// CurButton is the UI cursor kind used over buttons.
	CurButton = C.GOOME_CUR_BUTTON
	// CurInput is the UI cursor kind used over text inputs.
	CurInput = C.GOOME_CUR_INPUT
	// CurLink is the UI cursor kind used over links.
	CurLink = C.GOOME_CUR_LINK
	// CurDragValue is the UI cursor kind used over drag number widgets.
	CurDragValue = C.GOOME_CUR_DRAG_VALUE
	// CurSplitX is the UI cursor kind used over a vertical split handle.
	CurSplitX = C.GOOME_CUR_SPLIT_X
	// CurSplitY is the UI cursor kind used over a horizontal split handle.
	CurSplitY = C.GOOME_CUR_SPLIT_Y
	// CurGrip is the UI cursor kind used over the window resize grip.
	CurGrip = C.GOOME_CUR_GRIP
	// CurColumn is the UI cursor kind used over table column borders.
	CurColumn = C.GOOME_CUR_COLUMN
	// CurDragDrop is the UI cursor kind used during a drag and drop.
	CurDragDrop = C.GOOME_CUR_DRAG_DROP
	// CurCustomX is the UI cursor kind used over custom widgets that drag sideways.
	CurCustomX = C.GOOME_CUR_CUSTOM_X
	// CurCustomY is the UI cursor kind used over custom widgets that drag up and down.
	CurCustomY = C.GOOME_CUR_CUSTOM_Y
	// CurCustomXy is the UI cursor kind used over custom widgets that drag freely.
	CurCustomXy = C.GOOME_CUR_CUSTOM_XY
	// CurCount is the number of UI cursor kinds, not a kind itself.
	CurCount = C.GOOME_CUR_COUNT
	// RoleNone is the accessibility role for a node with no role.
	RoleNone = C.GOOME_ROLE_NONE
	// RoleWindow is the accessibility role for a window.
	RoleWindow = C.GOOME_ROLE_WINDOW
	// RoleButton is the accessibility role for a button.
	RoleButton = C.GOOME_ROLE_BUTTON
	// RoleCheckbox is the accessibility role for a checkbox.
	RoleCheckbox = C.GOOME_ROLE_CHECKBOX
	// RoleRadio is the accessibility role for a radio button.
	RoleRadio = C.GOOME_ROLE_RADIO
	// RoleSlider is the accessibility role for a slider.
	RoleSlider = C.GOOME_ROLE_SLIDER
	// RoleText is the accessibility role for plain text.
	RoleText = C.GOOME_ROLE_TEXT
	// RoleInput is the accessibility role for a text input.
	RoleInput = C.GOOME_ROLE_INPUT
	// RoleCombo is the accessibility role for a combo box.
	RoleCombo = C.GOOME_ROLE_COMBO
	// RoleList is the accessibility role for a list.
	RoleList = C.GOOME_ROLE_LIST
	// RoleItem is the accessibility role for a list item.
	RoleItem = C.GOOME_ROLE_ITEM
	// RoleTab is the accessibility role for a tab.
	RoleTab = C.GOOME_ROLE_TAB
	// RoleMenu is the accessibility role for a menu.
	RoleMenu = C.GOOME_ROLE_MENU
	// RoleMenuItem is the accessibility role for a menu item.
	RoleMenuItem = C.GOOME_ROLE_MENUITEM
	// RoleTable is the accessibility role for a table.
	RoleTable = C.GOOME_ROLE_TABLE
	// RoleTree is the accessibility role for a tree.
	RoleTree = C.GOOME_ROLE_TREE
	// RoleProgress is the accessibility role for a progress bar.
	RoleProgress = C.GOOME_ROLE_PROGRESS
	// RoleImage is the accessibility role for an image.
	RoleImage = C.GOOME_ROLE_IMAGE
	// RoleGroup is the accessibility role for a group of widgets.
	RoleGroup = C.GOOME_ROLE_GROUP
	// RoleDialog is the accessibility role for a dialog.
	RoleDialog = C.GOOME_ROLE_DIALOG
	// A11yFocused is set on an accessibility node that has keyboard focus.
	A11yFocused = C.GOOME_A11Y_FOCUSED
	// A11yChecked is set on an accessibility node that is checked.
	A11yChecked = C.GOOME_A11Y_CHECKED
	// A11yDisabled is set on an accessibility node that is disabled.
	A11yDisabled = C.GOOME_A11Y_DISABLED
	// A11yExpanded is set on an accessibility node that is expanded.
	A11yExpanded = C.GOOME_A11Y_EXPANDED
	// A11ySelected is set on an accessibility node that is selected.
	A11ySelected = C.GOOME_A11Y_SELECTED
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
	"GOOME_WINDOW_NO_TITLEBAR":       WindowNoTitlebar,
	"GOOME_WINDOW_NO_CONSOLE":        WindowNoConsole,
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
