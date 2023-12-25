package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	// Determine the base directory to use
	baseDir := ""
	if len(os.Args) > 1 {
		// Use the first argument as the base directory
		baseDir = os.Args[1]
	} else {
		// Fallback to the current working directory
		var err error
		baseDir, err = os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}
	}
	fmt.Println("Base directory:", baseDir)

	// Construct the path to the other executable
	executablePath := filepath.Join(baseDir, "resources", "go-web-server2")

	// Start the other executable
	cmd := exec.Command(executablePath)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start %v: %v", executablePath, err)
	}

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	// Notify sigChan on SIGINT or SIGTERM
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Set up HTTP server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	go func() {
		fmt.Println("Starting server at port 8090")
		if err := http.ListenAndServe(":8090", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Block until a signal is received
	<-sigChan
	fmt.Println("Shutting down server...")

	// Attempt to gracefully stop the other process
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		fmt.Println("Failed to send SIGTERM to go-web-server2:", err)
		cmd.Process.Kill()
	}

	// Optional: Wait for the other process to finish
	cmd.Wait()
}
