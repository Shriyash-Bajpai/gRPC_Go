package serializer

import (
	"fmt"
	"os"

	"google.golang.org/protobuf/proto"
)

func WriteProtobufToJSONFile(message proto.Message, filename string) error {

	data, err := ProtoToJSON(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message JSON:%v", err)
	}
	err = os.WriteFile(filename, []byte(data), 0644)
	return nil
}

func WriteProtoToBinaryFile(message proto.Message, filename string) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary")
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("cannot write binary data to file:%v", err)
	}
	return nil
}

func ReadProtoFromBinaryFile(filename string, message proto.Message) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("cannot read binary data from file:%v", err)
	}
	err = proto.Unmarshal(data, message)
	if err != nil {
		return fmt.Errorf("can not unmarshal the binary file to proto message:%v", err)
	}
	return nil
}
