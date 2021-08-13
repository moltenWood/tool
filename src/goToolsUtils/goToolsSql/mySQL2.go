package goToolsSql

//需要连接指定数据库后的有锁操作
import (
	"errors"
	"fmt"
	"goToolsUtils"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var MyRWLock sync.RWMutex

type Condition interface {
}

// GetMaxID 获得最大用户ID
func GetMaxID(dbName string, tableName string) int {
	MyRWLock.RLock()
	defer func() { MyRWLock.RUnlock() }()
	var maxUserID int
	err := DbPool[dbName].QueryRow("select id from " + tableName + " order by id desc limit 1").Scan(&maxUserID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			maxUserID = 0
		} else {
			panic(err)
		}
	}

	//regularTool.CheckErr(err)
	return maxUserID
}

func GetAllTableName(dbName string) []string {
	MyRWLock.RLock()
	defer func() { MyRWLock.RUnlock() }()
	rows, err := DbPool["emptyDb"].Query("select table_name from information_schema.tables where table_schema='" + dbName + "'")
	if err != nil {
		panic(err)
	}

	columns, _ := rows.Columns()

	length := len(columns)

	var tablesName []string
	//遍历返回结果
	for rows.Next() {
		result := make([]interface{}, length)

		for i := 0; i < length; i++ {
			result[i] = new(string)
		}
		if err := rows.Scan(result...); err != nil {
			fmt.Println("Scan=>", err)
			continue
		}
		for i := 0; i < length; i++ {
			tablesName = append(tablesName, *result[i].(*string))
		}
	}
	return tablesName
}

// 所有字段都获取
func FindRecords(dbName string, tableName string, conditionMap map[string]string, limit int, orderBy string, isASC bool, collumNames []string) []map[string]string {
	MyRWLock.RLock()
	defer func() { MyRWLock.RUnlock() }()
	var wherePart string
	if conditionMap == nil {
		wherePart = ""
	} else {
		var added = false
		for key, value := range conditionMap {
			if strings.HasPrefix(value, "gt:") {
				wherePart += " `" + key + "` > '" + value[3:] + "' and"
			} else if strings.HasPrefix(value, "gte:") {
				wherePart += " `" + key + "` >= '" + value[4:] + "' and"
			} else if strings.HasPrefix(value, "lt:") {
				wherePart += " `" + key + "` < '" + value[3:] + "' and"
			} else if strings.HasPrefix(value, "lte:") {
				wherePart += " `" + key + "` <= '" + value[4:] + "' and"
			} else if strings.HasPrefix(value, "e:") {
				wherePart += " `" + key + "` = '" + value[2:] + "' and"
			} else {
				wherePart += " `" + key + "` = '" + value + "' and"
			}
			added = true
		}
		if added {
			wherePart = " WHERE" + wherePart[:len(wherePart)-len(" and")]
		} else {
			wherePart = " WHERE" + wherePart
		}
	}

	var limitPart string
	if limit == -1 {
		limitPart = ""
	} else {
		limitPart = " LIMIT " + strconv.Itoa(limit)
	}

	var orderPart string = " ORDER BY " + orderBy
	if isASC {
		orderPart += " ASC"
	} else {
		orderPart += " DESC"
	}

	selectPart := ""
	if len(collumNames) > 0 {
		for _, collum := range collumNames {
			selectPart += collum + ","
		}
		selectPart = selectPart[:len(selectPart)-1]
	} else {
		selectPart = "*"
	}

	//fmt.Println("SELECT * FROM " + tableName + wherePart + orderPart + limitPart)
	rows, err := DbPool[dbName].Query("SELECT " + selectPart + " FROM " + tableName + wherePart + orderPart + limitPart)
	if err != nil {
		fmt.Println("Query=>", err.Error())
		return nil
	}

	columns, _ := rows.Columns()

	length := len(columns)

	var list []map[string]string

	//遍历返回结果
	for rows.Next() {
		var item = make(map[string]string)
		result := make([]interface{}, length)

		for i := 0; i < length; i++ {
			result[i] = new(string)
		}
		if err := rows.Scan(result...); err != nil {
			fmt.Println("Scan=>", err)
			continue
		}
		for i := 0; i < length; i++ {
			item[columns[i]] = *result[i].(*string)
		}
		list = append(list, item)
	}
	return list
}

