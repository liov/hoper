package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/kataras/iris/v12"
)

type ExportService struct{}

func (*ExportService) request(ctx iris.Context, api string, response interface{}) error {
	/*	req,_:=http.NewRequest(ctx.Method(),api,ctx.Request().Body)
		req.Header["Authorization"] = []string{ctx.GetHeader("Authorization")}
		req.Header["Content-Type"] = []string{"application/json"}*/
	req := ctx.Request() //原请求头可能会干扰返回
	req.URL, _ = url.Parse(api)
	req.RequestURI = ""
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		ctx.WriteString(err.Error())
		return err
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		log.Println(err)
		ctx.Write(buf.Bytes())
		return err
	}
	return nil
}

func (*ExportService) response(ctx iris.Context, filename string, f io.WriterTo) {
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment;filename="+filename)
	f.WriteTo(ctx.ResponseWriter())
}
