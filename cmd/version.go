package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// nolint: gochecknoinits
func init() {
	versionCmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "print Mass version",
		Long:    `Version prints the build information for Mass.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(os.Stdout, Version)

			if build, err := cmd.Flags().GetBool("build"); err == nil && build {
				fmt.Fprint(os.Stdout, BuildTime)
			}
		},
	}

	versionCmd.Flags().BoolP("build", "b", false, "print Mass build time")
	rootCmd.AddCommand(versionCmd)
}
