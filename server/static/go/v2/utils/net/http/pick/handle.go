package pick

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/liov/hoper/go/v2/utils/net/http/iris/response"
	"github.com/liov/hoper/go/v2/utils/net/http/pick/schema"
)

func commonHandler(w http.ResponseWriter, req *http.Request, handle reflect.Value, ps *Params) {
	handleTyp := handle.Type()
	handleNumIn := handleTyp.NumIn()
	if handleNumIn != 0 {
		params := make([]reflect.Value, handleNumIn)
		for i := 0; i < handleNumIn; i++ {
			params[i] = reflect.New(handleTyp.In(i).Elem())
			if handleTyp.In(i).Implements(contextType) {
				sess := params[i].Interface().(Claims)
				/*				cookie,err := req.Cookie("token")
								if err != nil {

								}
								value, _ := url.QueryUnescape(cookie.Value)
								if len(token) == 0 {
									token = req.Header.Get("Authorization")
								}*/

				sess.ParseToken(req)

			} else {
				if ps != nil || req.URL.RawQuery != "" {
					src := req.URL.Query()
					pathParam := *ps
					if len(pathParam) > 0 {
						for i := range pathParam {
							src.Set(pathParam[i].Key, pathParam[i].Value)
						}
					}
					decoder.PickDecode(params[i], src)
				}
				if req.Method != http.MethodGet {
					json.NewDecoder(req.Body).Decode(params[i].Interface())
				}
			}
		}
		result := handle.Call(params)
		if !result[1].IsNil() {
			json.NewEncoder(w).Encode(result[1].Interface())
			return
		}
		if info, ok := result[0].Interface().(*response.File); ok {
			header := w.Header()
			header.Set("Content-Type", "application/octet-stream")
			header.Set("Content-Disposition", "attachment;filename="+info.Name)
			io.Copy(w, info.File)
			if flusher, canFlush := w.(http.Flusher); canFlush {
				flusher.Flush()
			}
			info.File.Close()
			return
		}
		json.NewEncoder(w).Encode(response.ResData{
			Code:    0,
			Message: "success",
			Details: result[0].Interface(),
		})
	}
}

var decoder *schema.Decoder

func init() {
	decoder = schema.NewDecoder()
	decoder.SetAliasTag("json")
}
