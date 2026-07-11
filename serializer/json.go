package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtoToJSON(message proto.Message) (string, error) {

	// encoder
	marshaler := protojson.MarshalOptions{
		Multiline:       true,  // print pretty json
		Indent:          "  ",  // two spaces
		UseProtoNames:   false, // true: user_id , false: userId
		EmitUnpopulated: true,  // true means show default values
		UseEnumNumbers:  false,
	}

	// Json serialization ie. proto mssg to message
	data, err := marshaler.Marshal(message)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
