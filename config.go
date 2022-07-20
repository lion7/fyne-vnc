package fynevnc

import vnc "github.com/amitbet/vnc2video"

func CreateVncConfig(password string) *vnc.ClientConfig {
	cchServer := make(chan vnc.ServerMessage, 1)
	cchClient := make(chan vnc.ClientMessage, 1)
	errorCh := make(chan error, 1)
	var securityHandlers []vnc.SecurityHandler
	if password != "" {
		securityHandlers = []vnc.SecurityHandler{&vnc.ClientAuthVNC{Password: []byte(password)}, &vnc.ClientAuthNone{}}
	} else {
		securityHandlers = []vnc.SecurityHandler{&vnc.ClientAuthNone{}}
	}

	return &vnc.ClientConfig{
		SecurityHandlers: securityHandlers,
		DrawCursor:       true,
		PixelFormat:      vnc.PixelFormat16bit,
		ClientMessageCh:  cchClient,
		ServerMessageCh:  cchServer,
		ErrorCh:          errorCh,
		Messages:         vnc.DefaultServerMessages,
		Encodings: []vnc.Encoding{
			&vnc.TightEncoding{},
			&vnc.HextileEncoding{},
			&vnc.ZRLEEncoding{},
			&vnc.CopyRectEncoding{},
			&vnc.CursorPseudoEncoding{},
			&vnc.CursorPosPseudoEncoding{},
			&vnc.ZLibEncoding{},
			&vnc.RREEncoding{},
			&vnc.RawEncoding{},
		},
	}
}
