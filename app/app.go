package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/xuender/kit/los"
	"github.com/xuender/mass/pb"
)

const _dsnNotFound = "dsn not found: "

type App struct {
	cfg *pb.Config
}

func NewApp() *App {
	return &App{
		cfg: pb.NewConfig(),
	}
}

func (p *App) Exec(dsn, sql string) {
	slog.Info("Exec", "dsn", dsn, "sql", sql)

	url, has := p.cfg.GetDsn()[dsn]
	if !has {
		panic(_dsnNotFound + dsn)
	}

	gdb := NewDB(url)
	res := gdb.Exec(sql)
	los.Must0(res.Error)
	fmt.Fprintf(os.Stdout, "Rows Affected: %d\n", res.RowsAffected)
}

func (p *App) Raw(dsn, sql string) ([]string, []map[string]any) {
	slog.Info("Raw", "dsn", dsn, "sql", sql)

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

	return GetColumns(sql, dsn, gdb), data
}
