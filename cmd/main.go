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

	err := fynevnc.ConnectVncDisplay(addr, conf, w)
	if err != nil {
		panic(err)
	}

	w.SetPadded(false)
	w.CenterOnScreen()
	w.ShowAndRun()

	if err != nil {
		panic(err)
	}
}
