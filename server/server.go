package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/joshuarubin/grpc-helloworld-go/pb"
	"github.com/zvelo/zvelo-services/util"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	rpcPort  = 8080
	httpPort = 8081
)

func hello(name string) *pb.HelloReply {
	return &pb.HelloReply{Message: "Hello " + name}
}

type rpcServer struct{}

func (s *rpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return hello(in.Name), nil
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	util.Render(w, req, http.StatusOK, hello(req.FormValue("name")))
}

func startRPC() {
	addr := fmt.Sprintf(":%d", rpcPort)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &rpcServer{})
	log.Printf("rpc listening at %s\n", addr)
	log.Fatal(s.Serve(ln))
}

func startHTTP() {
	addr := fmt.Sprintf(":%d", httpPort)
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/hello.pb", helloHandler)
	http.HandleFunc("/hello.json", helloHandler)
	log.Printf("http listening at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func main() {
	go startRPC()
	startHTTP()
}
