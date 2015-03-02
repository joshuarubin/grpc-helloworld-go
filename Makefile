all: server/server client/client

.proto:	proto/helloworld.proto
	@mkdir -p pb
	protoc -I proto --go_out=plugins=grpc:pb proto/*.proto
	touch .proto

server/server: .proto server/server.go
	go build -v -o server/server ./server

client/client: .proto client/client.go
	go build -v -o client/client ./client

clean:
	rm -f server/server client/client
