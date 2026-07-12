package service

import (
	"context"
	"net"
	"testing"

	"github.com/Shriyash-Bajpai/gRPC_Go/pb"
	"github.com/Shriyash-Bajpai/gRPC_Go/sample"
	"github.com/Shriyash-Bajpai/gRPC_Go/serializer"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// Integration Testing

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, serverAdress := startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, serverAdress)
	laptop := sample.NewLaptop()
	expectedId := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedId, res.Id)

	// check if the laptop is really saved
	other, err := laptopServer.Store.Find(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	// check if laptop is same as we send
	requireSameLaptop(t, laptop, other)

}

func startTestLaptopServer(t *testing.T) (*LaptopServer, string) {

	laptopServer := NewLaptopServer(NewInMemoryLaptopStore())
	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go grpcServer.Serve(listener)
	return laptopServer, listener.Addr().String()
}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func requireSameLaptop(t *testing.T, laptop1 *pb.Laptop, laptop2 *pb.Laptop) {
	json1, err := serializer.ProtoToJSON(laptop1)
	require.NoError(t, err)
	json2, err := serializer.ProtoToJSON(laptop2)
	require.NoError(t, err)
	require.Equal(t, json1, json2)
}
