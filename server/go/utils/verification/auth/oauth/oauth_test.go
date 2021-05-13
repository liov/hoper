package oauth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func Test_Oauth(t *testing.T) {
	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID:     "222222",
		ClientSecret: "22222222",
		Scopes:       []string{"SCOPE1", "SCOPE2"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8070/oauth/authorize",
			TokenURL: "http://localhost:8070/oauth/access_token",
		},
		RedirectURL: "http://localhost:8080/auth/github/callback",
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
