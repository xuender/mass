package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/xuender/mass/app"
)

const _minRows = 10_000

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
				exec(cmd, args[0], args[1:]...)
			}
		},
	}

	execCmd.Flags().Int64P("min-rows", "m", _minRows, "min rows")
	rootCmd.AddCommand(execCmd)
}

func exec(cmd *cobra.Command, dsn string, sqls ...string) {
	mass := app.NewApp(cmd)

	for _, sql := range sqls {
		mass.Exec(dsn, sql)
	}
}
