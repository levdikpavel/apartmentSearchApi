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

func (m *MysqlManager) insertResidentalCompoundIfMissing (req Apartment) (int,error) {
	var whereConditions []string
	if req.City != "" {
		condition := fmt.Sprintf("city='%v'", req.City)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'city' is empty")
	}
	if req.District != "" {
		condition := fmt.Sprintf("district='%v'", req.District)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'district' is empty")
	}
	if req.Address != "" {
		condition := fmt.Sprintf("address='%v'", req.Address)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'address' is empty")
	}
	if req.ResidentalCompoundName != "" {
		condition := fmt.Sprintf("residental_compound_name='%v'", req.ResidentalCompoundName)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'residental_compound_name' is empty")
	}
	whereClause := strings.Join(whereConditions, " and ")
	countResidentalCompoundStatement := fmt.Sprintf("select count(1) as count from %v where %v;", m.TableResidentalCompound.FullName, whereClause)
	count, err := m.getMysqlCount(countResidentalCompoundStatement)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		insertStatement := fmt.Sprintf("insert into %v (city,district,address,residental_compound_name) VALUES ('%v','%v','%v','%v')",
			m.TableResidentalCompound.FullName, req.City, req.District, req.Address, req.ResidentalCompoundName)
		insert, err := m.DB.Query(insertStatement)
		if err != nil {
			return 0, err
		}
		defer insert.Close()
	}
	selectStatement := fmt.Sprintf("select residental_compound_id from %v where %v;", m.TableResidentalCompound.FullName, whereClause)
	var apartment Apartment
	err = m.DB.QueryRow(selectStatement).Scan(&apartment.ResidentalCompoundId)
	if err != nil {
		return 0, err
	}

	return apartment.ResidentalCompoundId, nil
}
func (m *MysqlManager) insertCorpusIfMissing (req Apartment) (int, error) {
	var whereConditions []string
	if req.ResidentalCompoundId != 0 {
		condition := fmt.Sprintf("residental_compound_id=%v", req.ResidentalCompoundId)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'residental_compound_id' is zero")
	}
	if req.CorpusName != "" {
		condition := fmt.Sprintf("corpus_name='%v'", req.CorpusName)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'corpus_name' is empty")
	}
	whereClause := strings.Join(whereConditions, " and ")
	countCorpuses := fmt.Sprintf("select count(1) as count from %v where %v;", m.TableCorpus.FullName, whereClause)
	count, err := m.getMysqlCount(countCorpuses)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		insertStatement := fmt.Sprintf("insert into %v (residental_compound_id,corpus_name,floors_count) VALUES (%v,'%v',%v)",
			m.TableCorpus.FullName, req.ResidentalCompoundId, req.CorpusName, req.FloorsCount)
		insert, err := m.DB.Query(insertStatement)
		if err != nil {
			return 0, err
		}
		defer insert.Close()
	}
	selectStatement := fmt.Sprintf("select corpus_id from %v where %v;", m.TableCorpus.FullName, whereClause)
	var apartment Apartment
	err = m.DB.QueryRow(selectStatement).Scan(&apartment.CorpusId)
	if err != nil {
		return 0, err
	}
	return apartment.CorpusId, nil
}
func (m *MysqlManager) insertApartmentIfMissing (req Apartment) (int, error) {
	var whereConditions []string
	if req.CorpusId != 0 {
		condition := fmt.Sprintf("corpus_id=%v", req.CorpusId)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'corpus_id' is zero")
	}
	if req.ApartmentName != "" {
		condition := fmt.Sprintf("apartment_name='%v'", req.ApartmentName)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'apartment_name' is empty")
	}
	if req.Floor != 0 {
		condition := fmt.Sprintf("floor=%v", req.Floor)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'floor' is zero")
	}
	if req.RoomsCount != 0 {
		condition := fmt.Sprintf("rooms_count=%v", req.RoomsCount)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'rooms_count' is zero")
	}
	if req.Square != 0 {
		condition := fmt.Sprintf("square=%v", req.Square)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'square' is zero")
	}
	if req.Cost != 0 {
		condition := fmt.Sprintf("cost=%v", req.Cost)
		whereConditions = append(whereConditions, condition)
	} else {
		return 0, errors.New("Request field 'cost' is zero")
	}
	whereClause := strings.Join(whereConditions, " and ")
	countCorpuses := fmt.Sprintf("select count(1) as count from %v where %v;", m.TableApartments.FullName, whereClause)
	count, err := m.getMysqlCount(countCorpuses)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		insertStatement := fmt.Sprintf("insert into %v (corpus_id,apartment_name,floor,rooms_count,square,cost) VALUES (%v,'%v',%v,%v,%v,%v)",
			m.TableApartments.FullName, req.CorpusId, req.ApartmentName, req.Floor, req.RoomsCount, req.Square, req.Cost)
		insert, err := m.DB.Query(insertStatement)
		if err != nil {
			return 0, err
		}
		defer insert.Close()
	}
	selectStatement := fmt.Sprintf("select apartment_id from %v where %v;", m.TableApartments.FullName, whereClause)
	var apartment Apartment
	err = m.DB.QueryRow(selectStatement).Scan(&apartment.ApartmentId)
	if err != nil {
		return 0, err
	}
	return apartment.ApartmentId, nil
}
func (m *MysqlManager) addApartment(req Apartment) (AparmentsApiResponse, error) {
	var result AparmentsApiResponse
	residentalCompoundId,err := m.insertResidentalCompoundIfMissing(req)
	if err != nil {
		return result, err
	}
	req.ResidentalCompoundId = residentalCompoundId

	corpusId,err := m.insertCorpusIfMissing(req)
	if err != nil {
		return result, err
	}
	req.CorpusId = corpusId

	apartmentId,err := m.insertApartmentIfMissing(req)
	if err != nil {
		return result, err
	}
	result.ApartmentId = apartmentId

	return result, nil
}
func (m *MysqlManager)searchApartments(req ApartmentSearchRequest) (AparmentsApiResponse, error) {
	var result AparmentsApiResponse
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
	if req.ResidentalCompoundName != "" {
		condition := fmt.Sprintf("residental_compound_name like '%%%v%%'", req.ResidentalCompoundName)
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
		return result, errors.New("Wrong or empty parameters")
	}
	whereClause := strings.Join(whereConditions, " and ")

	countStatement := fmt.Sprintf("select count(1) as count from %v where %v;", m.ViewApartments.FullName, whereClause)
	rowsCount, err := m.getMysqlCount(countStatement)
	if err != nil {
		return result, err
	}
	result.Count = rowsCount

	var limitClause string
	if req.Limit > 0 {
		limitClause = fmt.Sprintf("limit %v offset %v", req.Limit, req.Offset)
	}
	var orderByClause string
	if req.OrderBy != "" {
		orderByClause = fmt.Sprintf("order by %v", req.OrderBy)
	}
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
from %v where %v %v %v;`, m.ViewApartments.FullName, whereClause, orderByClause, limitClause)


	resultsSql, err := m.DB.Query(searchStatement)
	if err != nil {
		return result, err
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
			return result, errors.New(errText)
		}
		if apartment.CorpusName == "default" {
			apartment.CorpusName = ""
		}
		results = append(results, apartment)
	}
	result.Results = results
	return result, nil
}

func appendNumberWhereConditions(p NumberSearchParameters, columnName string, conditions *[]string) {
	condition, err := p.getWhereCondition(columnName)
	if err != nil {
		return
	}
	*conditions = append(*conditions, condition)
}

func (m *MysqlManager) getMysqlCount(countStatement string) (int, error) {
	var countStruct CountStruct
	err := m.DB.QueryRow(countStatement).Scan(&countStruct.Count)
	if err != nil {
		return 0, err
	}
	return countStruct.Count, nil
}
