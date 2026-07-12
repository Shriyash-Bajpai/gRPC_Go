package service

import (
	"errors"
	"fmt"
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
	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("cannot copy laptop data:%v", err)
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

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data:%w", err)
	}
	return other, nil
}
