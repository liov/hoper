package gql

import (
	"errors"
	"strconv"
)

type Uint64 uint64

func (Uint64) ImplementsGraphQLType(name string) bool {
	return name == "Uint64"
}

func (i *Uint64) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case uint64:
		*i = Uint64(input)
	case int64:
		*i = Uint64(input)
	case int:
		*i = Uint64(input)
	case uint32:
		*i = Uint64(input)
	case int32:
		*i = Uint64(input)
	case uint:
		*i = Uint64(input)
	default:
		err = errors.New("wrong type")
	}
	return err
}

func (i Uint64) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, int64(i), 10), nil
}

type Uint32 uint64

func (Uint32) ImplementsGraphQLType(name string) bool {
	return name == "Uint32"
}

func (i *Uint32) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case uint64:
		*i = Uint32(input)
	case int64:
		*i = Uint32(input)
	case int:
		*i = Uint32(input)
	case uint32:
		*i = Uint32(input)
	case int32:
		*i = Uint32(input)
	case uint:
		*i = Uint32(input)
	default:
		err = errors.New("wrong type")
	}
	return err
}

func (i Uint32) MarshalJSON() ([]byte, error) {
	return strconv.AppendInt(nil, int64(i), 10), nil
}
