package main

import (
	"github.com/halprin/grpc-streaming-example/pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal("Failed to listen")
	}

	grpcServer := grpc.NewServer()
	streamServerImplementation := StreamServerImplementation{}
	pb.RegisterStreamServer(grpcServer, streamServerImplementation)

	log.Println("Server starting up on port 8000")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Serving the gRPC server on the port failed")
	}
}
