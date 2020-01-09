package response

import (
	"context"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	model "github.com/liov/hoper/go/v2/protobuf/user"
)

func UserHook(tokenMaxAge int64) func(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
	return func(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
		switch rep := message.(type) {
		case *model.LoginRep:
			if rep.Details != nil {
				http.SetCookie(writer, &http.Cookie{
					Name:  "token",
					Value: rep.Details.Token,
					Path:  "/",
					//Domain:   "hoper.xyz",
					Expires:  time.Now().Add(time.Duration(tokenMaxAge) * time.Second),
					MaxAge:   int(time.Duration(tokenMaxAge) * time.Second),
					Secure:   false,
					HttpOnly: true,
				})
			}
		case *model.LogoutRep:
			http.SetCookie(writer, &http.Cookie{
				Name:  "token",
				Value: "del",
				Path:  "/",
				//Domain:   "hoper.xyz",
				Expires:  time.Now().Add(-1),
				MaxAge:   -1,
				Secure:   false,
				HttpOnly: true,
			})
		}
		return nil
	}
}
