package server

import (
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func StartServer(bucket string, region string, accessKeyID string, secretAccessKey string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// upload file
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		// CORS header
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// OPTION
		if r.Method == "OPTIONS" {
			return
		}

		if r.Method == "POST" {
			// file max size 1mb, if file size > 1mb, return error
			err := r.ParseMultipartForm(1 << 20)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			// get file stream from request
			file, fileHeader, err := r.FormFile("file")
			if err != nil {
				http.Error(w, "Unable to get file from form", http.StatusBadRequest)
				return
			}
			uploader := NewS3Uploader(
				bucket,
				region,
				accessKeyID,
				secretAccessKey,
			)
			// set filename to yyyyMMddHH + uuid + file extension
			// format current time to yyyyMMddHH
			currentTime := time.Now().Format("2006010215")
			fileHeader.Filename = "/assets/" + currentTime + "/" + uuid.New().String() + filepath.Ext(fileHeader.Filename)

			err = uploader.UploadFile(file, fileHeader)
			if err != nil {
				http.Error(w, "Failed to upload file", http.StatusInternalServerError)
				log.Printf("Failed to upload file: %v", err)
				return
			}
			// return file url
			w.Write([]byte(fileHeader.Filename))
		}
	})

	http.ListenAndServe(":9090", nil)
}
