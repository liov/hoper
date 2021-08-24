package middle

import (
	"log"
	"net/http"
)

func Log(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
}
