package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// nolint
var (
	cfgFile string
	//go:embed version.txt
	Version string
	//go:embed build.txt
	BuildTime string
)

type debugHandler struct {
	slog.Handler
}

func (p *debugHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= slog.LevelDebug
}

// nolint
var rootCmd = &cobra.Command{
	Use:   "mass",
	Short: "Mass",
	Long: fmt.Sprintf(`mass

Version: %sBuildTime: %s`, Version, BuildTime),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debug, _ := cmd.Flags().GetBool("debug"); debug {
			logger := slog.New(&debugHandler{Handler: slog.NewTextHandler(os.Stderr, nil)})
			slog.SetDefault(logger)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// nolint: gochecknoinits
func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config (default: $HOME/mass.toml)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("toml")
		viper.SetConfigName("mass")
	}

	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		slog.Info("Load config", "file", viper.ConfigFileUsed())
	}
}
