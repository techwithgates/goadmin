package config

import (
	"context"
	"goadmin/database"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

var mediaPath string
var Port int

func ConnectDb(url string) {
	db, dbErr := pgxpool.New(context.Background(), url)

	if dbErr != nil {
		log.Fatal(dbErr)
	} else {
		database.Db = db
	}
}

func DefineMediaPath(path string) {
	if _, err := os.Stat(path); err != nil {

		if os.IsNotExist(err) {
			workingDir, _ := os.Getwd()
			mediaPath = filepath.Join(workingDir, "media")
			err = os.Mkdir(mediaPath, os.ModeDir)

			if err != nil {
				log.Println(err)
			}

		} else {
			log.Println(err)
		}

	} else {
		mediaPath = path
	}
}

func SetPort(port int) { Port = port }

func GetMediaPath() string { return mediaPath }
