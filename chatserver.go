package chat_application

import (
	"fmt"
	"net"
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

	if chatConfig.logFile != "" {
		fmt.Printf("Log file for messages: %s\n", chatConfig.logFile)
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
	go func() {
		defer wg.Done()
		fmt.Println("goroutine telnet server")
		sock, err := net.Listen("tcp", fmt.Sprintf("%s:%d", chatConfig.Host, chatConfig.Port))
		if err != nil {
			fmt.Printf("Cannot bind to port %d", chatConfig.Port)
			os.Exit(1)
		}
		fmt.Printf("Started telnet server on port %d\n", chatConfig.Port)
		for {
			conn, err := sock.Accept()
			if err != nil {
				// TODO: need to handle error
			}

			go HandleUserConnection(conn)
		}
	}()

	// This goroutine accepts incoming  HTTP connections.
	// This will be the dispatch loop of the REST API.
	go func() {
		defer wg.Done()
		fmt.Printf("Started REST API endpoint on port %d\n", chatConfig.RestAPIPort)
	}()

	wg.Wait()
}
