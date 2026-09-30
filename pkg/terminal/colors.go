package terminal

import (
	"os"

	"forge.pi.home.arpa/govalds/bot/pkg/systemd"
)

var useColors bool

func init() {
	if os.Getenv("FORCE_COLOR") != "" {
		useColors = true
		return
	}
	if os.Getenv("NO_COLOR") != "" {
		useColors = false
		return
	}
	if os.Getenv("TERM") == "dumb" {
		useColors = false
		return
	}
	useColors = !systemd.IsSystemd()
}

func UseColors(override *bool) bool {
	if override != nil {
		return *override
	}
	return useColors
}
