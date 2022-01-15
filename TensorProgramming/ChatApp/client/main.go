package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Guilospanck/gRPC/TensorProgramming/ChatApp/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	srvAddr = flag.String("srvAddr", "localhost:8080", "Server address")
	port    = flag.Int("port", 4040, "Client port")
	name    = flag.String("name", "Anonymous", "Name of the client")

	client proto.BroadcastClient
	wait   *sync.WaitGroup
)

func init() {
	wait = &sync.WaitGroup{}
}

func connect(user *proto.User) error {
	var streamError error

	// gets stream
	stream, err := client.CreateStream(context.Background(), &proto.Connect{
		User:   user,
		Active: true,
	})
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}

	// spawns a new goroutine in order to receive stream of messages
	wait.Add(1)
	go func(str proto.Broadcast_CreateStreamClient) {
		defer wait.Done()

		for {
			message, err := str.Recv()
			if err != nil {
				streamError = fmt.Errorf("error reading message: %v", err)
				break
			}

			fmt.Printf("%v [%s]: %s\n", message.Id, message.Timestamp, message.Content)
		}

	}(stream)

	return streamError
}

func main() {
	flag.Parse()

	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithTimeout(10*time.Second))

	// dial to server
	conn, err := grpc.Dial(*srvAddr, opts...)
	if err != nil {
		log.Fatalf("Error trying to dial to server: %v", err)
	}
	defer conn.Close()

	// client stub
	client = proto.NewBroadcastClient(conn)

	timestamp := time.Now()
	done := make(chan int)

	id := sha256.Sum256([]byte(*name + timestamp.String()))

	user := &proto.User{
		Id:   hex.EncodeToString(id[:]),
		Name: *name,
	}

	connect(user)

	wait.Add(1)
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &proto.Message{
				Id:        user.Id,
				Timestamp: timestamp.String(),
				Content:   scanner.Text(),
			}

			_, err := client.BroadcastMessage(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error while trying to broadcast msg: %v", err)
				break
			}
		}

	}()

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
}
