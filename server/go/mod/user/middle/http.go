package middle

import (
	"net/http"

	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/gofiber/fiber/v2"
)

func Log(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.RequestURI)
}

func FiberLog(ctx *fiber.Ctx) error {
	log.Debug(ctx.BaseURL())
	return ctx.Next()
}
