package main

import (
	"context"
	"fmt"
	"grpc-app/proto"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type appServer struct {
	proto.UnimplementedAppServiceServer
}

func (s *appServer) Add(ctx context.Context, req *proto.AddRequest) (*proto.AddResponse, error) {
	time.Sleep(2 * time.Second)
	x := req.GetX()
	y := req.GetY()
	result := x + y
	res := &proto.AddResponse{
		Result: result,
	}
	return res, nil
}

func (s *appServer) GeneratePrimes(req *proto.PrimeRequest, stream proto.AppService_GeneratePrimesServer) error {
	start := req.GetStart()
	end := req.GetEnd()
	fmt.Printf("Generating primes from %d to %d\n", start, end)
	for no := start; no <= end; no++ {
		if isPrime(no) {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Sending prime no : %d\n", no)
			res := &proto.PrimeResponse{
				PrimeNo: no,
			}
			stream.Send(res)
		}
	}
	return nil
}

func isPrime(no int32) bool {
	for i := int32(2); i < no; i++ {
		if no%i == 0 {
			return false
		}
	}
	return true
}

func (s *appServer) CalculateAverage(stream proto.AppService_CalculateAverageServer) error {
	var sum, count int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			avg := sum / count
			res := &proto.AverageResponse{
				Result: avg,
			}
			stream.SendAndClose(res)
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		sum += req.GetNo()
		count++
	}
	return nil
}

func (s *appServer) GreetEveryone(stream proto.AppService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		user := req.GetUser()
		firstName := user.GetFirstName()
		lastName := user.GetLastName()
		fmt.Printf("Received req for greeting for %s and %s\n", firstName, lastName)
		res := &proto.GreetResponse{
			Greeting: fmt.Sprintf("Hello %s %s, Have a nice day!", firstName, lastName),
		}
		err = stream.Send(res)
		if err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func main() {

	s := &appServer{}
	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterAppServiceServer(grpcServer, s)
	grpcServer.Serve(listener)

}
