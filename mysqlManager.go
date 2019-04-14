package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

type MysqlManager struct {
	TableResidentalCompound mysqlTable
	TableCorpus             mysqlTable
	TableApartments         mysqlTable
	ViewApartments          mysqlTable
	DB                      *sql.DB
}

func (m *MysqlManager) Connect() error {
	m.TableResidentalCompound = createResidentalCompoundTable(GConfig.MysqlDb, GConfig.MySqlResidentalCompoundTable)
	m.TableCorpus = createCorpusTable(GConfig.MysqlDb, GConfig.MySqlCorpusTable)
	m.TableApartments = createApartmentTable(GConfig.MysqlDb, GConfig.MySqlApartmentTable)
	m.ViewApartments = createApartmentView(GConfig.MysqlDb, GConfig.MySqlApartmentView,
		m.TableResidentalCompound,
		m.TableCorpus,
		m.TableApartments)

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
	err = checkAndCreateTable(db, m.TableApartments)
	if err != nil {
		return err
	}
	err = checkAndCreateTable(db, m.ViewApartments)
	if err != nil {
		return err
	}
	m.DB = db
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

func (m *MysqlManager)searchApartments(req ApartmentSearchRequest) ([]Apartment, error) {
	var whereConditions []string
	if req.City != "" {
		condition := fmt.Sprintf("city like '%%%v%%'", req.City)
		whereConditions = append(whereConditions, condition)
	}
	if req.District != "" {
		condition := fmt.Sprintf("district like '%%%v%%'", req.District)
		whereConditions = append(whereConditions, condition)
	}
	if req.Address != "" {
		condition := fmt.Sprintf("address like '%%%v%%'", req.Address)
		whereConditions = append(whereConditions, condition)
	}
	if req.CorpusName != "" {
		condition := fmt.Sprintf("corpus_name like '%%%v%%'", req.CorpusName)
		whereConditions = append(whereConditions, condition)
	}
	appendNumberWhereConditions(req.FloorsCountRange, "floors_count", &whereConditions)
	appendNumberWhereConditions(req.FloorRange, "floor", &whereConditions)
	appendNumberWhereConditions(req.RoomsCountRange, "rooms_count", &whereConditions)
	appendNumberWhereConditions(req.SquareRange, "square", &whereConditions)
	appendNumberWhereConditions(req.CostRange, "cost", &whereConditions)

	if len(whereConditions) == 0 {
		return nil, errors.New("Wrong or empty parameters")
	}
	whereClause := strings.Join(whereConditions, " and ")

	searchStatement := fmt.Sprintf(`select 
city, 
district, 
address, 
residental_compound_name, 
corpus_name, 
floors_count, 
floor, 
rooms_count, 
square, 
cost
from %v where %v;`, m.ViewApartments.FullName, whereClause)


	resultsSql, err := m.DB.Query(searchStatement)
	if err != nil {
		return nil, err
	}

	var results []Apartment
	for resultsSql.Next() {
		var apartment Apartment
		err = resultsSql.Scan(
			&apartment.City,
			&apartment.District,
			&apartment.Address,
			&apartment.ResidentalCompoundName,
			&apartment.CorpusName,
			&apartment.FloorsCount,
			&apartment.Floor,
			&apartment.RoomsCount,
			&apartment.Square,
			&apartment.Cost)
		if err != nil {
			errText := fmt.Sprintf("Error while parsing next item from DB. %v", err)
			return nil, errors.New(errText)
		}
		results = append(results, apartment)
	}
	return results, nil
}

func appendNumberWhereConditions(p NumberSearchParameters, columnName string, conditions *[]string) {
	condition, err := p.getWhereCondition(columnName)
	if err != nil {
		return
	}
	*conditions = append(*conditions, condition)
}