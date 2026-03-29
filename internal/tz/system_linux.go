//go:build linux

package tz

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jakobnielsen/tztui/internal/sudo"
)

// GetSystemTZ returns the currently active system timezone IANA identifier.
// It reads the output of `timedatectl show --property=Timezone --value`.
// Falls back to reading the /etc/localtime symlink if timedatectl is unavailable.
func GetSystemTZ() (string, error) {
	out, err := exec.Command("timedatectl", "show", "--property=Timezone", "--value").Output()
	if err == nil {
		if tz := strings.TrimSpace(string(out)); tz != "" {
			return tz, nil
		}
	}

	// Fallback: /etc/localtime symlink (works without systemd).
	return tzFromLocaltimeSymlink()
}

// tzFromLocaltimeSymlink resolves /etc/localtime to an IANA timezone string.
func tzFromLocaltimeSymlink() (string, error) {
	dest, err := os.Readlink("/etc/localtime")
	if err != nil {
		return "", fmt.Errorf("cannot determine system timezone: %w", err)
	}
	for _, prefix := range []string{
		"/usr/share/zoneinfo/",
		"/usr/lib/zoneinfo/",
		"/usr/share/lib/zoneinfo/",
		"/etc/zoneinfo/",
	} {
		if after, ok := strings.CutPrefix(dest, prefix); ok && after != "" {
			return after, nil
		}
	}
	return "", fmt.Errorf("unrecognised localtime path: %s", dest)
}

// SetSystemTZ sets the system timezone to the given IANA identifier.
// Returns ErrPermission if the process lacks root privileges.
func SetSystemTZ(iana string) error {
	cmd := exec.Command("timedatectl", "set-timezone", iana)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if isPermissionMsg(msg) {
			return ErrPermission
		}
		if msg != "" {
			return fmt.Errorf("timedatectl set-timezone: %s", msg)
		}
		return fmt.Errorf("timedatectl set-timezone: %w", err)
	}
	return nil
}

// SetSystemTZWithSudo sets the system timezone by re-running the command
// under sudo, reading the password from the provided string via stdin.
func SetSystemTZWithSudo(iana, password string) error {
	return sudo.Run(password, "timedatectl", "set-timezone", iana)
}

func isPermissionMsg(msg string) bool {
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "access denied") ||
		strings.Contains(lower, "permission denied") ||
		strings.Contains(lower, "not allowed")
}

// parseTimedatectl parses `timedatectl show` output (key=value pairs) into a map.
// Exported for testing without invoking the real binary.
func parseTimedatectl(output string) map[string]string {
	result := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if k, v, ok := strings.Cut(line, "="); ok {
			result[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	return result
}
