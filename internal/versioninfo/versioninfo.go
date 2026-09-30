package versioninfo

import (
	"fmt"
	"strconv"
	"time"
)

var botVersion, buildTimestampText, commit, repository string
var buildTimestamp time.Time

func BotVersion() string   { return botVersion }
func BuildTime() time.Time { return buildTimestamp }
func Commit() string       { return commit }
func Repository() string   { return repository }

func init() {
	if buildTimestampText == "" {
		fmt.Println("linus[init]: warning: buildTimestampText is empty!")
		return
	}

	bt, err := strconv.ParseInt(buildTimestampText, 10, 64)
	if err != nil {
		fmt.Println("linus[init]: warning: unable to parse buildTimestampText as uint64!")
		return
	}

	buildTimestamp = time.Unix(bt, 0)
}
