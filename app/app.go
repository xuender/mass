package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/xuender/kit/los"
	"github.com/xuender/kit/types"
	"github.com/xuender/mass/pb"
)

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
		panic(fmt.Sprintf("dsn not found: %s", dsn))
	}

	gdb := NewDB(url)
	res := gdb.Exec(sql)
	los.Must0(res.Error)
	fmt.Fprintf(os.Stdout, "Rows Affected: %d\n", res.RowsAffected)
}

func (p *App) Raw(dsn, sql string) {
	slog.Info("Raw", "dsn", dsn, "sql", sql)

	url, has := p.cfg.GetDsn()[dsn]
	if !has {
		panic(fmt.Sprintf("dsn not found: %s", dsn))
	}

	data := []map[string]any{}

	gdb := NewDB(url)
	los.Must0(gdb.Raw(sql).Scan(&data).Error)

	if len(data) == 0 {
		return
	}

	titles := []string{"#"}

	for key := range data[0] {
		titles = append(titles, key)
	}

	vals := [][]string{}

	for idx, row := range data {
		val := make([]string, len(titles))

		for num, key := range titles {
			if num == 0 {
				val[0] = types.Itoa(idx + 1)
			} else {
				val[num] = fmt.Sprintf("%v", row[key])
			}
		}

		vals = append(vals, val)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(titles)
	// footer := make([]string, len(titles))
	// footer[0] = fmt.Sprintf("%d", len(data))
	// table.SetFooter(footer)
	table.AppendBulk(vals)
	table.SetBorder(false)
	table.Render()
}
