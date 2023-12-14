package routes

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"

	"github.com/techwithgates/goadmin/database"
	"github.com/techwithgates/goadmin/utils"
)

type Paginator struct {
	Offset       int
	Objects      []any
	TableName    string
	TotalObjects int
	ShowCheckBox bool
	HasMore      bool
}

var htmlTemplate, baseTemplate *template.Template
var tablesFile, objectsFile []byte
var dbContext = context.Background()
var embedder embed.FS

// variables for pagination
var pkField string
var totalObjects int
var pageLimit int = 20

// var offset = 0

func SetEmbedder(_embedder *embed.FS) {
	embedder = *_embedder

	baseFile, _ := embedder.ReadFile("template/base.html")
	tblFile, _ := embedder.ReadFile("template/tables.html")
	objFile, _ := embedder.ReadFile("template/objects.html")

	tablesFile = tblFile
	objectsFile = objFile

	baseTmpl, _ := template.New("eap").Parse(string(baseFile))
	baseTemplate, _ = baseTmpl.Clone()
}

// function that retrieves & returns table names from the database
func ListTables(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	data := []string{}

	// iterate & append table names
	for _, name := range utils.GetTables() {
		data = append(data, name)
	}

	htmlTemplate, _ = baseTemplate.Parse(string(tablesFile))
	htmlTemplate.Execute(w, data)
}

// function that returns objects of a table
func ListObjects(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	db := database.Db
	tableName := ps.ByName("tableName")

	// get the offset query param
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	// to check and prevent additional queries
	if offset == 0 {
		// verifiy if the table name exists in the database
		err := db.QueryRow(dbContext, database.TableVerifyStmt, tableName).Scan(&tableName)

		if err != nil {
			if err == pgx.ErrNoRows {
				log.Println(err)
				return
			} else {
				log.Println(err)
			}
		}

		// query the total number of objects
		db.QueryRow(dbContext, fmt.Sprintf(database.GetTotalObjStmt, tableName)).Scan(&totalObjects)

		// query the pk field of the table
		_pkField := utils.GetPkField(tableName)

		// save the pk field for efficient data retrieval
		pkField = _pkField
	}

	// select pk fields from the table
	rows, err := db.Query(dbContext, fmt.Sprintf(database.RetrievePksStmt, pkField, tableName, pkField, pageLimit, offset))

	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	idList := []any{}
	var id any

	// iterate the table rows and append data into map
	for rows.Next() {
		err := rows.Scan(&id)

		if err != nil {
			return
		}

		idList = append(idList, id)
	}

	showCb, hasMore := false, false

	// if there are any data record, select all checkbox will be displayed
	if len(idList) > 0 {
		showCb = true
	}

	// check if there more data to load
	if offset+pageLimit < totalObjects {
		hasMore = true
	}

	// set data output
	paginator := Paginator{
		Offset:       offset,
		Objects:      idList,
		TotalObjects: totalObjects,
		ShowCheckBox: showCb,
		TableName:    tableName,
		HasMore:      hasMore,
	}

	if offset > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(paginator)
	} else {
		htmlTemplate, _ = baseTemplate.Parse(string(objectsFile))
		htmlTemplate.Execute(w, paginator)
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
