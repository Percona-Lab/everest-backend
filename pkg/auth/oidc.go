package auth

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type OIDC struct {
	IDTokenVerifier *oidc.IDTokenVerifier
	Oauth2Config    oauth2.Config
	provider        *oidc.Provider
}

func New(ctx context.Context) (*OIDC, error) {
	provider, err := oidc.NewProvider(ctx, "http://127.0.0.1:5556/dex")
	if err != nil {
		return nil, err
	}

	oauth2Config := oauth2.Config{
		ClientID:     "foo",
		ClientSecret: "bar",

		// The redirectURL.
		RedirectURL: "http://localhost:3000/login/callback",

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		// Other scopes, such as "groups" can be requested.
		Scopes: []string{oidc.ScopeOpenID, "profile", "email", "groups"},
	}

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: "foo"})

	return &OIDC{
		IDTokenVerifier: idTokenVerifier,
		Oauth2Config:    oauth2Config,
		provider:        provider,
	}, nil
}

func (o *OIDC) HandleCallback(ctx context.Context, code, state string) (string, error) {
	// TODO: verify state
	oauth2Token, err := o.Oauth2Config.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("could not verify token. %w", err)
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("missing token")
	}

	// Parse and verify ID Token payload.
	idToken, err := o.IDTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return "", fmt.Errorf("could not verify ID token. %w", err)
	}

	// Extract custom claims.
	var claims struct {
		Email    string   `json:"email"`
		Verified bool     `json:"email_verified"`
		Groups   []string `json:"groups"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return "", fmt.Errorf("could not parse claims. %w", err)
	}

	return rawIDToken, nil
}
