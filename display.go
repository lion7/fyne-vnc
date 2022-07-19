package fynevnc

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	vnc "github.com/amitbet/vnc2video"
)

const (
	minWidth  = 32
	minHeight = 24
)

type VncDisplay struct {
	widget.BaseWidget
	keyboardHandler
	mouseHandler

	Client  *vnc.ClientConn
	Display *canvas.Image
}

func (v *VncDisplay) MinSize() fyne.Size {
	v.ExtendBaseWidget(v)
	return fyne.Size{
		Width:  minWidth,
		Height: minHeight,
	}
}

func (v *VncDisplay) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.Display)
}

// NewVncDisplay creates a new VncDisplay and does all the heavy lifting, setting up all event handlers
func NewVncDisplay(addr string, config *vnc.ClientConfig) *VncDisplay {
	client := connectVnc(addr, config)
	screenImage := client.Canvas.Image

	for _, encoding := range config.Encodings {
		renderer, ok := encoding.(vnc.Renderer)
		if ok {
			renderer.SetTargetImage(screenImage)
		}
	}

	err := client.SetEncodings([]vnc.EncodingType{
		//vnc.EncCursorPseudo,
		//vnc.EncPointerPosPseudo,
		//vnc.EncCopyRect,
		vnc.EncTight,
		//vnc.EncZRLE,
		vnc.EncRaw,
		//vnc.EncHextile,
		//vnc.EncZlib,
		//vnc.EncRRE,
	})
	if err != nil {
		fmt.Printf("error setting encodings: %v\n", err)
	}

	// Create a fyne canvas image from our screen image
	display := canvas.NewImageFromImage(screenImage)
	display.FillMode = canvas.ImageFillContain

	// Instantiate the VncDisplay
	viewer := &VncDisplay{Client: client, Display: display}

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
	go PeriodicallyRequestFramebufferUpdate(client, 10)

	// Refresh the display when we receive a framebuffer update
	go ExecuteOnFramebufferUpdate(config, display.Refresh)

	return viewer
}
