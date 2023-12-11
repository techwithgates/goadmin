package routes

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	// "os"
	// "path/filepath"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"

	"github.com/techwithgates/goadmin/database"
	"github.com/techwithgates/goadmin/utils"
)

type Data struct {
	Tables  []string
	Objects []any
}

type Navigator struct {
	TableName       string
	ShowObjectsPath bool
	ShowCheckBox    bool
}

type Output struct {
	Data      Data
	Navigator Navigator
}

// var htmlTemplate, _ = template.ParseFiles("template/base.html", "template/content.html")
var htmlTemplate *template.Template
var dbContext = context.Background()
var embedder embed.FS

func SetEmbedder(_embedder *embed.FS) {
	embedder = *_embedder
	baseFile, err := embedder.ReadFile("template/base.html")
	contentFile, err := embedder.ReadFile("template/content.html")

	if err != nil {
		log.Println(err)
	}

	tmpl, err := template.New("eap").Parse(string(baseFile))
	tmpl, err = tmpl.Clone()
	htmlTemplate, err = tmpl.Parse(string(contentFile))
}

// function that retrieves & returns table names from the database
func ListTables(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := Data{}

	// iterate & append table names
	for _, name := range utils.GetTables() {
		data.Tables = append(data.Tables, name)
	}

	// set the output format
	output := Output{
		Data:      data,
		Navigator: Navigator{ShowObjectsPath: false},
	}

	err := htmlTemplate.Execute(w, output)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// function that returns objects of a table
func ListTableObjects(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.Db
	tableName := ps.ByName("tableName")

	// verifiy if the table name exists in the database
	err := db.QueryRow(dbContext, database.TableVerifyStmt, tableName).Scan(&tableName)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println(err)
			http.NotFound(w, r)
			return
		} else {
			log.Println(err)
		}
	}

	// query the pk field of the table
	pkField := utils.GetPkField(tableName)

	// select pk fields from the table
	rows, queryErr := db.Query(dbContext, fmt.Sprintf(database.RetrievePksStmt, pkField, tableName, pkField))

	if queryErr != nil {
		log.Println(queryErr)
		return
	}
	defer rows.Close()

	idList := []interface{}{}
	var id any

	// iterate the table rows and append data into map
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
			return
		}
		idList = append(idList, id)
	}

	hasMany := false

	// if there are any data record, select all checkbox will be displayed
	if len(idList) > 0 {
		hasMany = true
	}

	// set output format
	output := Output{
		Data:      Data{Objects: idList},
		Navigator: Navigator{TableName: tableName, ShowObjectsPath: true, ShowCheckBox: hasMany},
	}

	err = htmlTemplate.Execute(w, output)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// function that returns form template and handles create operation
func AddObject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		// call function to generate HTML add form template
		formTemplate, err := utils.GenerateAddForm(ps.ByName("tableName"))

		if err != nil {
			log.Println(err)
		}

		// set json response format
		output := struct {
			Content string `json:"title"`
			Form    string `json:"form"`
		}{Content: "Add Form", Form: formTemplate}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(output)
	} else {
		// set form multipart size to 10 MB max
		err := r.ParseMultipartForm(10 << 20)

		if err != nil {
			fmt.Println("Form parsing error!")
		}

		bodyData := make(map[string]interface{})

		// iterate through uploaded media files and save them into file system
		for key, files := range r.MultipartForm.File {
			for _, file := range files {
				utils.UploadFile(&bodyData, key, file)
			}
		}

		// retrieve the json data
		jsonData := r.FormValue("jsonData")

		// append key and value pairs to body data
		err = json.Unmarshal([]byte(jsonData), &bodyData)

		if err != nil {
			log.Println(err)
		}

		// call the function to add data for a specific model
		err = utils.AddData(ps.ByName("tableName"), bodyData)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("Insertion successful!")
	}
}

// function that returns form template and handles update operation
func EditObject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		// call the function to generate HTML edit form template
		formTemplate, err := utils.GenerateEditForm(ps.ByName("tableName"), ps.ByName("id"))

		if err != nil {
			log.Println(err)
		}

		// set json response format
		output := struct {
			Content string `json:"title"`
			Form    string `json:"form"`
		}{Content: "Edit Form", Form: formTemplate}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(output)
	} else {
		// set form multipart size to 10 MB max
		err := r.ParseMultipartForm(10 << 20)

		if err != nil {
			fmt.Println("Form parsing error!")
		}

		// retrieve stringified json data
		jsonData := r.FormValue("jsonData")

		// set variable to append keys & values of formData
		bodyData := make(map[string]interface{})

		// iterate through uploaded media files and save them into file system
		for key, files := range r.MultipartForm.File {
			for _, file := range files {
				utils.UploadFile(&bodyData, key, file)
			}
		}

		// append key and value pairs to body data
		err = json.Unmarshal([]byte(jsonData), &bodyData)

		if err != nil {
			log.Println(err)
		}

		// call the function to update data
		err = utils.EditData(ps.ByName("tableName"), ps.ByName("id"), bodyData)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		json.NewEncoder(w).Encode("Update successful!")
	}
}

// function that handles delete operation
func DeleteObject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.Db
	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println(err)
	}

	// parse json data and assign to a slice of strings
	var objectsIds []string
	err = json.Unmarshal(body, &objectsIds)

	if err != nil {
		log.Println(err)
	}

	// assign variable to append formatted data
	var dataFormat string

	// format to sql syntax
	for _, val := range objectsIds {
		dataFormat += fmt.Sprintf(val + ",")
	}

	// remove comma at the end to prevent sql syntax error
	dataFormat = strings.TrimSuffix(dataFormat, ",")

	// get the pk field of the table
	pkField := utils.GetPkField(ps.ByName("tableName"))

	// execute query to delete data
	_, err = db.Exec(dbContext, fmt.Sprintf(database.DeleteDataStmt, ps.ByName("tableName"), pkField, dataFormat))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("Deletion successful!"))
}
