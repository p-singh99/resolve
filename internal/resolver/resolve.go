package resolver

import (
	"github.com/pawanjot/resolve/internal/adapter"
	"github.com/pawanjot/resolve/internal/distro"
)

// Resolve picks the best package source for the given app name and returns a Decision.
// Resolution is distro-aware: Ubuntu/derivatives prefer Snap; other distros prefer Flatpak.
func Resolve(app string, prefer PreferredSource) Decision {
	d := distro.Detect()

	flatpak := &adapter.FlatpakAdapter{}
	apt := &adapter.AptAdapter{}
	snap := &adapter.SnapAdapter{}

	aptAvail := apt.Available()
	snapAvail := snap.Available()

	// Explicit user preference overrides the algorithm.
	if prefer == PreferSnap {
		if snapAvail {
			if snapName, found := snap.Search(app); found {
				return Decision{Source: "snap", AppID: snapName, Reason: "user preference: --prefer snap"}
			}
		}
		return Decision{Source: "none", Reason: "snap preferred but app not found in Snap Store"}
	}

	if prefer == PreferFlatpak {
		if appID, found := flatpak.Search(app); found {
			return Decision{Source: "flatpak", AppID: appID, Reason: "user preference: --prefer flatpak"}
		}
		return Decision{Source: "none", Reason: "flatpak preferred but app not found on Flathub"}
	}

	if prefer == PreferNative {
		if aptAvail {
			if pkg, found := apt.Search(app); found {
				return Decision{Source: "apt", AppID: pkg, Reason: "user preference: --prefer native"}
			}
		}
		return Decision{Source: "none", Reason: "native preferred but app not found in apt"}
	}

	// Auto mode: apply distro-aware hierarchy.
	appType := Classify(app)

	if d.IsUbuntuFamily() {
		return resolveUbuntu(app, appType, snap, snapAvail, apt, aptAvail, flatpak)
	}
	return resolveGeneric(app, appType, flatpak, apt, aptAvail)
}

// resolveUbuntu applies the Ubuntu hierarchy: Snap → apt → Flatpak (last resort).
func resolveUbuntu(app string, appType AppType, snap *adapter.SnapAdapter, snapAvail bool, apt *adapter.AptAdapter, aptAvail bool, flatpak *adapter.FlatpakAdapter) Decision {
	if appType == AppTypeCLI {
		// CLI tools: apt first, then Snap
		if aptAvail {
			if pkg, found := apt.Search(app); found && !apt.IsSnapTransitional(pkg) {
				return Decision{Source: "apt", AppID: pkg, Reason: "CLI tool — apt preferred for system integration (Ubuntu)"}
			}
		}
		if snapAvail {
			if snapName, found := snap.Search(app); found {
				return Decision{Source: "snap", AppID: snapName, Reason: "CLI tool — not in apt, falling back to Snap (Ubuntu)"}
			}
		}
		return Decision{Source: "none", Reason: "not found in apt or Snap"}
	}

	// GUI app: Snap → apt (skip transitional stubs) → Flatpak
	if snapAvail {
		if snapName, found := snap.Search(app); found {
			return Decision{Source: "snap", AppID: snapName, Reason: "GUI app — Snap preferred on Ubuntu (native universal format)"}
		}
	}
	if aptAvail {
		if pkg, found := apt.Search(app); found && !apt.IsSnapTransitional(pkg) {
			return Decision{Source: "apt", AppID: pkg, Reason: "GUI app — not in Snap Store, falling back to apt (Ubuntu)"}
		}
	}
	// Last resort: Flathub (Flatpak will be bootstrapped at install time if needed)
	if appID, found := flatpak.Search(app); found {
		return Decision{Source: "flatpak", AppID: appID, Reason: "GUI app — not in Snap or apt, falling back to Flatpak (Ubuntu last resort)"}
	}
	return Decision{Source: "none", Reason: "not found in Snap, apt, or Flatpak"}
}

// resolveGeneric applies the non-Ubuntu hierarchy: Flatpak → native (GUI) / native → Flatpak (CLI).
func resolveGeneric(app string, appType AppType, flatpak *adapter.FlatpakAdapter, apt *adapter.AptAdapter, aptAvail bool) Decision {
	if appType == AppTypeCLI {
		if aptAvail {
			if pkg, found := apt.Search(app); found {
				return Decision{Source: "apt", AppID: pkg, Reason: "CLI tool — native apt preferred for system integration"}
			}
		}
		if appID, found := flatpak.Search(app); found {
			return Decision{Source: "flatpak", AppID: appID, Reason: "CLI tool — not in apt, falling back to Flatpak"}
		}
		return Decision{Source: "none", Reason: "not found in apt or Flatpak"}
	}

	// GUI app: Flatpak first (Flatpak will be bootstrapped at install time if needed)
	if appID, found := flatpak.Search(app); found {
		return Decision{Source: "flatpak", AppID: appID, Reason: "GUI app — Flatpak preferred: sandboxed, newer version, cross-distro consistent"}
	}
	if aptAvail {
		if pkg, found := apt.Search(app); found {
			return Decision{Source: "apt", AppID: pkg, Reason: "GUI app — not on Flathub, falling back to apt"}
		}
	}
	return Decision{Source: "none", Reason: "not found in Flatpak or apt"}
}
