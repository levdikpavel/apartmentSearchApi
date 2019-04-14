package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type MysqlManager struct {
	TableResidentalCompound mysqlTable
	TableCorpus             mysqlTable
	TableAppartments        mysqlTable
}

func (m *MysqlManager) Connect() error {
	m.TableResidentalCompound = createResidentalCompoundTable(GConfig.MysqlDb, GConfig.MySqlResidentalCompoundTable)
	m.TableCorpus = createCorpusTable(GConfig.MysqlDb, GConfig.MySqlCorpusTable)
	m.TableAppartments = createAppartmentTable(GConfig.MysqlDb, GConfig.MySqlAppartmentTable)
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
	err = checkAndCreateTable(db, m.TableResidentalCompound)
	if err != nil {
		return err
	}
	err = checkAndCreateTable(db, m.TableCorpus)
	if err != nil {
		return err
	}
	err = checkAndCreateTable(db, m.TableAppartments)
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

func checkAndCreateTable(db *sql.DB, table mysqlTable) error {
	// Check if table exists
	err := checkTable(db, table)
	if err != nil {
		// Create if not exist
		log.Printf("Table '%v' is missing. Create table...", table.TableName)
		err = createTable(db,table)
		if err != nil {
			return err
		}
	}
	// Check if table exist now
	err = checkTable(db, table)
	if err != nil {
		return err
	}
	log.Printf("Table '%v' exists", table.TableName)
	return nil
}

func checkTable(db *sql.DB, table mysqlTable) error {
	query := fmt.Sprintf(`SELECT 'something' FROM %v.%v LIMIT 1;`, table.DbName, table.TableName)
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

func createTable(db *sql.DB, table mysqlTable) error {
	_, err := db.Exec(table.CreateStatement)
	if err != nil {
		return err
	}
	log.Printf("Table '%v' is created", table.TableName)
	return nil
}
