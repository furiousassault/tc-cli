package configuration

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func fileContentString(filePath string) (content string, err error) {
	t, err := ioutil.ReadFile(filePath)
	return string(t), err
}

func configurationFromYaml(filePath string) (config *Configuration, err error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err // todo wrap
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err // todo wrap
	}

	return
}
