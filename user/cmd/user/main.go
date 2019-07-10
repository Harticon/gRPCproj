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
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)

func runHttpReverse() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := user.RegisterRouteGuideHandlerFromEndpoint(ctx, mux, ":9090", opts)
	if err != nil {
		return err
	}

	fmt.Println("Reverse proxy is running")
	return http.ListenAndServe(":80", mux)

}

func runGrpcServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9090))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	grpcServer := grpc.NewServer()

	viper.SetDefault("connection", "host=db-svc user=goo dbname=goo sslmode=disable password=goo port=5432")
	viper.SetDefault("secret", "secret")
	viper.SetDefault("hashSecret", "salt&peper")

	db, err := gorm.Open("postgres", viper.GetString("connection"))
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	db.AutoMigrate(&user.User{})

	//access := order.NewAccess(&gorm.DB{})
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

	fmt.Println("before sleep")

	//for i := 0; i < 10; i++ {
	//	fmt.Println((i + 1) * 5)
	//	time.Sleep(5 * time.Second)
	//}

	time.Sleep(5 * time.Second)

	defer glog.Flush()
	if err := runHttpReverse(); err != nil {
		glog.Fatal(err)
	}

}
