package templatei

import (
	"io"
	"strings"
	"text/template"

	"github.com/liov/hoper/v2/utils/log"
)

var CommonTemp = template.New("all")

func Parse(tpl string) *template.Template {
	t, err := CommonTemp.Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}
	t.Funcs(template.FuncMap{"join": strings.Join})
	return t
}

func Execute(wr io.Writer, name string, data interface{}) error {
	return CommonTemp.ExecuteTemplate(wr, name, data)
}
