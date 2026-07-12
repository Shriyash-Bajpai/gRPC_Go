package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/Shriyash-Bajpai/gRPC_Go/pb"
	"github.com/Shriyash-Bajpai/gRPC_Go/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This is the Client appln
// It connects to the server, create a client stub, etc etc
func main() {

	// receive the address of the server from the command line using flag
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	// Create a connection to the server
	// It does nto send the request
	// conn is the gRPC Client Connection
	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server:", err)
	}

	// Create the Client Stub
	// Stub is a local object that represents the remote service.
	laptopClient := pb.NewLaptopServiceClient(conn)

	// Create a new sample request to send in the request
	laptop := sample.NewLaptop()
	// Test cases
	// laptop.Id=""												// no id sent
	// laptop.Id = "280ae34a-85c5-4a90-ba80-8941eb9fc519"       // laptop already exists
	// laptop.Id = "invalid" // laptop id invalid

	// Make a new request with the above created sample laptop
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	// Context carries deadline, cancellation, request metadata etc
	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Calls the RPC
	res, err := laptopClient.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Printf("laptop already exists")
		} else {
			log.Fatal("connot create laptop:", err)
		}
		return
	}
	log.Printf("create laptop with id:%s", res.Id)

}
