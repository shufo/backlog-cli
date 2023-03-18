package alias

import "github.com/shufo/backlog-cli/internal/config"

func Set(alias string, expansion string) {
	config.SetAlias(alias, expansion)
}
