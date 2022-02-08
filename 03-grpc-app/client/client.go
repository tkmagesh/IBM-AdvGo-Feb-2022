package main

import (
	"context"
	"fmt"
	"log"

	"grpc-app/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	clientConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := proto.NewAppServiceClient(clientConn)
	ctx := context.Background()
	addRequest := &proto.AddRequest{
		X: 100,
		Y: 200,
	}
	res, err := client.Add(ctx, addRequest)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(res.GetResult())
}
