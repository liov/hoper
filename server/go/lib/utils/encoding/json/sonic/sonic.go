package sonic

import "github.com/bytedance/sonic"

func Marshal(v interface{}) ([]byte, error) {
	return sonic.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return sonic.Unmarshal(data, v)
}
