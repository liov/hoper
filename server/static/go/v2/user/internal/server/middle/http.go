package middle

import (
	"context"
	"net/http"

	"github.com/liov/hoper/go/v2/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/http/auth"
	"github.com/liov/hoper/go/v2/utils/http/iris/response"
	"github.com/liov/hoper/go/v2/utils/json"
)

func HttpAuth(w http.ResponseWriter, r *http.Request) {
	tokens := r.Header["Authorization"]
	authErr := response.ResData{Code: int(errorcode.Auth), Message: errorcode.Auth.Error()}
	resp, _ := json.Json.Marshal(&authErr)
	if len(tokens) == 0 || tokens[0] == "" {
		w.Write(resp)
		return
	}
	token, err := ParseToken(tokens[0])
	if err != nil {
		w.Write(resp)
		return
	}
	user, err := auth.UserFromRedis(token.UserID)
	if err != nil {
		w.Write(resp)
		return
	}
	r = r.WithContext(context.WithValue(r.Context(), "auth", user))
}
