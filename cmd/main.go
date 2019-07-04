package main

import (
	"context"
	"fmt"
	user2 "github.com/Harticon/gRPCproj/user/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	if err != nil {
		fmt.Printf("cant connect to server\n")
	}

	defer conn.Close()

	client := user2.NewRouteGuideClient(conn)

	feature, err := client.SignUp(context.Background(), &user2.User{
		Email:    "hroj@seznam.cz",
		Password: "vojta",
	},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("probehlo")
	fmt.Println(feature)

}
