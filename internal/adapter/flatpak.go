package adapter

import (
	"fmt"
	"os/exec"
	"strings"
)

// FlatpakAdapter handles search and install via Flatpak/Flathub.
type FlatpakAdapter struct{}

// Available returns true if flatpak is installed on the system.
func (f *FlatpakAdapter) Available() bool {
	_, err := exec.LookPath("flatpak")
	return err == nil
}

// Search returns true if the app is available on Flathub.
// Returns ("", false) if flatpak is not installed or the app is not found.
func (f *FlatpakAdapter) Search(app string) (appID string, found bool) {
	if !f.Available() {
		return "", false
	}
	out, err := exec.Command("flatpak", "search", "--columns=application", app).Output()
	if err != nil {
		return "", false
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "Application ID" {
			continue
		}
		return line, true
	}
	return "", false
}

// Version returns the remote version string for an app ID from Flathub.
func (f *FlatpakAdapter) Version(appID string) string {
	out, err := exec.Command("flatpak", "remote-info", "flathub", appID, "--columns=version").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(out))
}

// Install installs the given Flatpak app ID from Flathub (user install, no sudo).
func (f *FlatpakAdapter) Install(appID string) error {
	cmd := exec.Command("flatpak", "install", "--noninteractive", "flathub", appID)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("flatpak install failed: %w", err)
	}
	return nil
}
