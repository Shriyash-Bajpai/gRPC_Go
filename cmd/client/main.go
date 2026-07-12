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

func main() {

	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server:", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)
	laptop := sample.NewLaptop()
	// Test cases
	// laptop.Id=""												// no id sent
	// laptop.Id = "280ae34a-85c5-4a90-ba80-8941eb9fc519"       // laptop already exists
	// laptop.Id = "invalid" // laptop id invalid

	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

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
