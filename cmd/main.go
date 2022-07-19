package main

import (
	"fmt"
	fynevnc "github.com/lion7/fyne-vnc"
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

	err := fynevnc.OpenVncViewer(addr, fynevnc.CreateVncConfig(pass))
	if err != nil {
		panic(err)
	}
}
