package app

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/xuender/kit/los"
	"github.com/xuender/mass/pb"
	"gorm.io/gorm"
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

func (p *App) Delete(dsnKey, sql string, notExec bool) {
	slog.Debug("Delete", "dsn", dsnKey, "sql", sql)

	dsn, has := p.cfg.GetDsn()[dsnKey]
	if !has {
		panic(_dsnNotFound + dsnKey)
	}

	gdb := NewDB(dsn)
	table := NewTable(sql)

	one := false

	count := table.Explain(dsn, gdb)
	if count < p.minRows {
		one = true
	}

	// count, _ := table.Count(dsn, gdb).(int64)
	// if !one && count < p.minRows {
	// 	one = true
	// }

	if one {
		fmt.Fprintf(os.Stdout, "%s;\n", sql)

		if !notExec {
			los.Must0(p.exec(gdb, sql))
		}

		return
	}

	delSQL := fmt.Sprintf("%s LIMIT %d", sql, p.minRows)

	slog.Debug("del", "sql", delSQL, "count", count)
	fmt.Fprintf(os.Stdout, "%s;\n", delSQL)

	if !notExec {
		var sum int64

		for {
			row := lo.Must(p.execAffected(gdb, delSQL))
			sum += row
			fmt.Fprintf(os.Stdout,
				"%s Rows Affected: %d, Count: %d, Min Rows: %d\n",
				time.Now().Format("15:04:05"), row, sum, p.minRows)

			if row < p.minRows {
				return
			}
		}
	}
}

func (p *App) Exec(dsnKey, sql string, notExec bool) {
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
		fmt.Fprintf(os.Stdout, "%s;\n", sql)

		if notExec {
			los.Must0(p.exec(gdb, sql))
		}

		return
	}

	pks := table.Pks(dsn, gdb)
	update := fmt.Sprintf("%s ORDER BY `%s` LIMIT %d", sql, strings.Join(pks, "`,`"), p.minRows)
	slog.Debug("exec", "update", update, "count", count)
	fmt.Fprintf(os.Stdout, "%s;\n", update)

	if !notExec {
		los.Must0(p.exec(gdb, update))
	}

	table.Iterator(dsn, gdb, p.minRows, func(where string) error {
		exe := sql
		if strings.Contains(strings.ToLower(exe), "where") {
			exe += where
		} else {
			exe += " WHERE" + where[4:]
		}

		slog.Debug("exec", "update", exe)
		fmt.Fprintf(os.Stdout, "%s;\n", exe)

		if !notExec {
			los.Must0(p.exec(gdb, exe))
		}

		return nil
	})
}

func (p *App) exec(gdb *gorm.DB, sql string) error {
	ctx := gdb.Exec(sql)
	if ctx.Error != nil {
		return ctx.Error
	}

	fmt.Fprintf(os.Stdout, "Rows Affected: %d, Min Rows: %d\n", ctx.RowsAffected, p.minRows)

	return nil
}

func (p *App) execAffected(gdb *gorm.DB, sql string) (int64, error) {
	ctx := gdb.Exec(sql)

	return ctx.RowsAffected, ctx.Error
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
