package utils

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/techwithgates/goadmin/config"
)

func UploadFile(bodyData *map[string]interface{}, key string, file *multipart.FileHeader) {
	(*bodyData)[key] = file.Filename
	// join the specified path with the uploaded file name
	filePath := filepath.Join(config.GetMediaPath(), file.Filename)
	uploadedFile, err := file.Open()

	if err != nil {
		log.Println(err)
	}
	defer uploadedFile.Close()

	// create a new empty file
	newFile, err := os.Create(filePath)

	if err != nil {
		log.Println(err)
	}
	defer newFile.Close()

	// copy the content of the file to the empty file
	_, err = io.Copy(newFile, uploadedFile)

	if err != nil {
		log.Println(err)
	}
}
