package alias

import "github.com/shufo/backlog-cli/config"

func Delete(alias string) {
	config.DeleteAlias(alias)
}
