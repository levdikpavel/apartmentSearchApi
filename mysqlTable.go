package main

import "fmt"

type mysqlTable struct {
	DbName          string
	TableName       string
	CreateStatement string
}

func createResidentalCompoundTable(dbName string, tableName string) mysqlTable {
	createTableTemplate := `CREATE TABLE %v.%v (
  residental_compound_id INT AUTO_INCREMENT,
  city VARCHAR(45),
  district VARCHAR(45),
  address VARCHAR(100),
  residental_compound_name VARCHAR(45),
  PRIMARY KEY (residental_compound_id))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;`
	createTableStatement := fmt.Sprintf(createTableTemplate, dbName, tableName)
	table := mysqlTable{
		DbName:          dbName,
		TableName:       tableName,
		CreateStatement: createTableStatement,
	}
	return table
}
func createCorpusTable(dbName string, tableName string) mysqlTable {
	createTableTemplate := `CREATE TABLE %v.%v (
  corpus_id INT AUTO_INCREMENT,
  residental_compound_id INT,
  corpus_name VARCHAR(10),
  floors_count INT,
  PRIMARY KEY (corpus_id))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;`
	createTableStatement := fmt.Sprintf(createTableTemplate, dbName, tableName)
	table := mysqlTable{
		DbName:          dbName,
		TableName:       tableName,
		CreateStatement: createTableStatement,
	}
	return table
}
func createAppartmentTable(dbName string, tableName string) mysqlTable {
	createTableTemplate := `CREATE TABLE %v.%v (
  appartment_id INT AUTO_INCREMENT,
  corpus_id INT,
  number INT,
  floor INT NULL,
  rooms_count INT,
  square DOUBLE,
  cost DOUBLE,
  PRIMARY KEY (appartment_id))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;`
	createTableStatement := fmt.Sprintf(createTableTemplate, dbName, tableName)
	table := mysqlTable{
		DbName:          dbName,
		TableName:       tableName,
		CreateStatement: createTableStatement,
	}
	return table
}
