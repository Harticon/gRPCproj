package main

import (
	"flag"
	"fmt"
	order "github.com/Harticon/gRPCproj/user"
	user "github.com/Harticon/gRPCproj/user/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	user.RegisterRouteGuideServer(grpcServer, &order.RouteGuideServer{})
	fmt.Println("server is running")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("couldn't start server")
	}
}
