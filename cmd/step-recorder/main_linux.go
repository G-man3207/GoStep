//go:build !windows
// +build !windows

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("GoStep")
	fmt.Println("----------------")
	fmt.Println("This application is designed for Windows and requires Windows-specific libraries.")
	fmt.Println("The Linux version is provided only for build compatibility and CI/CD purposes.")
	fmt.Println("")
	fmt.Println("To use this application, please build and run it on a Windows system.")
	fmt.Println("")
	fmt.Println("For more information, see the README.md file.")

	os.Exit(0)
}
