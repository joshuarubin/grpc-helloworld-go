package main

import (
	"log"
	"os"

	"github.com/joshuarubin/grpc-helloworld-go/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	addr        = "localhost:8080"
	defaultName = "world"
)

func main() {
	conn, err := grpc.Dial(addr)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.Message)
}
