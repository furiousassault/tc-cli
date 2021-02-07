package authorization

import "github.com/furiousassault/tc-cli/pkg/configuration"

const APIPathSuffixAuthGuest = "guestAuth"

// AuthorizerGuest should be used for guest authorization mode.
type AuthorizerGuest struct{}

func (a AuthorizerGuest) ProvideURLAuthSuffix() string {
	return APIPathSuffixAuthGuest
}

func (a AuthorizerGuest) Credentials() configuration.Authorization {
	return configuration.Authorization{}
}

// NewAuthorizerGuest creates an Authorizer for non-authorized access.
func NewAuthorizerGuest() AuthorizerGuest {
	return AuthorizerGuest{}
}
