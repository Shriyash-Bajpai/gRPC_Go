package main

import (
	"flag" // flag is used to take command line arguement -port+8080
	"fmt"
	"log"
	"net"

	"github.com/Shriyash-Bajpai/gRPC_Go/pb"
	"github.com/Shriyash-Bajpai/gRPC_Go/service"
	"google.golang.org/grpc"
)

// Entry file of our server appl.
func main() {

	// We first read the server port from the command line (using flags)
	port := flag.Int("port", 0, "the server port")
	// Processes the command line args, without it port remains at default value 0
	flag.Parse()
	log.Printf("start server on port:%d", *port)

	// Create a app-specific server that implements RPC
	// The passing of an interface helps by Dependency Injection
	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)
	// Creates gRPC runtime, just an empty gRPC object
	grpcServer := grpc.NewServer()
	// Now we connect the implementation to the gRPC runtime
	// The idea is that the runtime now knows which object should handle incoming CreateLaptop RPCs.
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	// Construct the network address to listen to
	// 0.0.0.0 means listen on all network interfaces
	address := fmt.Sprintf("0.0.0.0:%d", *port)

	// Opens a TCP socket and binds it to the chosen address
	// We separeate these 2 to provide kinds of listeners without changing gRPC runtime.
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server")
	}

	// Server becomes live
	// Remember that serve is blocking,
	// never returns during normal operations
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
