package chat_application

import (
	"fmt"
	"os"
	"sync"
)

func ChatServerStart() {
	// Initialize the commands which a user can send to the server
	// by the server.  See commands.go.
	RegisterCommands()

	err := ParseConfigFile()
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		os.Exit(1)
	}

	if chatConfig.LogFile != "" {
		fmt.Printf("Log file for messages: %s\n", chatConfig.LogFile)
	}

	// When we start up, create a room #default
	// that everyone joins by default.
	defaultRoom, err := NewRoom("#default")

	if err != nil {
		fmt.Printf("Error creating room: %v\n", err)
		os.Exit(1)
	}
	go defaultRoom.Loop()

	var wg sync.WaitGroup
	wg.Add(2)

	// This goroutine accepts incoming telnet connections.
	// This is the core dispatch loop for the chat server.
	go TelnetServer(&wg)

	// This goroutine accepts incoming  HTTP connections.
	// This will be the dispatch loop of the REST API.
	go RestAPIServer(&wg)

	wg.Wait()
}
