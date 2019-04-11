package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MysqlManager struct {
}

func (m *MysqlManager) Connect() error {
	config := mysql.Config{
		Net:                  "tcp",
		Addr:                 GConfig.MysqlHost,
		// DBName:               GConfig.MysqlDb,
		User:                 GConfig.MysqlUser,
		Passwd:               GConfig.MysqlPassword,
		AllowNativePasswords: true,
	}
	connectionString := config.FormatDSN()
	log.Printf("Connecting to Server %v", connectionString)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	err = createDatabase(db, GConfig.MysqlDb)
	if err != nil {
		return err
	}
	err = checkAndCreateTable(db, GConfig.MySqlAppartmentsTable)
	if err != nil {
		return err
	}
	return nil
}

func createDatabase(conn *sql.DB, dbName string) error {
	createDatabaseStatement := fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %v DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`, dbName)
	_, err := conn.Exec(createDatabaseStatement)
	if err != nil {
		return err
	}
	useDbStatement := fmt.Sprintf("use %v;", dbName)
	_, err = conn.Exec(useDbStatement)
	if err != nil {
		return err
	}
	log.Printf("Connected to DB '%v'", dbName)
	return nil
}

func checkAndCreateTable(db *sql.DB, tableName string) error {
	// Check if table exists
	err := checkTable(db, tableName)
	if err != nil {
		// Create if not exist
		log.Printf("Table '%v' is missing. Create table...", tableName)
		err = createTable(db,tableName)
		if err != nil {
			return err
		}
	}
	// Check if table exist now
	err = checkTable(db, tableName)
	if err != nil {
		return err
	}
	log.Printf("Table '%v' exists", tableName)
	return nil
}

func checkTable(db *sql.DB, tableName string) error {
	query := fmt.Sprintf(`SELECT 'something' FROM %v LIMIT 1;`, tableName)
	results, err := db.Query(query)
	if err != nil {
		return err
	}
	type somethingType struct {
		something string
	}
	for results.Next() {
		var smth somethingType
		err = results.Scan(&smth.something)
		if err != nil {
			return err
		}
	}
	return nil
}

func createTable(db *sql.DB, tableName string) error {
	createTableStatement := fmt.Sprintf(
`CREATE TABLE realty.%v (
  id INT AUTO_INCREMENT,
  city VARCHAR(45) NULL,
  district VARCHAR(45) NULL,
  address VARCHAR(100) NULL,
  residental_compound VARCHAR(45) NULL,
  corpus VARCHAR(10) NULL,
  floors_count INT NULL,
  floor INT NULL,
  rooms_count INT NULL,
  square DOUBLE NULL,
  cost DOUBLE NULL,
  PRIMARY KEY (id))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;
`, tableName)
	_, err := db.Exec(createTableStatement)
	if err != nil {
		return err
	}
	log.Printf("Table '%v' is created", tableName)
	return nil
}
