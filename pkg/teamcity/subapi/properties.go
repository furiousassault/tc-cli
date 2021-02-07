package subapi

// Property represents a key/value/type structure of properties present in several other entities
type Property struct {
	Name      string `json:"name,omitempty" xml:"name"`
	Value     string `json:"value" xml:"value"`
	Type      *Type  `json:"type,omitempty"`
	Inherited *bool  `json:"inherited,omitempty" xml:"inherited"`
}

// Type represents a parameter type.
type Type struct {
	RawValue string `json:"rawValue,omitempty" xml:"rawValue"`
}

// Properties represents a collection of key/value properties for a resource
type Properties struct {
	Count int32       `json:"count,omitempty" xml:"count"`
	Href  string      `json:"href,omitempty" xml:"href"`
	Items []*Property `json:"property"`
}
