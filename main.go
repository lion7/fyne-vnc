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

	var channelErr error
	go func() {
		if err := <-config.ErrorCh; err != nil {
			channelErr = err
			w.Close()
		}
	}()

	v, err := ConnectVncDisplay(addr, config)
	if err != nil {
		return err
	}
	defer v.Close()

	w.CenterOnScreen()
	w.Resize(v.Size())
	w.SetContent(v)
	w.ShowAndRun()

	return channelErr
}
