# grpc-streaming-example
An example of gRPC streaming messages in Go.

## Prerequisites

You need `protobuf`, `protoc-gen-go`, and `protoc-gen-go-grpc`.

You can install these on the Mav via `brew`:

1. `brew install protobuf`
1. `brew install protoc-gen-go`
1. `brew install protoc-gen-go-grpc`

## Compile

To compile the program, run `make compile`.

## Run

### Server

Run the server first.

```shell
$ ./server
```

The server will listen for a stream of people.  For each received person, it will send an immediate message back and
then send another message five seconds later.

To stop the server, press control + C.

### Client

Then run the client.

```shell
$ ./client
```

The client will query details about a person from the console.  Provide a name, location, and a number with spaces in
between.  E.g...
`halprin Internet 768`

The client will then stream the person to the server five times with 10 seconds in between.  In between these 10
seconds, you will receive both messages back from the server between sends.

After sending the person the fifth time, you will be able to type in details for another person like before.

To stop the client, press control + C or press return without providing any person details when the client is querying
for a person.
