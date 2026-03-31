package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pawanjot/resolve/internal/resolver"
)

var dryRunCmd = &cobra.Command{
	Use:   "dry-run <app>",
	Short: "Show how an application would be resolved without installing it",
	Long: `Alias for 'resolve install <app> --dry-run'.
Resolves the best source for an application and prints the decision without installing.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app := args[0]
		prefer := resolvePreference(preferFlag)

		fmt.Printf("Resolving %q...\n", app)
		decision := resolver.Resolve(app, prefer)

		printDecision(app, decision, true)
	},
}

func init() {
	rootCmd.AddCommand(dryRunCmd)
	dryRunCmd.Flags().StringVar(&preferFlag, "prefer", "", "prefer a specific source: flatpak | snap | native")
}
