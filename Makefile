compile: compileServer compileClient

compileServer: compileProtobuf
	go build -o server -v ./cmd/server/

compileClient: compileProtobuf
	go build -o client -v ./cmd/client/

compileProtobuf:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./pb/stream_service.proto
