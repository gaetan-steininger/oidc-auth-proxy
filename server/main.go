package server

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	sessionCookieName string
	sessionStore      *sessions.CookieStore
	sessionMaxAge     int
	oauth2Config      oauth2.Config
	oidcVerifier      *oidc.IDTokenVerifier
	allowedGroups     []string
	allowedUsers      []string
	userClaim         string
	groupsClaim       string
	backendUrl        *url.URL
)

func Run(
	sessionCookieNameArg string,
	sessionSecret string,
	maxAge int,
	backendUrlString string,
	oidcScopes []string,
	oidcRedirectUrl string,
	oidcClientId string,
	oidcClientSecret string,
	oidcProviderUrl string,
	oidcUserClaim string,
	oidcGroupsClaim string,
	allowedUsersList []string,
	allowedGroupsList []string,
) {
	var err error
	backendUrl, err = url.Parse(backendUrlString)
	if err != nil {
		log.Fatalf("invalid backend url : %v", err)
	}
	sessionCookieName = sessionCookieNameArg
	sessionStore = sessions.NewCookieStore([]byte(sessionSecret))
	sessionMaxAge = maxAge
	userClaim = oidcUserClaim
	groupsClaim = oidcGroupsClaim
	allowedGroups = allowedGroupsList
	allowedUsers = allowedUsersList

	log.Printf("backend url : %v\n", backendUrl)
	log.Printf("session cookie name : %s\n", sessionCookieName)
	log.Printf("oidc user claim : %s\n", userClaim)
	log.Printf("oidc groups claim : %s\n", groupsClaim)
	log.Printf("allowed users : %v", allowedUsers)
	log.Printf("allowed groups : %v", allowedGroups)

	provider, err := oidc.NewProvider(context.Background(), oidcProviderUrl)
	if err != nil {
		log.Fatalf("failed to get provider: %v", err)
	}

	parsedOidcRedirectUrl, err := url.Parse(oidcRedirectUrl)
	if err != nil {
		log.Fatalf("failed to parse oidc redirect url : %v", err)
	}

	oauth2Config = oauth2.Config{
		ClientID:     oidcClientId,
		ClientSecret: oidcClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  oidcRedirectUrl,
		Scopes:       oidcScopes,
	}

	oidcVerifier = provider.Verifier(&oidc.Config{ClientID: oidcClientId})

	http.HandleFunc(parsedOidcRedirectUrl.Path, callbackHandler)
	http.HandleFunc("/", proxyHandler)

	log.Println("oidc proxy server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
