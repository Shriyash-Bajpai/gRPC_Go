package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtoToJSON(message proto.Message) (string, error) {

	marshaler := protojson.MarshalOptions{
		Multiline:       true,
		Indent:          "  ",
		UseProtoNames:   false,
		EmitUnpopulated: true,
		UseEnumNumbers:  false,
	}

	data, err := marshaler.Marshal(message)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
