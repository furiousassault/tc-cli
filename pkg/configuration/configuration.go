package configuration

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	API API
}

type API struct {
	Authorization Authorization `yaml:"authorization"`
	HTTP          HTTP          `yaml:"http"`

	URL string `yaml:"url" envconfig:"TC_CLI_URL"` // including schema, url and port
}

type Authorization struct {
	Username      string `yaml:"username" envconfig:"TC_CLI_USERNAME"`
	Password      string `yaml:"password" envconfig:"TC_CLI_PASSWORD"`
	Token         string `yaml:"token" envconfig:"TC_CLI_TOKEN"`
	TokenFilePath string `yaml:"token_path" envconfig:"TC_CLI_TOKEN_FILE_PATH"`
}

type HTTP struct {
	RequestTimeout time.Duration `yaml:"request_timeout" envconfig:"TC_CLI_REQUEST_TIMEOUT"`
}

func ConfigFromYAML(configPath string) (config *Configuration, err error) {
	configPath, err = configPathWithDefault(configPath)
	if err != nil {
		return nil, err
	}

	config, err = configurationFromYaml(configPath)
	if err != nil {
		return nil, err
	}

	// fill unset config fields with envs
	err = envconfig.Process("", config)
	if err != nil {
		fmt.Println("Error processing environment variables in configuration init")
	}

	initializeAuthParameters(config)
	return config, config.validate()
}

func initializeAuthParameters(config *Configuration) {
	// if it's the httpAuth, no point in reading token path
	if config.API.Authorization.Username != "" && config.API.Authorization.Password != "" {
		return
	}

	// if token is set explicitly, the same
	if config.API.Authorization.Token != "" {
		return
	}

	// try to evaluate token from token path
	tokenFilePath, err := tokenPathWithDefault(config.API.Authorization.TokenFilePath)
	if err != nil {
		return
	}

	token, errTokenFromFile := fileContentString(tokenFilePath)
	if errTokenFromFile != nil {
		return
	}

	config.API.Authorization.Token = token
	return
}
