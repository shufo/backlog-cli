package alias

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/shufo/backlog-cli/config"
	"gopkg.in/yaml.v3"
)

func FindAlias(alias string) (string, error) {
	configPath := config.CommandConfigPath()

	configBytes, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}

	var commandConfig config.CommandConfig

	err = yaml.Unmarshal(configBytes, &commandConfig)

	if err != nil {
		log.Fatalf("parse error: %s at %s\n", err, configPath)
	}

	expansion, ok := commandConfig.Aliases[alias]

	if ok {
		return expansion, nil
	}

	return "", fmt.Errorf("alias not found")
}
