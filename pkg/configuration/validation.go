package configuration

import "github.com/pkg/errors"

func (c *Configuration) validate() error {
	if c.API.URL == "" {
		return errors.Wrap(ErrConfigValidation, "Teamcity URL is not provided")
	}

	return nil
}
