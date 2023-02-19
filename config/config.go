package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type BacklogSettings struct {
	Project      string `json:"project"`
	Organization string `json:"organization"`
}

type Config map[string]HostConfig

type HostConfig struct {
	Hostname     string `yaml:"hostname"`
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token"`
}

func ReadConfig() (Config, error) {
	configBytes, err := ioutil.ReadFile(configPath())

	if err != nil {
		return Config{"": HostConfig{}}, err
	}

	var hostConfig Config
	err = yaml.Unmarshal(configBytes, &hostConfig)

	if err != nil {
		log.Fatalf("parse error: %s at %s\n", err, configPath)
	}

	return hostConfig, nil
}

func WriteConfig(space string, update *HostConfig) {
	configPath := configPath()
	configBytes, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}

	var hostConfig Config
	err = yaml.Unmarshal(configBytes, &hostConfig)

	if err != nil {
		log.Fatalf("parse error: %s at %s\n", err, configPath)
	}

	host := fmt.Sprintf("%s.backlog.com", space)

	_, ok := hostConfig[host]

	if ok {
		for key, config := range hostConfig {
			if key == host {
				if update.Hostname != "" {
					config.Hostname = update.Hostname
				}
				if update.AccessToken != "" {
					config.AccessToken = update.AccessToken

				}
				if update.RefreshToken != "" {
					config.RefreshToken = update.RefreshToken
				}

				hostConfig[key] = config
			}
		}

	} else {
		hostConfig[host] = *update
	}

	newData, err := yaml.Marshal(&hostConfig)

	if err != nil {
		log.Fatalf("error marshaling YAML: %v", err)
	}

	// Write the updated YAML data back to the file
	err = ioutil.WriteFile(configPath, newData, 0644)
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
}

func GetAccessToken(space string) (string, error) {
	configPath := configPath()
	configBytes, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}

	var hostConfig Config
	err = yaml.Unmarshal(configBytes, &hostConfig)

	if err != nil {
		log.Fatalf("parse error: %s at %s\n", err, configPath)
	}

	host := fmt.Sprintf("%s.backlog.com", space)

	_, ok := hostConfig[host]

	if ok && hostConfig[host].AccessToken != "" {
		return hostConfig[host].AccessToken, nil
	} else {
		return "", fmt.Errorf("access token not found")
	}
}

func configPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homeDir, ".config", "bl", "hosts.yml")
}
