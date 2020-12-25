package middle

import (
	"net/http"

	"github.com/liov/hoper/go/v2/utils/log"
)

func Log(w http.ResponseWriter,r *http.Request) {
	log.Debug(r.RequestURI)
}
