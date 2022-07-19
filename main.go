package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	vnc "github.com/amitbet/vnc2video"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		cmd := filepath.Base(os.Args[0])
		fmt.Printf("Usage  : %s address password", cmd)
		fmt.Printf("Example: %s localhost:5900 secret", cmd)
		os.Exit(1)
	}

	addr := os.Args[1]
	pass := os.Args[2]

	err := OpenVncViewer(addr, CreateVncConfig(pass))
	if err != nil {
		panic(err)
	}
}

func OpenVncViewer(addr string, config *vnc.ClientConfig) error {
	a := app.New()
	defer a.Quit()

	w := a.NewWindow("VNC")
	defer w.Close()

	w.CenterOnScreen()

	v := NewVncDisplay(addr, config)
	defer v.Client.Close()

	go func() {
		if err := <-config.ErrorCh; err != nil {
			w.Close()
		}
	}()

	w.Resize(v.Size())
	w.SetContent(v)
	w.ShowAndRun()

	return <-config.ErrorCh
}
