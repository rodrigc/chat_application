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
		sock, err := net.Listen("tcp", ":8080")
		if err != nil {
			fmt.Println("Cannot bind to port 8080")
			os.Exit(1)
		}
		for {
			conn, err := sock.Accept()
			if err != nil {
				// handle error
			}

			go HandleUserConnection(conn)
		}
	}()

	// This goroutine accepts incoming  HTTP connections.
	// This will be the dispatch loop of the REST API.
	go func() {
		defer wg.Done()
		fmt.Println("goroutine REST API")
	}()

	wg.Wait()
}
