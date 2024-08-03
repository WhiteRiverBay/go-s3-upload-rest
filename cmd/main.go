package main

import (
	"flag"
	"log"

	"github.com/WhiteRiverBay/go-s3-upload-rest/server"
)

func main() {
	bucket := flag.String("bucket", "", "S3 bucket name")
	region := flag.String("region", "", "S3 region")
	accessKeyID := flag.String("access", "", "AWS access key ID")
	secretAccessKey := flag.String("secret", "", "AWS secret access key")
	width := flag.Int("width", 300, "Image width")
	height := flag.Int("height", 300, "Image height")
	bind := flag.String("bind", ":9090", "Server bind address")
	dailyLimit := flag.Int("daily", 30, "Daily limit")
	minuteLimit := flag.Int("minute", 5, "Minute limit")

	flag.Parse()

	if *bucket == "" || *region == "" || *accessKeyID == "" || *secretAccessKey == "" {
		log.Fatal("Missing required flags")
	}

	// show start text in console
	log.Println("Starting server on :9090")
	server.StartServer(*bucket, *region, *accessKeyID, *secretAccessKey, *width, *height, *bind, *dailyLimit, *minuteLimit)
}
