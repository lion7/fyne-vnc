package fynevnc

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	vnc "github.com/amitbet/vnc2video"
)

const (
	framerate = 10
	minWidth  = 32
	minHeight = 24
)

type VncDisplay struct {
	widget.BaseWidget
	keyboardHandler
	mouseHandler

	closed  bool
	client  *vnc.ClientConn
	config  *vnc.ClientConfig
	display *canvas.Image
}

func (v *VncDisplay) Close() {
	v.closed = true
	v.client.Close()
}

func (v *VncDisplay) MinSize() fyne.Size {
	v.ExtendBaseWidget(v)
	return fyne.Size{
		Width:  minWidth,
		Height: minHeight,
	}
}

func (v *VncDisplay) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.display)
}

// ConnectVncDisplay creates a new VncDisplay and does all the heavy lifting, setting up all event handlers
func ConnectVncDisplay(addr string, config *vnc.ClientConfig) (*VncDisplay, error) {
	client, err := connectVnc(addr, config)
	if err != nil {
		return nil, err
	}

	for _, encoding := range config.Encodings {
		renderer, ok := encoding.(vnc.Renderer)
		if ok {
			renderer.SetTargetImage(client.Canvas)
		}
	}

	err = client.SetEncodings([]vnc.EncodingType{
		vnc.EncCursorPseudo,
		vnc.EncPointerPosPseudo,
		vnc.EncCopyRect,
		//vnc.EncTight,
		vnc.EncZRLE,
		vnc.EncRaw,
		vnc.EncHextile,
		vnc.EncZlib,
		vnc.EncRRE,
	})
	if err != nil {
		return nil, fmt.Errorf("error setting encodings: %v\n", err)
	}

	// Create a fyne canvas image from our screen image
	display := canvas.NewImageFromImage(client.Canvas)
	display.FillMode = canvas.ImageFillContain

	// Instantiate the VncDisplay
	viewer := &VncDisplay{client: client, config: config, display: display}

	// Set the initial size equal to the framebuffer size
	viewer.Resize(fyne.NewSize(float32(client.Width()), float32(client.Height())))

	// Add handler for keyboard events
	viewer.handleKeyEvent = func(msg vnc.KeyEvent) {
		err := msg.Write(client)
		if err != nil {
			fmt.Printf("error sending key event: %v\n", err)
		}
	}

	// Add handler for mouse / pointer events
	viewer.handlePointerEvent = func(msg vnc.PointerEvent) {
		err := msg.Write(client)
		if err != nil {
			fmt.Printf("error sending pointer event: %v\n", err)
		}
		display.Refresh()
	}

	// Request framebuffer updates 10 times per second
	go viewer.PeriodicallyRequestFramebufferUpdate(framerate)

	// Refresh the display when we receive a framebuffer update
	go viewer.RefreshOnFramebufferUpdate()

	// Record a video
	//go viewer.RecordVideo(framerate)

	return viewer, nil
}
