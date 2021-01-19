package pick

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/liov/hoper/go/v2/utils/net/http"
	"github.com/liov/hoper/go/v2/utils/net/http/request/binding"
)

func commonHandler(w http.ResponseWriter, req *http.Request, handle *reflect.Value, ps *Params) {
	handleTyp := handle.Type()
	handleNumIn := handleTyp.NumIn()
	if handleNumIn != 0 {
		params := make([]reflect.Value, handleNumIn)
		for i := 0; i < handleNumIn; i++ {
			if handleTyp.In(i).Implements(contextType) {
				params[i] = reflect.ValueOf(req.Context())
			} else {
				params[i] = reflect.New(handleTyp.In(i).Elem())
				if ps != nil || req.URL.RawQuery != "" {
					src := req.URL.Query()
					if ps != nil {
						pathParam := *ps
						if len(pathParam) > 0 {
							for i := range pathParam {
								src.Set(pathParam[i].Key, pathParam[i].Value)
							}
						}
					}
					binding.PickDecode(params[i], src)
				}
				if req.Method != http.MethodGet {
					json.NewDecoder(req.Body).Decode(params[i].Interface())
				}
			}
		}
		result := handle.Call(params)
		header := w.Header()
		header.Set("Content-Type", "application/json")
		if !result[1].IsNil() {
			json.NewEncoder(w).Encode(result[1].Interface())
			return
		}

		if info, ok := result[0].Interface().(*httpi.File); ok {
			header.Set("Content-Type", "application/octet-stream")
			header.Set("Content-Disposition", "attachment;filename="+info.Name)
			io.Copy(w, info.File)
			if flusher, canFlush := w.(http.Flusher); canFlush {
				flusher.Flush()
			}
			info.File.Close()
			return
		}
		json.NewEncoder(w).Encode(httpi.ResData{
			Code:    0,
			Message: "success",
			Details: result[0].Interface(),
		})
	}
}