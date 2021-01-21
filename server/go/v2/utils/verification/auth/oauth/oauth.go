package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/liov/hoper/go/v2/protobuf/utils/oauth"
	"github.com/liov/hoper/go/v2/protobuf/utils/response"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/server"
)

type Server struct {
	Config                   *server.Config
	Manager                  oauth2.Manager
	ClientAuthorizedHandler  ClientAuthorizedHandler
	ClientScopeHandler       ClientScopeHandler
	UserAuthorizationHandler UserAuthorizationHandler
	RefreshingScopeHandler   RefreshingScopeHandler
	ResponseErrorHandler     ResponseErrorHandler
	InternalErrorHandler     InternalErrorHandler
	ExtensionFieldsHandler   ExtensionFieldsHandler
}

func NewDefaultServer(manager oauth2.Manager) *Server {
	return NewServer(server.NewConfig(), manager)
}

// NewServer create authorization server
func NewServer(cfg *server.Config, manager oauth2.Manager) *Server {
	srv := &Server{
		Config:  cfg,
		Manager: manager,
	}

	return srv
}

func (s *Server) GetRedirectURI(req *oauth.OauthReq, data map[string]interface{}) (uri string, err error) {
	u, err := url.Parse(req.RedirectURI)
	if err != nil {
		return
	}
	if req.LoginURI != "" {
		u = &url.URL{Path: req.LoginURI}
	}
	q := u.Query()
	if req.LoginURI != "" {
		q.Set("client_id", req.ClientID)
		q.Set("access_type", req.AccessType)
		q.Set("redirect_uri", req.RedirectURI)
		q.Set("response_type", req.ResponseType)
		q.Set("scope", req.Scope)
	} else {
		if req.State != "" {
			q.Set("state", req.State)
		}
	}

	for k, v := range data {
		q.Set(k, fmt.Sprint(v))
	}

	switch oauth2.ResponseType(req.ResponseType) {
	case oauth2.Code:
		u.RawQuery = q.Encode()
	case oauth2.Token:
		u.RawQuery = ""
		u.Fragment, err = url.QueryUnescape(q.Encode())
		if err != nil {
			return
		}
	}

	uri = u.String()
	return
}

func (s *Server) CheckResponseType(rt oauth2.ResponseType) bool {
	for _, art := range s.Config.AllowedResponseTypes {
		if art == rt {
			return true
		}
	}
	return false
}

func (s *Server) ValidationAuthorizeRequest(req *oauth.OauthReq) error {
	if req.ClientID == "" || req.RedirectURI == "" {
		return errors.ErrInvalidRequest
	}

	if req.ResponseType == "" {
		return errors.ErrUnsupportedResponseType
	} else if allowed := s.CheckResponseType(oauth2.ResponseType(req.ResponseType)); !allowed {
		return errors.ErrUnauthorizedClient
	}

	return nil
}

func (s *Server) GetAuthorizeToken(req *oauth.OauthReq) (ti oauth2.TokenInfo, err error) {
	// check the client allows the grant type
	if fn := s.ClientAuthorizedHandler; fn != nil {
		gt := oauth2.AuthorizationCode

		if oauth2.ResponseType(req.ResponseType) == oauth2.Token {
			gt = oauth2.Implicit
		}

		allowed, verr := fn(req.ClientID, gt)
		if verr != nil {
			err = verr
			return
		} else if !allowed {
			err = errors.ErrUnauthorizedClient
			return
		}
	}

	// check the client allows the authorized scope
	if fn := s.ClientScopeHandler; fn != nil {

		allowed, verr := fn(req.ClientID, req.Scope)
		if verr != nil {
			err = verr
			return
		} else if !allowed {
			err = errors.ErrInvalidScope
			return
		}
	}

	tgr := &oauth2.TokenGenerateRequest{
		ClientID:       req.ClientID,
		UserID:         req.UserID,
		RedirectURI:    req.RedirectURI,
		Scope:          req.Scope,
		AccessTokenExp: time.Duration(req.AccessTokenExp),
		Request:        nil,
	}

	ti, err = s.Manager.GenerateAuthToken(oauth2.ResponseType(req.ResponseType), tgr)
	return
}

