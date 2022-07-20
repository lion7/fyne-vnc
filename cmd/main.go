package main

import (
	"fmt"
	"fyne.io/fyne/v2/app"
	fynevnc "github.com/lion7/fyne-vnc"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 3 {
		cmd := filepath.Base(os.Args[0])
		fmt.Printf("Usage  : %s address password\n", cmd)
		fmt.Printf("Example: %s localhost:5900 secret\n", cmd)
		os.Exit(1)
	}

	addr := os.Args[1]
	pass := os.Args[2]
	conf := fynevnc.CreateVncConfig(pass)

	a := app.New()
	defer a.Quit()

	w := a.NewWindow("VNC")
	defer w.Close()

	go func() {
		if err := <-conf.ErrorCh; err != nil {
			w.Close()
		}
	}()

	v, err := fynevnc.ConnectVncDisplay(addr, conf)
	if err != nil {
		panic(err)
	}
	defer v.Close()

	// Add keyboard handler. Hopefully not needed in later versions of Fyne...
	w.Canvas().SetOnTypedKey(v.TypedKey)

	// Remove default padding to get a seamless viewer.
	w.SetPadded(false)

	// Center the window on screen.
	w.CenterOnScreen()

	// Initially resize the window to fully fit the viewer.
	w.Resize(v.Size())

	// Set the viewer as the content of the window.
	w.SetContent(v)

	// Show the window.
	w.ShowAndRun()

	if err != nil {
		panic(err)
	}
}
