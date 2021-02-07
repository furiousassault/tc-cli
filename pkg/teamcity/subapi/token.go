package subapi

import (
	"fmt"
	"net/url"

	"github.com/dghubble/sling"
)

type Token struct {
	Name  string `json:"name" xml:"name"`
	Value string `json:"value" xml:"value"`
}

type Tokens struct {
	Count int      `json:"count,omitempty" xml:"count"`
	Items []*Token `json:"token" xml:"token"`
}

type TokenService struct {
	sling         *sling.Sling
	httpClient    sling.Doer
	requestsMaker *requestsMaker
}

func NewTokenService(base *sling.Sling, client sling.Doer) *TokenService {
	s := base.Path("users/")
	return &TokenService{
		sling:         s,
		httpClient:    client,
		requestsMaker: newRequestsMakerWithSling(client, s),
	}
}

func (s *TokenService) TokenCreate(username, tokenName string) (token Token, err error) {
	path := fmt.Sprintf("%s/tokens", username)
	data := Token{
		Name: tokenName,
	}

	err = s.requestsMaker.post(path, data, &token, "tokenCreate")
	return
}

func (s *TokenService) TokenRemove(userID string, tokenName string) (err error) {
	path := fmt.Sprintf("%s/tokens/%s", userID, tokenName)
	return s.requestsMaker.delete(path, "tokenRemove")
}

func (s *TokenService) TokenList(userID string) (tokens Tokens, err error) {
	path := fmt.Sprintf("%s/tokens", url.PathEscape(userID))
	err = s.requestsMaker.get(path, &tokens, "tokenList")
	return
}