// AppendOne 向数据库末端追加一条数据
func AppendOne(dbName string, tableName string, condition Condition) {
	MyRWLock.Lock()
	defer func() { MyRWLock.Unlock() }()
	targetDb := DbPool[dbName]
	var conditionMap map[string]string
	conditionMap = goToolsUtils.StructConvertMapByTag(condition)

	if conditionMap["id"] != "" {
		delete(conditionMap, "id")
	}

	var keys string
	var values string
	//fmt.Print(conditionMap)
	for key, value := range conditionMap {
		keys += "`" + key + "`" + ","
		values += "'" + value + "'" + ","
		//values += "" + value + "" + ","
	}
	keys = keys[:len(keys)-1]
	values = values[:len(values)-1]

	fmt.Println("INSERT INTO " + tableName + " (" + keys + ") VALUES (" + values + ")")
	_, err := targetDb.Exec("INSERT INTO " + tableName + " (" + keys + ") VALUES (" + values + ")")
	if err != nil {
		panic("targetDb.Exec failed err:" + err.Error())
	}
}

// AppendMany 向数据库末端追加多条数据
func AppendMany(dbName string, tableName string, conditions interface{}) {
	MyRWLock.Lock()
	defer func() { MyRWLock.Unlock() }()
	targetDb := DbPool[dbName]

	//事件操作
	conn, err := targetDb.Begin()
	if err != nil {
		panic(errors.New(fmt.Sprintf(" targetDb.Begin() failed:%v ", dbName)))
	}

	t := reflect.TypeOf(conditions)
	v := reflect.ValueOf(conditions)

	length := v.Len()
	kind := t.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			condition := v.Index(i).Interface()

			t2 := reflect.TypeOf(condition).Kind()
			var conditionMap map[string]string
			if t2 == reflect.Struct {
				conditionMap = goToolsUtils.StructConvertMapByTag(condition)
			} else if t2 == reflect.Map {
				conditionMap = condition.(map[string]string)
			}

			if len(conditionMap) == 0 {
				panic("数据为空")
			}

			if conditionMap["Id"] != "" {
				delete(conditionMap, "Id")
			}
			if conditionMap["id"] != "" {
				delete(conditionMap, "id")
			}

			var keys string
			var values string
			//fmt.Print(conditionMap)
			for key, value := range conditionMap {
				keys += "`" + key + "`" + ","
				values += "'" + value + "'" + ","
				//values += "" + value + "" + ","
			}

			keys = keys[:len(keys)-1]
			values = values[:len(values)-1]

			//fmt.Println("INSERT INTO " + tableName + " (" + keys + ") VALUES (" + values + ")")
			_, err := targetDb.Exec("INSERT INTO " + tableName + " (" + keys + ") VALUES (" + values + ")")
			if err != nil {
				fmt.Println("INSERT INTO " + tableName + " (" + keys + ") VALUES (" + values + ")")
				panic("targetDb.Exec failed err:" + err.Error())
			}
		}
		//提交事务
		conn.Commit()
		fmt.Printf("AppendMany tableNmae:%25s len: %v successfully\n", tableName, length)
	} else {
		conn.Rollback()
		panic("必须传入切片/数组")
	}
}

func UpdateOne(dbName string, tableName string, condition Condition, id string) {
	MyRWLock.Lock()
	defer func() { MyRWLock.Unlock() }()

	var conditionMap map[string]string
	conditionMap = goToolsUtils.StructConvertMapByTag(condition)
	if conditionMap["id"] != "" {
		delete(conditionMap, "id")
	}
	if conditionMap["Id"] != "" {
		delete(conditionMap, "Id")
	}

	var keys string
	for key, value := range conditionMap {
		keys += "`" + key + "`" + " = '" + value + "', "
	}
	keys = keys[0 : len(keys)-2]

	//fmt.Println("UPDATE " + tableName + " SET " + keys + " where id = " +id)
	_, err := DbPool[dbName].Exec("UPDATE " + tableName + " SET " + keys + " where id = " + id)
	if err != nil {
		panic(err)
	}
}

