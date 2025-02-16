package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IslomSobirov/tcp-with-pow/internal/config"
	"github.com/IslomSobirov/tcp-with-pow/internal/pkg/pow"
	"github.com/IslomSobirov/tcp-with-pow/internal/protocol"
	"io"
	"net"
	"time"
)

// RunClient function to start the client
func RunClient(ctx context.Context, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	defer conn.Close()

	// Client makes request every 3 seconds
	for {
		message, err := HandleConnection(ctx, conn, conn)
		if err != nil {
			return err
		}
		fmt.Println("Message Received: ", message)
		time.Sleep(3 * time.Second)
	}

}

func HandleConnection(ctx context.Context, readerConn io.Reader, writerConn io.WriteCloser) (string, error) {
	reader := bufio.NewReader(readerConn)
	err := sendMessage(protocol.Message{
		Header: protocol.RequestChallenge,
	}, writerConn)

	if err != nil {
		return "", err
	}

	messageStr, err := readConnMessage(reader)
	if err != nil {
		return "", err
	}

	message, err := protocol.ParseMessage(messageStr)
	if err != nil {
		return "", fmt.Errorf("failed while parsing the message %w", err)
	}

	var HashCash pow.HashCash

	err = json.Unmarshal([]byte(message.Payload), &HashCash)
	if err != nil {
		return "", fmt.Errorf("failed while unmarshalling the message %w", err)
	}

	// Compute the challenge
	fmt.Println("HashCash data is ready ", HashCash)
	conf := ctx.Value("config").(*config.Config)
	HashCash, err = HashCash.ComputeHashCash(conf.HashcashMaxIterations)

	if err != nil {
		return "", fmt.Errorf("failed while computing the hash cash %w", err)
	}

	fmt.Println("hashcash computed:", HashCash)
	// marshal solution to json
	byteData, err := json.Marshal(HashCash)
	if err != nil {
		return "", fmt.Errorf("err marshal hashcash: %w", err)
	}
	// Sending solution back to the server
	err = sendMessage(protocol.Message{
		Header:  protocol.RequestResource,
		Payload: string(byteData),
	}, writerConn)

	if err != nil {
		return "", fmt.Errorf("failed while sending message %w", err)
	}

	// Get response from the server
	messageStr, err = readConnMessage(reader)
	if err != nil {
		return "", fmt.Errorf("err read msg: %w", err)
	}
	message, err = protocol.ParseMessage(messageStr)
	if err != nil {
		return "", fmt.Errorf("err parse msg: %w", err)
	}
	return message.Payload, nil
}

// readConnMessage read connection message
func readConnMessage(reader *bufio.Reader) (string, error) {
	return reader.ReadString('\n')
}

// sendMessage writing message
func sendMessage(message protocol.Message, conn io.Writer) error {
	messageStr := fmt.Sprintf("%s\n", message.Stringify())
	_, err := conn.Write([]byte(messageStr))
	return err
}
