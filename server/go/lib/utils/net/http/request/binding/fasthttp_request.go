// Copyright 2019 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"errors"
	"mime/multipart"
	"net/http"
	"reflect"

	"github.com/valyala/fasthttp"
)

type multipartFasthttpRequest fasthttp.Request

var _ setter = (*multipartFasthttpRequest)(nil)

// TrySet tries to set a value by the multipart request with the binding a form file
func (r *multipartFasthttpRequest) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (isSetted bool, err error) {
	req := (*fasthttp.Request)(r)
	form, err := req.MultipartForm()
	if err != nil {
		return false, err
	}
	if files := form.File[key]; len(files) != 0 {
		return setByMultipartFormFile(value, field, files)
	}

	return setByKV(value, field, formSource(form.Value), key, opt)
}

type multipartRequest http.Request

var _ setter = (*multipartRequest)(nil)

// TrySet tries to set a value by the multipart request with the binding a form file
func (r *multipartRequest) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (isSetted bool, err error) {
	if files := r.MultipartForm.File[key]; len(files) != 0 {
		return setByMultipartFormFile(value, field, files)
	}

	return setByKV(value, field, formSource(r.MultipartForm.Value), key, opt)
}

func setByMultipartFormFile(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSetted bool, err error) {
	switch value.Kind() {
	case reflect.Ptr:
		switch value.Interface().(type) {
		case *multipart.FileHeader:
			value.Set(reflect.ValueOf(files[0]))
			return true, nil
		}
	case reflect.Struct:
		switch value.Interface().(type) {
		case multipart.FileHeader:
			value.Set(reflect.ValueOf(*files[0]))
			return true, nil
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(value.Type(), len(files), len(files))
		isSetted, err = setArrayOfMultipartFormFiles(slice, field, files)
		if err != nil || !isSetted {
			return isSetted, err
		}
		value.Set(slice)
		return true, nil
	case reflect.Array:
		return setArrayOfMultipartFormFiles(value, field, files)
	}
	return false, errors.New("unsupported field type for multipart.FileHeader")
}

func setArrayOfMultipartFormFiles(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSetted bool, err error) {
	if value.Len() != len(files) {
		return false, errors.New("unsupported len of array for []*multipart.FileHeader")
	}
	for i := range files {
		setted, err := setByMultipartFormFile(value.Index(i), field, files[i:i+1])
		if err != nil || !setted {
			return setted, err
		}
	}
	return true, nil
}
