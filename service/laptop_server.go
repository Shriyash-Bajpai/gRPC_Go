package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Shriyash-Bajpai/gRPC_Go/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This has the actual appln. logic
// This is the service layer i.e. it is given a req and it has to decide how to res

// This struct defines the concrete type that will implement the gRPC interface.
// Dependency Injection: a dependency on the storage layer
type LaptopServer struct {
	Store LaptopStore

	// Returns an "unimplemented" error for every new added RPC
	pb.UnimplementedLaptopServiceServer
}

// Create a new instance of LaptopServer
// Accepts a laptopStore and return Server instance
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{Store: store}
}

// CreateLaptop is a unary RPC to create a new laptop
func (server *LaptopServer) CreateLaptop(

	// Context: carries request-scoped information liek cancellation, deadlines, metadata etc
	ctx context.Context,
	// the req from the client
	req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {

	// Extract the laptop instance from the incoming request
	laptop := req.GetLaptop()
	log.Printf("received a create laptop request with id:%s", laptop.Id)

	// Verify laptopID, if no ID->assign one, if invalidID, return error
	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop Id is not a valid UUID:%v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Cannot generate a new laptop ID:%v", err)
		}
		laptop.Id = id.String()
	}

	// some heavy processing to test timeout
	time.Sleep(6 * time.Second)

	if err := ctx.Err(); err != nil {
		log.Printf("ctx.Err() = %v", err)
		return nil, status.Error(codes.DeadlineExceeded, err.Error())
	}

	//if ctx.Err() == context.DeadlineExceeded {
	//	log.Print("deadline is exceeded")
	//	return nil, status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	//}

	// save the laptop to in memory storage
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}

		return nil, status.Errorf(code, "cannot save laptop to the store:%v", err)
	}
	log.Printf("saved laptop with id:%s", laptop.Id)
	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}
