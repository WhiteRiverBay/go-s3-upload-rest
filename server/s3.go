package server

import (
	"io"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Uploader handles file uploads to S3
type S3Uploader struct {
	s3Client *s3.S3
	bucket   string
}

func NewS3Uploader(bucket string, region string, accessKeyID string, secretAccessKey string) *S3Uploader {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	return &S3Uploader{
		s3Client: s3.New(sess),
		bucket:   bucket,
	}
}

// UploadFile uploads a file to S3
func (u *S3Uploader) UploadFile(file io.ReadSeeker, fileHeader *multipart.FileHeader) error {
	// defer file.Close()

	_, err := u.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(u.bucket),
		Key:           aws.String(fileHeader.Filename),
		Body:          file,
		ContentLength: aws.Int64(fileHeader.Size),
		ContentType:   aws.String(fileHeader.Header.Get("Content-Type")),
		ACL:           aws.String("public-read"),
	})
	return err
}
