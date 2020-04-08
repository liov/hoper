package oauth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/kataras/iris/v12"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func Test_Oauth(t *testing.T) {
	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_KEY"),
		ClientSecret: os.Getenv("GITHUB_SECRET"),
		Scopes:       []string{"SCOPE1", "SCOPE2"},
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8070/auth/github/callback",
	}
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}

	// Use the custom HTTP client when requesting a token.
	httpClient := &http.Client{Timeout: 2 * time.Second}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := conf.Client(ctx, tok)
	_ = client
}

func Test_Oauth2(t *testing.T) {
	app := iris.New()
	app.Get("/", iris.FromStd(Index))
	app.Get("/auth/{provider}", CallBack)
	app.Get("/auth/{provider}/callback", CallBack)
	app.Get("/logout/{provider}", Logout)
	app.Run(iris.Addr(":8070"))
}
