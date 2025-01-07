package m2m

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/auth0/go-auth0/authentication"
)

// Assertion of the proper interface
var _ M2MGenerator = (*M2M)(nil)

const (
	audience          string = "https://career-cue-auth-api.com"
	clientCredentials string = "client_credentials"
)

type M2MGenerator interface {
	GetToken() (Token, error)
}

type M2M struct {
	domain       string
	clientID     string
	clientSecret string
	auth         *authentication.Authentication
	mu           *sync.RWMutex
	token        Token
}

type tokenRequestBody struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

type tokenResponseBody struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Token struct {
	accessToken string
	expiresIn   time.Time
	tokenType   string
}

func (t Token) IsExpired() bool {
	now := time.Now().UTC()
	return t.expiresIn.Before(now)
}

func (t Token) IsZero() bool {
	return t.expiresIn.IsZero()
}

func (t Token) GetHeaderValue() string {
	return fmt.Sprintf("%s %s", t.tokenType, t.accessToken)
}

func NewM2M(auth0Domain, auth0ClientID, auth0ClientSecret string) (*M2M, error) {
	ctx := context.Background()
	a, err := authentication.New(ctx, auth0Domain, authentication.WithClientID(auth0ClientID), authentication.WithClientSecret(auth0ClientSecret))
	if err != nil {
		return &M2M{}, fmt.Errorf("cannot connect to auth0 Management API: %w", err)
	}

	return &M2M{
		domain:       auth0Domain,
		clientID:     auth0ClientID,
		clientSecret: auth0ClientSecret,
		auth:         a,
		mu:           new(sync.RWMutex),
		token:        Token{},
	}, nil
}

func (g *M2M) GetToken() (Token, error) {
	if g.token.IsExpired() {
		token, err := g.fetchToken()
		if err != nil {
			return Token{}, errors.New("trouble fetching new token from auth0")
		}

		g.mu.Lock()
		g.token = token
		g.mu.Unlock()
	}

	// Read token from the m2m struct
	var foundToken Token
	g.mu.RLock()
	foundToken = g.token
	g.mu.RUnlock()

	return foundToken, nil
}

func (g *M2M) fetchToken() (Token, error) {
	reqBody := tokenRequestBody{
		ClientID:     g.clientID,
		ClientSecret: g.clientSecret,
		Audience:     audience,
		GrantType:    clientCredentials,
	}

	req, err := g.auth.NewRequest(context.Background(), http.MethodPost, fmt.Sprintf("%s/oauth/token", g.domain), reqBody)
	if err != nil {
		return Token{}, err
	}

	res, err := g.auth.Do(req)
	if err != nil {
		return Token{}, err
	}
	defer res.Body.Close()

	var tokenResp tokenResponseBody
	if err := json.NewDecoder(res.Body).Decode(&tokenResp); err != nil {
		return Token{}, err
	}

	token := Token{
		accessToken: tokenResp.AccessToken,
		tokenType:   tokenResp.TokenType,
		expiresIn:   time.Unix(tokenResp.ExpiresIn, 0),
	}

	return token, nil
}
