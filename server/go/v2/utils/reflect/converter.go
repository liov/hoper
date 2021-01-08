// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package reflecti

import (
	"reflect"
	"strconv"
)

type Converter func(string) reflect.Value

var (
	invalidValue = reflect.Value{}
)

// Default converters for basic types.
var Converters = map[reflect.Kind]Converter{
	reflect.Bool:    convertBool,
	reflect.Float32: convertFloat32,
	reflect.Float64: convertFloat64,
	reflect.Int:     convertInt,
	reflect.Int8:    convertInt8,
	reflect.Int16:   convertInt16,
	reflect.Int32:   convertInt32,
	reflect.Int64:   convertInt64,
	reflect.String:  convertString,
	reflect.Uint:    convertUint,
	reflect.Uint8:   convertUint8,
	reflect.Uint16:  convertUint16,
	reflect.Uint32:  convertUint32,
	reflect.Uint64:  convertUint64,
}

func convertBool(value string) reflect.Value {
	if value == "on" {
		return reflect.ValueOf(true)
	} else if v, err := strconv.ParseBool(value); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertFloat32(value string) reflect.Value {
	if v, err := strconv.ParseFloat(value, 32); err == nil {
		return reflect.ValueOf(float32(v))
	}
	return invalidValue
}

func convertFloat64(value string) reflect.Value {
	if v, err := strconv.ParseFloat(value, 64); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertInt(value string) reflect.Value {
	if v, err := strconv.ParseInt(value, 10, 0); err == nil {
		return reflect.ValueOf(int(v))
	}
	return invalidValue
}

func convertInt8(value string) reflect.Value {
	if v, err := strconv.ParseInt(value, 10, 8); err == nil {
		return reflect.ValueOf(int8(v))
	}
	return invalidValue
}

func convertInt16(value string) reflect.Value {
	if v, err := strconv.ParseInt(value, 10, 16); err == nil {
		return reflect.ValueOf(int16(v))
	}
	return invalidValue
}

func convertInt32(value string) reflect.Value {
	if v, err := strconv.ParseInt(value, 10, 32); err == nil {
		return reflect.ValueOf(int32(v))
	}
	return invalidValue
}

func convertInt64(value string) reflect.Value {
	if v, err := strconv.ParseInt(value, 10, 64); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}

func convertString(value string) reflect.Value {
	return reflect.ValueOf(value)
}

func convertUint(value string) reflect.Value {
	if v, err := strconv.ParseUint(value, 10, 0); err == nil {
		return reflect.ValueOf(uint(v))
	}
	return invalidValue
}

func convertUint8(value string) reflect.Value {
	if v, err := strconv.ParseUint(value, 10, 8); err == nil {
		return reflect.ValueOf(uint8(v))
	}
	return invalidValue
}

func convertUint16(value string) reflect.Value {
	if v, err := strconv.ParseUint(value, 10, 16); err == nil {
		return reflect.ValueOf(uint16(v))
	}
	return invalidValue
}

func convertUint32(value string) reflect.Value {
	if v, err := strconv.ParseUint(value, 10, 32); err == nil {
		return reflect.ValueOf(uint32(v))
	}
	return invalidValue
}

func convertUint64(value string) reflect.Value {
	if v, err := strconv.ParseUint(value, 10, 64); err == nil {
		return reflect.ValueOf(v)
	}
	return invalidValue
}
