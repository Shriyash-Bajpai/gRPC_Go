package service

import (
	"context"
	"io"
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

	laptopServer, serverAdress := startTestLaptopServer(t, NewInMemoryLaptopStore())
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

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()

	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	store := NewInMemoryLaptopStore()
	expectedIds := make(map[string]bool)

	for i := 0; i < 6; i++ {

		laptop := sample.NewLaptop()
		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &pb.Memory{Value: 4096, Unit: pb.Memory_MEGABYTE}
		case 4:
			laptop.PriceUsd = 1999
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.5
			laptop.Cpu.MaxGhz = 4.5
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedIds[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberCores = 6
			laptop.Cpu.MinGhz = 2.8
			laptop.Cpu.MaxGhz = 5.0
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedIds[laptop.Id] = true
		}
		err := store.Save(laptop)
		require.NoError(t, err)
	}

	_, serverAdress := startTestLaptopServer(t, store)
	laptopClient := newTestLaptopClient(t, serverAdress)

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)

	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		require.Contains(t, expectedIds, res.GetLaptop().GetId())
		found += 1
	}
	require.Equal(t, found, len(expectedIds))
}

func startTestLaptopServer(t *testing.T, store LaptopStore) (*LaptopServer, string) {

	laptopServer := NewLaptopServer(store)
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
