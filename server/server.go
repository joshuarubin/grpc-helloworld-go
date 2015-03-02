package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joshuarubin/grpc-helloworld-go/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = 8080
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err = s.Serve(ln); err != nil {
		log.Fatal(err)
	}
}
