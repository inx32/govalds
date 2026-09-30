//go:build !linux

package systemd

func IsSystemd() bool {
	return false
}
