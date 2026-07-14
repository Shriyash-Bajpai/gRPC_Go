To view the lecture series, visit:  
https://youtube.com/playlist?list=PLy_6D98if3UJd5hxWNfAqKMr15HZqFnqf&si=7-By6G0D98jYmxAL

## Protocol Buffer Files

To finalize the protobuf part: ``make gen``  
To test run the ``file.test.go`` file

```azure
                        Development Time (One-time)
                        ───────────────────────────

 .proto file
      │
      │ protoc + protoc-gen-go
      ▼
Generated Go struct (.pb.go) (also called proto message)
      │
      │ 
      ▼
  Used in Go program
      
 Runtime (Writing)
      
   Go Struct
       │
       │ Serialization
       │ Marshal()
       ▼
Binary (Protocol Buffer encoded bytes)
      │
      │ os.WriteFile()
      ▼
     Disk (raw bytes)


            Runtime (Reading)

     Disk (raw bytes)
      │
      │ os.ReadFile()
      ▼
Binary (Protocol Buffer encoded bytes)
      │
      │ Deserialization
      │ Unmarshal()
      ▼
   Go Struct
      
            Convert to JSON  (here we have to create an intermediate marshaller)
       
Generated Go Struct (.pb.go)
        │
        │ (protoc + protoc-gen-go)
        ▼
  Create protobuf message
person := &pb.Person{...}
        │
        │ protojson.MarshalOptions{...}
        ▼
    MarshalOptions object
(stores JSON formatting options)
        │
        │ marshaler.Marshal(message)
        ▼
    JSON []byte
        │
        │ string(data)
        ▼
    JSON string  

```

### Note
The generated functions in pb.go file work only with **Pointers to instances of structs** and not the struct itself.  
eg. *Person will work  
but Person will not work  
See the func definition in ``.pb.go files for further details``  
  
  
## Remote Procedure Calls (RPC)

**The .proto IDL defines the service and its RPC methods. The generated gRPC code expects your server type to implement methods with those exact names and signatures. At runtime, gRPC uses the incoming service name + method name from the request to dispatch to the corresponding implementation.**

We define a service that expose the functions in a server (proto file) and implement all the functions.  
Those are the exact functions that are served by our gRPC server.  
See ``laptop_service.proto`` for further info.  
  
  
  
### What is a gRPC RPC?
An RPC defines a remote function exposed by a service.
**eg:**
```
service BankService {

    rpc CreateAccount(CreateAccountRequest)
        returns (CreateAccountResponse);

    rpc CheckBalance(CheckBalanceRequest)
        returns (CheckBalanceResponse);

    rpc TransferMoney(TransferMoneyRequest)
        returns (TransferMoneyResponse);
}
```

#### How to start a gRPC server?
Read the file ``./cmd/server/server_main.go`` to read about actual process.  
The general flow is : Decide port -> start a gRPC server and a LapServerService Instance -> Register the server -> start listening of the decided address.  
Basic code to get an idea.
```azure
func startTestLaptopServer(t *testing.T) (*LaptopServer, string) {

	laptopServer := NewLaptopServer(NewInMemoryLaptopStore())   // create new appln. server
	grpcServer := grpc.NewServer()                              // create new grpc server 

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)    // register the server

	listener, err := net.Listen("tcp", ":0")                    // start listening on the given port
	require.NoError(t, err)

	go grpcServer.Serve(listener)
	return laptopServer, listener.Addr().String()
}
```

### Service Layer
Service layer's purpose is: Given a req, decide what should happen.  
```azure
  Client
    │
    ▼
gRPC Runtime
    │
    ▼
LaptopServer.CreateLaptop()
    │
    ├── Validate request
    ├── Generate ID if needed
    ├── Check deadline
    ├── Save to store
    └── Build response
    │
    ▼
gRPC Runtime
    │
    ▼
  Client
```

### Storage
In this project, we are using a InMemLaptopStore (a slice of Laptop objects).  
See ``laptop_store.go`` for further info.


## Client Side
The client sends req to the server  
Sample code to crearte a client:
```azure
    conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
    laptopClient := pb.NewLaptopServiceClient(conn)
    req := &pb.CreateLaptopRequest{
        Laptop: laptop,
    }
    res, err := laptopClient.CreateLaptop(ctx, req)
```
```azure
         Client Program
              │
         grpc.Dial()
              │
          TCP Connection
              │
           HTTP/2
              │
        Generated Client Stub
              │
        CreateLaptop()
              │
        Serialize Request
              │
──────────── Network ────────────
              │
            Server
              │
          CreateLaptop()
              │
        Serialize Response
              │
──────────── Network ────────────
                │
        Client receives Response
```


#### Calling the remote RPC Server from Client
```azure
        Client Program
                │
            grpc.Dial()
                │
           TCP Connection
                │
              HTTP/2
                │
        Generated Client Stub
                │
          CreateLaptop()
                │
        Serialize Request
                │
──────────── Network ────────────
                │
              Server
                │
        CreateLaptop()
                │
        Serialize Response
                │
──────────── Network ────────────
                │
        Client receives Response
```


In this project:
Unary- Create a laptop in store
Server Streaming - Search for laptops with a filter in store



# gRPC Go Quick Reference

## Server Setup

```go
type Server struct {
    pb.UnimplementedLaptopServiceServer
}

grpcServer := grpc.NewServer()
pb.RegisterLaptopServiceServer(grpcServer, &Server{})

listener, _ := net.Listen("tcp", ":50051")
grpcServer.Serve(listener)
```

---

## Client Setup

```go
conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
client := pb.NewLaptopServiceClient(conn)
```

---

# Unary RPC

## Client (Send Request)

```go
res, err := client.Method(ctx, req)
```

## Server (Receive Request & Send Response)

```go
func (s *Server) Method(
    ctx context.Context,
    req *pb.Request,
) (*pb.Response, error) {

    return &pb.Response{}, nil
}
```

---

# Server Streaming RPC

## Client (Receive Stream)

```go
stream, err := client.Method(ctx, req)

for {
    res, err := stream.Recv()
    if err == io.EOF {
        break
    }
}
```

## Server (Receive Unary Request)

```go
func (s *Server) Method(
    req *pb.Request,
    stream pb.Service_MethodServer,
) error {
```

## Server (Send Stream)

```go
for {
    stream.Send(res)
}

return nil
```

---

# Client Streaming RPC

## Client (Open Stream)

```go
stream, err := client.Method(ctx)
```

## Client (Send Stream)

```go
stream.Send(req)
stream.Send(req)
stream.Send(req)
```

## Client (Close & Receive Response)

```go
res, err := stream.CloseAndRecv()
```

## Server (Receive Stream)

```go
func (s *Server) Method(
    stream pb.Service_MethodServer,
) error {

    for {
        req, err := stream.Recv()
        if err == io.EOF {
            break
        }
    }

    return stream.SendAndClose(&pb.Response{})
}
```

---

# Bidirectional Streaming RPC

## Client (Open Stream)

```go
stream, err := client.Method(ctx)
```

## Client (Send)

```go
stream.Send(req)
```

## Client (Receive)

```go
res, err := stream.Recv()
```

> Typically, receiving runs in a separate goroutine while sending continues.

## Client (Close Sending)

```go
stream.CloseSend()
```

## Server (Receive & Send)

```go
func (s *Server) Method(
    stream pb.Service_MethodServer,
) error {

    for {
        req, err := stream.Recv()
        if err == io.EOF {
            break
        }

        stream.Send(res)
    }

    return nil
}
```
