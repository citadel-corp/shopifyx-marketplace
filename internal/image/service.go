package image

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	bucket = os.Getenv("S3_BUCKET_NAME")
	region = os.Getenv("S3_REGION")
	key    = os.Getenv("S3_SECRET_KEY")
)

type Service interface {
	UploadToS3(ctx context.Context, readSeeker io.ReadSeeker, filename string) (string, error)
}

type imageService struct {
	awsSession *session.Session
}

func NewService(awsSession *session.Session) Service {
	return &imageService{
		awsSession: awsSession,
	}
}

func (s *imageService) UploadToS3(ctx context.Context, readSeeker io.ReadSeeker, filename string) (string, error) {
	svc := s3.New(s.awsSession)

	// This uploads the contents of the buffer to S3
	_, err := svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(filename),
		ACL:         aws.String("public-read"),
		Body:        readSeeker,
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", err
	}
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	}

	req, _ := svc.GetObjectRequest(params)
	req.Build()
	return req.HTTPRequest.URL.String(), nil
}
