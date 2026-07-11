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