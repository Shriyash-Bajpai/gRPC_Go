package serializer

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
)

// Entire journey of file
// .proto -> go struct / proto message (by protoc) -> Binary (by serialization) -> file (writing) -> Binary (reading) -> Go Struct(Unmarshal)

// accept a proto message and convert it into JSON String
// see that we make it generic by using proto.Message and not specifying any particular message type eg. Laptop,Screen etc
func WriteProtobufToJSONFile(message proto.Message, filename string) error {

	data, err := ProtoToJSON(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message JSON:%v", err)
	}
	// Write the bytes received as data into a slice of bytes and store it into the file "filename"
	// 0644 is the octal Unix file permission
	err = os.WriteFile(filename, []byte(data), 0644)
	return nil
}

// serializes a protobuf message into protobuf wire format
// i.e. convert proto message into binary
func WriteProtoToBinaryFile(message proto.Message, filename string) error {

	// serializes the protobuf message into binary
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary")
	}

	// write the serialized (binary byts) into the "filename"
	// note that data is a slice of bytes too
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file:%v", err)
	}
	return nil
}

// reads serialized protobuf bytes and reconstructs the message
// File->Bytes->Go Object
func ReadProtoFromBinaryFile(filename string, message proto.Message) error {

	// read data from filename in the form of bytes
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file:%v", err)
	}

	// deserializes protobuf binary into go message
	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("can not unmarshal the binary file to proto message:%v", err)
	}
	return nil
}
