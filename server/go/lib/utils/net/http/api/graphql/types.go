package gql

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/types/known/anypb"
	"log"

	"io"
	"strconv"
)

type DummyResolver struct{}

func (r *DummyResolver) Dummy(ctx context.Context) (*bool, error) { return nil, nil }

func MarshalBytes(b []byte) Marshaler {
	return WriterFunc(func(w io.Writer) {
		_, _ = fmt.Fprintf(w, "%q", string(b))
	})
}

func UnmarshalBytes(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case string:
		return []byte(v), nil
	case *string:
		return []byte(*v), nil
	case []byte:
		return v, nil
	case json.RawMessage:
		return v, nil
	default:
		return nil, fmt.Errorf("%T is not []byte", v)
	}
}

func MarshalAny(any anypb.Any) Marshaler {
	return WriterFunc(func(w io.Writer) {
		d, err := any.UnmarshalNew()
		if err != nil {
			log.Println("unable to unmarshal any: ", err)
			return
		}

		if err := json.NewEncoder(w).Encode(d); err != nil {
			log.Println("unable to encode json: ", err)
		}
	})
}

func UnmarshalAny(v interface{}) (anypb.Any, error) {
	switch v := v.(type) {
	case []byte:
		return anypb.Any{}, nil //TODO add an unmarshal mechanism
	case json.RawMessage:
		return anypb.Any{}, nil
	default:
		return anypb.Any{}, fmt.Errorf("%T is not json.RawMessage", v)
	}
}

func MarshalInt32(any int32) Marshaler {
	return WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(any))))
	})
}

func UnmarshalInt32(v interface{}) (int32, error) {
	switch v := v.(type) {
	case int:
		return int32(v), nil
	case int32:
		return v, nil
	case json.Number:
		i, err := v.Int64()
		return int32(i), err
	default:
		return 0, fmt.Errorf("%T is not int32", v)
	}
}

func MarshalInt64(any int64) Marshaler {
	return WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(any))))
	})
}

func UnmarshalInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case json.Number:
		i, err := v.Int64()
		return i, err
	default:
		return 0, fmt.Errorf("%T is not int32", v)
	}
}

func MarshalUint8(any uint8) Marshaler {
	return WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(any))))
	})
}

func UnmarshalUint8(v interface{}) (uint8, error) {
	switch v := v.(type) {
	case int:
		return uint8(v), nil
	case uint8:
		return v, nil //TODO add an unmarshal mechanism
	case json.Number:
		i, err := v.Int64()
		return uint8(i), err
	default:
		return 0, fmt.Errorf("%T is not uint64", v)
	}
}

func MarshalUint32(any uint32) Marshaler {
	return WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(any))))
	})
}

func UnmarshalUint32(v interface{}) (uint32, error) {
	switch v := v.(type) {
	case int:
		return uint32(v), nil
	case uint32:
		return v, nil
	case json.Number:
		i, err := v.Int64()
		return uint32(i), err
	default:
		return 0, fmt.Errorf("%T is not int32", v)
	}
}

func MarshalUint64(any uint64) Marshaler {
	return WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(any))))
	})
}

func UnmarshalUint64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case int:
		return uint64(v), nil
	case uint64:
		return v, nil //TODO add an unmarshal mechanism
	case json.Number:
		i, err := v.Int64()
		return uint64(i), err
	default:
		return 0, fmt.Errorf("%T is not uint64", v)
	}
}

func MarshalFloat32(any float32) Marshaler {
	return WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Itoa(int(any))))
	})
}

func UnmarshalFloat32(v interface{}) (float32, error) {
	switch v := v.(type) {
	case int:
		return float32(v), nil
	case float32:
		return v, nil
	case json.Number:
		f, err := v.Float64()
		return float32(f), err
	default:
		return 0, fmt.Errorf("%T is not float32", v)
	}
}

func Int64(v interface{}) int64 {
	switch v := v.(type) {
	case int:
		return int64(v)
	case int64:
		return v
	case json.Number:
		f, _ := v.Int64()
		return f
	default:
		return 0
	}
}

type Marshaler interface {
	MarshalGQL(w io.Writer)
}

type WriterFunc func(writer io.Writer)

func (f WriterFunc) MarshalGQL(w io.Writer) {
	f(w)
}

func MarshalHttpResponse_HeaderEntry(val map[string]string) Marshaler {
	return WriterFunc(func(w io.Writer) {
		err := json.NewEncoder(w).Encode(val)
		if err != nil {
			panic(err)
		}
	})
}

func UnmarshalHttpResponse_HeaderEntry(v interface{}) (map[string]string, error) {
	if m, ok := v.(map[string]string); ok {
		return m, nil
	}

	return nil, fmt.Errorf("%T is not a map", v)
}
