package jira

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/qmuntal/stateless"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/exp/slog"
	"golang.org/x/oauth2"
)

const (
	// states
	stateInit                        = "Init"
	stateOauthConfigSetup            = "OauthConfigSetup"
	stateTokenObtained               = "TokenObtained"
	stateAccessibleResourcesObtained = "AccessibleResourcesObtained"

	// events
	triggerSetupOauthConfig         = "SetupOauthConfig"
	triggerCodeExchange             = "CodeExchange"
	triggerFetchAccessibleResources = "FetchAccessibleResources"
)

type Authenticator struct {
	loginWorkflow *stateless.StateMachine

	oauthToken *oauth2.Token
	conf       *oauth2.Config
	verifier   string
	ctx        context.Context
	callback   chan *oauth2.Token
}

func NewAuthenticator() *Authenticator {
	a := &Authenticator{}

	loginWorkflow := stateless.NewStateMachine(stateInit)
	loginWorkflow.Configure(stateInit).
		Permit(triggerSetupOauthConfig, stateOauthConfigSetup)

	loginWorkflow.Configure(stateOauthConfigSetup).
		OnEntry(a.setupOauthConfig).
		Permit(triggerCodeExchange, stateTokenObtained, a.isVerifierPresent)

	loginWorkflow.Configure(stateTokenObtained).
		OnEntry(a.exchangeCode).
		Permit(triggerFetchAccessibleResources, stateAccessibleResourcesObtained, a.isTokenValid, a.isOauthContextPresent)

	loginWorkflow.Configure(stateAccessibleResourcesObtained).
		OnEntry(a.fetchAccessibleResources)

	a.loginWorkflow = loginWorkflow

	return a
}

func (a *Authenticator) Login(ctx context.Context) error {
	a.loginWorkflow.Fire(triggerSetupOauthConfig)

	if err := a.loginWorkflow.FireCtx(ctx, triggerCodeExchange); err != nil {
		return err
	}

	if err := a.loginWorkflow.FireCtx(ctx, triggerFetchAccessibleResources); err != nil {
		return err
	}
	return nil
}

func (a *Authenticator) isVerifierPresent(ctx context.Context, args ...any) bool {
	return a.verifier != ""
}

func (a *Authenticator) isTokenValid(ctx context.Context, args ...any) bool {
	return a.oauthToken != nil && a.oauthToken.Valid()
}

func (a *Authenticator) isOauthContextPresent(ctx context.Context, args ...any) bool {
	return a.ctx != nil
}

func (a *Authenticator) setupOauthConfig(ctx context.Context, args ...any) error {
	a.conf = &oauth2.Config{
		ClientID:     "OOGf9PTJL0hGGC5hWD17G6OkiGKjO0FG",
		ClientSecret: "ATOAhihA9MN3TOWAJEC4DxxPZMxGyjmA_mH8rUtSGXRIoUP6WQ3UvjCk5Mtx9TUBH6JF089B37D6",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://auth.atlassian.com/authorize",
			TokenURL: "https://auth.atlassian.com/oauth/token",
		},
		RedirectURL: "http://127.0.0.1:9999/gh-jira-changelog/oauth/callback",
		Scopes:      []string{"read:jira-work"},
	}

	a.verifier = oauth2.GenerateVerifier()

	slog.Info("Configured oauth")
	return nil
}

func (a *Authenticator) exchangeCode(ctx context.Context, args ...any) error {
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

	// sping up server for callback from RedirectURL and shut it down once we get response
	a.callback = make(chan *oauth2.Token)
	mux := http.NewServeMux()
	mux.HandleFunc("/gh-jira-changelog/oauth/callback", http.HandlerFunc(a.callbackHandler))
	svr := http.Server{Addr: "127.0.0.1:9999", Handler: mux}
	go func() { svr.ListenAndServe() }()

	a.oauthToken = <-a.callback
	svr.Shutdown(a.ctx)
	return nil
}

func (a *Authenticator) fetchAccessibleResources(ctx context.Context, args ...any) error {
	resp, err := a.Client().Get("https://api.atlassian.com/oauth/token/accessible-resources")
	if err != nil || resp.StatusCode != http.StatusOK {
		slog.Error("Failed to fetch accessible-resource from jira", "error", err)
		slog.Error("Response information", "response_status", resp.Status)
		return err
	}

	accessibleResources, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read accessibleResources", "error", err)
		return err
	}
	slog.Debug("Retrieved accessible resources successfully", "resources", string(accessibleResources))

	return nil
}

func (a *Authenticator) Client() *http.Client {
	return a.conf.Client(a.ctx, a.oauthToken)
}

func (a *Authenticator) callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := a.conf.Exchange(a.ctx, code, oauth2.VerifierOption(a.verifier))
	if err != nil {
		log.Fatal(err)
	}

	msg := "<h1>Authentication successful!</h1>"
	msg = msg + "<p>You are authenticated and can now return to the CLI.</p>"
	fmt.Fprintln(w, msg)

	a.callback <- token
}
