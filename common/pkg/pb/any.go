package pbhelp

import (
	"encoding/json"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// https://ravina01997.medium.com/converting-interface-to-any-proto-and-vice-versa-in-golang-27badc3e23f1
func AnyToAnyPB(v any) (*anypb.Any, error) {
	anyValue := &anypb.Any{}
	bytes, _ := json.Marshal(v)
	bytesValue := &wrappers.BytesValue{
		Value: bytes,
	}

	err := anypb.MarshalFrom(anyValue, bytesValue, proto.MarshalOptions{})
	return anyValue, err
}

func AnyPBToAny(anyValue *anypb.Any) (any, error) {
	var value any

	bytesValue := &wrappers.BytesValue{}
	err := anypb.UnmarshalTo(anyValue, bytesValue, proto.UnmarshalOptions{})
	if err != nil {
		return value, err
	}

	uErr := json.Unmarshal(bytesValue.Value, &value)
	if uErr != nil {
		return value, uErr
	}

	return value, nil
}
