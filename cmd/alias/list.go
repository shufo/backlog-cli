package alias

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"github.com/shufo/backlog-cli/internal/config"
	"gopkg.in/yaml.v3"
)

func List() {
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

	// sort
	aliases := sortMapStringSliceByKey(commandConfig.Aliases)

	for _, v := range aliases {
		for i, w := range v {
			fmt.Printf("%s: %s\n", i, w)
		}
	}
}

func sortMapStringSliceByKey(m map[string]string) []map[string]string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	sorted := make([]map[string]string, len(keys))
	for i, k := range keys {
		sorted[i] = map[string]string{k: m[k]}
	}
	return sorted
}
