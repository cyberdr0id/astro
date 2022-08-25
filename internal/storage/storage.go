package storage

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 presents type for work with AWS Simple storage service.
type S3 struct {
	s3     *s3.S3
	bucket string
}

var (
	errEmptyAccessKey   = errors.New("empty AWS access key")
	errEmptyAccessKeyID = errors.New("empty AWS access key ID")
	errEmptyRegion      = errors.New("empty AWS region")
	errEmptyBucketName  = errors.New("empty AWS bucket name")
)

// New creates a new instance of S3.
func New() (*S3, error) {
	if os.Getenv("AWS_BUCKET") == "" {
		return nil, errEmptyBucketName
	}

	if os.Getenv("AWS_REGION") == "" {
		return nil, errEmptyRegion
	}

	if os.Getenv("AWS_ACCESS_KEY") == "" {
		return nil, errEmptyAccessKey
	}

	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		return nil, errEmptyAccessKeyID
	}

	sn, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_ACCESS_KEY"), ""),
	})
	if err != nil {
		return &S3{}, fmt.Errorf("unable to create session: %w", err)
	}

	return &S3{
		s3:     s3.New(sn),
		bucket: os.Getenv("AWS_BUCKET"),
	}, nil
}

// Upload uploads file to the bucket.
func (s *S3) Upload(file io.ReadSeeker, fileName string) error {
	_, err := s.s3.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})
	if err != nil {
		return fmt.Errorf("unable to put file to object storage: %w", err)
	}

	return nil
}
