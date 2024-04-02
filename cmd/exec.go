package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/xuender/mass/app"
)

// nolint: gochecknoinits
func init() {
	execCmd := &cobra.Command{
		Use:   "exec",
		Short: "Exec sql",
		Long:  `Exec sql.`,
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				slog.Error("miss dsn")

				return
			case 1:
				slog.Debug("open db", "dsn", args[0])

			default:
				exec(args[0], args[1:]...)
			}
		},
	}

	rootCmd.AddCommand(execCmd)
}

func exec(dsn string, sqls ...string) {
	mass := app.NewApp()

	for _, sql := range sqls {
		mass.Exec(dsn, sql)
	}
}
