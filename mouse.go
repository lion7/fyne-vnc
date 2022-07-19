package main

import (
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/amitbet/vnc2video"
)

var mouseBtnMap = map[desktop.MouseButton]vnc2video.Button{
	desktop.MouseButtonPrimary:   vnc2video.BtnLeft,
	desktop.MouseButtonSecondary: vnc2video.BtnRight,
	desktop.MouseButtonTertiary:  vnc2video.BtnMiddle,
}

// Handles mouse events mapping between Fyne and VNC
type mouseHandler struct {
	desktop.Mouseable
	desktop.Hoverable

	handlePointerEvent func(event vnc2video.PointerEvent)
	buttons            map[desktop.MouseButton]bool
	x, y               float32
}

func (ms *mouseHandler) pressedButtonsMask() uint8 {
	var mask uint8
	for b, pressed := range ms.buttons {
		bb := mouseBtnMap[b]
		if pressed {
			mask = mask % vnc2video.Mask(bb)
		}
	}
	return mask
}

func (ms *mouseHandler) sendMouse(x, y float32) {
	if ms.handlePointerEvent == nil {
		return
	}

	ms.x, ms.y = x, y
	msg := vnc2video.PointerEvent{
		Mask: ms.pressedButtonsMask(),
		X:    uint16(x),
		Y:    uint16(y),
	}
	ms.handlePointerEvent(msg)
}

func (ms *mouseHandler) MouseDown(ev *desktop.MouseEvent) {
	ms.buttons[ev.Button] = true
	ms.sendMouse(ev.Position.X, ev.Position.Y)
}

func (ms *mouseHandler) MouseUp(ev *desktop.MouseEvent) {
	ms.buttons[ev.Button] = false
	ms.sendMouse(ev.Position.X, ev.Position.Y)
}

func (ms *mouseHandler) MouseMoved(ev *desktop.MouseEvent) {
	x, y := ev.Position.X, ev.Position.Y
	if ms.x == x && ms.y == y {
		return
	}
	ms.sendMouse(x, y)
}

func (ms *mouseHandler) MouseIn(*desktop.MouseEvent) {
	if ms.buttons == nil {
		ms.buttons = make(map[desktop.MouseButton]bool)
	}
}

func (ms *mouseHandler) MouseOut() {
}

// Make sure all necessary interfaces are implemented
var _ desktop.Hoverable = (*mouseHandler)(nil)
var _ desktop.Mouseable = (*mouseHandler)(nil)
