package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/xuender/kit/los"
	"github.com/xuender/mass/pb"
)

// nolint: gochecknoinits
func init() {
	initCmd := &cobra.Command{
		Use:     "init",
		Aliases: []string{"i"},
		Short:   "initialize the configuration file",
		Long: `Initialize the configuration file

		like:

		`,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := pb.NewConfig()
			if len(cfg.GetDsn()) > 0 && !los.Must(cmd.Flags().GetBool("force")) {
				slog.Error("Config is exists.", "file", viper.ConfigFileUsed())

				return
			}

			if cfg.GetDsn() == nil {
				cfg.Dsn = map[string]string{}
			}

			cfg.Dsn["dbname1"] = "username:password@tcp(host:port)/schema?charset=utf8mb4&parseTime=True&loc=Local"
			cfg.Dsn["dbname2"] = "username:password@tcp(host:port)/schema?charset=utf8mb4&parseTime=True&loc=Local"

			cfg.Save()
		},
	}

	initCmd.Flags().Bool("force", false, "forced to initialize the config")
	rootCmd.AddCommand(initCmd)
}
