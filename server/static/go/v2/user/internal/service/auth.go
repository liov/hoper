package service

import (
	"context"
	"net/http"

	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/http/iris/response"
	"github.com/liov/hoper/go/v2/utils/http/token"
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	tokens := r.Header["Authorization"]
	errHandle := func(w http.ResponseWriter) {
		authErr := response.ResData{Code: uint32(errorcode.Auth), Message: errorcode.Auth.Error()}
		resp, _ := json.Json.Marshal(&authErr)
		w.Write(resp)
	}

	if len(tokens) == 0 || tokens[0] == "" {
		errHandle(w)
		return
	}
	claims, err := token.ParseToken(tokens[0], config.Conf.Server.TokenSecret)
	if err != nil {
		errHandle(w)
		return
	}
	user, err := UserFromRedis(claims.UserID)
	if err != nil {
		log.Error(err)
		errHandle(w)
		return
	}

	r.WithContext(context.WithValue(r.Context(), "auth", user))
}