func (s *Server) redirectError(req *oauth.OauthReq, err error) *response.HttpResponse {
	data, _, _ := s.GetErrorData(err)
	return s.redirect(req, data)
}

func (s *Server) redirect(req *oauth.OauthReq, data map[string]interface{}) *response.HttpResponse {
	w := &response.HttpResponse{Header: make(map[string]string)}
	uri, err := s.GetRedirectURI(req, data)
	if err != nil {
		w.Body = []byte(err.Error())
		return w
	}
	if req.LoginURI != "" {
		w.Body = []byte("未登录")
	}
	w.Header["HeaderLocation"] = uri
	w.StatusCode = 302
	return w
}

func (s *Server) HandleAuthorizeRequest(req *oauth.OauthReq, token string) (w *response.HttpResponse) {
	err := s.ValidationAuthorizeRequest(req)
	if err != nil {
		return s.redirectError(req, err)
	}

	// user authorization
	req.UserID, err = s.UserAuthorizationHandler(token)

	if err != nil || req.UserID == "" {
		return s.redirect(req, nil)
	}
	req.LoginURI = ""
	// specify the expiration time of access token

	ti, verr := s.GetAuthorizeToken(req)
	if verr != nil {
		return s.redirectError(req, verr)
	}

	return s.redirect(req, s.GetAuthorizeData(oauth2.ResponseType(req.ResponseType), ti))
}

func (s *Server) GetAuthorizeData(rt oauth2.ResponseType, ti oauth2.TokenInfo) map[string]interface{} {
	if rt == oauth2.Code {
		return map[string]interface{}{
			"code": ti.GetCode(),
		}
	}
	return s.GetTokenData(ti)
}

func (s *Server) GetTokenData(ti oauth2.TokenInfo) map[string]interface{} {
	data := map[string]interface{}{
		"access_token": ti.GetAccess(),
		"token_type":   s.Config.TokenType,
		"expires_in":   int64(ti.GetAccessExpiresIn() / time.Second),
	}

	if scope := ti.GetScope(); scope != "" {
		data["scope"] = scope
	}

	if refresh := ti.GetRefresh(); refresh != "" {
		data["refresh_token"] = refresh
	}

	if fn := s.ExtensionFieldsHandler; fn != nil {
		ext := fn(ti)
		for k, v := range ext {
			if _, ok := data[k]; ok {
				continue
			}
			data[k] = v
		}
	}
	return data
}

func (s *Server) GetErrorData(err error) (map[string]interface{}, int, http.Header) {
	var re errors.Response
	if v, ok := errors.Descriptions[err]; ok {
		re.Error = err
		re.Description = v
		re.StatusCode = errors.StatusCodes[err]
	} else {
		if fn := s.InternalErrorHandler; fn != nil {
			if v := fn(err); v != nil {
				re = *v
			}
		}

		if re.Error == nil {
			re.Error = errors.ErrServerError
			re.Description = errors.Descriptions[errors.ErrServerError]
			re.StatusCode = errors.StatusCodes[errors.ErrServerError]
		}
	}

	if fn := s.ResponseErrorHandler; fn != nil {
		fn(&re)
	}

	data := make(map[string]interface{})
	if err := re.Error; err != nil {
		data["error"] = err.Error()
	}

	if v := re.ErrorCode; v != 0 {
		data["error_code"] = v
	}

	if v := re.Description; v != "" {
		data["error_description"] = v
	}

	if v := re.URI; v != "" {
		data["error_uri"] = v
	}

	statusCode := http.StatusInternalServerError
	if v := re.StatusCode; v > 0 {
		statusCode = v
	}

	return data, statusCode, re.Header
}

func (s *Server) ValidationTokenRequest(r *oauth.OauthReq) (*oauth2.TokenGenerateRequest, error) {

	if r.GrantType == "" {
		return nil, errors.ErrUnsupportedGrantType
	}

	tgr := &oauth2.TokenGenerateRequest{
		ClientID:     r.ClientID,
		ClientSecret: r.ClientSecret,
		Request:      nil,
	}

	switch oauth2.GrantType(r.GrantType) {
	case oauth2.AuthorizationCode:
		tgr.RedirectURI = r.RedirectURI
		tgr.Code = r.Code
		if tgr.RedirectURI == "" ||
			tgr.Code == "" {
			return nil, errors.ErrInvalidRequest
		}
	case oauth2.PasswordCredentials:
		return nil, errors.ErrInvalidGrant
	case oauth2.ClientCredentials:
		tgr.Scope = r.Scope
	case oauth2.Refreshing:
		tgr.Refresh = r.RefreshToken
		tgr.Scope = r.Scope
		if tgr.Refresh == "" {
			return nil, errors.ErrInvalidRequest
		}
	}
	return tgr, nil
}

