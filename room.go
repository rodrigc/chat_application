package chat_application

import (
	"fmt"
	"strings"
	"sync"
)

var roomMap map[string]*Room = make(map[string]*Room)
var roomMtx sync.Mutex

type Room struct {
	name     string
	messages chan []byte
	users    map[*User]bool
	usersMtx sync.Mutex
}

func (r *Room) WriteTimestampedMessage(msg []byte) {
	newMsg := NewTimestampedMessage(msg).ToBytes()
	newMsg = append(newMsg, "\n"...)
	r.messages <- newMsg
}

// NewRoom will create a new room, and enter
// it into the roomMap
func NewRoom(name string) (*Room, error) {
	// Name must begin with '#'
	if !strings.HasPrefix(name, "#") {
		name = "#" + name
	}
	roomMtx.Lock()
	defer roomMtx.Unlock()
	if _, present := roomMap[name]; present {
		return nil, fmt.Errorf("Room %s exists", name)
	}

	roomMap[name] = &Room{name: name, messages: make(chan []byte), users: map[*User]bool{}}

	return roomMap[name], nil
}

// Join will enter a User into the Room
func (r *Room) Join(u *User) error {
	var present bool
	r.usersMtx.Lock()
	if _, present := r.users[u]; !present {
		oldRoom := u.room
		if oldRoom != nil && oldRoom != r {
			r.usersMtx.Unlock()
			oldRoom.Leave(u)
			r.usersMtx.Lock()
		}
		r.users[u] = true
		u.room = r
	}
	r.usersMtx.Unlock()

	if !present {
		r.WriteTimestampedMessage([]byte(fmt.Sprintf("%s has entered room %s\n", u.name, r.name)))
	}
	return nil
}

// Leave will exit a User from the Room
func (r *Room) Leave(u *User) error {
	r.usersMtx.Lock()
	_, present := r.users[u]
	if present {
		delete(r.users, u)
		u.room = nil
	}
	r.usersMtx.Unlock()
	if present {
		// Send message to all other users in room that user has left
		msg := []byte(fmt.Sprintf("%s has left %s\n", u.name, r.name))
		r.WriteTimestampedMessage(msg)

		// Need to also send message directly to user, since they have already left
		// the room.
		u.WriteTimestampedMessage(msg)
	}
	return nil
}

// Loop will take all messages sent to the channel for this Room,
// and for each User in the room, it will write the message to that
// User's channel.
func (r *Room) Loop() {
	for msg := range r.messages {
		r.usersMtx.Lock()
		for user, _ := range r.users {
			user.msgToUser <- msg
		}
		r.usersMtx.Unlock()
	}
}
