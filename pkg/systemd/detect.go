//go:build linux

package systemd

import (
	"bytes"
	"fmt"
	"os"
)

var isSystemd bool

func init() {
	ppid := os.Getppid()

	cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", ppid))
	if err != nil {
		return
	}

	fields := bytes.Split(cmdline, []byte{'\x00'})
	if len(fields) == 0 {
		return
	}

	if string(fields[0]) == "/usr/lib/systemd/systemd" {
		isSystemd = true
	}
}

func IsSystemd() bool {
	return isSystemd
}
