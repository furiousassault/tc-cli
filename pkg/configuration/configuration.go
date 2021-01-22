package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const (
	DefaultPathToken         = ".tc-client/access.token"
	DefaultPathConfiguration = ".tc-client/configuration.yaml"
)

var (
	ConfigPath string
	config     = &Configuration{}
)

type Configuration struct {
	API API
}

type API struct {
	Authorization Authorization `yaml:"authorization"`
	HTTP          HTTP          `yaml:"http"`

	URL string `yaml:"url" envconfig:"TC_URL"` // including schema, url and port
}

type Authorization struct {
	Username      string `yaml:"username" envconfig:"TC_USERNAME"`
	Password      string `yaml:"password" envconfig:"TC_PASSWORD"`
	Token         string `yaml:"token" envconfig:"TC_TOKEN"`
	TokenFilePath string `yaml:"token_path" envconfig:"TC_TOKEN_FILE_PATH"`
}

type HTTP struct {
	RequestTimeout time.Duration `yaml:"request_timeout" envconfig:"TC_REQUEST_TIMEOUT"`
}

func InitConfigFromYAML() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error during user's home directory evaluation")
	}

	if ConfigPath == "" {
		defaultPathConfiguration := filepath.Join(homeDir, DefaultPathConfiguration)
		ConfigPath = defaultPathConfiguration
		fmt.Printf(
			"Configuration file path is not specified, trying to find it in '%s'\n",
			defaultPathConfiguration,
		)
	}

	file, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(file, config); err != nil {
		return err
	}

	if err = envconfig.Process("", config); err != nil {
		// todo global variables overuse with cobra everywhere... pass normal logger somehow.
		fmt.Println("Error processing environment variables in configuration init")
	}

	if config.API.Authorization.Username == "" || config.API.Authorization.Password == "" {
		if config.API.Authorization.Token == "" {
			if config.API.Authorization.TokenFilePath == "" {
				defaultTokenPath := filepath.Join(homeDir, DefaultPathConfiguration)
				fmt.Printf(
					"Token file path is not specified, trying to find it in '%s'\n",
					defaultTokenPath,
				)

				config.API.Authorization.TokenFilePath = defaultTokenPath
			}

			fileToken, err := ioutil.ReadFile(config.API.Authorization.TokenFilePath)
			if err != nil {
				fmt.Println("Token file read unsuccessful")
				return nil
			}
			
			config.API.Authorization.Token = string(fileToken)
		}
	}

	return nil
}

func GetConfig() Configuration {
	return *config
}
