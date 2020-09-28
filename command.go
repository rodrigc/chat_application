package chat_application

import "fmt"

type CommandFunc func(u *User, c []string) error

type Command struct {
	name        string
	help        string
	commandFunc CommandFunc
}

var ServerCommands map[string]Command

func displayCommandHelp(u *User, c []string) error {
	commandHelp := "\n"

	for _, command := range ServerCommands {
		commandHelp += command.help + "\n"
	}
	commandHelp += "\n"

	u.msgToUser <- []byte(commandHelp)

	return nil
}

func displayRooms(u *User, c []string) error {
	roomMtx.Lock()
	defer roomMtx.Unlock()

	roomList := "\nROOMS available:\n"

	for room, _ := range roomMap {
		roomList = fmt.Sprintf("%s\n%s\n", roomList, room)
	}
	roomList += "\n"

	u.msgToUser <- []byte(roomList)
	return nil
}

func joinRoom(u *User, c []string) error {
	if len(c) < 2 {
		u.msgToUser <- []byte("Room not specified")
		return nil
	}

	var err error
	roomMtx.Lock()
	room, ok := roomMap[c[1]]
	roomMtx.Unlock()
	if !ok {
		// Room does not exist, create it before joining
		room, err = NewRoom(c[1])

		if err != nil {
			return err
		}
		go room.Loop()
	}
	room.Join(u)
	return nil
}

func leaveRoom(u *User, c []string) error {
	if u.room != nil {
		u.room.Leave(u)
	}
	return nil
}

func newNick(u *User, c []string) error {
	if len(c) < 2 {
		u.msgToUser <- []byte("New nickname not specified")
		return nil
	}

	msg := []byte(fmt.Sprintf("%s has been renamed to %s", u.name, c[1]))
	u.name = c[1]
	if u.room != nil {
		u.room.WriteTimestampedMessage(msg)
		return nil
	}
	// User is not in a room, so send message directly to user
	u.WriteTimestampedMessage(msg)
	return nil
}

func usersInRoom(u *User, c []string) error {
	if u.room == nil {
		return nil
	}

	u.room.usersMtx.Lock()
	defer u.room.usersMtx.Unlock()

	msg := fmt.Sprintf("\nUsers in room: %s\n\n", u.room.name)
	for user, _ := range u.room.users {
		msg = fmt.Sprintf("%s\n%s", msg, user.name)
	}
	msg += "\n\n"
	u.msgToUser <- []byte(msg)
	return nil
}

func RegisterCommands() {
	ServerCommands = map[string]Command{
		"/help":  Command{name: "/help", help: "/help , display help output of all commands", commandFunc: displayCommandHelp},
		"/join":  Command{name: "/join", help: "/join #room1 , enter #room1 , create it if it does not exist", commandFunc: joinRoom},
		"/nick":  Command{name: "/nick", help: "/nick newName , change nickname to newName", commandFunc: newNick},
		"/leave": Command{name: "/leave", help: "/leave , leave the current room", commandFunc: leaveRoom},
		"/rooms": Command{name: "/rooms", help: "/rooms , display all available rooms", commandFunc: displayRooms},
		"/who":   Command{name: "/who", help: "/who , display the users in the current room", commandFunc: usersInRoom},
	}
}
