package s3storage

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Fetcher interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

func GetS3Data(fetcher S3Fetcher, bucket string, key string) ([]byte, error) {
	results, err := fetcher.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}
	defer results.Body.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, results.Body); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
