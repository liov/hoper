package oauth

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/liov/hoper/go/v2/utils/net/http/auth/oauth/provider"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/amazon"
	"github.com/markbates/goth/providers/auth0"
	"github.com/markbates/goth/providers/bitbucket"
	"github.com/markbates/goth/providers/dropbox"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/intercom"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/markbates/goth/providers/onedrive"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/paypal"
	"github.com/markbates/goth/providers/soundcloud"
	"github.com/markbates/goth/providers/spotify"
	"github.com/markbates/goth/providers/steam"
	"github.com/markbates/goth/providers/twitch"
	"github.com/markbates/goth/providers/twitter"
	"github.com/markbates/goth/providers/uber"
	"github.com/markbates/goth/providers/yandex"
)

var sessionsManager *sessions.Sessions

func init() {
	sessionsManager = sessions.New(sessions.Config{
		Cookie:  ".cookiesession.id",
		Expires: time.Minute,
	})
}

// These are some function helpers that you may use if you want

// GetProviderName is a function used to get the name of a provider
// for a given request. By default, this provider is fetched from
// the URL query string. If you provide it in a different way,
// assign your own function to this variable that returns the provider
// name for your request.
var GetProviderName = func(ctx iris.Context) (string, error) {
	// try to get it from the url param "provider"
	if p := ctx.URLParam("provider"); p != "" {
		return p, nil
	}

	// try to get it from the url PATH parameter "{provider} or :provider or {provider:string} or {provider:alphabetical}"
	if p := ctx.Params().Get("provider"); p != "" {
		return p, nil
	}

	// try to get it from context's per-request storage
	if p := ctx.Values().GetString("provider"); p != "" {
		return p, nil
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
func BeginAuthHandler(ctx iris.Context) {
	url, err := GetAuthURL(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.Writef("%v", err)
		return
	}

	ctx.Redirect(url, iris.StatusTemporaryRedirect)
}

/*
GetAuthURL starts the authentication process with the requested provided.
It will return a URL that should be used to send users to.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider" or from the context's value of "provider" key.
I would recommend using the BeginAuthHandler instead of doing all of these steps
yourself, but that's entirely up to you.
*/
func GetAuthURL(ctx iris.Context) (string, error) {
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
	session := sessionsManager.Start(ctx)
	session.Set(providerName, sess.Marshal())
	return url, nil
}

// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
var SetState = func(ctx iris.Context) string {
	state := ctx.URLParam("state")
	if len(state) > 0 {
		return state
	}

	return "state"

}

// GetState gets the state returned by the provider during the callback.
// This is used to prevent CSRF attacks, see
// http://tools.ietf.org/html/rfc6749#section-10.12
var GetState = func(ctx iris.Context) string {
	return ctx.URLParam("state")
}

/*
CompleteUserAuth does what it says on the tin. It completes the authentication
process and fetches all of the basic information about the user from the provider.
It expects to be able to get the name of the provider from the query parameters
as either "provider" or ":provider".
See https://github.com/markbates/goth/examples/main.go to see this in action.
*/
var CompleteUserAuth = func(ctx iris.Context) (goth.User, error) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		return goth.User{}, err
	}

	provider, err := goth.GetProvider(providerName)
	if err != nil {
		return goth.User{}, err
	}
	session := sessionsManager.Start(ctx)
	value := session.GetString(providerName)
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
	_, err = sess.Authorize(provider, ctx.Request().URL.Query())
	if err != nil {
		return goth.User{}, err
	}

	session.Set(providerName, sess.Marshal())
	return provider.FetchUser(sess)
}

// Logout invalidates a user session.
func Logout(ctx iris.Context) {
	providerName, err := GetProviderName(ctx)
	if err != nil {
		ctx.WriteString(err.Error())
		return
	}
	session := sessionsManager.Start(ctx)
	session.Delete(providerName)
	ctx.Header("Location", "/")
	ctx.StatusCode(http.StatusTemporaryRedirect)
}

// End of the "some function helpers".

