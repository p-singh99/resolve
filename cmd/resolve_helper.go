package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/pawanjot/resolve/internal/resolver"
)

// resolvePreference returns the effective PreferredSource from the flag, config, or default.
func resolvePreference(flag string) resolver.PreferredSource {
	prefer := resolver.PreferredSource(flag)
	if prefer == "" {
		prefer = resolver.PreferredSource(viper.GetString("preferred_source"))
	}
	if prefer == "" {
		prefer = resolver.PreferAuto
	}
	return prefer
}

// printDecision prints the resolution decision. If dryRun is true it shows
// what would happen without installing; if false it shows what is about to happen.
func printDecision(app string, decision resolver.Decision, dryRun bool) {
	if decision.Source == "none" {
		fmt.Fprintf(os.Stderr, "✗  Could not find %q: %s\n", app, decision.Reason)
		return
	}
	if dryRun {
		fmt.Printf("→  Would install %s from %s (%s)\n", decision.AppID, decision.Source, decision.Reason)
		fmt.Println("   No changes made.")
	} else {
		fmt.Printf("→  Installing %s from %s (%s)\n", decision.AppID, decision.Source, decision.Reason)
	}
}
