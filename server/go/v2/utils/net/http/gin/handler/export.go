package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type ExportService struct{}

func (*ExportService) request(ctx *gin.Context, api string, response interface{}) error {
	/*	req,_:=http.NewRequest(ctx.Method(),api,ctx.Request().Body)
		req.Header["HeaderAuthorization"] = []string{ctx.GetHeader("HeaderAuthorization")}
		req.Header["Content-Type"] = []string{"application/json"}*/
	req := ctx.Request //原请求头可能会干扰返回
	req.URL, _ = url.Parse(api)
	req.RequestURI = ""
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		ctx.Writer.WriteString(err.Error())
		return err
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)

	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		log.Println(err)
		ctx.Writer.Write(buf.Bytes())
		return err
	}
	return nil
}

func (*ExportService) response(ctx *gin.Context, filename string, f io.WriterTo) {
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment;filename="+filename)
	f.WriteTo(ctx.Writer)
}