// UpdateMany 向数据库多条数据update
//reference 参照列 必须在conditions里面有
func UpdateMany(dbName string, tableName string, conditions interface{}, reference string) {
	MyRWLock.Lock()
	defer func() { MyRWLock.Unlock() }()
	targetDb := DbPool[dbName]

	//事件操作
	conn, err := targetDb.Begin()
	if err != nil {
		panic(errors.New(fmt.Sprintf(" targetDb.Begin() failed:%v ", dbName)))
	}

	t := reflect.TypeOf(conditions)
	v := reflect.ValueOf(conditions)

	kind := t.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			condition := v.Index(i).Interface()
			t2 := reflect.TypeOf(condition).Kind()
			var conditionMap map[string]string
			if t2 == reflect.Struct {
				conditionMap = goToolsUtils.StructConvertMapByTag(condition)
			} else if t2 == reflect.Map {
				conditionMap = condition.(map[string]string)
			}


			if len(conditionMap) == 0 {
				panic("数据为空")
			}
			//fmt.Println(conditionMap)
			if conditionMap[reference] == "" {
				panic(fmt.Sprintf("update needs %v", reference))
			}
			referenceValue := conditionMap[reference]
			delete(conditionMap, reference)

			var keys string
			for key, value := range conditionMap {
				keys += "`" + key + "`" + " = '" + value + "', "
			}
			keys = keys[0 : len(keys)-2]

			//fmt.Println("UPDATE " + tableName + " SET " + keys + " where id = " +id)
			_, err := DbPool[dbName].Exec("UPDATE " + tableName + " SET " + keys + " where `"+ reference +"` = " + referenceValue)
			if err != nil {
				panic(err)
			}
		}
		//提交事务
		conn.Commit()
	} else {
		conn.Rollback()
		panic("必须传入切片/数组")
	}
}

// UpdateMany temp
func UpdateMany2(dbName string, tableName string, conditions interface{}) {
	MyRWLock.Lock()
	defer func() { MyRWLock.Unlock() }()
	targetDb := DbPool[dbName]

	//事件操作
	conn, err := targetDb.Begin()
	if err != nil {
		panic(errors.New(fmt.Sprintf(" targetDb.Begin() failed:%v ", dbName)))
	}

	t := reflect.TypeOf(conditions)
	v := reflect.ValueOf(conditions)

	kind := t.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			condition := v.Index(i).Interface()
			t2 := reflect.TypeOf(condition).Kind()
			var conditionMap map[string]string
			if t2 == reflect.Struct {
				conditionMap = goToolsUtils.StructConvertMapByTag(condition)
			} else if t2 == reflect.Map {
				conditionMap = condition.(map[string]string)
			}


			if len(conditionMap) == 0 {
				panic("数据为空")
			}
			//fmt.Println(conditionMap)
			//if conditionMap["Id"] == "" && conditionMap["id"] == "" {
			//	panic("update needs id")
			//}
			openTime := conditionMap["OpenTime"] 
			//delete(conditionMap, "Id")
			//delete(conditionMap, "id")

			var keys string
			for key, value := range conditionMap {
				keys += "`" + key + "`" + " = '" + value + "', "
			}
			keys = keys[0 : len(keys)-2]

			//fmt.Println("UPDATE " + tableName + " SET " + keys + " where id = " +id)
			_, err := DbPool[dbName].Exec("UPDATE " + tableName + " SET " + keys + " where OpenTime = " + openTime)
			if err != nil {
				panic(err)
			}
		}
		//提交事务
		conn.Commit()
	} else {
		conn.Rollback()
		panic("必须传入切片/数组")
	}
}

// 并将该表格的所有数据 （除了指定需要删除的以外）添加到新的表格，并删除该表格，重命名新表格为该表格
//todo insert 功能没有做
func DeleteAndReplace(dbName string, tableName string, deleteWhere string, creatTable func(dbName string, tableName string)) {
	if len(FindRecords2(dbName, tableName, nil, -1, "id", true, []string{}, deleteWhere))==0{
		fmt.Printf("%v is not needed to DeleteAndReplace\n", tableName)
		return}
	notDeleteWhere := ReverseWhere(deleteWhere)
	allRecords := FindRecords2(dbName, tableName, nil, -1, "id", true, []string{}, notDeleteWhere)

	creatTable(dbName, tableName+"temp")
	AppendMany(dbName, tableName+"temp", allRecords)
	DropTable(dbName, tableName)
	RenameTable(dbName, tableName+"temp", tableName)
}

