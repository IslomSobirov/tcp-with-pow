package server

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IslomSobirov/tcp-with-pow/internal/config"
	"github.com/IslomSobirov/tcp-with-pow/internal/pkg/pow"
	"github.com/IslomSobirov/tcp-with-pow/internal/protocol"
	"math/rand"
	"net"
	"strconv"
	"time"
)

var Quotes = [10]string{
	"Give me liberty, or give me death!",
	"Injustice anywhere is a threat to justice everywhere.",
	"I came, I saw, I conquered",
	"The only thing we have to fear is fear itself.",
	"An eye for an eye will only make the whole world blind.",
	"Speak softly and carry a big stick; you will go far.",
	"History will be kind to me, for I intend to write it.",
	"That’s one small step for man, one giant leap for mankind.",
	"Power tends to corrupt, and absolute power corrupts absolutely.",
	"Let them eat cake.",
}

type Clock interface {
	Now() time.Time
}

var ErrQuit = errors.New("client requests to close connection")

type Cache interface {
	// Add - add rand value with expiration (in seconds) to cache
	Add(int, int64) error
	// Get - check existence of int key in cache
	Get(int) (bool, error)
	// Delete - delete key from cache
	Delete(int)
}

func RunServer(ctx context.Context, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	defer listener.Close()
	fmt.Println("listening", listener.Addr())
	for {
		// Listen for an incoming connection.
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accept connection: %w", err)
		}
		// Handle connections in a new goroutine.
		go handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	fmt.Println("new client:", conn.RemoteAddr())
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		req, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("err read connection:", err)
			return
		}
		msg, err := ProcessRequest(ctx, req, conn.RemoteAddr().String())
		if err != nil {
			fmt.Println("err process request:", err)
			return
		}
		if msg != nil {
			err := sendMessage(*msg, conn)
			if err != nil {
				fmt.Println("err send message:", err)
			}
		}
	}
}

// ProcessRequest - process request from client
// returns not-nil pointer to Message if needed to send it back to client
func ProcessRequest(ctx context.Context, msgStr string, clientInfo string) (*protocol.Message, error) {
	msg, err := protocol.ParseMessage(msgStr)
	if err != nil {
		return nil, err
	}
	// switch by header of msg
	switch msg.Header {
	case protocol.Quit:
		return nil, ErrQuit
	case protocol.RequestChallenge:
		fmt.Printf("client %s requests challenge\n", clientInfo)
		// create new challenge for client
		conf := ctx.Value("config").(*config.Config)
		clock := ctx.Value("clock").(Clock)
		cache := ctx.Value("cache").(Cache)
		date := clock.Now()

		// add new created rand value to cache to check it later on RequestResource stage
		// with duration in seconds
		randValue := rand.Intn(100000)
		err := cache.Add(randValue, conf.HashcashDuration)
		if err != nil {
			return nil, fmt.Errorf("err add rand to cache: %w", err)
		}

		hashcash := pow.HashCash{
			Ver:      1,
			Bits:     conf.HashcashZerosCount,
			Date:     date.Unix(),
			Resource: clientInfo,
			Rand:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", randValue))),
			Counter:  0,
		}
		hashcashMarshaled, err := json.Marshal(hashcash)
		if err != nil {
			return nil, fmt.Errorf("err marshal hashcash: %v", err)
		}
		msg := protocol.Message{
			Header:  protocol.ResponseChallenge,
			Payload: string(hashcashMarshaled),
		}
		return &msg, nil
	case protocol.RequestResource:
		fmt.Printf("client %s requests resource with payload %s\n", clientInfo, msg.Payload)
		// parse client's solution
		var hashcash pow.HashCash
		err := json.Unmarshal([]byte(msg.Payload), &hashcash)
		if err != nil {
			return nil, fmt.Errorf("err unmarshal hashcash: %w", err)
		}
		// validate hashcash params
		if hashcash.Resource != clientInfo {
			return nil, fmt.Errorf("invalid hashcash resource")
		}
		conf := ctx.Value("config").(*config.Config)
		clock := ctx.Value("clock").(Clock)
		cache := ctx.Value("cache").(Cache)

		// decoding rand from base64 field in received client's hashcash
		randValueBytes, err := base64.StdEncoding.DecodeString(hashcash.Rand)
		if err != nil {
			return nil, fmt.Errorf("err decode rand: %w", err)
		}
		randValue, err := strconv.Atoi(string(randValueBytes))
		if err != nil {
			return nil, fmt.Errorf("err decode rand: %w", err)
		}

		// if rand exists in cache, it means, that hashcash is valid and really challenged by this server in past
		exists, err := cache.Get(randValue)
		if err != nil {
			return nil, fmt.Errorf("err get rand from cache: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("challenge expired or not sent")
		}

		// sent solution should not be outdated
		if clock.Now().Unix()-hashcash.Date > conf.HashcashDuration {
			return nil, fmt.Errorf("challenge expired")
		}
		//to prevent indefinite computing on server if client sent hashcash with 0 counter
		maxIter := hashcash.Counter
		if maxIter == 0 {
			maxIter = 1
		}
		_, err = hashcash.ComputeHashCash(maxIter)
		if err != nil {
			return nil, fmt.Errorf("invalid hashcash")
		}
		//get random quote
		fmt.Printf("client %s succesfully computed hashcash %s\n", clientInfo, msg.Payload)
		msg := protocol.Message{
			Header:  protocol.ResponseResource,
			Payload: Quotes[rand.Intn(4)],
		}
		// delete rand from cache to prevent duplicated request with same hashcash value
		cache.Delete(randValue)
		return &msg, nil
	default:
		return nil, fmt.Errorf("unknown header")
	}
}

func sendMessage(message protocol.Message, conn net.Conn) error {
	messageStr := fmt.Sprintf("%s\n", message.Stringify())
	_, err := conn.Write([]byte(messageStr))
	return err
}
