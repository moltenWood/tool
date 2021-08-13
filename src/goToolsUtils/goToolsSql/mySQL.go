package goToolsSql

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

var DbPool map[string]*sql.DB

const dsn = "root:my36254796@tcp(127.0.0.1:3306)/{[dbName]}?charset=utf8"

func init() {
	emptyDb := connectDb("", dsn)
	DbPool = map[string]*sql.DB{
		"emptyDb": emptyDb,
	}

}

func ConnectDb(dbNames ...string) {
	for _, dbName := range dbNames {
		binanceDB := connectDb(dbName, dsn)
		DbPool[dbName] = binanceDB
	}
}

func connectDb(dbName string, dsn string) *sql.DB {
	newDsn := strings.Replace(dsn, "{[dbName]}", dbName, 1)
	//fmt.Println(newDsn)
	db, err := sql.Open("mysql", newDsn)
	if err != nil {
		panic(errors.New("sql.Open failed : " + err.Error()))
	}
	if len(dbName) != 0 {
		_, err = db.Exec("USE " + dbName)
		if err != nil {
			panic(errors.New("db.Exec failed : " + err.Error()))
		}
	}
	return db
}

func Find_tables(dbName string) map[string]string {
	db := DbPool[dbName]

	sqlStr := `SELECT table_name tableName,TABLE_COMMENT tableDesc
			FROM INFORMATION_SCHEMA.TABLES
			WHERE UPPER(table_type)='BASE TABLE'
			AND LOWER(table_schema) = ?
			ORDER BY table_name asc`

	var result = make(map[string]string)

	rows, err := db.Query(sqlStr, "test_dbName")
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var tableName, tableDesc string
		err = rows.Scan(&tableName, &tableDesc)

		if len(tableDesc) == 0 {
			tableDesc = tableName
		}
		result[tableName] = tableDesc
	}
	return result
}

//删除某个数据库的表
func DropTable(dbName string, tableName string) {
	db := DbPool[dbName]
	//defer db.Close()

	_, err := db.Exec("USE " + dbName)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("drop table " + tableName)
	if err != nil {
		panic(err)
	} else {
		fmt.Print("drop table " + tableName + " successfully\n")
	}
}

//rename table UNIUSDT2021Emas170m to UNIUSDT2021Emas170m1;
func RenameTable(dbName string, tableName string, newName string) {
	db := DbPool[dbName]
	_, err := db.Exec("USE " + dbName)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("rename table " + tableName + " to " + newName)
	if err != nil {
		panic(err)
	} else {
		fmt.Print("drop table " + tableName + " successful")
	}
}

//创建数据库
func CreateDB(dbName string) {
	_, err := DbPool["emptyDb"].Exec("CREATE DATABASE `" + dbName + "`")
	if err != nil {
		panic(errors.New(fmt.Sprintf("new db named %v failed\n", dbName)))
	} else {
		fmt.Printf("new db named %v successfully\n", dbName)
	}
}

//创建表
func CreateTable(dbName string, tableName string) {
	db := DbPool[dbName]

	allTables := GetAllTableName(dbName)
	var tableNameIsNotExisted = true
	for _, table := range allTables {
		if table == tableName {
			tableNameIsNotExisted = false
			break
		}
	}
	if tableNameIsNotExisted {
		_, err := db.Exec("USE " + dbName)
		if err != nil {
			panic(err)
		}

		//fmt.Print("CREATE TABLE test_tableName (id integer, key varchar(32))")
		_, err = db.Exec(`create table if not exists ` + tableName + `(
id bigint primary key AUTO_INCREMENT UNIQUE,
OpenTime bigint NOT NULL DEFAULT 0,
OpenTimeStr varchar(19) DEFAULT 0,
Open decimal(15,8) NOT NULL DEFAULT 0,
High decimal(15,8) NOT NULL DEFAULT 0,
Low decimal(15,8) NOT NULL DEFAULT 0,
Close decimal(15,8) NOT NULL DEFAULT 0,
Volume decimal(22,6) NOT NULL DEFAULT 0,
CloseTime bigint NOT NULL DEFAULT 0,
CloseTimeStr varchar(19) DEFAULT 0,
Amount decimal(25,8) NOT NULL DEFAULT 0,
Count bigint NOT NULL DEFAULT 0,
BUYVolume decimal(22,6) NOT NULL DEFAULT 0,
BUYAmount decimal(25,8) NOT NULL DEFAULT 0,
CreateTime TIMESTAMP NOT NULL DEFAULT current_timestamp ,
UpdateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)`)

		if err != nil {
			panic(err)
		}
		fmt.Printf("new table in %v named: %v create successfully\n", dbName, tableName)

	} else {
		//fmt.Printf("table in %v named: %v is already exisited\n", dbName, tableName)
	}
}