//反转where条件
func ReverseWhere(extraWhere string) string {
	extraWhere = strings.ReplaceAll(extraWhere, "  ", " ")
	splited := strings.Split(extraWhere, " ")
	var tar = [][]string{{"=", "<>"}, {"<", ">="}, {">", "<="}, {">=", "<"}, {"<=", ">"}, {"and", "or"}, {"or", "and"}}
	output := ""
	for _, part := range splited {
		startIndex := 0
		for _, dui := range tar {
			if strings.Contains(part, dui[0]) {
				startIndex = strings.Index(part, dui[0])
				if dui[0] == "=" && (part[startIndex-1:startIndex] == "<" || part[startIndex-1:startIndex] == ">") {
					continue
				} else if (dui[0] == "<" || dui[0] == ">") && part[startIndex+1:startIndex+2] == "=" {
					continue
				}
				part = strings.Replace(part, dui[0], dui[1], 1)
				break
			} else if strings.Contains(part, dui[1]) {
				startIndex = strings.Index(part, dui[1])
				part = strings.Replace(part, dui[1], dui[0], 1)
				break
			}
		}
		output += part + " "
	}
	return output
}

// 这个可以加额外的where条件 更复杂
func FindRecords2(dbName string, tableName string, conditionMap map[string]string, limit int, orderBy string, isASC bool, collumNames []string, extraWhere string) []map[string]string {
	MyRWLock.RLock()
	defer func() { MyRWLock.RUnlock() }()
	var wherePart string
	if conditionMap == nil {
		wherePart = ""
	} else {
		var added = false
		for key, value := range conditionMap {
			if strings.HasPrefix(value, "gt:") {
				wherePart += " `" + key + "` > '" + value[3:] + "' and"
			} else if strings.HasPrefix(value, "gte:") {
				wherePart += " `" + key + "` >= '" + value[4:] + "' and"
			} else if strings.HasPrefix(value, "lt:") {
				wherePart += " `" + key + "` < '" + value[3:] + "' and"
			} else if strings.HasPrefix(value, "lte:") {
				wherePart += " `" + key + "` <= '" + value[4:] + "' and"
			} else if strings.HasPrefix(value, "e:") {
				wherePart += " `" + key + "` = '" + value[2:] + "' and"
			} else {
				wherePart += " `" + key + "` = '" + value + "' and"
			}
			added = true
		}
		if added {
			wherePart = " WHERE" + wherePart[:len(wherePart)-len(" and")]
		} else {
			wherePart = " WHERE" + wherePart
		}
	}

	if conditionMap == nil {
		wherePart += " where " + extraWhere + " "
	}else {
		wherePart += " and " + extraWhere + " "
	}

	var limitPart string
	if limit == -1 {
		limitPart = ""
	} else {
		limitPart = " LIMIT " + strconv.Itoa(limit)
	}

	var orderPart string = " ORDER BY " + orderBy
	if isASC {
		orderPart += " ASC"
	} else {
		orderPart += " DESC"
	}

	selectPart := ""
	if len(collumNames) > 0 {
		for _, collum := range collumNames {
			selectPart += collum + ","
		}
		selectPart = selectPart[:len(selectPart)-1]
	} else {
		selectPart = "*"
	}

	//fmt.Println("SELECT * FROM " + tableName + wherePart + orderPart + limitPart)
	rows, err := DbPool[dbName].Query("SELECT " + selectPart + " FROM " + tableName + wherePart + orderPart + limitPart)
	if err != nil {
		fmt.Printf("Query=>%v\n%v\n", err.Error(),"SELECT " + selectPart + " FROM " + tableName + wherePart + orderPart + limitPart)
		return nil
	}

	columns, _ := rows.Columns()

	length := len(columns)

	var list []map[string]string

	//遍历返回结果
	for rows.Next() {
		var item = make(map[string]string)
		result := make([]interface{}, length)

		for i := 0; i < length; i++ {
			result[i] = new(string)
		}
		if err := rows.Scan(result...); err != nil {
			fmt.Println("Scan=>", err)
			continue
		}
		for i := 0; i < length; i++ {
			item[columns[i]] = *result[i].(*string)
		}
		list = append(list, item)
	}
	return list
}
