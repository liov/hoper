package oauth

import (
	"encoding/base64"
	"strings"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
)

type (
	// ClientAuthorizedHandler check the client allows to use this authorization grant type
	ClientAuthorizedHandler func(clientID string, grant oauth2.GrantType) (allowed bool, err error)

	// ClientScopeHandler check the client allows to use scope
	ClientScopeHandler func(clientID, scope string) (allowed bool, err error)

	// UserAuthorizationHandler get user id from request authorization
	UserAuthorizationHandler func(token string) (userID string, err error)

	// PasswordAuthorizationHandler get user id from username and password
	PasswordAuthorizationHandler func(username, password, verifyCode string) (userID string, err error)

	// RefreshingScopeHandler check the scope of the refreshing token
	RefreshingScopeHandler func(newScope, oldScope string) (allowed bool, err error)

	// ResponseErrorHandler response error handing
	ResponseErrorHandler func(re *errors.Response)

	// InternalErrorHandler internal error handing
	InternalErrorHandler func(err error) (re *errors.Response)

	// ExtensionFieldsHandler in response to the access token with the extension of the field
	ExtensionFieldsHandler func(ti oauth2.TokenInfo) (fieldsValue map[string]interface{})
)

func ParseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	// Case insensitive prefix match. See Issue 22736.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
