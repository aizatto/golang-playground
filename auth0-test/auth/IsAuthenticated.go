package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"session"

	oidc "github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	store, err := session.Store.Get(r, session.DEFAULT)
	if err != nil {
		http.Error(w, "cannot get session", http.StatusInternalServerError)
		return
	}

	if shouldAuthenticate(store) {
		authenticator, err := NewAuthenticator()
		if authenticator == nil {
			http.Error(w, "authenticator is nil", http.StatusInternalServerError)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// next(w, r)
		b := make([]byte, 32)
		_, err = rand.Read(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		state := base64.StdEncoding.EncodeToString(b)
		store.Values["state"] = state
		err = store.Save(r, w)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, authenticator.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
		return
	}

	next(w, r)
}

func shouldAuthenticate(store *sessions.Session) bool {
	if _, ok := store.Values["profile"]; !ok {
		return true
	}

	// check validity of profile
	valid := true

	return valid != true
}

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, "https://"+os.Getenv("AUTH0_DOMAIN")+"/")
	if err != nil {
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	fmt.Println(conf.RedirectURL)

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}
