package adapter

import (
	"fmt"
	"os/exec"
)

const flathubURL = "https://dl.flathub.org/repo/flathub.flatpakrepo"

// EnsureFlatpak installs Flatpak and adds the Flathub remote if not already present.
// Should only be called on non-Ubuntu distros when the resolution decision is Flatpak.
// nativeInstaller is the package manager command to use, e.g. "apt-get" or "dnf".
func EnsureFlatpak(nativeInstaller string) error {
	// Check if flatpak binary already exists.
	if _, err := exec.LookPath("flatpak"); err == nil {
		// Flatpak is installed — ensure Flathub remote is present.
		return ensureFlathub()
	}

	fmt.Println("   Flatpak is not installed. Setting it up now...")

	var installCmd *exec.Cmd
	switch nativeInstaller {
	case "apt-get":
		installCmd = exec.Command("pkexec", "apt-get", "install", "-y", "flatpak")
	case "dnf":
		installCmd = exec.Command("pkexec", "dnf", "install", "-y", "flatpak")
	default:
		return fmt.Errorf("unsupported installer for Flatpak bootstrap: %s", nativeInstaller)
	}

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("failed to install flatpak: %w", err)
	}
	fmt.Println("   ✓  Flatpak installed.")

	if err := ensureFlathub(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("   ⚠  Note: a session restart may be needed before Flatpak GUI apps")
	fmt.Println("      launch correctly from your desktop environment.")
	fmt.Println()
	return nil
}

func ensureFlathub() error {
	// Check if flathub remote is already configured.
	out, _ := exec.Command("flatpak", "remote-list").Output()
	for _, line := range splitLines(string(out)) {
		if line == "flathub" {
			return nil // already present
		}
	}

	fmt.Println("   Adding Flathub repository...")
	cmd := exec.Command("flatpak", "remote-add", "--if-not-exists", "flathub", flathubURL)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add Flathub remote: %w", err)
	}
	fmt.Println("   ✓  Flathub configured.")
	return nil
}

func splitLines(s string) []string {
	var lines []string
	for _, l := range splitNewlines(s) {
		if l != "" {
			lines = append(lines, l)
		}
	}
	return lines
}

func splitNewlines(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		result = append(result, s[start:])
	}
	return result
}