//创建表
func CreateFakeTrades(dbName string, tableName string) {
	db := DbPool[dbName]

	allTables := GetAllTableName(dbName)
	var tableNameIsNotExisted = true
	for _, table := range allTables {
		if table == tableName {
			tableNameIsNotExisted = false
			break
		}
	}
	if tableNameIsNotExisted {
		_, err := db.Exec("USE " + dbName)
		if err != nil {
			panic(err)
		}

		//fmt.Print("CREATE TABLE test_tableName (id integer, key varchar(32))")
		_, err = db.Exec(`create table if not exists ` + tableName + `(
id bigint primary key AUTO_INCREMENT UNIQUE,
Symbol varchar(19) DEFAULT 0,
TimeInterval SMALLINT UNSIGNED NOT NULL DEFAULT 0,
SmallPeriod SMALLINT UNSIGNED NOT NULL DEFAULT 0,
BigPeriod SMALLINT UNSIGNED NOT NULL DEFAULT 0,
TransactTime bigint NOT NULL DEFAULT 0,
TransactTimeStr varchar(19) DEFAULT 0,
Price decimal(15,8) NOT NULL DEFAULT 0,
Type varchar(10) DEFAULT 0,
Side         varchar(10) DEFAULT 0,
CreateTime TIMESTAMP NOT NULL DEFAULT current_timestamp ,
UpdateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)`)

		if err != nil {
			panic(err)
		}
		fmt.Printf("new table in %v named: %v create successfully\n", dbName, tableName)

	} else {
		//fmt.Printf("table in %v named: %v is already exisited\n", dbName, tableName)
	}
}

//创建表
func CreateTableForEmas(dbName string, tableName string, emaPeriodStart int, emaPeriodEnd int, emaPeriodStep int) {
	db := DbPool[dbName]

	allTables := GetAllTableName(dbName)
	var tableNameIsNotExisted = true
	for _, table := range allTables {
		if table == tableName {
			tableNameIsNotExisted = false
			break
		}
	}
	if tableNameIsNotExisted {
		_, err := db.Exec("USE " + dbName)
		if err != nil {
			panic(err)
		}
		var emaPart string
		for period := emaPeriodStart; period < emaPeriodEnd+1; period += emaPeriodStep {
			emaPart += fmt.Sprintf("Ema%d decimal(15,8) NOT NULL DEFAULT 0,", period)
		}
		//fmt.Print("CREATE TABLE test_tableName (id integer, key varchar(32))")
		_, err = db.Exec(`create table if not exists ` + tableName + `(
id bigint primary key AUTO_INCREMENT UNIQUE,
OpenTime bigint NOT NULL DEFAULT 0,
OpenTimeStr varchar(19) DEFAULT 0,` + emaPart +
			`CreateTime TIMESTAMP NOT NULL DEFAULT current_timestamp ,
UpdateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)`)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(fmt.Sprintf("new db named %v in %v successfully\n", tableName, dbName))
		}
	}
}

//创建表
func CreateTableForTradeRecord(dbName string, tableName string) {
	db := DbPool[dbName]

	allTables := GetAllTableName(dbName)
	var tableNameIsNotExisted = true
	for _, table := range allTables {
		if table == tableName {
			tableNameIsNotExisted = false
			break
		}
	}
	if tableNameIsNotExisted {
		_, err := db.Exec("USE " + dbName)
		if err != nil {
			panic(err)
		}

		//fmt.Print("CREATE TABLE test_tableName (id integer, key varchar(32))")
		//EventType option["BUY","SELL","LongStopLoss0","LongStopLoss1","LongStopLoss2","ShortStopLoss0","ShortStopLoss1","ShortStopLoss2"]
		_, err = db.Exec(`create table if not exists ` + tableName + `(
id bigint primary key AUTO_INCREMENT UNIQUE,
Symbol varchar(10) NOT NULL,
OrderId             bigint NOT NULL      ,     
OrderListId         bigint NOT NULL      ,     
ClientOrderId       varchar(25) NOT NULL,
TransactTime        bigint NOT NULL      ,     
Price                decimal(15,8) NOT NULL ,
OrigQty             decimal(15,8) NOT NULL ,      
ExecutedQty          decimal(15,8) NOT NULL ,      
CummulativeQuoteQty  decimal(15,8) NOT NULL ,      
Status              varchar(10) NOT NULL,
TimeInForce         varchar(10) NOT NULL,
Type                varchar(10) NOT NULL,
Side                varchar(10) NOT NULL,
CreateTime TIMESTAMP NOT NULL DEFAULT current_timestamp ,
UpdateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)`)

		if err != nil {
			panic(err)
		}
		fmt.Printf("new table in %v named: %v create successfully\n", dbName, tableName)

	} else {
		//fmt.Printf("table in %v named: %v is already exisited\n", dbName, tableName)

	}
}

func CreateTableForResults(dbName string, tableName string) {
	db := DbPool[dbName]

	allTables := GetAllTableName(dbName)
	var tableNameIsNotExisted = true
	for _, table := range allTables {
		if table == tableName {
			tableNameIsNotExisted = false
			break
		}
	}
	if tableNameIsNotExisted {
		_, err := db.Exec("USE " + dbName)
		if err != nil {
			panic(err)
		}

		//fmt.Print("CREATE TABLE test_tableName (id integer, key varchar(32))")
		_, err = db.Exec(`create table if not exists ` + tableName + `(
id bigint primary key AUTO_INCREMENT UNIQUE,
EndTimeStr varchar(19) DEFAULT 0,
TimeInterval SMALLINT UNSIGNED NOT NULL DEFAULT 0,
SmallPeriod SMALLINT UNSIGNED NOT NULL DEFAULT 0,
BigPeriod SMALLINT UNSIGNED NOT NULL DEFAULT 0,
Value decimal(15,8) NOT NULL DEFAULT 0,
WinRate float(5,4) NOT NULL DEFAULT 0,
CreateTime TIMESTAMP NOT NULL DEFAULT current_timestamp ,
UpdateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)`)

		if err != nil {
			panic(err)
		} else {
			fmt.Println(fmt.Sprintf("new db named %v in %v successfully\n", tableName, dbName))
		}
	}
	//defer db.Close()
	fmt.Printf("new table in %v named: %v create successfully\n", dbName, tableName)
}
