package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/techwithgates/goadmin/config"
	"github.com/techwithgates/goadmin/database"
)

var intField = `
	<div class="form-field">
		<label>%s</label>
		<input name="%s" type="number" value="%v" %s>
	</div>
`

var numericField = `
	<div class="form-field">
		<label>%s</label>
		<input name="%s" type="number" value="%v" step="any" %s>
	</div>
`

var charField = `
	<div class="form-field">
		<label>%s</label>
		<input name="%s" type="text" value="%s" %s>
	</div>
`

var passwdField = `
	<div class="form-field">
		<label>%s</label>
		<input name="%s" type="password" value="%s" %s>
	</div>
`

var emailField = `
	<div class="form-field">
		<label>%s</label>
		<input name="%s" type="email" value="%s" %s>
	</div>
`

var dateField = `
	<div class="form-field">
		<label>%s</label>
		<input name="%s" type="date" value="%s" %s>
	</div>
`

var dateTimeField = `
	<div class="form-field">
		<label>%s</label>
		<input name="%s" type="datetime-local" value="%s" %s>
	</div>
`

var timeField = `
	<div class="form-field">
		<label>%s</label>
		<input name="%s" type="time" value="%s" %s>
	</div>
`

var boolField = `
	<div class="checkbox-field">
		<label>%s</label>
		<input id="%s" name="%s" type="checkbox" %s>
	</div>
`

var textField = `
	<div class="form-field">
		<label>%s</label>
		<textarea name="%s" %s>%s</textarea>
	</div>
`

var urlField = `
	<div class="form-field">
		<label>%s</label>
		<input name="%s" type="url" value="%s" %s>
	</div>
`

var fileField = `
	<div class="form-field">
		<label>%s</label>
		<input name="%s" type="file" %s>
	</div>
`

var editFileField = `
	<div class="form-field">
		<label>%s (Current File: <a href="%s" target="_blank">%s<a>)</label>
		<input name="%s" type="file">
	</div>
`

var fieldMap = map[string]string{
	"character":                   charField,
	"character varying":           charField,
	"smallint":                    intField,
	"integer":                     intField,
	"bigint":                      intField,
	"numeric":                     numericField,
	"text":                        textField,
	"email":                       emailField,
	"url":                         urlField,
	"password":                    passwdField,
	"boolean":                     boolField,
	"date":                        dateField,
	"time without time zone":      timeField,
	"timestamp without time zone": dateTimeField,
	"json":                        textField,
	"jsonb":                       textField,
	"ARRAY":                       textField,
	"file":                        fileField,
	"edit_file":                   editFileField,
}

var dateTimeMap = map[string]int{
	"date":                        1,
	"time without time zone":      1,
	"timestamp without time zone": 1,
}

var nilMap = map[string]string{
	"YES": "",
	"NO":  "required",
}

var boolMap = map[bool]string{
	true:  "checked",
	false: "",
}

func GenerateAddForm(tableName string) (string, error) {
	db := database.Db
	var defaultPk string

	// get the pk field of the table
	pkField := GetPkField(tableName)

	// query the pk type
	row := db.QueryRow(dbContext, database.PkTypeStmt, tableName, pkField)
	row.Scan(nil, &defaultPk)

	// check if the PK field is SERIAL type or not
	if defaultPk == "" {
		pkField = ""
	}

	// fetch the column name, data type and nullable constraint
	rows, err := db.Query(dbContext, database.GetMetaDataStmt, tableName, pkField)

	if err != nil {
		return "", err
	}
	defer rows.Close()

	var colName, dataType, nullable string
	form := fmt.Sprintf(`<form id="formId" onsubmit="submitForm(event, '%s')">`, tableName)

	for rows.Next() {
		// get the column name, data type and nullable constraint
		if err := rows.Scan(&colName, &dataType, &nullable); err != nil {
			return "", err
		}

		_, ok := fieldMap[dataType]

		// if the table field is not supported by EAP, use <input type='text'> by default
		if !ok {
			form += fmt.Sprintf(fieldMap["character varying"], strings.Title(colName), colName, "", nilMap[nullable])
			continue
		}

		if dataType == "boolean" {
			form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, colName, "")
			continue
		}

		if dataType == "text" || dataType == "json" {
			form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, nilMap[nullable], "")
			continue
		}

		if dataType == "character varying" {
			// retrieve the appropriate field type
			dataType = verifyCharField(tableName, colName)

			if dataType == "file" {
				form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, nilMap[nullable])
			} else {
				form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, "", nilMap[nullable])
			}
			continue
		}

		form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, "", nilMap[nullable])
	}

	form += `<div class="button-area"><button type="submit">Add</button></div></form>`
	return form, nil
}

