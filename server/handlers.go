package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/gorilla/sessions"
)

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionCookieName)

	origPath, ok := session.Values["path"].(string)
	if !ok {
		origPath = "/"
	}

	state, ok := session.Values["state"].(string)
	if !ok || r.URL.Query().Get("state") != state {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}
	token, err := oauth2Config.Exchange(context.Background(), r.URL.Query().Get("code"))
	if err != nil {
		log.Printf("failed to exchange token : %v", err)
		http.Error(w, "failed to exchange token", http.StatusInternalServerError)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Print("missing ID token")
		http.Error(w, "missing ID token", http.StatusInternalServerError)
		return
	}

	idToken, err := oidcVerifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		log.Printf("failed to verify ID token : %v", err)
		http.Error(w, "failed to verify ID token", http.StatusInternalServerError)
		return
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("failed to parse claims : %v", err)
		http.Error(w, "failed to parse claims", http.StatusInternalServerError)
		return
	}

	username, groups, err := extractUserGroupsClaims(claims, userClaim, groupsClaim)
	if err != nil {
		log.Printf("failed to extract user / groups claims : %v", err)
		http.Error(w, "failed to extract user / groups claims", http.StatusInternalServerError)
		return
	}

	if len(allowedGroups) > 0 || len(allowedUsers) > 0 {
		authorized := false

		if len(allowedGroups) > 0 {
			for _, userGroup := range groups {
				for _, allowedGroup := range allowedGroups {
					if userGroup == allowedGroup {
						authorized = true
						break
					}
				}
			}
		}

		if len(allowedUsers) > 0 {
			for _, allowedUser := range allowedUsers {
				if allowedUser == username {
					authorized = true
					break
				}
			}
		}

		if !authorized {
			msg := fmt.Sprintf("user %s unauthorized", username)
			log.Print(msg)
			http.Error(w, msg, http.StatusForbidden)
			return
		}
	}

	log.Printf("%s - successful login", username)

	session.Values["authenticated"] = true
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   sessionMaxAge,
		HttpOnly: true,
		Secure:   true,
	}
	session.Save(r, w)

	http.Redirect(w, r, origPath, http.StatusFound)
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionStore.Get(r, sessionCookieName)
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		state := generateState()
		session, _ := sessionStore.Get(r, sessionCookieName)
		session.Values["state"] = state
		session.Values["path"] = r.URL.Path
		session.Save(r, w)

		http.Redirect(w, r, oauth2Config.AuthCodeURL(state), http.StatusFound)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backendUrl)
	proxy.ServeHTTP(w, r)
}
