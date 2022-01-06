package oauth

import (
	"errors"
	"html/template"
	"net/http"
	"os"

	"github.com/actliboy/hoper/server/go/lib/utils/verification/auth/oauth/provider"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
)

var store *sessions.CookieStore

const SessionName = "session-name"

func init() {
	store = sessions.NewCookieStore([]byte(os.Getenv(".cookiesession.id")))
}

// These are some function helpers that you may use if you want

// GetProviderName is a function used to get the name of a provider
// for a given request. By default, this provider is fetched from
// the URL query string. If you provide it in a different way,
// assign your own function to this variable that returns the provider
// name for your request.
var GetProviderName = func(ctx *gin.Context) (string, error) {
	// try to get it from the url param "provider"
	if p, ok := ctx.GetQuery("provider"); ok {
		return p, nil
	}

	// try to get it from the url PATH parameter "{provider} or :provider or {provider:string} or {provider:alphabetical}"
	if p := ctx.PostForm("provider"); p != "" {
		return p, nil
	}

	// try to get it from context's per-request storage
	if p, _ := ctx.Get("provider"); p.(string) != "" {
		return p.(string), nil
	}
	// if not found then return an empty string with the corresponding error
	return "", errors.New("you must select a provider")
}

/*
BeginAuthHandler is a convenience handler for starting the authentication process.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".
BeginAuthHandler will redirect the user to the appropriate authentication end-point
for the requested provider.
See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
func BeginAuthHandler(ctx *gin.Context) {
	url, err := GetAuthURL(ctx)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Writer.WriteString(err.Error())
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider" or from the context's value of "provider" key.
I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/
func GetAuthURL(ctx *gin.Context) (string, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return "", err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return "", err
	}
	sess, err := provider.BeginAuth(SetState(ctx))
	if err != nil {
		return "", err
	}

	url, err := sess.GetAuthURL()
	if err != nil {
		return "", err
	}
	session, _ := store.Get(ctx.Request, SessionName)
	session.Values[providerName] = sess.Marshal()
	return url, nil
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var SetState = func(ctx *gin.Context) string {
	if state, ok := ctx.GetQuery("state"); ok {
		return state
	}

	return "state"

}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(ctx *gin.Context) string {
	return ctx.Query("state")
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all of the basic information about the user from the provider.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".
See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
var CompleteUserAuth = func(ctx *gin.Context) (goth.User, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}
	session, _ := store.Get(ctx.Request, SessionName)
	value := session.Values[providerName].(string)
	if value == "" {
		return goth.User{}, errors.New("session value for " + providerName + " not found")
	}

	sess, err := provider.UnmarshalSession(value)
	if err != nil {
		return goth.User{}, err
	}

	user, err := provider.FetchUser(sess)
	if err == nil {
		// user can be found with existing session data
		return user, err
	}

	// get new token and retry fetch
	_, err = sess.Authorize(provider, ctx.Request.URL.Query())
	if err != nil {
		return goth.User{}, err
	}

	session.Values[providerName] = sess.Marshal()
	return provider.FetchUser(sess)
}

// Logout invalidates a user session.
func Logout(ctx *gin.Context) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		ctx.Writer.WriteString(err.Error())
		return
	}
	session, _ := store.Get(ctx.Request, SessionName)
	delete(session.Values, providerName)
	ctx.Header("HeaderLocation", "/")
	ctx.Status(http.StatusTemporaryRedirect)
}

// End of the "some function helpers".

func CallBack(ctx *gin.Context) {
	user, err := CompleteUserAuth(ctx)
	if err != nil {
		BeginAuthHandler(ctx)
		return
	}
	t, _ := template.New("foo").Parse(userTemplate)
	t.Execute(ctx.Writer, user)
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">ErrorLog in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>HeaderLocation: {{.HeaderLocation}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

func Index(res http.ResponseWriter, req *http.Request) {
	t, _ := template.New("foo").Parse(indexTemplate)
	t.Execute(res, provider.NewAuth())
}
