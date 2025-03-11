//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/gustaf/go-test/pkg/config"
	"github.com/gustaf/go-test/pkg/gui"
)

func init() {
	// Lock the OS Thread since this is a GUI application
	runtime.LockOSThread()
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			fmt.Printf("An unexpected error occurred: %v\n", r)
			os.Exit(1)
		}
	}()

	// Create and show the main window
	if err := gui.NewRecorderWindow(config.DefaultSettings()); err != nil {
		log.Printf("Failed to create main window: %v", err)
		fmt.Printf("Failed to create application window: %v\n", err)
		os.Exit(1)
	}
}
