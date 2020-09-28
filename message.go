package chat_application

import (
	"fmt"
	"time"
)

type Message struct {
	sender    string
	timestamp string
	msg       []byte
}

type TimestampedMessage struct {
	timestamp string
	msg       []byte
}

func NewMessage(sender string, msg []byte) *Message {
	return &Message{sender: sender, timestamp: time.Now().Format(time.RFC3339), msg: msg}
}

func (m *Message) ToBytes() []byte {
	sendMsg := []byte{}
	sendMsg = append(sendMsg, []byte(fmt.Sprintf("[%s] %s: ", m.timestamp, m.sender))...)
	sendMsg = append(sendMsg, m.msg...)
	return sendMsg
}

func NewTimestampedMessage(msg []byte) *TimestampedMessage {
	return &TimestampedMessage{timestamp: time.Now().Format(time.RFC3339), msg: msg}
}

func (m *TimestampedMessage) ToBytes() []byte {
	sendMsg := []byte{}
	sendMsg = append(sendMsg, []byte(fmt.Sprintf("[%s] ", m.timestamp))...)
	sendMsg = append(sendMsg, m.msg...)
	return sendMsg
}
