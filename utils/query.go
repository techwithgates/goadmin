package utils

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/techwithgates/goadmin/database"

	"golang.org/x/crypto/bcrypt"
)

var dbContext = context.Background()

// retrieve table names,
func GetTables() []string {
	tableList := []string{}
	rows, err := database.Db.Query(context.Background(), database.RetrieveTableStmt)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tableName string

	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			panic(err)
		}
		tableList = append(tableList, tableName)
	}

	return tableList
}

// fetch PK field of a table
func GetPkField(tableName string) string {
	db := database.Db
	var pkField string

	row := db.QueryRow(dbContext, database.FindPkStmt, tableName, tableName+"_pkey")

	if scanErr := row.Scan(&pkField, nil); scanErr != nil {
		log.Println(scanErr)
		return ""
	}

	return pkField
}

// handle data create operation
func AddData(tableName string, bodyData map[string]interface{}) error {
	db := database.Db
	var fields, values string

	// iterate and format SQL syntax
	for key, val := range bodyData {
		fields += key + ","

		switch v := val.(type) {
		case string:
			fieldType := verifyCharField(tableName, key)

			if fieldType == "password" {
				passwd, err := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)

				if err != nil {
					log.Println(err)
				}

				v = string(passwd)
			}

			values += fmt.Sprintf("'%s',", v)
		case nil:
			values += fmt.Sprintf("%v,", "null")
		default:
			values += fmt.Sprintf("%v,", v)
		}
	}

	// remove commas at the end to prevent sql syntax error
	fields = strings.TrimSuffix(fields, ",")
	values = strings.TrimSuffix(values, ",")

	// execute query to add data
	_, err := db.Exec(dbContext, fmt.Sprintf(database.InsertDataStmt, tableName, fields, values))

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// handle data update operation
func EditData(tableName string, id string, bodyData map[string]interface{}) error {
	db := database.Db
	var values string

	// iterate and format SQL syntax
	for key, val := range bodyData {
		switch v := val.(type) {
		case string:
			fieldType := verifyCharField(tableName, key)

			if fieldType == "password" {
				passwd, err := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)

				if err != nil {
					log.Println(err)
				}

				v = string(passwd)
			}

			values += fmt.Sprintf("%s='%s',", key, v)
		case nil:
			values += fmt.Sprintf("%s=null,", key)
		default:
			values += fmt.Sprintf("%s=%v,", key, v)
		}
	}

	// remove the comma at the end to prevent sql syntax error
	values = strings.TrimSuffix(values, ",")

	// get pk field
	pkField := GetPkField(tableName)

	// execute query to update data
	_, err := db.Exec(dbContext, fmt.Sprintf(database.UpdateDataStmt, tableName, values, pkField), id)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
