package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"github.com/Shriyash-Bajpai/gRPC_Go/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This has the actual appln. logic
// This is the service layer i.e. it is given a req and it has to decide how to res

// 1 MB
const maxImageSize = 1 << 20

// This struct defines the concrete type that will implement the gRPC interface.
// Dependency Injection: a dependency on the storage layer
type LaptopServer struct {
	LaptopStore LaptopStore
	ImageStore  ImageStore

	// Returns an "unimplemented" error for every new added RPC
	pb.UnimplementedLaptopServiceServer
}

// Create a new instance of LaptopServer
// Accepts a laptopStore and return Server instance
func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore) *LaptopServer {
	return &LaptopServer{LaptopStore: laptopStore, ImageStore: imageStore}
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
	//time.Sleep(6 * time.Second)
	if err := contextError(ctx); err != nil {
		return nil, err
	}

	// save the laptop to in memory storage
	err := server.LaptopStore.Save(laptop)
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

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is cancelled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
	//if err := ctx.Err(); err != nil {
	//	log.Printf("ctx.Err() = %v", err)
	//	return nil, status.Error(codes.DeadlineExceeded, err.Error())
	//}
}

func (server *LaptopServer) SearchLaptop(
	req *pb.SearchLaptopRequest,
	stream pb.LaptopService_SearchLaptopServer) error {

	filter := req.GetFilter()
	log.Printf("receive a search-laptop request with filter:%v", filter)

	// some heavy processing

	err := server.LaptopStore.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {

			res := &pb.SearchLaptopResponse{Laptop: laptop}
			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("sent laptop with id:%s", laptop.GetId())
			return nil
		},
	)

	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error:%v", err)
	}
	return nil
}

// client streaming
func (server *LaptopServer) UploadImage(stream pb.LaptopService_UploadImageServer) error {

	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "cannot receive image info"))
	}

	laptopID := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	log.Printf("receive an upload image req for laptop %s with image type %s", laptopID, imageType)

	laptop, err := server.LaptopStore.Find(laptopID)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot find laptop:%v", err))
	}

	if laptop == nil {
		return logError(status.Errorf(codes.InvalidArgument, "laptop %s does not exist", laptopID))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {

		if err := contextError(stream.Context()); err != nil {
			return err
		}
		log.Print("waiting for more data to come")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("no more data")
			break
		}

		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot receive chunk data:%v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)
		log.Printf("received a chunk with size:%d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return logError(status.Errorf(codes.InvalidArgument, "image is too large:%v", err))
		}

		// writes very slowly
		//time.Sleep(time.Second)

		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot write chunk data:%v", err))
		}
	}
	imageID, err := server.ImageStore.Save(laptopID, imageType, imageData)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot save image to the store:%s", err))
	}

	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}
	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Errorf(codes.Internal, "cannot send response:%v", err))
	}
	log.Printf("the image is saved with id:%s and size:%v", imageID, imageSize)
	return nil

	return nil
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
