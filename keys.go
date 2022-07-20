package fynevnc

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	vnc "github.com/amitbet/vnc2video"
)

// Handles keyboard events mapping between Fyne and VNC
type keyboardHandler struct {
	desktop.Keyable

	config *vnc.ClientConfig
	shift  bool
	caps   bool
}

func (ks *keyboardHandler) TypedKey(ev *fyne.KeyEvent) {
	keyName := ev.Name
	k, ok := keyMap[keyName]
	if !ok && len(keyName) == 1 && keyName[0] < 128 {
		k = vnc.Key(keyName[0])
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
		if keyCode == vnc.ShiftLeft || keyCode == vnc.ShiftRight {
			ks.shift = pressed
		}
		if keyCode == vnc.CapsLock {
			ks.caps = pressed
		}
		ks.sendKey(keyCode, pressed)
	}
}

func (ks *keyboardHandler) sendKey(key vnc.Key, pressed bool) {
	if ks.config == nil {
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
	ks.config.ClientMessageCh <- &vnc.KeyEvent{
		Down: down,
		Key:  key,
	}
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
	keyMap        map[fyne.KeyName]vnc.Key
	desktopKeyMap map[fyne.KeyName]vnc.Key
)

func init() {
	desktopKeyMap = map[fyne.KeyName]vnc.Key{
		desktop.KeyAltLeft:      vnc.AltLeft,
		desktop.KeyAltRight:     vnc.AltRight,
		desktop.KeyControlLeft:  vnc.ControlLeft,
		desktop.KeyControlRight: vnc.ControlRight,
		desktop.KeyShiftLeft:    vnc.ShiftLeft,
		desktop.KeyShiftRight:   vnc.ShiftRight,
		desktop.KeySuperLeft:    vnc.SuperLeft,
		desktop.KeySuperRight:   vnc.SuperRight,
		desktop.KeyCapsLock:     vnc.CapsLock,
	}

	keyMap = map[fyne.KeyName]vnc.Key{
		fyne.KeySpace:     vnc.Space,
		fyne.KeyBackspace: vnc.BackSpace,
		fyne.KeyDelete:    vnc.Delete,
		fyne.KeyDown:      vnc.Down,
		fyne.KeyEnd:       vnc.End,
		fyne.KeyEnter:     vnc.Return,
		fyne.KeyReturn:    vnc.Return,
		fyne.KeyEscape:    vnc.Escape,
		fyne.KeyF1:        vnc.F1,
		fyne.KeyF2:        vnc.F2,
		fyne.KeyF3:        vnc.F3,
		fyne.KeyF4:        vnc.F4,
		fyne.KeyF5:        vnc.F5,
		fyne.KeyF6:        vnc.F6,
		fyne.KeyF7:        vnc.F7,
		fyne.KeyF8:        vnc.F8,
		fyne.KeyF9:        vnc.F9,
		fyne.KeyF10:       vnc.F10,
		fyne.KeyF11:       vnc.F11,
		fyne.KeyF12:       vnc.F12,
		fyne.KeyHome:      vnc.Home,
		fyne.KeyInsert:    vnc.Insert,
		fyne.KeyLeft:      vnc.Left,
		fyne.KeyPageDown:  vnc.PageDown,
		fyne.KeyPageUp:    vnc.PageUp,
		fyne.KeyRight:     vnc.Right,
		fyne.KeyTab:       vnc.Tab,
		fyne.KeyUp:        vnc.Up,
		//:        vnc.Pause,
		//:  vnc.PrintScreen,
		//:      vnc.NumLock,
		//:         vnc.Meta,
		//:          vnc.Win,
	}
}

// Make sure all necessary interfaces are implemented
var _ desktop.Keyable = (*keyboardHandler)(nil)
