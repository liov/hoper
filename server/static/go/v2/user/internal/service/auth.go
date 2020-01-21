package service

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/user/internal/config"

	"github.com/liov/hoper/go/v2/utils/http/iris/response"
	"github.com/liov/hoper/go/v2/utils/http/token"
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	var auth string
	cookie, _ := r.Cookie("token")
	value, _ := url.QueryUnescape(cookie.Value)
	if value == "" {
		auth = r.Header.Get("authorization")
	} else {
		auth = value
	}
	errHandle := func(w http.ResponseWriter) {
		authErr := response.ResData{Code: uint32(errorcode.Auth), Message: errorcode.Auth.Error()}
		resp, _ := json.Json.Marshal(&authErr)
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: "del",
			Path:  "/",
			//Domain:  "hoper.xyx",
			Expires:  time.Now().Add(-1),
			MaxAge:   -1,
			Secure:   false,
			HttpOnly: true,
		})
		w.Write(resp)
	}

	if auth == "" {
		errHandle(w)
		return
	}
	claims, err := token.ParseToken(auth, config.Conf.Server.TokenSecret)
	if err != nil {
		errHandle(w)
		return
	}
	user, err := UserHashFromRedis(claims.UserID)
	if err != nil {
		log.Error(err)
		errHandle(w)
		return
	}

	r.WithContext(context.WithValue(r.Context(), "auth", user))
}
