package db

import (
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/mestvv/NorthBridgeBackend/internal/config"
)

const DuplicateEntry = 1062

func New(cfg config.Database) (*sqlx.DB, error) {
	location, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("time load location failed: %v", err)
	}
	conf := mysql.NewConfig()
	conf.Net = cfg.Net
	conf.Addr = cfg.Server
	conf.User = cfg.User
	conf.Passwd = cfg.Password
	conf.DBName = cfg.DBName
	conf.Timeout = cfg.Timeout
	conf.Loc = location
	conf.ParseTime = true

	dbConn, err := sqlx.Connect("mysql", conf.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %v", err)
	}

	dbConn.SetMaxIdleConns(cfg.MaxIdleConnections)
	dbConn.SetMaxOpenConns(cfg.MaxOpenConnections)

	if err := dbConn.Ping(); err != nil {
		return nil, err
	}

	return dbConn, nil
}
