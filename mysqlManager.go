package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlManager struct {
}

func (m *MysqlManager) Connect() error {
	config := mysql.Config{
		Net:    "tcp",
		Addr:   GConfig.MysqlHost,
		DBName: GConfig.MysqlDb,
		User:   GConfig.MysqlUser,
		Passwd: GConfig.MysqlPassword,
	}
	connectionString := config.FormatDSN()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}
