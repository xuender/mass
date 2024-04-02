package cmd

import (
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

// nolint
var rootCmd = &cobra.Command{
	Use:   "mass",
	Short: "Mass",
	Long: fmt.Sprintf(`mass

Version: %sBuildTime: %s`, Version, BuildTime),
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
