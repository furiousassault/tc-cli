package configuration

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	DefaultPathToken              = ".tc-client/access.token"
	DefaultPathConfiguration      = ".tc-client/configuration.yaml"
	DefaultPathArtifactsDirectory = "/tmp/tc-client/artifacts/"
)

func configPathWithDefault(path string) (configPath string, err error) {
	if path != "" {
		return path, nil
	}

	return pathInHomeDir(DefaultPathConfiguration)
}

func tokenPathWithDefault(path string) (tokenPath string, err error) {
	if path != "" {
		return path, nil
	}

	return pathInHomeDir(DefaultPathToken)
}

func pathInHomeDir(pathSuffix string) (path string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrapf(err, "error during user's home directory evaluation")
	}

	return filepath.Join(homeDir, pathSuffix), nil
}
