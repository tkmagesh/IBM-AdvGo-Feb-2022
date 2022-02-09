package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"grpc-app/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {

	clientConn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := proto.NewAppServiceClient(clientConn)
	ctx := context.Background()

	//doRequestResponse(ctx, client)
	//doServerStreaming(ctx, client)
	//doClientStreaming(ctx, client)
	//doBidirectionalStreaming(ctx, client)
	doRequestResponseWithTimeout(ctx, client)
}

func doRequestResponse(ctx context.Context, client proto.AppServiceClient) {
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

func doServerStreaming(ctx context.Context, client proto.AppServiceClient) {
	req := &proto.PrimeRequest{
		Start: 5,
		End:   100,
	}
	stream, err := client.GeneratePrimes(ctx, req)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("All prime numbers are received")
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Received Prime No : %d\n", res.GetPrimeNo())
	}

}

func doClientStreaming(ctx context.Context, client proto.AppServiceClient) {
	var nos []int32 = []int32{3, 1, 4, 2, 5, 8, 6, 7, 9}
	stream, err := client.CalculateAverage(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	for _, no := range nos {
		time.Sleep(500 * time.Millisecond)
		req := &proto.AverageRequest{
			No: no,
		}
		fmt.Printf("Sending %d\n", no)
		stream.Send(req)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	result := res.GetResult()
	fmt.Printf("Average = %d\n", result)
}

func doBidirectionalStreaming(ctx context.Context, client proto.AppServiceClient) {
	stream, err := client.GreetEveryone(ctx)
	if err != nil {
		log.Fatalf("failed to greet everyone: %v", err)
	}
	done := make(chan bool)
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				done <- true
				break
			}
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("Greeting: %v\n", res.GetGreeting())
		}
	}()

	users := []proto.UserName{
		proto.UserName{
			FirstName: "Magesh",
			LastName:  "Kuppan",
		},
		proto.UserName{
			FirstName: "Suresh",
			LastName:  "Kannan",
		},
		proto.UserName{
			FirstName: "Ramesh",
			LastName:  "Jayaraman",
		},
		proto.UserName{
			FirstName: "Rajesh",
			LastName:  "Pandit",
		},
		proto.UserName{
			FirstName: "Ganesh",
			LastName:  "Kumar",
		},
	}

	for _, user := range users {
		log.Printf("Sending user: %v\n", user)
		time.Sleep(1 * time.Second)
		req := &proto.GreetRequest{
			User: &user,
		}
		if err := stream.Send(req); err != nil {
			log.Fatalln(err)
		}
	}

	<-done
}

func doRequestResponseWithTimeout(ctx context.Context, client proto.AppServiceClient) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()
	addRequest := &proto.AddRequest{
		X: 100,
		Y: 200,
	}
	res, err := client.Add(timeoutCtx, addRequest)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.DeadlineExceeded {
			fmt.Println("Timeout occurred")
		} else {
			log.Fatalln(err)
		}
		log.Fatalln(err)
	}
	fmt.Println(res.GetResult())
}
