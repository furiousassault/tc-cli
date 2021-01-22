package authorization

const APIPathSuffixAuthToken = ""

type AuthorizerToken struct {
	Token string
}

func (a AuthorizerToken) ProvideURLAuthSuffix() string {
	return APIPathSuffixAuthToken
}

func AuthToken(token string) AuthorizerToken {
	return AuthorizerToken{Token: token}
}
