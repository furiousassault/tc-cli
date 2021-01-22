package subapi

import (
	"net/url"
)

// Locator represents an arbitrary locator to be used when querying resources
type Locator string

// LocatorID creates a locator by ID
func LocatorID(id string) Locator {
	return Locator(url.QueryEscape("id:") + id)
}

func (l Locator) String() string {
	return string(l)
}
