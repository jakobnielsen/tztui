//go:build darwin

package tz

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jakobnielsen/tztui/internal/sudo"
)

// GetSystemTZ returns the currently active system timezone IANA identifier.
// It reads the /etc/localtime symlink, which works without root privileges.
// Falls back to systemsetup if the symlink is absent or unresolvable.
func GetSystemTZ() (string, error) {
	// Primary: /etc/localtime -> /var/db/timezone/zoneinfo/<iana> or
	//          /usr/share/zoneinfo/<iana> — no root required.
	if tz, err := tzFromLocaltimeSymlink(); err == nil {
		return tz, nil
	}

	// Fallback: systemsetup (requires root; may fail for regular users).
	out, err := exec.Command("systemsetup", "-gettimezone").Output()
	if err != nil {
		return "", fmt.Errorf("cannot determine system timezone: %w", err)
	}
	line := strings.TrimSpace(string(out))
	_, after, found := strings.Cut(line, "Time Zone: ")
	if !found {
		return "", fmt.Errorf("unexpected systemsetup output: %q", line)
	}
	tz := strings.TrimSpace(after)
	if tz == "" {
		return "", fmt.Errorf("systemsetup returned empty timezone")
	}
	return tz, nil
}

// tzFromLocaltimeSymlink resolves /etc/localtime to an IANA timezone string.
func tzFromLocaltimeSymlink() (string, error) {
	dest, err := os.Readlink("/etc/localtime")
	if err != nil {
		return "", err
	}
	// Possible prefixes on macOS:
	//   /var/db/timezone/zoneinfo/<iana>
	//   /usr/share/zoneinfo/<iana>
	for _, prefix := range []string{
		"/var/db/timezone/zoneinfo/",
		"/usr/share/zoneinfo/",
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
	cmd := exec.Command("systemsetup", "-settimezone", iana)
	out, err := cmd.CombinedOutput()
	msg := strings.TrimSpace(string(out))
	// systemsetup exits 0 even when it lacks permission — check output first.
	if isPermissionMsg(msg) {
		return ErrPermission
	}
	if err != nil {
		if msg != "" {
			return fmt.Errorf("systemsetup -settimezone: %s", msg)
		}
		return fmt.Errorf("systemsetup -settimezone: %w", err)
	}
	return nil
}

// SetSystemTZWithSudo sets the system timezone by re-running the command
// under sudo, reading the password from the provided string via stdin.
func SetSystemTZWithSudo(iana, password string) error {
	return sudo.Run(password, "systemsetup", "-settimezone", iana)
}

func isPermissionMsg(msg string) bool {
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "administrator") ||
		strings.Contains(lower, "permission denied") ||
		strings.Contains(lower, "not allowed")
}
