package chat_application

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/rs/xid"
)

type User struct {
	// name of user
	name string

	// socket connection for user
	conn net.Conn

	// channel for writing messages to user
	msgToUser chan []byte

	// user can only be in one room at a time
	room *Room
}

var userMap map[string]*User = make(map[string]*User)
var userMapMtx sync.Mutex

func generateRandomUserName() string {
	return xid.New().String()
}

// NewUser takes in a socket connection, and
// allocates a new User with:
//  - a randomly generated name,
//  - a channel, which will be used to write messages back to the User
func NewUser(conn net.Conn) *User {
	u := User{name: generateRandomUserName(), conn: conn, msgToUser: make(chan []byte)}
	userMapMtx.Lock()
	// TODO: avoid collisions in map, even though we are generating a random name
	userMap[u.name] = &u
	userMapMtx.Unlock()
	return &u
}

// PrintBanner prints a banner message which is sent back to the User
// when they connect to the chat server.
func (u *User) PrintBanner() {
	msg := fmt.Sprintf("Welcome to Chat Server!\nTo see a list of commands type: /help\n\nYou are %s\n\n", u.name)
	u.msgToUser <- []byte(msg)
}

// WriteLoop checks for messages on a channel that are meant for the User,
// and writes them back to the User's socket connection.
func (u *User) WriteLoop() {
	for msg := range u.msgToUser {
		u.conn.Write(msg)
	}
}

func (u *User) WriteTimestampedMessage(msg []byte) {
	newMsg := NewTimestampedMessage(msg).ToBytes()
	newMsg = append(newMsg, "\n"...)
	u.msgToUser <- newMsg
}

// ProcessMessage parses an incoming message sent from a User to the chat server.
// If the message begins with '/', it might be a special command that needs to
// be executed.  If it is not a command, pass the message along to the server,
// where it will be broadcast to the room.
func (u *User) ProcessMessage(msg []byte) *Message {
	if msg[0] == '/' {
		userCommand := strings.Split(strings.TrimSpace(string(msg)), " ")
		if _, ok := ServerCommands[userCommand[0]]; ok {
			ServerCommands[userCommand[0]].commandFunc(u, userCommand)
			return nil
		}
		fmt.Printf("Command not found %v\n", userCommand)
	}
	return NewMessage(u.name, msg)
}

// HandleUserConnection is called when a client connects
// to the chat server.  It allocates a new User, and then
// listens on the socket for data from the user, and then
// sends the data to the room that the user is in.
func HandleUserConnection(conn net.Conn) {
	user := NewUser(conn)

	go user.WriteLoop()
	user.PrintBanner()

	defaultRoom, _ := roomMap["#default"]
	defaultRoom.Join(user)

	bufReader := bufio.NewReader(conn)
	for {
		bytes, err := bufReader.ReadBytes('\n')
		if err != nil {
			continue
		}
		message := user.ProcessMessage(bytes)
		if message != nil && user.room != nil {
			user.room.messages <- message.ToBytes()
		}
	}
}
