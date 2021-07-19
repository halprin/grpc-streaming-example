package main

import (
	"fmt"
	"github.com/halprin/grpc-streaming-example/grpc"
	"io"
	"log"
	"time"
)

type streamServer struct {
	grpc.UnimplementedStreamServer
}

func (receiver *streamServer) HelloWorld(stream grpc.Stream_HelloWorldServer) error {
	for {
		//repeatedly get a message from a client
		receivedPerson, err := stream.Recv()
		if err == io.EOF{
			//the client is done and therefore so are we
			log.Println("Client is done, wrapping up")
			return nil
		} else if err != nil {
			//there was some other actual error, bubble up the error
			log.Println("There was an error receiving a message")
			return err
		}

		//grab the elements from the protobuf
		personName := receivedPerson.GetName()
		personLocation := receivedPerson.GetLocation()
		personDistanceToDc := receivedPerson.GetDistanceWashingtonDc()

		//construct our response
		responseMessage := grpc.HelloMessage{
			Message: generateResponseMessage(personName, personLocation, personDistanceToDc),
		}

		//and send our response
		err = sendResponseToStream(stream, &responseMessage)
		if err != nil {
			//bubble up error
			return err
		}

		go waitAndSendResponse(stream, &responseMessage)
	}
}

func generateResponseMessage(name string, location string, distance int64) string {
	return fmt.Sprintf("Hello World, %s!  You are located in %s which is %d miles from Washington D.C.", name, location, distance)
}

func sendResponseToStream(stream grpc.Stream_HelloWorldServer, message *grpc.HelloMessage) error {
	err := stream.Send(message)
	if err != nil {
		//something happened while trying to send a response
		log.Println("Sending a response failed")
		return err
	}

	return nil
}

func waitAndSendResponse(stream grpc.Stream_HelloWorldServer, message *grpc.HelloMessage) {
	time.Sleep(5 * time.Second)
	_ = sendResponseToStream(stream, message)
}
