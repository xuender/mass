package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xuender/kit/los"
	"github.com/xuender/mass/pb"
)

const _dsnNotFound = "dsn not found: "

type App struct {
	cfg     *pb.Config
	minRows int64
}

func NewApp(cmd *cobra.Command) *App {
	minRows, err := cmd.Flags().GetInt64("min-rows")
	if err != nil {
		minRows = 10_000
	}

	return &App{
		cfg:     pb.NewConfig(),
		minRows: minRows,
	}
}

func (p *App) Exec(dsnKey, sql string) {
	slog.Debug("Exec", "dsn", dsnKey, "sql", sql)

	dsn, has := p.cfg.GetDsn()[dsnKey]
	if !has {
		panic(_dsnNotFound + dsnKey)
	}

	gdb := NewDB(dsn)
	table := NewTable(sql)
	one := false

	if rows := table.Explain(dsn, gdb); rows < p.minRows {
		one = true
	}

	count, _ := table.Count(dsn, gdb).(int64)
	if !one && count < p.minRows {
		one = true
	}

	if one {
		res := gdb.Exec(sql)
		los.Must0(res.Error)
		slog.Debug("exec", "one", one, "min", p.minRows)
		fmt.Fprintf(os.Stdout, "Rows Affected: %d, min: %d\n", res.RowsAffected, p.minRows)

		return
	}

	slog.Debug("exec", "count", count)

	pks := table.Pks(dsn, gdb)
	update := fmt.Sprintf("%s ORDER BY `%s` LIMIT %d", sql, strings.Join(pks, "`,`"), p.minRows)
	slog.Debug("update", "sql", update)

	ctx := gdb.Exec(update)
	if ctx.Error != nil {
		panic(ctx.Error)
	}

	fmt.Fprintf(os.Stdout, "Rows Affected: %d, Min Rows: %d\n", ctx.RowsAffected, p.minRows)

	table.Iterator(dsn, gdb, p.minRows, func(where string) error {
		exe := sql
		if strings.Contains(strings.ToLower(exe), "where") {
			exe += where
		} else {
			exe += " WHERE" + where[4:]
		}

		slog.Debug("row", "where", where, "sql", sql, "exe", exe)

		ctx := gdb.Exec(exe)

		if ctx.Error == nil {
			fmt.Fprintf(os.Stdout, "Rows Affected: %d, Min Rows: %d\n", ctx.RowsAffected, p.minRows)
		}

		return ctx.Error
	})
}

func (p *App) Raw(dsn, sql string) ([]string, []map[string]any) {
	slog.Debug("Raw", "dsn", dsn, "sql", sql)

	url, has := p.cfg.GetDsn()[dsn]
	if !has {
		panic(_dsnNotFound + dsn)
	}

	data := []map[string]any{}

	gdb := NewDB(url)
	los.Must0(gdb.Raw(sql).Scan(&data).Error)

	if len(data) == 0 {
		return nil, nil
	}

	table := NewTable(sql)

	return table.Columns(dsn, gdb), data
}