func NewAuth() *ProviderIndex {
	goth.UseProviders(
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost/api/auth/twitter/callback"),
		// If you'd like to use authenticate instead of authorize in Twitter provider, use this instead.
		// twitter.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost/api/auth/twitter/callback"),

		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), "http://localhost/api/auth/facebook/callback"),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost/api/auth/google/callback"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), "http://localhost:8070/auth/github/callback"),
		spotify.New(os.Getenv("SPOTIFY_KEY"), os.Getenv("SPOTIFY_SECRET"), "http://localhost/api/auth/spotify/callback"),
		twitch.New(os.Getenv("TWITCH_KEY"), os.Getenv("TWITCH_SECRET"), "http://localhost/api/auth/twitch/callback"),
		dropbox.New(os.Getenv("DROPBOX_KEY"), os.Getenv("DROPBOX_SECRET"), "http://localhost/api/auth/dropbox/callback"),
		bitbucket.New(os.Getenv("BITBUCKET_KEY"), os.Getenv("BITBUCKET_SECRET"), "http://localhost/api/auth/bitbucket/callback"),
		instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), "http://localhost/api/auth/instagram/callback"),
		intercom.New(os.Getenv("INTERCOM_KEY"), os.Getenv("INTERCOM_SECRET"), "http://localhost/api/auth/intercom/callback"),
		amazon.New(os.Getenv("AMAZON_KEY"), os.Getenv("AMAZON_SECRET"), "http://localhost/api/auth/amazon/callback"),
		onedrive.New(os.Getenv("ONEDRIVE_KEY"), os.Getenv("ONEDRIVE_SECRET"), "http://localhost/api/auth/onedrive/callback"),
		microsoftonline.New(os.Getenv("MICROSOFTONLINE_KEY"), os.Getenv("MICROSOFTONLINE_SECRET"), "http://localhost/api/auth/microsoftonline/callback"),

		paypal.New(os.Getenv("PAYPAL_KEY"), os.Getenv("PAYPAL_SECRET"), "http://localhost/api/auth/paypal/callback"),
		steam.New(os.Getenv("STEAM_KEY"), "http://localhost/api/auth/steam/callback"),
		uber.New(os.Getenv("UBER_KEY"), os.Getenv("UBER_SECRET"), "http://localhost/api/auth/uber/callback"),
		soundcloud.New(os.Getenv("SOUNDCLOUD_KEY"), os.Getenv("SOUNDCLOUD_SECRET"), "http://localhost/api/auth/soundcloud/callback"),
		gitlab.New(os.Getenv("GITLAB_KEY"), os.Getenv("GITLAB_SECRET"), "http://localhost/api/auth/gitlab/callback"),

		//Auth0 allocates domain per customer, a domain must be provided for auth0 to work
		auth0.New(os.Getenv("AUTH0_KEY"), os.Getenv("AUTH0_SECRET"), "http://localhost/api/auth/auth0/callback", os.Getenv("AUTH0_DOMAIN")),
		yandex.New(os.Getenv("YANDEX_KEY"), os.Getenv("YANDEX_SECRET"), "http://localhost/api/auth/yandex/callback"),
		provider.New(os.Getenv("HOPER_KEY"), os.Getenv("HOPER_SECRET"), "http://localhost:8080/auth/hoper/callback"),
	)

	// OpenID Connect is based on OpenID Connect Auto Discovery URL (https://openid.net/specs/openid-connect-discovery-1_0-17.html)
	// because the OpenID Connect provider initialize it self in the New(), it can return an error which should be handled or ignored
	// ignore the error for now
	openidConnect, _ := openidConnect.New(os.Getenv("OPENID_CONNECT_KEY"), os.Getenv("OPENID_CONNECT_SECRET"), "http://localhost/api/auth/openid-connect/callback", os.Getenv("OPENID_CONNECT_DISCOVERY_URL"))
	if openidConnect != nil {
		goth.UseProviders(openidConnect)
	}

	m := make(map[string]string)
	m["amazon"] = "Amazon"
	m["bitbucket"] = "Bitbucket"
	m["dropbox"] = "Dropbox"
	m["facebook"] = "Facebook"
	m["github"] = "Github"
	m["gitlab"] = "Gitlab"
	m["google"] = "Google"
	m["gplus"] = "Google Plus"
	m["steam"] = "Steam"
	m["twitch"] = "Twitch"
	m["uber"] = "Uber"
	m["wepay"] = "Wepay"
	m["instagram"] = "Instagram"
	m["onedrive"] = "Onedrive"
	m["azuread"] = "Azure AD"
	m["microsoftonline"] = "Microsoft Online"
	m["paypal"] = "Paypal"
	m["twitter"] = "Twitter"
	m["auth0"] = "Auth0"
	m["hoper"] = "hoper"

	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return &ProviderIndex{Providers: keys, ProvidersMap: m}

}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

func CallBack(ctx iris.Context) {
	user, err := CompleteUserAuth(ctx)
	if err != nil {
		BeginAuthHandler(ctx)
		return
	}
	t, _ := template.New("foo").Parse(userTemplate)
	t.Execute(ctx.ResponseWriter(), user)
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`

func Index(res http.ResponseWriter, req *http.Request) {
	t, _ := template.New("foo").Parse(indexTemplate)
	t.Execute(res, NewAuth())
}
