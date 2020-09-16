// Package storage is the responsible of fetching data (single files) from S3
package storage

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Fetcher is an interface that defines a method to obtain the data.
type S3Fetcher interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
}

// GetS3Data returns the corresponding S3 element using a fetcher Interface passed as a parameter.
// Every object that fulfills the interface S3Fetcher can be used to fetch the data. It is injected as a
// parameter (dependency injection).
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
