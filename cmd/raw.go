package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/xuender/mass/app"
)

// nolint: gochecknoinits
func init() {
	rawCmd := &cobra.Command{
		Use:   "raw",
		Short: "Raw sql",
		Long:  `Raw sql`,
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				slog.Error("miss dsn")

				return
			case 1:
				slog.Debug("open db", "dsn", args[0])

			default:
				raw(cmd, args[0], args[1:]...)
			}
		},
	}

	rawCmd.Flags().StringP("type", "t", "grid", "grid, csv, json, toml, yaml")
	rootCmd.AddCommand(rawCmd)
}

func raw(cmd *cobra.Command, dsn string, sqls ...string) {
	mass := app.NewApp(cmd)

	for _, sql := range sqls {
		titles, data := mass.Raw(dsn, sql)

		if name, err := cmd.Flags().GetString("type"); err == nil {
			switch name {
			case "json":
			case "toml":
			case "yaml":
			case "csv":
				app.Csv(titles, data)
			default:
				app.Grid(titles, data)
			}
		}
	}
}
