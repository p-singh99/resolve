package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pawanjot/resolve/internal/adapter"
)

var searchCmd = &cobra.Command{
	Use:   "search <app>",
	Short: "Search for an application across all sources",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app := args[0]
		found := false

		sa := &adapter.SnapAdapter{}
		if sa.Available() {
			if snapName, ok := sa.Search(app); ok {
				version := sa.Version(snapName)
				installed := sa.InstalledVersion(app)
				if installed != "" {
					fmt.Printf("[snap]    %s  (version: %s, installed: %s)\n", snapName, version, installed)
				} else {
					fmt.Printf("[snap]    %s  (version: %s)\n", snapName, version)
				}
				found = true
			}
		}

		fa := &adapter.FlatpakAdapter{}
		if fa.Available() {
			if appID, ok := fa.Search(app); ok {
				version := fa.Version(appID)
				fmt.Printf("[flatpak] %s  (version: %s)\n", appID, version)
				found = true
			}
		}

		aa := &adapter.AptAdapter{}
		if aa.Available() {
			if pkg, ok := aa.Search(app); ok {
				version := aa.Version(pkg)
				if aa.IsSnapTransitional(pkg) {
					fmt.Printf("[apt]     %s  (version: %s) [snap transitional stub — installs via Snap]\n", pkg, version)
				} else {
					fmt.Printf("[apt]     %s  (version: %s)\n", pkg, version)
				}
				found = true
			}
		}

		if !found {
			fmt.Printf("No results found for %q.\n", app)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
