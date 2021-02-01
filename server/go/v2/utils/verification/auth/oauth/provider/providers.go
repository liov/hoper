package provider

import (
	"os"

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


type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}


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
		New(os.Getenv("HOPER_KEY"), os.Getenv("HOPER_SECRET"), "http://localhost:8080/auth/hoper/callback"),
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

