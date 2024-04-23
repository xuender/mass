package app

import (
	"fmt"
	"log/slog"

	"github.com/xuender/kit/los"
	"gorm.io/gorm"
)

type explain struct {
	ID           int32  `gorm:"column:id"`
	SelectType   string `gorm:"column:select_type"`
	Table        string `gorm:"column:table"`
	PossibleKeys string `gorm:"column:possible_keys"`
	Key          string `gorm:"column:key"`
	KeyLen       int    `gorm:"column:key_len"`
	Ref          string `gorm:"column:ref"`
	Rows         int64  `gorm:"column:rows"`
	Extra        string `gorm:"column:Extra"`
}

func (p *Table) Explain(dsn string, gdb *gorm.DB) int64 {
	exp := fmt.Sprintf("EXPLAIN SELECT * FROM %s.%s", p.Schema(dsn), p.Table())

	if where := p.Where(); where != "" {
		exp += " WHERE " + where
	}

	explains := []explain{}
	los.Must0(gdb.Raw(exp).Scan(&explains).Error)

	var rows int64

	for _, exp := range explains {
		rows += exp.Rows
	}

	slog.Debug("explain", "rows", rows, "explains", explains, "sql", p.sql, "exp", exp)

	return rows
}
