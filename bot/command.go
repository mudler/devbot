package bot

import (
	"strings"

	util "github.com/mudler/devbot/shared/utils"
)

func HandleCommand(msg, command string, f func(args string)) {
	if strings.Contains(msg, Config.CommandPrefix+command) {
		args := strings.TrimSpace(util.StripPluginCommand(msg, Config.CommandPrefix, command))
		f(args)
	}
}
