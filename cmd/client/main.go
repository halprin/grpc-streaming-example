package main

import (
	"context"
	"fmt"
	"github.com/halprin/grpc-streaming-example/pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
	"time"
)

func main() {
	//connect to the server listening on port 8000
	dialOptions := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	connection, err := grpc.Dial("localhost:8000", dialOptions...)
	if err != nil {
		log.Fatalf("Unable to dial gRPC server: %+v", err)
	}
	defer connection.Close()

	//create a new client using the connection
	client := pb.NewStreamClient(connection)

	//call the HelloWorld endpoint to start the stream
	stream, err := client.HelloWorld(context.Background())
	if err != nil {
		log.Fatal("Unable to establish HelloWorld stream")
	}

	//spin-off the listener for received messages from the server
	go streamListener(stream)

	//send the person details
	_ = continuallyReadPeople(stream)

	log.Println("Closing down the client")
}

func streamListener(stream pb.Stream_HelloWorldClient) {
	//loop forever
	for {
		//wait to receive a message from the server
		receivedMessage, err := stream.Recv()
		if err == io.EOF {
			//the server closed the connection, so just return to finish this separate listener
			log.Println("Server is done, wrapping up")
			return
		} else if err != nil {
			//an actual error happened, return because no need to keep listening
			log.Println("There was an error receiving a message")
			return
		}

		log.Printf("Received message: %s", receivedMessage.GetMessage())
	}
}

func continuallyReadPeople(stream pb.Stream_HelloWorldClient) error {
	for {
		//construct a person given the details provided on the console
		person, err := getNextPerson()
		if err != nil {
			log.Println("Failed to get a person, done sending people")
			break
		}

		log.Printf("Sending a person to the server: %s", person.String())

		//send the person to the server via the stream some times
		err = sendPersonDetails(stream, person)
		if err != nil {
			return err
		}
	}

	//we're done sending people so close down the stream
	log.Println("Closing down the sending of people")
	err := stream.CloseSend()
	if err != nil {
		log.Println("Failed to finish the stream")
	}

	return nil
}

func sendPersonDetails(stream pb.Stream_HelloWorldClient, person *pb.People) error {
	//send the person to the server 5 times and wait 10 seconds in between
	for i := 0; i < 5; i++ {
		log.Println("Sending...")
		err := stream.Send(person)
		if err != nil {
			log.Println("Failed to send person to server")
			return err
		}

		time.Sleep(10 * time.Second)
	}

	return nil
}

func getNextPerson() (*pb.People, error) {
	//read strings from the console
	personName, personLocation, personDistanceString, err := readDetailsFromConsole()
	if err != nil {
		return nil, err
	}

	personDistance, err := strconv.ParseInt(personDistanceString, 10, 64)
	if err != nil {
		return nil, err
	}

	//construct the person from the provided strings/int64
	return &pb.People{
		Name:                 personName,
		Location:             personLocation,
		DistanceWashingtonDc: personDistance,
	}, nil
}

func readDetailsFromConsole() (string, string, string, error) {
	log.Println("Enter a name, location, and distance")

	var name string
	var location string
	var distance string

	_, err := fmt.Scanln(&name, &location, &distance)
	if err != nil {
		log.Println("Unable to read from the console")
		return "", "", "", err
	}

	return name, location, distance, nil
}
