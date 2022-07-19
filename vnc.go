package main

import (
	"context"
	"fmt"
	vnc "github.com/amitbet/vnc2video"
	"net"
	"time"
)

func connectVnc(addr string, config *vnc.ClientConfig) *vnc.ClientConn {
	// Establish TCP connection to VNC server.
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to VNC host. %v", err))
	}

	client, err := vnc.Connect(context.Background(), conn, config)
	if err != nil {
		panic(fmt.Sprintf("Error negotiating connection to VNC host. %v", err))
	}

	return client
}

func RequestFramebufferUpdate(client *vnc.ClientConn) {
	reqMsg := vnc.FramebufferUpdateRequest{Inc: 1, X: 0, Y: 0, Width: client.Width(), Height: client.Height()}
	if err := reqMsg.Write(client); err != nil {
		fmt.Printf("error requesting framebuffer update: %v\n", err)
	}
}

func PeriodicallyRequestFramebufferUpdate(client *vnc.ClientConn, framerate int) {
	for {
		timeStart := time.Now()
		RequestFramebufferUpdate(client)
		timeTarget := timeStart.Add((1000 / time.Duration(framerate)) * time.Millisecond)
		timeLeft := timeTarget.Sub(time.Now())
		if timeLeft > 0 {
			time.Sleep(timeLeft)
		}
	}
}

func ExecuteOnFramebufferUpdate(config *vnc.ClientConfig, onFramebufferUpdate func()) {
	for {
		msg := <-config.ServerMessageCh
		if msg.Type() == vnc.FramebufferUpdateMsgType {
			onFramebufferUpdate()
		}
	}
}

func LogVncMessages(config *vnc.ClientConfig) {
	for {
		select {
		case err := <-config.ErrorCh:
			fmt.Printf("Received error message: %s\n", err.Error())
		case msg := <-config.ClientMessageCh:
			fmt.Printf("Received client message type:%v msg:%v\n", msg.Type(), msg)
		case msg := <-config.ServerMessageCh:
			fmt.Printf("Received server message type:%v msg:%v\n", msg.Type(), msg)
		}
	}
}
