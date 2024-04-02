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
				raw(args[0], args[1:]...)
			}
		},
	}

	rawCmd.Flags().StringP("type", "t", "table", "table, csv")
	rootCmd.AddCommand(rawCmd)
}

func raw(dsn string, sqls ...string) {
	mass := app.NewApp()

	for _, sql := range sqls {
		mass.Raw(dsn, sql)
	}
}
