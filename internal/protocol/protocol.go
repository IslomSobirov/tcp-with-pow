package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	Quit = iota
	RequestChallenge
	ResponseChallenge
	RequestResource
	ResponseResource
)

type Message struct {
	Header  int
	Payload string
}

func (m Message) Stringify() string {
	return fmt.Sprintf("%d|%s", m.Header, m.Payload)
}

func ParseMessage(data string) (*Message, error) {
	cleanData := strings.TrimSpace(data)
	var messageType int
	parts := strings.SplitN(cleanData, "|", 2)
	if len(parts) < 1 || len(parts) > 2 { //we need at least 1 or 2 parts
		return nil, fmt.Errorf("message doesn't match protocol")
	}
	messageType, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid message type: %s", parts[0])
	}

	message := Message{
		Header: messageType,
	}

	if len(parts) == 2 {
		message.Payload = parts[1]
	}

	return &message, nil
}
