package app

import (
	"log/slog"

	"github.com/xuender/kit/los"
	"gorm.io/gorm"
)

type column struct {
	OrdinalPosition int32  `gorm:"column:ORDINAL_POSITION"`
	ColumnName      string `gorm:"column:COLUMN_NAME"`
	DataType        string `gorm:"column:DATA_TYPE"`
	ColumnKey       string `gorm:"column:COLUMN_KEY"`
}

func (p *Table) getColumns() []string {
	ret := make([]string, len(p.columns))

	for idx, c := range p.columns {
		ret[idx] = c.ColumnName
	}

	return ret
}

func (p *Table) Columns(dsn string, gormDB *gorm.DB) []string {
	if p.columns != nil {
		return p.getColumns()
	}

	schema := p.Schema(dsn)
	table := p.Table()
	p.columns = []*column{}

	slog.Debug("colums", "schema", schema, "table", table)
	// nolint: lll
	los.Must0(gormDB.Raw(
		"SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY, ORDINAL_POSITION FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? ORDER BY ORDINAL_POSITION",
		schema, table,
	).Scan(&p.columns).Error)

	return p.getColumns()
}
