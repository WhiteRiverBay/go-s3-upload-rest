package main

import (
	"bytes"
	"flag"
	"go-s3-upload-rest/server"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// start sever
	log.Println("Starting server on :9090")
	bucket := flag.String("bucket", "test-bucket", "S3 bucket name")
	region := flag.String("region", "us-west-2", "S3 region")
	accessKeyID := flag.String("access", "test-access", "AWS access key ID")
	secretAccessKey := flag.String("secret", "test-secret", "AWS secret access key")
	width := flag.Int("width", 300, "Image width")
	height := flag.Int("height", 300, "Image height")

	flag.Parse()
	go server.StartServer(*bucket, *region, *accessKeyID, *secretAccessKey, *width, *height)

	// sleep 5
	time.Sleep(5 * time.Second)
	code := m.Run()
	os.Exit(code)
}

func TestUpload(t *testing.T) {
	// exePath, err := os.Executable()
	// if err != nil {
	// 	t.Fatalf("failed to get executable path: %v", err)
	// }

	// exePath, err = filepath.EvalSymlinks(exePath)
	// if err != nil {
	// 	t.Fatalf("failed to get executable path: %v", err)
	// }

	// exeDir := filepath.Dir(exePath)
	// log.Println(exeDir)

	exeDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	// test cases
	tests := []struct {
		filename string
		status   int
	}{

		{exeDir + "/../test/beaver.jpg", http.StatusOK},
		{exeDir + "/../test/beaver.png", http.StatusOK},
		{exeDir + "/../test/beaver.gif", http.StatusOK},
		{exeDir + "/../test/beaver.txt", http.StatusBadRequest},
	}
	for i, tt := range tests {
		file, err := os.Open(tt.filename)
		if err != nil {
			t.Fatalf("failed to open file: %v", err)
		}
		defer file.Close()

		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		// not using filepath.Base(tt.filename) because it will return the full path
		fileName := filepath.Base(tt.filename)

		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatalf("failed to copy file: %v", err)
		}
		writer.Close()

		req, err := http.NewRequest("POST", "http://localhost:9090/upload", &buf)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("X-Forwarded-For", "1.1.1."+strconv.Itoa(i))

		rr := httptest.NewRecorder()

		http.DefaultServeMux.ServeHTTP(rr, req)

		if status := rr.Code; status != tt.status {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, tt.status)
		}
	}
}

func TestRateLimit(t *testing.T) {

	exeDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	for i := 0; i < 3; i++ {
		file, err := os.Open(exeDir + "/../test/beaver.jpg")
		if err != nil {
			t.Fatalf("failed to open file: %v", err)
		}
		defer file.Close()

		var buf bytes.Buffer
		writer := multipart.NewWriter(&buf)
		fileName := filepath.Base(exeDir + "/../test/beaver.jpg")

		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			t.Fatalf("failed to copy file: %v", err)
		}
		writer.Close()

		req, err := http.NewRequest("POST", "http://localhost:9090/upload", &buf)
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("X-Forwarded-For", "2.2.2.2")

		rr := httptest.NewRecorder()

		http.DefaultServeMux.ServeHTTP(rr, req)

		if status := rr.Code; i < 2 && status != http.StatusOK || i == 2 && status != http.StatusTooManyRequests {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	}

}
