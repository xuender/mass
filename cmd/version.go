package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// nolint: gochecknoinits
func init() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "version",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(os.Stdout, Version)
		},
	}

	rootCmd.AddCommand(versionCmd)
}
