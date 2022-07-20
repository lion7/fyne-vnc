package fynevnc

import (
	"context"
	"fmt"
	vnc "github.com/amitbet/vnc2video"
	"github.com/amitbet/vnc2video/encoders"
	"net"
	"time"
)

func connectVnc(addr string, config *vnc.ClientConfig) (*vnc.ClientConn, error) {
	// Establish TCP connection to VNC server.
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error connecting to VNC host. %v", err)
	}

	// Set a read/write deadline of 5 seconds.
	conn.SetDeadline(time.Now().Add(5 * time.Second))

	// Attempt to negotiate the VNC connection
	client, err := vnc.Connect(context.Background(), conn, config)
	if err != nil {
		return nil, fmt.Errorf("error negotiating connection to VNC host. %v", err)
	}

	// Remove the deadline for future read/writes
	conn.SetDeadline(time.Time{})

	return client, nil
}

func (v *VncDisplay) RequestFramebufferUpdate() {
	reqMsg := vnc.FramebufferUpdateRequest{Inc: 1, X: 0, Y: 0, Width: v.client.Width(), Height: v.client.Height()}
	if err := reqMsg.Write(v.client); err != nil {
		fmt.Printf("error requesting framebuffer update: %v\n", err)
		v.config.ErrorCh <- err
	}
}

func (v *VncDisplay) PeriodicallyRequestFramebufferUpdate(framerate int) {
	for !v.closed {
		timeStart := time.Now()
		v.RequestFramebufferUpdate()
		timeTarget := timeStart.Add((1000 / time.Duration(framerate)) * time.Millisecond)
		timeLeft := timeTarget.Sub(time.Now())
		if timeLeft > 0 {
			time.Sleep(timeLeft)
		}
	}
}

func (v *VncDisplay) RefreshOnFramebufferUpdate() {
	for !v.closed {
		msg := <-v.config.ServerMessageCh
		if msg.Type() == vnc.FramebufferUpdateMsgType {
			v.display.Refresh()
		}
	}
}

func (v *VncDisplay) LogVncMessages() {
	for !v.closed {
		select {
		case err := <-v.config.ErrorCh:
			fmt.Printf("Received error message: %s\n", err.Error())
		case msg := <-v.config.ClientMessageCh:
			fmt.Printf("Received client message type:%v msg:%v\n", msg.Type(), msg)
		case msg := <-v.config.ServerMessageCh:
			fmt.Printf("Received server message type:%v msg:%v\n", msg.Type(), msg)
		}
	}
}

func (v *VncDisplay) RecordVideo(framerate int) {
	codec := &encoders.MJPegImageEncoder{Quality: 60, Framerate: int32(framerate)}
	go codec.Run("./output")
	for !v.closed {
		timeStart := time.Now()

		codec.Encode(v.client.Canvas)

		timeTarget := timeStart.Add((1000 / time.Duration(framerate)) * time.Millisecond)
		timeLeft := timeTarget.Sub(time.Now())
		if timeLeft > 0 {
			time.Sleep(timeLeft)
		}
	}
	codec.Close()
}
