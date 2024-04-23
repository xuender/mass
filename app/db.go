package app

import (
	"log"
	"os"
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

func IsNumber(val any) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}
