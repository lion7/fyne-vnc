package fynevnc

import (
	"fyne.io/fyne/v2/app"
	vnc "github.com/amitbet/vnc2video"
)

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
