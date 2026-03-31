package adapter

import (
	"fmt"
	"os/exec"
	"strings"
)

// AptAdapter handles search and install via apt on Debian/Ubuntu systems.
type AptAdapter struct{}

// Available returns true if apt is installed on the system.
func (a *AptAdapter) Available() bool {
	_, err := exec.LookPath("apt-cache")
	return err == nil
}

// IsSnapTransitional returns true if the apt package is a transitional stub
// that exists only to install a Snap (version string contains "snap").
// Example: firefox version "1:1snap1-0ubuntu5" is a Snap stub, not a real apt package.
func (a *AptAdapter) IsSnapTransitional(pkg string) bool {
	version := a.Version(pkg)
	return strings.Contains(strings.ToLower(version), "snap")
}

// Search returns true if the package name is available in apt.
// It also returns the exact package name.
func (a *AptAdapter) Search(app string) (pkg string, found bool) {
	out, err := exec.Command("apt-cache", "show", app).Output()
	if err != nil {
		return "", false
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "Package: ") {
			return strings.TrimPrefix(line, "Package: "), true
		}
	}
	return "", false
}

// Version returns the candidate version of the package from apt-cache.
func (a *AptAdapter) Version(pkg string) string {
	out, err := exec.Command("apt-cache", "policy", pkg).Output()
	if err != nil {
		return "unknown"
	}
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Candidate:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Candidate:"))
		}
	}
	return "unknown"
}

// Install installs the given apt package using pkexec (polkit, no raw sudo).
func (a *AptAdapter) Install(pkg string) error {
	cmd := exec.Command("pkexec", "apt-get", "install", "-y", pkg)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("apt install failed: %w", err)
	}
	return nil
}
