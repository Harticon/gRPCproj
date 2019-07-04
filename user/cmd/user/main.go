package main

import (
	"context"
	"flag"
	"fmt"
	order "github.com/Harticon/gRPCproj/user"
	user "github.com/Harticon/gRPCproj/user/proto"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint1 = flag.String("/signup", "localhost:9090", "gRPC server endpoint")
	grpcServerEndpoint2 = flag.String("/signin", "localhost:9090", "gRPC server endpoint")
)

func runHttpReverse() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := user.RegisterRouteGuideHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint1, opts)
	if err != nil {
		return err
	}
	err = user.RegisterRouteGuideHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint2, opts)
	if err != nil {
		return err
	}

	fmt.Println("Reverse proxy is running")
	return http.ListenAndServe(":8080", mux)

}

func runGrpcServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9090))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	grpcServer := grpc.NewServer()

	viper.SetDefault("connection", "grpc.db")
	viper.SetDefault("secret", "secret")
	viper.SetDefault("hashSecret", "salt&peper")

	db, err := gorm.Open("sqlite3", viper.GetString("connection"))
	if err != nil {
		panic("failed to connect to database	")
	}

	db.AutoMigrate(&user.User{})

	access := order.NewAccess(db)
	service := order.NewRouteGuideServer(access)

	user.RegisterRouteGuideServer(grpcServer, service)
	fmt.Println("server is running")
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func main() {

	flag.Parse()

	//todo ??? maybe error handling
	go runGrpcServer()

	time.Sleep(2 * time.Second)

	defer glog.Flush()
	if err := runHttpReverse(); err != nil {
		glog.Fatal(err)
	}

}
