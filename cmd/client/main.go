package main

import (
	"context"
	"github.com/halprin/grpc-streaming-example/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	dialOptions := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	connection, err := grpc.Dial("localhost:8000", dialOptions...)
	if err != nil {
		log.Fatalf("Unable to dial gRPC server: %+v", err)
	}
	defer connection.Close()

	client := pb.NewStreamClient(connection)

	stream, err := client.HelloWorld(context.Background())
	if err != nil {
		log.Fatal("Unable to establish HelloWorld stream")
	}

	//spin-off the listener for received messages from the server
	go streamListener(stream)

	//send the person details
	_ = sendPersonDetails(stream)

	log.Println("Closing down the client")
}

func streamListener(stream pb.Stream_HelloWorldClient) {
	for {
		receivedMessage, err := stream.Recv()
		if err == io.EOF {
			log.Println("Server is done, wrapping up")
			return
		} else if err != nil {
			log.Println("There was an error receiving a message")
			return
		}

		log.Printf("Received message: %s", receivedMessage.GetMessage())
	}
}

func sendPersonDetails(stream pb.Stream_HelloWorldClient) error {
	person := pb.People{
		Name:                 "halprin",
		Location:             "The Internet",
		DistanceWashingtonDc: int64(74),
	}

	log.Printf("Sending a person to the server: %s", person.String())

	err := stream.Send(&person)
	if err != nil {
		log.Println("Failed to send person to server")
		return err
	}

	time.Sleep(10 * time.Second)

	log.Println("Closing down the sending of people")
	err = stream.CloseSend()
	if err != nil {
		log.Println("Failed to finish the stream")
	}

	return nil
}
