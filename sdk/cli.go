package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ghostnet <command>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	switch cmd {
	case "connect":
		fmt.Println("Connecting to GhostNet...")
	case "status":
		fmt.Println("GhostNet Status: Connected")
	default:
		fmt.Printf("Unknown command: %s\n", cmd)
	}
}
