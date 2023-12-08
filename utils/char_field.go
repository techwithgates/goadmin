package utils

import (
	"admin_panel/database"
	"fmt"
	"log"
	"strings"
)

var charMap = map[string]string{
	"email":    "email",
	"url":      "url",
	"file":     "file",
	"password": "password",
}

func verifyCharField(tableName string, colName string) string {
	db := database.Db
	var colNumber int
	var comment *string

	// query the column number
	row := db.QueryRow(dbContext, fmt.Sprintf(database.GetColNumber, tableName, colName))

	// get the column number value
	err := row.Scan(&colNumber)

	if err != nil {
		log.Println(err)
	}

	// query the comment value of a field
	row = db.QueryRow(dbContext, fmt.Sprintf(database.GetCommentStmt, tableName, colNumber))

	// get the comment value
	err = row.Scan(&comment)

	if err != nil {
		log.Println(err)
	}

	if comment != nil {
		value, ok := charMap[strings.ToLower(*comment)]

		if ok {
			return value
		}
	}

	return "character varying"
}
