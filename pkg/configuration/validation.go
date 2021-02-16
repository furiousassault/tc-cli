package configuration

import (
	"fmt"
)

func (c *Configuration) validate() error {
	if c.API.URL == "" {
		return fmt.Errorf("%w: Teamcity URL is not provided", ErrConfigValidation)
	}

	return nil
}
