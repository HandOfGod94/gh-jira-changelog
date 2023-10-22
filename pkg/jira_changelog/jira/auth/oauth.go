package auth

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/qmuntal/stateless"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/exp/slog"
	"golang.org/x/oauth2"
)

const (
	// states
	stateInit                       = "Init"
	stateTokenConfigured            = "TokenConfigured"
	stateAllowedResourcesConfigured = "AllowedResourcesConfigured"

	// events
	triggerCodeExchange          = "CodeExchange"
	triggerSaveToken             = "SaveToken"
	triggerFetchAllowedResources = "FetchAccessibleResources"
	triggerPersistResourcesInfo  = "PersistResourcesInfo"
)

type oauthAuthenticator struct {
	loginWorkflow *stateless.StateMachine

	oauthToken               *oauth2.Token
	conf                     *oauth2.Config
	verifier                 string
	ctx                      context.Context
	callback                 chan *oauth2.Token
	allowedResourcesResponse []byte
}

func NewAuthenticator() *oauthAuthenticator {
	a := &oauthAuthenticator{}

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

	loginWorkflow := stateless.NewStateMachine(stateInit)
	loginWorkflow.Configure(stateInit).
		Permit(triggerCodeExchange, stateTokenConfigured, a.isVerifierPresent)

	loginWorkflow.Configure(stateTokenConfigured).
		OnEntryFrom(triggerCodeExchange, a.exchangeCode).
		OnEntryFrom(triggerSaveToken, a.saveToken).
		Permit(triggerFetchAllowedResources, stateAllowedResourcesConfigured, a.isTokenValid, a.isOauthContextPresent).
		PermitReentry(triggerSaveToken, a.isTokenValid)

	loginWorkflow.Configure(stateAllowedResourcesConfigured).
		OnEntryFrom(triggerFetchAllowedResources, a.fetchAccessibleResources).
		OnEntryFrom(triggerPersistResourcesInfo, a.saveJiraConfig).
		PermitReentry(triggerPersistResourcesInfo)

	a.loginWorkflow = loginWorkflow

	return a
}

func (a *oauthAuthenticator) Login(ctx context.Context) error {
	if err := a.loginWorkflow.FireCtx(ctx, triggerCodeExchange); err != nil {
		return err
	}

	if err := a.loginWorkflow.FireCtx(ctx, triggerSaveToken); err != nil {
		return err
	}

	if err := a.loginWorkflow.FireCtx(ctx, triggerFetchAllowedResources); err != nil {
		return err
	}

	if err := a.loginWorkflow.FireCtx(ctx, triggerPersistResourcesInfo); err != nil {
		return err
	}
	return nil
}

func (a *oauthAuthenticator) Client() *http.Client {
	return a.conf.Client(a.ctx, a.oauthToken)
}

// ActionFuncs
func (a *oauthAuthenticator) exchangeCode(ctx context.Context, args ...any) error {
	url := a.conf.AuthCodeURL("state", oauth2.S256ChallengeOption(a.verifier),
		oauth2.SetAuthURLParam("response_type", "code"),
		oauth2.SetAuthURLParam("prompt", "consent"),
		oauth2.SetAuthURLParam("audience", "api.atlassian.com"),
	)

	fmt.Println(color.CyanString("You will now be taken to browser for authentication."))
	fmt.Println(color.CyanString("Please grant permission to access jira issues"))

	time.Sleep(1 * time.Second)
	open.Run(url)
	time.Sleep(1 * time.Second)

	a.ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})

	// spin up server for callback from RedirectURL and shut it down once we get response
	a.callback = make(chan *oauth2.Token)
	mux := http.NewServeMux()
	mux.HandleFunc("/gh-jira-changelog/oauth/callback", a.callbackHandler)
	svr := http.Server{Addr: "127.0.0.1:9999", Handler: mux}
	go func() { svr.ListenAndServe() }()

	a.oauthToken = <-a.callback
	svr.Shutdown(a.ctx)
	return nil
}

func (a *oauthAuthenticator) saveToken(ctx context.Context, args ...any) error {

	token := Token{
		AccessToken:  a.oauthToken.AccessToken,
		RefreshToken: a.oauthToken.RefreshToken,
		ExpiresIn:    a.oauthToken.Expiry,
	}

	return token.Save()
}

func (a *oauthAuthenticator) fetchAccessibleResources(ctx context.Context, args ...any) error {
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

	a.allowedResourcesResponse = accessibleResources

	return nil
}

func (a *oauthAuthenticator) saveJiraConfig(ctx context.Context, args ...any) error {
	resources, err := parseResources(a.allowedResourcesResponse)
	if err != nil {
		return err
	}

	err = resources[0].Save()
	if err != nil {
		return err
	}
	return nil
}

// GaurdFuncs
func (a *oauthAuthenticator) isVerifierPresent(ctx context.Context, args ...any) bool {
	return a.verifier != ""
}

func (a *oauthAuthenticator) isTokenValid(ctx context.Context, args ...any) bool {
	return a.oauthToken != nil && a.oauthToken.Valid()
}

func (a *oauthAuthenticator) isOauthContextPresent(ctx context.Context, args ...any) bool {
	return a.ctx != nil
}

// oauth RedirectURL callback handler
func (a *oauthAuthenticator) callbackHandler(w http.ResponseWriter, r *http.Request) {
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
