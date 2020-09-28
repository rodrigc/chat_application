package chat_application

import (
	"fmt"
	"net"
	"os"
	"sync"
)

// This function accepts incoming telnet connections.
// This is the core dispatch loop for the chat server.
func TelnetServer(wg sync.WaitGroup) {
	defer wg.Done()
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
}
