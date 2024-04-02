package app

import (
	"log"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/xuender/kit/los"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(dsn string) *gorm.DB {
	return los.Must(gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:   true,
		DisableNestedTransaction: true,
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}),
	}))
}

func GetSchema(sql string) string {
	if table := GetAllTable(sql); strings.Contains(table, ".") {
		ts := strings.Split(GetAllTable(sql), ".")

		return ts[0]
	}

	return ""
}

func GetAllTable(sql string) string {
	sql = strings.ToLower(sql)
	sql = strings.ReplaceAll(sql, "`", "")

	if strings.Contains(sql, "select") || strings.Contains(sql, "delete") {
		next := false
		for _, s := range strings.Split(sql, " ") {
			if next {
				return s
			}

			if s == "from" {
				next = true
			}
		}
	}

	return ""
}

func GetTable(sql string) string {
	ts := strings.Split(GetAllTable(sql), ".")

	return ts[len(ts)-1]
}

type Column struct {
	OrdinalPosition int32  `gorm:"column:ORDINAL_POSITION"`
	ColumnName      string `gorm:"column:COLUMN_NAME"`
	DataType        string `gorm:"column:DATA_TYPE"`
	ColumnKey       string `gorm:"column:COLUMN_KEY"`
}

func GetColumns(sql, dsn string, gormDB *gorm.DB) []string {
	schema := GetSchema(sql)
	table := GetTable(sql)
	columns := []*Column{}

	if schema == "" {
		s := strings.Index(dsn, "/") + 1
		e := strings.Index(dsn, "?")
		schema = dsn[s:e]
	}

	slog.Debug("colums", "schema", schema, "table", table)
	// nolint: lll
	gormDB.Raw(
		"SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY, ORDINAL_POSITION FROM `information_schema`.`COLUMNS` WHERE `TABLE_SCHEMA` = ? AND `TABLE_NAME` = ? ORDER BY ORDINAL_POSITION",
		schema, table,
	).Scan(&columns)

	ret := make([]string, len(columns))

	for idx, c := range columns {
		ret[idx] = c.ColumnName
	}

	slog.Debug("columns", "schema", schema, "table", table, "columns", ret)

	return ret
}
