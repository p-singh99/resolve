package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pawanjot/resolve/internal/adapter"
	"github.com/pawanjot/resolve/internal/distro"
	"github.com/pawanjot/resolve/internal/resolver"
)

var preferFlag string
var dryRunFlag bool

var installCmd = &cobra.Command{
	Use:   "install <app>",
	Short: "Install an application",
	Long:  `Resolve the best source for an application and install it.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app := args[0]
		prefer := resolvePreference(preferFlag)

		fmt.Printf("Resolving %q...\n", app)
		decision := resolver.Resolve(app, prefer)

		printDecision(app, decision, dryRunFlag)

		if decision.Source == "none" {
			os.Exit(1)
		}
		if dryRunFlag {
			return
		}

		var installErr error
		switch decision.Source {
		case "flatpak":
			fa := &adapter.FlatpakAdapter{}
			if !fa.Available() {
				d := distro.Detect()
				nativeInstaller := "apt-get"
				if d.Family == distro.FamilyFedora {
					nativeInstaller = "dnf"
				}
				if bootstrapErr := adapter.EnsureFlatpak(nativeInstaller); bootstrapErr != nil {
					fmt.Fprintf(os.Stderr, "✗  Could not set up Flatpak: %v\n", bootstrapErr)
					os.Exit(1)
				}
			}
			installErr = fa.Install(decision.AppID)
		case "snap":
			sa := &adapter.SnapAdapter{}
			installErr = sa.Install(decision.AppID)
		case "apt":
			aa := &adapter.AptAdapter{}
			installErr = aa.Install(decision.AppID)
		}

		if installErr != nil {
			fmt.Fprintf(os.Stderr, "✗  Install failed: %v\n", installErr)
			os.Exit(1)
		}

		fmt.Printf("✓  %s installed successfully.\n", app)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVar(&preferFlag, "prefer", "", "prefer a specific source: flatpak | snap | native")
	installCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "show what would be installed without installing")
}
