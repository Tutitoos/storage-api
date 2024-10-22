package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	r2 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"io"
	"log"
	"storage-api/src/domain"
	"time"
)

type ICloudflareService struct {
	Client     *r2.Client
	BucketName string
	Context    context.Context
}

func CloudflareService() *ICloudflareService {
	ctx := context.TODO()

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{URL: domain.CONFIG.BucketUrl}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(domain.CONFIG.BucketRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(domain.CONFIG.CloudflareAccessKeyId, domain.CONFIG.CloudflareSecretAccessKey, "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Fatalf("Error loading config " + err.Error())

		return nil
	}

	return &ICloudflareService{
		Client:     r2.NewFromConfig(cfg),
		BucketName: domain.CONFIG.BucketName,
		Context:    ctx,
	}
}

func (s *ICloudflareService) GetFiles(folder string) ([]types.Object, error) {
	input := &r2.ListObjectsV2Input{
		Prefix: &folder,
		Bucket: &s.BucketName,
	}

	resp, err := s.Client.ListObjectsV2(s.Context, input)
	if err != nil {
		return nil, err
	}
	return resp.Contents, nil
}

func (s *ICloudflareService) GetFile(base64Filename string) (*r2.GetObjectOutput, error) {
	filenameEncoded := base64.StdEncoding.EncodeToString([]byte(base64Filename))
	filename, _ := base64.StdEncoding.DecodeString(filenameEncoded)

	resp, err := s.Client.GetObject(s.Context, &r2.GetObjectInput{
		Bucket: &s.BucketName,
		Key:    aws.String(string(filename)),
	})
	if err != nil {
		var awsErr *types.NoSuchKey
		if errors.As(err, &awsErr) {
			return nil, fmt.Errorf("File is not exist")
		}

		var smithyErr smithy.APIError
		if errors.As(err, &smithyErr) && smithyErr.ErrorCode() == "NoSuchKey" {
			return nil, fmt.Errorf("File is not exist")
		}

		return nil, err
	}

	return resp, nil
}

func (s *ICloudflareService) UploadFile(fileReader io.Reader, folderName string, filename string, contentType string) (*r2.PutObjectOutput, error) {
	filePath := folderName + "/" + filename
	resp, err := s.Client.PutObject(s.Context, &r2.PutObjectInput{
		Bucket:      &s.BucketName,
		Key:         aws.String(filePath),
		Body:        fileReader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ICloudflareService) DeleteFile(filename string) (*r2.DeleteObjectOutput, error) {
	resp, err := s.Client.DeleteObject(s.Context, &r2.DeleteObjectInput{
		Bucket: &s.BucketName,
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ICloudflareService) GenerateSignedURL(filename string) (string, error) {
	resignClient := r2.NewPresignClient(s.Client)
	input := &r2.GetObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(filename),
	}
	resp, err := resignClient.PresignGetObject(s.Context, input, r2.WithPresignExpires(1*time.Hour))
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}
