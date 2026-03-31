package adapter

import (
	"fmt"
	"os/exec"
	"strings"
)

// SnapAdapter handles search and install via Snap.
type SnapAdapter struct{}

// Available returns true if snap is installed on the system.
func (s *SnapAdapter) Available() bool {
	_, err := exec.LookPath("snap")
	return err == nil
}

// InstalledVersion returns the installed version of a snap, or "" if not installed.
func (s *SnapAdapter) InstalledVersion(app string) string {
	out, err := exec.Command("snap", "list", app).Output()
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		// snap list output: Name  Version  Rev  Tracking  Publisher  Notes
		if len(fields) >= 2 && fields[0] == app {
			return fields[1]
		}
	}
	return ""
}

// Search returns true if the app has an exact name match in the Snap Store.
// It does NOT fall back to fuzzy/partial matches to avoid routing installs
// to the wrong package (e.g. searching "git" returning "gitkraken").
func (s *SnapAdapter) Search(app string) (snapName string, found bool) {
	out, err := exec.Command("snap", "find", app).Output()
	if err != nil {
		return "", false
	}
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		// snap find output: Name  Version  Publisher  Notes  Summary
		// Skip header line
		if len(fields) == 0 || fields[0] == "Name" {
			continue
		}
		// Only return an exact case-insensitive name match
		if strings.EqualFold(fields[0], app) {
			return fields[0], true
		}
	}
	return "", false
}

// Version returns the store version of a snap.
func (s *SnapAdapter) Version(snapName string) string {
	out, err := exec.Command("snap", "info", snapName).Output()
	if err != nil {
		return "unknown"
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "version:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "version:"))
		}
	}
	return "unknown"
}

// Install installs the given snap. Snap handles privilege escalation internally.
func (s *SnapAdapter) Install(snapName string) error {
	cmd := exec.Command("snap", "install", snapName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("snap install failed: %w", err)
	}
	return nil
}
