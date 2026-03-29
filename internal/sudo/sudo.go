// Package sudo provides a helper for running a command under sudo,
// reading the password via stdin (sudo -S).
package sudo

import (
	"fmt"
	"os/exec"
	"strings"
)

// Run executes args under `sudo -S`, writing password to stdin.
// It returns a descriptive error on failure, including when the
// password is wrong.
func Run(password string, args ...string) error {
	cmdArgs := append([]string{"-S", "--"}, args...)
	cmd := exec.Command("sudo", cmdArgs...)
	cmd.Stdin = strings.NewReader(password + "\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		lower := strings.ToLower(msg)
		if strings.Contains(lower, "incorrect password") ||
			strings.Contains(lower, "sorry") ||
			strings.Contains(lower, "authentication failure") {
			return fmt.Errorf("incorrect password")
		}
		if msg != "" {
			// Strip sudo's "password:" prompt line from output.
			lines := strings.Split(msg, "\n")
			var clean []string
			for _, l := range lines {
				if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(l)), "password") {
					clean = append(clean, l)
				}
			}
			if trimmed := strings.TrimSpace(strings.Join(clean, "\n")); trimmed != "" {
				return fmt.Errorf("%s", trimmed)
			}
		}
		return fmt.Errorf("sudo: %w", err)
	}
	return nil
}
