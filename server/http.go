package server

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/WhiteRiverBay/go-s3-upload-rest/util"

	"github.com/google/uuid"
)

func StartServer(bucket string,
	region string,
	accessKeyID string,
	secretAccessKey string, height int, width int, bind string,
	dailyLimit int, minuteLimit int) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	rl := util.NewRateLimiter(minuteLimit, dailyLimit)

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
			// TODO rateof limit 1 IP 1 request per 60 second
			ipHeader := r.Header.Get("X-Forwarded-For")
			ip := strings.Split(ipHeader, ",")[0]

			if !rl.Allow(ip) {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				// log ip
				log.Printf("Rate limit exceeded for IP: %s", ip)
				return
			}

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
			defer file.Close()

			// check if file is supported format
			_, err = util.IsSupportedFormat(file)
			if err != nil {
				http.Error(w, "Unsupported file format", http.StatusBadRequest)
				return
			}

			// get dimensions of the image
			rWidth, rHeight, err := util.GetImageDimensions(file)
			if err != nil {
				http.Error(w, "Failed to get image dimensions", http.StatusInternalServerError)
				log.Printf("Failed to get image dimensions: %v", err)
				return
			}
			var uploadData io.ReadSeeker

			uploadData = file
			// resize image if it's larger than the specified dimensions
			// gif ignore
			if (rWidth > width || rHeight > height) && filepath.Ext(fileHeader.Filename) != ".gif" {
				var contentLength int64
				uploadData, contentLength, err = util.ResizeImage(file, width, height, filepath.Ext(fileHeader.Filename))
				if err != nil {
					http.Error(w, "Failed to resize image", http.StatusInternalServerError)
					log.Printf("Failed to resize image: %v", err)
					return
				}
				fileHeader.Header.Set("Content-Length", strconv.Itoa(int(contentLength)))
				// how to change file size in fileHeader?
				fileHeader.Size = contentLength
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
			err = uploader.UploadFile(uploadData, fileHeader)
			if err != nil {
				http.Error(w, "Failed to upload file", http.StatusInternalServerError)
				log.Printf("Failed to upload file: %v", err)
				return
			}
			// return file url
			w.Write([]byte(fileHeader.Filename))
		}
	})

	http.ListenAndServe(bind, nil)
}
