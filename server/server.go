package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/bradfitz/http2"
	"github.com/joshuarubin/grpc-helloworld-go/pb"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/graceful"
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

func helloHandler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	util.Render(w, req, http.StatusOK, hello(req.FormValue("name")))
}

func startRPC(msgCh chan<- string, errCh chan<- error) {
	addr := fmt.Sprintf(":%d", rpcPort)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		errCh <- err
		return
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &rpcServer{})

	msgCh <- fmt.Sprintf("rpc listening at %s", addr)
	errCh <- s.Serve(ln)
}

func startHTTP(msgCh chan<- string, errCh chan<- error) {
	mux := httprouter.New()
	mux.GET("/hello", helloHandler)
	mux.GET("/hello.pb", helloHandler)
	mux.GET("/hello.json", helloHandler)

	addr := fmt.Sprintf(":%d", httpPort)

	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	http2.ConfigureServer(s, nil)

	msgCh <- fmt.Sprintf("http listening at %s", addr)
	errCh <- graceful.ListenAndServe(s, 3*time.Minute)
}

func main() {
	msgCh := make(chan string)
	errCh := make(chan error)

	go startHTTP(msgCh, errCh)
	go startRPC(msgCh, errCh)

	for {
		select {
		case err := <-errCh:
			log.Fatal(err)
		case msg := <-msgCh:
			log.Println(msg)
		}
	}
}
