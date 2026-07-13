package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Shriyash-Bajpai/gRPC_Go/pb"
	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")

// Data Access Logic

// Create the interface defining the storage.
// Interface bcoz that we want to define methods too.
type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

// In mem implementation of LaptopStore
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

// Constructor to create a new instance of InMemLapStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore {

	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

// We need to implement the func defined in the interface above.
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	// deep copy
	// not just storing the reference
	other, err := deepCopy(laptop)
	if err != nil {
		return fmt.Errorf("cannot copy laptop data:%w", err)
	}
	store.data[other.Id] = other
	return nil
}

func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	return deepCopy(laptop)
}

func (store *InMemoryLaptopStore) Search(
	ctx context.Context,
	filter *pb.Filter,
	found func(laptop *pb.Laptop) error,
) error {

	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {

		// Some heavy task being done
		//time.Sleep(time.Second)
		log.Print("checking laptop id:", laptop.Id)

		if isQualified(filter, laptop) {

			if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
				log.Print("context is cancelled")
				return errors.New("context is cancelled")
			}
			// deep copy and not the reference
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {

	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}
	if laptop.GetCpu().GetNumberCores() < filter.GetMinCpuCores() {
		return false
	}
	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}
	if toBit(laptop.GetRam()) < toBit(filter.MinRam) {
		return false
	}
	return true
}

func toBit(memory *pb.Memory) uint64 {
	if memory == nil {
		return 0
	}

	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value

	case pb.Memory_BYTE:
		return value * 8

	case pb.Memory_KILOBYTE:
		return value * 8 * 1024

	case pb.Memory_MEGABYTE:
		return value * 8 * 1024 * 1024

	case pb.Memory_GIGABYTE:
		return value * 8 * 1024 * 1024 * 1024

	case pb.Memory_TERABYTE:
		return value * 8 * 1024 * 1024 * 1024 * 1024

	default:
		return 0
	}
}

func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data:%w", err)
	}
	return other, nil
}
