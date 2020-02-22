package cmd

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// Config includes target manager and user info
type Config struct {
	ManagerAddr string `yaml:"ManagerAddr"`
	User string `yaml:"User"`
	Passwd string `yaml:"Passwd"`
}

func readConfigFromFile(filepath string) (*Config, error) {
	c := &Config{}
	if filepath == "~/.jf/config.yml" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		filepath = path.Join(home, ".jf/config.yml")
	}
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	err = yaml.UnmarshalStrict(file, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}