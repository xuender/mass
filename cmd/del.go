package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/xuender/mass/app"
)

// nolint: gochecknoinits
func init() {
	delCmd := &cobra.Command{
		Use:   "del",
		Short: "批量删除",
		Long:  `含量数据删除.`,
		Run: func(cmd *cobra.Command, args []string) {
			switch len(args) {
			case 0:
				slog.Error("miss dsn")

				return
			case 1:
				slog.Debug("open db", "dsn", args[0])

			default:
				execDel(cmd, args[0], args[1:]...)
			}
		},
	}

	rootCmd.AddCommand(delCmd)
	delCmd.Flags().Int64P("min-rows", "m", _minRows, "min rows")
	delCmd.Flags().BoolP("not-exec", "n", false, "generate sql but do not exec")
}

func execDel(cmd *cobra.Command, dsn string, sqls ...string) {
	mass := app.NewApp(cmd)
	not, _ := cmd.Flags().GetBool("not-exec")

	for _, sql := range sqls {
		mass.Delete(dsn, sql, not)
	}
}
