package middle

import (
	"net/http"

	"github.com/hopeio/zeta/utils/log"
)

func Log(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.RequestURI)
}
