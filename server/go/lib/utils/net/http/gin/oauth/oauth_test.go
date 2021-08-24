package oauth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/server/go/lib/utils/net/http/gin/handler"
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
	r := gin.New()
	r.GET("/", handler.Wrap(Index))
	r.GET("/auth/{provider}", CallBack)
	r.GET("/auth/{provider}/callback", CallBack)
	r.GET("/logout/{provider}", Logout)
	r.Run(":8080")
}
