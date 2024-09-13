package gorm

import (
	"context"
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dbPrefix = "MYSQL_"
)

func NewMySQLDB(ctx context.Context) (*gorm.DB, error) {
	conf, err := config.NewDBConfig(ctx, dbPrefix)
	if err != nil {
		log.Error("Failed to load database config", log.Ferror(err))
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}) // ping is automatically called
	if err != nil {
		return nil, err
	}

	return db, nil
}
