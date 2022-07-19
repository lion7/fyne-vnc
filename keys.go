package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/amitbet/vnc2video"
)

// Handles keyboard events mapping between Fyne and VNC
type keyboardHandler struct {
	desktop.Keyable

	handleKeyEvent func(event vnc2video.KeyEvent)
	shift          bool
	caps           bool
}

func (ks *keyboardHandler) TypedKey(ev *fyne.KeyEvent) {
	keyName := ev.Name
	k, ok := keyMap[keyName]
	if !ok && len(keyName) == 1 && keyName[0] < 128 {
		k = vnc2video.Key(keyName[0])
	}
	if k > 0 {
		ks.sendKey(k, true)
		ks.sendKey(k, false)
	}
}

func (ks *keyboardHandler) KeyDown(ev *fyne.KeyEvent) {
	ks.handleDesktopKey(ev.Name, true)
}

func (ks *keyboardHandler) KeyUp(ev *fyne.KeyEvent) {
	ks.handleDesktopKey(ev.Name, false)
}

func (ks *keyboardHandler) handleDesktopKey(keyName fyne.KeyName, pressed bool) {
	if keyCode, ok := desktopKeyMap[keyName]; ok {
		if keyCode == vnc2video.ShiftLeft || keyCode == vnc2video.ShiftRight {
			ks.shift = pressed
		}
		if keyCode == vnc2video.CapsLock {
			ks.caps = pressed
		}
		ks.sendKey(keyCode, pressed)
	}
}

func (ks *keyboardHandler) sendKey(key vnc2video.Key, pressed bool) {
	if ks.handleKeyEvent == nil {
		return
	}

	if key >= 'A' && key <= 'Z' {
		if !ks.shift && !ks.caps {
			key = key + 32
		}
	}
	var down uint8
	if pressed {
		down = 1
	}
	msg := vnc2video.KeyEvent{
		Down: down,
		Key:  key,
	}
	ks.handleKeyEvent(msg)
}

func (ks *keyboardHandler) Focused() bool {
	return true
}

func (ks *keyboardHandler) FocusGained() {
}

func (ks *keyboardHandler) FocusLost() {
}

func (ks *keyboardHandler) TypedRune(ch rune) {
}

var (
	keyMap        map[fyne.KeyName]vnc2video.Key
	desktopKeyMap map[fyne.KeyName]vnc2video.Key
)

func init() {
	desktopKeyMap = map[fyne.KeyName]vnc2video.Key{
		desktop.KeyAltLeft:      vnc2video.AltLeft,
		desktop.KeyAltRight:     vnc2video.AltRight,
		desktop.KeyControlLeft:  vnc2video.ControlLeft,
		desktop.KeyControlRight: vnc2video.ControlRight,
		desktop.KeyShiftLeft:    vnc2video.ShiftLeft,
		desktop.KeyShiftRight:   vnc2video.ShiftRight,
		desktop.KeySuperLeft:    vnc2video.SuperLeft,
		desktop.KeySuperRight:   vnc2video.SuperRight,
		desktop.KeyCapsLock:     vnc2video.CapsLock,
	}

	keyMap = map[fyne.KeyName]vnc2video.Key{
		fyne.KeySpace:     vnc2video.Space,
		fyne.KeyBackspace: vnc2video.BackSpace,
		fyne.KeyDelete:    vnc2video.Delete,
		fyne.KeyDown:      vnc2video.Down,
		fyne.KeyEnd:       vnc2video.End,
		fyne.KeyEnter:     vnc2video.Return,
		fyne.KeyReturn:    vnc2video.Return,
		fyne.KeyEscape:    vnc2video.Escape,
		fyne.KeyF1:        vnc2video.F1,
		fyne.KeyF2:        vnc2video.F2,
		fyne.KeyF3:        vnc2video.F3,
		fyne.KeyF4:        vnc2video.F4,
		fyne.KeyF5:        vnc2video.F5,
		fyne.KeyF6:        vnc2video.F6,
		fyne.KeyF7:        vnc2video.F7,
		fyne.KeyF8:        vnc2video.F8,
		fyne.KeyF9:        vnc2video.F9,
		fyne.KeyF10:       vnc2video.F10,
		fyne.KeyF11:       vnc2video.F11,
		fyne.KeyF12:       vnc2video.F12,
		fyne.KeyHome:      vnc2video.Home,
		fyne.KeyInsert:    vnc2video.Insert,
		fyne.KeyLeft:      vnc2video.Left,
		fyne.KeyPageDown:  vnc2video.PageDown,
		fyne.KeyPageUp:    vnc2video.PageUp,
		fyne.KeyRight:     vnc2video.Right,
		fyne.KeyTab:       vnc2video.Tab,
		fyne.KeyUp:        vnc2video.Up,
		//:        vnc2video.Pause,
		//:  vnc2video.PrintScreen,
		//:      vnc2video.NumLock,
		//:         vnc2video.Meta,
		//:          vnc2video.Win,
	}
}

// Make sure all necessary interfaces are implemented
var _ desktop.Keyable = (*keyboardHandler)(nil)
