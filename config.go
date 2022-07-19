package main

import vnc "github.com/amitbet/vnc2video"

func CreateVncConfig(password string) *vnc.ClientConfig {
	cchServer := make(chan vnc.ServerMessage)
	cchClient := make(chan vnc.ClientMessage)
	errorCh := make(chan error)
	var securityHandlers []vnc.SecurityHandler
	if password != "" {
		securityHandlers = []vnc.SecurityHandler{&vnc.ClientAuthVNC{Password: []byte(password)}}
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