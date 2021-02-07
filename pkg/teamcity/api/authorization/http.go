package authorization

const APIPathSuffixAuthHTTP = "httpAuth"

type AuthorizerHTTP struct {
	Username, Password string
}

func (a AuthorizerHTTP) ProvideURLAuthSuffix() string {
	return APIPathSuffixAuthHTTP
}

func NewAuthorizerHTTP(username, password string) AuthorizerHTTP {
	return AuthorizerHTTP{Username: username, Password: password}
}