func GenerateEditForm(tableName string, id string) (string, error) {
	db := database.Db
	var defaultPk, colName, dataType, nullable string
	var value any

	// get the pk field of the table
	pkField := GetPkField(tableName)
	_pkField := pkField

	// query the pk type
	row := db.QueryRow(dbContext, database.PkTypeStmt, tableName, pkField)
	row.Scan(nil, &defaultPk)

	// check if the PK field is SERIAL type or not
	if defaultPk == "" {
		_pkField = ""
	}

	// query data type, column name and nullable constraint
	rows, err := db.Query(dbContext, database.GetMetaDataStmt, tableName, _pkField)

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	// construct form with dynamic parameters
	form := fmt.Sprintf(`<form id="formId" onsubmit="submitForm(event, '%s', '%s')">`, tableName, id)

	for rows.Next() {
		// get the column name, data type and nullable constraint
		if err := rows.Scan(&colName, &dataType, &nullable); err != nil {
			log.Println(err)
		}

		// query the value of a single field/column
		row := db.QueryRow(dbContext, fmt.Sprintf(database.GetFieldStmt, colName, tableName, pkField), id)

		// get the value of the field
		if err := row.Scan(&value); err != nil {
			log.Println(err)
		}

		if value == nil {
			value = ""
		}

		_, ok := fieldMap[dataType]

		// if the table field is not supported by EAP, use <input type='text'> by default
		if !ok {
			form += fmt.Sprintf(fieldMap["character varying"], strings.Title(colName), colName, value, nilMap[nullable])
			continue
		}

		_, isDateTime := dateTimeMap[dataType]

		if isDateTime {
			// call the function to convert date, time or datetime to string format
			value = dateTimeToString(dataType, value)
		}

		if result, ok := value.(pgtype.Numeric); ok {
			// call the function to convert numeric type to string format
			value = numericToString(result)
		}

		if dataType == "boolean" {
			if val, ok := value.(bool); ok {
				form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, colName, boolMap[val])
			}
			continue
		}

		if dataType == "ARRAY" {
			value = arrayToString(value)
			form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, nilMap[nullable], value)
			continue
		}

		if dataType == "text" {
			form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, nilMap[nullable], value)
			continue
		}

		if dataType == "json" || dataType == "jsonb" {
			if jsonMap, ok := value.(map[string]any); ok {
				value, err = json.Marshal(jsonMap)
			}

			form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, nilMap[nullable], value)
			continue
		}

		if dataType == "character varying" {
			// retrieve the appropriate char field type
			dataType = verifyCharField(tableName, colName)

			if dataType == "file" {
				mediaUrl := fmt.Sprintf("http://localhost:%d/media/%s", config.Port, value)
				form += fmt.Sprintf(fieldMap["edit_"+dataType], strings.Title(colName), mediaUrl, value, colName)
			} else {
				form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, value, nilMap[nullable])
			}
			continue
		}

		form += fmt.Sprintf(fieldMap[dataType], strings.Title(colName), colName, value, nilMap[nullable])
	}

	form += `<div class="button-area"><button>Edit</button></div></form>`
	return form, nil
}
