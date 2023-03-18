package alias

import "github.com/shufo/backlog-cli/internal/config"

func Delete(alias string) {
	config.DeleteAlias(alias)
}
