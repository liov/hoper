package middle

import (
	"net/http"

	"github.com/hopeio/lemon/utils/log"
)

func Log(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.RequestURI)
}
