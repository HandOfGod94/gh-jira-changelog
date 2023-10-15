package jira

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/skratchdot/open-golang/open"
	"golang.org/x/exp/slog"
	"golang.org/x/oauth2"
)

type Authenticator struct {
	oauthToken *oauth2.Token
	conf       *oauth2.Config
	verifier   string
	ctx        context.Context
}

func NewAuthenticator() *Authenticator {
	conf := &oauth2.Config{
		ClientID:     "OOGf9PTJL0hGGC5hWD17G6OkiGKjO0FG",
		ClientSecret: "ATOAhihA9MN3TOWAJEC4DxxPZMxGyjmA_mH8rUtSGXRIoUP6WQ3UvjCk5Mtx9TUBH6JF089B37D6",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://auth.atlassian.com/authorize",
			TokenURL: "https://auth.atlassian.com/oauth/token",
		},
		RedirectURL: "http://127.0.0.1:9999/gh-jira-changelog/oauth/callback",
		Scopes:      []string{"read:jira-work"},
	}

	return &Authenticator{
		conf: conf,
	}
}

func (a *Authenticator) Login(ctx context.Context) error {
	a.verifier = oauth2.GenerateVerifier()
	url := a.conf.AuthCodeURL("state", oauth2.S256ChallengeOption(a.verifier),
		oauth2.SetAuthURLParam("response_type", "code"),
		oauth2.SetAuthURLParam("prompt", "consent"),
		oauth2.SetAuthURLParam("audience", "api.atlassian.com"),
	)

	slog.Info("You will now be taken to browser for authentication.")
	slog.Info("Please grant permission to access jira issues")

	time.Sleep(1 * time.Second)
	open.Run(url)
	time.Sleep(1 * time.Second)

	a.ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})

	http.HandleFunc("/gh-jira-changelog/oauth/callback", http.HandlerFunc(a.callbackHandler))
	http.ListenAndServe(":9999", nil)

	return nil
}

func (a *Authenticator) Client() *http.Client {
	return a.conf.Client(a.ctx, a.oauthToken)
}

func (a *Authenticator) callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	tok, err := a.conf.Exchange(a.ctx, code, oauth2.VerifierOption(a.verifier))
	if err != nil {
		log.Fatal(err)
	}

	a.oauthToken = tok

	resp, err := a.Client().Get("https://api.atlassian.com/oauth/token/accessible-resources")
	if err != nil || resp.StatusCode != 200 {
		slog.Error("Failed to fetch accessible-resource from jira", "error", err)
		slog.Error("Response information", "response_status", resp.Status)
		panic(err)
	}

	accessibleResources, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read accessibleResources", "error", err)
		panic(err)
	}
	slog.Debug("Retrieved accessible resources successfully", "resources", string(accessibleResources))


	msg := "<h1>Authentication successful!</h1>"
	msg = msg + "<p>You are authenticated and can now return to the CLI.</p>"
	fmt.Fprintln(w, msg)
}
