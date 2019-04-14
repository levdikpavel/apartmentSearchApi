package main

import "fmt"

type mysqlTable struct {
	DbName          string
	TableName       string
	FullName        string
	CreateStatement string
}

func createResidentalCompoundTable(dbName string, tableName string) mysqlTable {
	fullName := fmt.Sprintf("%v.%v", dbName, tableName)
	createTableTemplate := `CREATE TABLE %v (
  residental_compound_id INT NOT NULL AUTO_INCREMENT,
  city VARCHAR(45),
  district VARCHAR(45),
  address VARCHAR(100),
  residental_compound_name VARCHAR(45),
  PRIMARY KEY (residental_compound_id))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;`
	createTableStatement := fmt.Sprintf(createTableTemplate, fullName)
	table := mysqlTable{
		DbName:          dbName,
		TableName:       tableName,
		CreateStatement: createTableStatement,
		FullName:        fullName,
	}
	return table
}
func createCorpusTable(dbName string, tableName string) mysqlTable {
	fullName := fmt.Sprintf("%v.%v", dbName, tableName)
	createTableTemplate := `CREATE TABLE %v (
  corpus_id INT NOT NULL AUTO_INCREMENT,
  residental_compound_id INT,
  corpus_name VARCHAR(10),
  floors_count INT,
  KEY residental_compound_id_key (residental_compound_id),
  PRIMARY KEY (corpus_id))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;`
	createTableStatement := fmt.Sprintf(createTableTemplate, fullName)
	table := mysqlTable{
		DbName:          dbName,
		TableName:       tableName,
		CreateStatement: createTableStatement,
		FullName:        fullName,
	}
	return table
}
func createApartmentTable(dbName string, tableName string) mysqlTable {
	fullName := fmt.Sprintf("%v.%v", dbName, tableName)
	createTableTemplate := `CREATE TABLE %v (
  apartment_id INT NOT NULL AUTO_INCREMENT,
  corpus_id INT,
  floor INT NULL,
  rooms_count INT,
  square DOUBLE,
  cost DOUBLE,
  KEY corpus_id_key (corpus_id),
  PRIMARY KEY (apartment_id))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8;`
	createTableStatement := fmt.Sprintf(createTableTemplate, fullName)
	table := mysqlTable{
		DbName:          dbName,
		TableName:       tableName,
		CreateStatement: createTableStatement,
		FullName:        fullName,
	}
	return table
}

func createApartmentView(dbName string, viewName string,
	tableResidentalCompound mysqlTable,
	tableCorpus mysqlTable,
	tableApartments mysqlTable) mysqlTable {
	fullName := fmt.Sprintf("%v.%v", dbName, viewName)
	createTableTemplate := `CREATE VIEW %v as
select 
  rc.city,
  rc.district,
  rc.address,
  rc.residental_compound_name,

  c.corpus_name,
  c.floors_count,

  a.apartment_id,
  a.floor,
  a.rooms_count,
  a.square,
  a.cost
from %v rc
join %v c on rc.residental_compound_id=c.residental_compound_id
join %v a on c.corpus_id=a.corpus_id;`
	createTableStatement := fmt.Sprintf(createTableTemplate, fullName,
		tableResidentalCompound.FullName,
		tableCorpus.FullName,
		tableApartments.FullName)
	table := mysqlTable{
		DbName:          dbName,
		TableName:       viewName,
		CreateStatement: createTableStatement,
		FullName:        fullName,
	}
	return table
}

