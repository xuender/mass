package app

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/xuender/kit/los"
	"gorm.io/gorm"
)

type Table struct {
	table    string
	where    string
	schema   string
	allTable string
	sql      string
	columns  []*column
	pks      []string
}

func NewTable(sql string) *Table {
	return &Table{
		allTable: getAllTable(sql),
		sql:      sql,
	}
}

func (p *Table) Min(dsn string, gdb *gorm.DB) any {
	return p.exec("MIN", dsn, gdb)
}

func (p *Table) Max(dsn string, gdb *gorm.DB) any {
	return p.exec("MAX", dsn, gdb)
}

func (p *Table) Count(dsn string, gdb *gorm.DB) any {
	return p.exec("COUNT", dsn, gdb)
}

func (p *Table) exec(cmd, dsn string, gdb *gorm.DB) any {
	key := "*"

	if cmd != "COUNT" {
		pks := p.Pks(dsn, gdb)
		if len(pks) < 1 {
			panic("no primary key")
		}

		key = pks[0]
	}

	schema := p.Schema(dsn)
	table := p.Table()

	sql := fmt.Sprintf("SELECT %s(%s) AS %s FROM `%s`.`%s`", cmd, key, cmd, schema, table)
	slog.Debug(cmd, "sql", sql)

	data := []map[string]any{}

	los.Must0(gdb.Raw(sql).Scan(&data).Error)

	if len(data) == 0 {
		return 0
	}

	return data[0][cmd]
}

func (p *Table) Pks(dsn string, gdb *gorm.DB) []string {
	if p.pks != nil {
		return p.pks
	}

	schema := p.Schema(dsn)
	table := p.Table()
	columns := []*column{}

	slog.Debug("pks", "schema", schema, "table", table)
	// nolint: lll
	los.Must0(gdb.Raw(
		"SELECT COLUMN_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE CONSTRAINT_NAME = 'PRIMARY' AND TABLE_SCHEMA = ? AND TABLE_NAME = ? order by ORDINAL_POSITION",
		schema, table,
	).Scan(&columns).Error)

	p.pks = make([]string, len(columns))
	for idx, column := range columns {
		p.pks[idx] = column.ColumnName
	}

	return p.pks
}

func (p *Table) Iterator(dsn string, gdb *gorm.DB, limit int64, call func(where string) error) {
	pks := p.Pks(dsn, gdb)
	schema := p.Schema(dsn)
	table := p.Table()
	data := []map[string]any{}
	keys := strings.Join(pks, "`,`")

	one := fmt.Sprintf(
		"SELECT `%s` FROM `%s`.`%s` ORDER BY `%s` LIMIT %d,1",
		keys,
		schema, table,
		keys,
		limit,
	)

	slog.Debug("sql", "one", one)
	los.Must0(gdb.Raw(one).Scan(&data).Error)

	for len(data) > 0 {
		row := make([]any, len(pks))
		for idx, key := range pks {
			row[idx] = data[0][key]
		}

		where := IteratorWhere(pks, row)
		if err := call(fmt.Sprintf("%s ORDER BY `%s` LIMIT %d", where, keys, limit)); err != nil {
			slog.Error("call", "err", err)

			return
		}

		next := fmt.Sprintf(
			"SELECT `%s` FROM `%s`.`%s` WHERE %s ORDER BY `%s` LIMIT %d,1",
			keys,
			schema, table,
			where[4:],
			keys,
			limit,
		)

		slog.Debug("iter", "next", next)

		newdata := []map[string]any{}
		los.Must0(gdb.Raw(next).Scan(&newdata).Error)

		data = newdata
	}
}

func IteratorWhere(pks []string, row []any) string {
	list := []string{}

	for idx, key := range pks {
		if idx == len(pks)-1 {
			if IsNumber(row[idx]) {
				list = append(list, fmt.Sprintf(" AND `%s` > %v", key, row[idx]))
			} else {
				list = append(list, fmt.Sprintf(" AND `%s` > '%v'", key, row[idx]))
			}

			break
		}

		if IsNumber(row[idx]) {
			list = append(list, fmt.Sprintf(" AND (`%s` > %v OR `%s` = %v)", key, row[idx], key, row[idx]))
		} else {
			list = append(list, fmt.Sprintf(" AND (`%s` > '%v' OR `%s` = '%v')", key, row[idx], key, row[idx]))
		}
	}

	ret := ""

	for _, str := range list {
		idx := strings.LastIndex(ret, "OR ")
		if idx < 0 {
			ret = str
		} else {
			tmp := ret[:idx+3] + "(" + strings.TrimRight(ret[idx+3:], ")") + str
			left := strings.Count(tmp, "(")
			right := strings.Count(tmp, ")")
			ret = tmp + strings.Repeat(")", left-right)
		}
	}

	return ret
}

func (p *Table) Where() string {
	if p.where != "" {
		return p.where
	}

	idx := strings.LastIndex(strings.ToLower(p.sql), "where")
	if idx < 0 {
		return ""
	}

	idx += 5

	p.where = strings.TrimSpace(p.sql[idx:])

	return p.where
}

func (p *Table) Schema(dsn string) string {
	if p.schema != "" {
		return p.schema
	}

	if strings.Contains(p.allTable, ".") {
		list := strings.Split(p.allTable, ".")

		if list[0] != "" {
			p.schema = list[0]

			return p.schema
		}
	}

	left := strings.Index(dsn, "/") + 1
	right := strings.Index(dsn, "?")
	p.schema = dsn[left:right]

	return p.schema
}

func (p *Table) Table() string {
	if p.table != "" {
		return p.table
	}

	list := strings.Split(p.allTable, ".")
	p.table = list[len(list)-1]

	return p.table
}

func getAllTable(sql string) string {
	sql = strings.ToLower(sql)
	sql = strings.ReplaceAll(sql, "`", "")

	if strings.Contains(sql, "select") || strings.Contains(sql, "delete") {
		next := false
		for _, str := range strings.Split(sql, " ") {
			if next {
				return str
			}

			if str == "from" {
				next = true
			}
		}
	}

	if strings.Contains(sql, "update") {
		table := ""

		for _, str := range strings.Split(sql, " ") {
			if str == "set" {
				return table
			}

			table = str
		}
	}

	return ""
}