func (s *Server) HandleTokenRequest(r *oauth.OauthReq) (*response.HttpResponse, error) {
	tgr, err := s.ValidationTokenRequest(r)
	if err != nil {
		return s.tokenError(err)
	}

	ti, err := s.GetAccessToken(oauth2.GrantType(r.GrantType), tgr)
	if err != nil {
		return s.tokenError(err)
	}

	return s.token(s.GetTokenData(ti), nil)
}

func (s *Server) tokenError(err error) (*response.HttpResponse, error) {
	data, statusCode, header := s.GetErrorData(err)
	return s.token(data, header, statusCode)
}
func (s *Server) token(data map[string]interface{}, header http.Header, statusCode ...int) (*response.HttpResponse, error) {
	res := &response.HttpResponse{Header: map[string]string{}}
	res.Header["Content-Type"] = "application/json;charset=UTF-8"
	res.Header["Cache-Control"] = "no-store"
	res.Header["Pragma"] = "no-cache"

	for key := range header {
		res.Header[key] = header.Get(key)
	}

	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}

	res.StatusCode = uint32(status)
	res.Body, _ = json.Marshal(data)
	return res, nil
}

func (s *Server) GetAccessToken(gt oauth2.GrantType, tgr *oauth2.TokenGenerateRequest) (oauth2.TokenInfo, error) {
	if allowed := s.CheckGrantType(gt); !allowed {
		return nil, errors.ErrUnauthorizedClient
	}

	if fn := s.ClientAuthorizedHandler; fn != nil {
		allowed, err := fn(tgr.ClientID, gt)
		if err != nil {
			return nil, err
		} else if !allowed {
			return nil, errors.ErrUnauthorizedClient
		}
	}

	switch gt {
	case oauth2.AuthorizationCode:
		ti, err := s.Manager.GenerateAccessToken(gt, tgr)
		if err != nil {
			switch err {
			case errors.ErrInvalidAuthorizeCode:
				return nil, errors.ErrInvalidGrant
			case errors.ErrInvalidClient:
				return nil, errors.ErrInvalidClient
			default:
				return nil, err
			}
		}
		return ti, nil
	case oauth2.PasswordCredentials, oauth2.ClientCredentials:
		if fn := s.ClientScopeHandler; fn != nil {
			allowed, err := fn(tgr.ClientID, tgr.Scope)
			if err != nil {
				return nil, err
			} else if !allowed {
				return nil, errors.ErrInvalidScope
			}
		}
		return s.Manager.GenerateAccessToken(gt, tgr)
	case oauth2.Refreshing:
		// check scope
		if scope, scopeFn := tgr.Scope, s.RefreshingScopeHandler; scope != "" && scopeFn != nil {
			rti, err := s.Manager.LoadRefreshToken(tgr.Refresh)
			if err != nil {
				if err == errors.ErrInvalidRefreshToken || err == errors.ErrExpiredRefreshToken {
					return nil, errors.ErrInvalidGrant
				}
				return nil, err
			}

			allowed, err := scopeFn(scope, rti.GetScope())
			if err != nil {
				return nil, err
			} else if !allowed {
				return nil, errors.ErrInvalidScope
			}
		}

		ti, err := s.Manager.RefreshAccessToken(tgr)
		if err != nil {
			if err == errors.ErrInvalidRefreshToken || err == errors.ErrExpiredRefreshToken {
				return nil, errors.ErrInvalidGrant
			}
			return nil, err
		}
		return ti, nil
	}

	return nil, errors.ErrUnsupportedGrantType
}

func (s *Server) CheckGrantType(gt oauth2.GrantType) bool {
	for _, agt := range s.Config.AllowedGrantTypes {
		if agt == gt {
			return true
		}
	}
	return false
}
