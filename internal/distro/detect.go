package distro

import (
	"bufio"
	"os"
	"strings"
)

// Family represents the detected distro family.
type Family int

const (
	FamilyUnknown Family = iota
	FamilyUbuntu         // Ubuntu and derivatives (Mint, Pop!_OS, elementaryOS, etc.)
	FamilyFedora         // Fedora, RHEL, CentOS
	FamilyArch           // Arch, Manjaro
	FamilyDebian         // Debian (non-Ubuntu)
)

// Info holds the parsed fields from /etc/os-release.
type Info struct {
	ID     string // e.g. "ubuntu", "fedora", "arch"
	IDLike string // e.g. "ubuntu debian"
	Name   string // e.g. "Ubuntu 24.04 LTS"
	Family Family
}

// Detect reads /etc/os-release and returns distro information.
func Detect() Info {
	info := parse("/etc/os-release")
	info.Family = classify(info)
	return info
}

func parse(path string) Info {
	f, err := os.Open(path)
	if err != nil {
		return Info{}
	}
	defer f.Close()

	info := Info{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		value = strings.Trim(value, `"`)
		switch key {
		case "ID":
			info.ID = strings.ToLower(value)
		case "ID_LIKE":
			info.IDLike = strings.ToLower(value)
		case "PRETTY_NAME":
			info.Name = value
		}
	}
	return info
}

func classify(info Info) Family {
	// Check ID and ID_LIKE for family membership.
	fields := strings.Fields(info.ID + " " + info.IDLike)
	for _, f := range fields {
		switch f {
		case "ubuntu":
			return FamilyUbuntu
		case "fedora", "rhel", "centos":
			return FamilyFedora
		case "arch":
			return FamilyArch
		case "debian":
			return FamilyDebian
		}
	}
	return FamilyUnknown
}

// IsUbuntuFamily returns true for Ubuntu and all derivatives that set ID_LIKE=ubuntu.
func (i Info) IsUbuntuFamily() bool {
	return i.Family == FamilyUbuntu
}
