package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	config "github.com/shufo/find-config"
	"gopkg.in/yaml.v3"
)

var defaultConfigName = "backlog.json"

type BacklogSettings struct {
	BacklogDomain string `json:"backlog_domain"`
	Organization  string `json:"organization"`
	Project       string `json:"project"`
}

type Config map[string]HostConfig

type HostConfig struct {
	Hostname     string `yaml:"hostname"`
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token"`
}

func ConfigExists() bool {
	file, err := os.Open(defaultConfigName)

	if err != nil {
		return false
	}

	defer file.Close()

	return true
}

func GetBacklogSetting() (BacklogSettings, error) {
	configBytes, err := ioutil.ReadFile(defaultConfigName)

	if err != nil {
		return BacklogSettings{}, err
	}

	var settings BacklogSettings
	err = json.Unmarshal(configBytes, &settings)

	if err != nil {
		log.Fatalf("parse error: %s at %s\n", err, defaultConfigName)
	}

	return settings, nil
}

func ReadConfig() (Config, error) {
	configBytes, err := ioutil.ReadFile(configPath())

	if err != nil {
		return Config{"": HostConfig{}}, err
	}

	var hostConfig Config
	err = yaml.Unmarshal(configBytes, &hostConfig)

	if err != nil {
		log.Fatalf("parse error: %s at %s\n", err, configPath())
	}

	return hostConfig, nil
}

func WriteConfig(space string, update *HostConfig) {
	configPath := configPath()

	var hostConfig Config = make(Config)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		createDefaultHostsConfig()
	}

	configBytes, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configBytes, &hostConfig)

	if err != nil {
		log.Fatalf("parse error: %s at %s\n", err, configPath)
	}

	host := update.Hostname

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

func GetAccessToken(setting BacklogSettings) (string, error) {
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

	host := fmt.Sprintf("%s.%s", setting.Organization, setting.BacklogDomain)

	_, ok := hostConfig[host]

	if ok && hostConfig[host].AccessToken != "" {
		return hostConfig[host].AccessToken, nil
	} else {
		return "", fmt.Errorf("access token not found")
	}
}

func GetRefreshToken(setting BacklogSettings) (string, error) {
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

	host := fmt.Sprintf("%s.%s", setting.Organization, setting.BacklogDomain)

	_, ok := hostConfig[host]

	if ok && hostConfig[host].RefreshToken != "" {
		return hostConfig[host].RefreshToken, nil
	} else {
		return "", fmt.Errorf("refresh token not found")
	}
}

func CreateDefaultConfig(p *BacklogSettings) {
	// Convert the Person struct to JSON
	jsonBytes, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		panic(err)
	}

	file, err := os.Create(defaultConfigName)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write the JSON to the file
	_, err = file.Write(jsonBytes)

	if err != nil {
		panic(err)
	}
}

func configPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(configDir, "bl", "hosts.yml")
}

func FindConfig(name string) (string, error) {
	return config.Find(name, ".")
}

func ShowConfigNotFound() {
	color.Red("backlog.json not found")
	color.White("You can set by `bk auth login`")
}

type CommandConfig struct {
	Aliases Aliases `json:"aliases"`
}

type Aliases map[string]string

func SetAlias(alias string, expansion string) {
	configPath := CommandConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		createDefaultCommandConfig()
	}

	configBytes, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}

	var commandConfig CommandConfig
	err = yaml.Unmarshal(configBytes, &commandConfig)

	if err != nil {
		log.Fatalf("parse error: %s at %s\n", err, configPath)
	}

	commandConfig.Aliases[alias] = expansion

	newData, err := yaml.Marshal(&commandConfig)

	if err != nil {
		log.Fatalf("error marshaling YAML: %v", err)
	}

	// Write the updated YAML data back to the file

	fmt.Printf("- Adding alias for %s: %s\n", alias, expansion)

	err = ioutil.WriteFile(configPath, newData, 0644)
	if err != nil {
		log.Fatalf("error writing file: %v", err)
	}
	fmt.Printf("%s Added alias.\n", color.GreenString("✓"))
}

func DeleteAlias(alias string) error {
	configPath := CommandConfigPath()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		createDefaultCommandConfig()
	}

	configBytes, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}

	var commandConfig CommandConfig
	err = yaml.Unmarshal(configBytes, &commandConfig)

	if err != nil {
		log.Fatalf("parse error: %s at %s\n", err, configPath)
	}

	was, ok := commandConfig.Aliases[alias]

	if ok {
		delete(commandConfig.Aliases, alias)
		newData, err := yaml.Marshal(&commandConfig)

		if err != nil {
			log.Fatalf("error marshaling YAML: %v", err)
		}

		// Write the updated YAML data back to the file
		err = ioutil.WriteFile(configPath, newData, 0644)
		if err != nil {
			log.Fatalf("error writing file: %v", err)
		}

		fmt.Printf("%s Deleted alias %s; was %s\n", color.RedString("✓"), alias, was)

	} else {
		fmt.Println("target alias not found")
		os.Exit(1)
	}

	return nil
}

func createDefaultCommandConfig() {
	file, err := os.Create(CommandConfigPath())

	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write([]byte{})

	if err != nil {
		panic(err)
	}
}

func createDefaultHostsConfig() {
	userConfigDir, err := os.UserConfigDir()

	if err != nil {
		panic(err)
	}

	appConfigDir := filepath.Join(userConfigDir, "bl")

	if _, err := os.Stat(appConfigDir); os.IsNotExist(err) {
		os.MkdirAll(appConfigDir, 0755)
	}

	file, err := os.Create(HostsConfigPath())

	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write([]byte{})

	if err != nil {
		panic(err)
	}
}

func CommandConfigPath() string {
	configDir, err := os.UserConfigDir()

	if err != nil {
		panic(err)
	}

	return filepath.Join(configDir, "bl", "config.yml")
}

func HostsConfigPath() string {
	configDir, err := os.UserConfigDir()

	if err != nil {
		panic(err)
	}

	return filepath.Join(configDir, "bl", "hosts.yml")
}
