package main

import (
	"flag"
	"go-s3-upload-rest/server"
	"log"
)

func main() {
	bucket := flag.String("bucket", "", "S3 bucket name")
	region := flag.String("region", "", "S3 region")
	accessKeyID := flag.String("access", "", "AWS access key ID")
	secretAccessKey := flag.String("secret", "", "AWS secret access key")
	flag.Parse()

	if *bucket == "" || *region == "" || *accessKeyID == "" || *secretAccessKey == "" {
		log.Fatal("Missing required flags")
	}

	// show start text in console
	log.Println("Starting server on :9090")
	server.StartServer(*bucket, *region, *accessKeyID, *secretAccessKey)
}
