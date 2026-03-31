package resolver

// AppType classifies an application as GUI or CLI to inform source preference.
type AppType int

const (
	AppTypeUnknown AppType = iota
	AppTypeGUI
	AppTypeCLI
)

// knownCLITools is a set of well-known CLI utilities that should prefer apt over Flatpak.
var knownCLITools = map[string]bool{
	"git": true, "curl": true, "wget": true, "vim": true, "neovim": true,
	"htop": true, "tmux": true, "zsh": true, "fish": true, "bash": true,
	"ffmpeg": true, "imagemagick": true, "python3": true, "node": true,
	"nodejs": true, "npm": true, "ripgrep": true, "fzf": true, "jq": true,
	"tree": true, "rsync": true, "ssh": true, "nmap": true, "strace": true,
}

// Classify returns AppTypeGUI or AppTypeCLI for a given app name.
func Classify(app string) AppType {
	if knownCLITools[app] {
		return AppTypeCLI
	}
	return AppTypeGUI // default: assume GUI desktop app
}

// Decision represents the resolver's chosen source and its rationale.
type Decision struct {
	Source string // "flatpak" | "apt" | "none"
	AppID  string // resolved package/app identifier
	Reason string
}

// PreferredSource enumerates valid user-configured source preferences.
type PreferredSource string

const (
	PreferAuto    PreferredSource = "auto"
	PreferFlatpak PreferredSource = "flatpak"
	PreferNative  PreferredSource = "native"
	PreferSnap    PreferredSource = "snap"
)
