package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/xuender/mass/app"
	"github.com/xuender/mass/pb"
)

// nolint: gochecknoinits
func init() {
	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Dsn check",
		Long:  `Dsn check`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := pb.NewConfig()

			for key, url := range cfg.GetDsn() {
				slog.Info("check", "name", key)
				app.NewDB(url)
			}
		},
	}

	rootCmd.AddCommand(checkCmd)
}
