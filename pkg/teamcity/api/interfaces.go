package api

type Authorizer interface {
	URLAuthSuffixProvider
}

type URLAuthSuffixProvider interface {
	ProvideURLAuthSuffix() string
}
