package storage

import (
	"context"
	"fmt"
	"os"
	"bytes"

	l "github.com/CarlosCaravanTsz/imgAI/internal/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	_ "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	_ "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type FotoUpload struct {
	Filename string
	Path     string
	Buffer   []byte
}

type S3Client struct {
	Client *s3.Client
	Bucket string
	Endpoint string
}

func NewS3Client() (*S3Client, error) {
	if err := godotenv.Load(); err != nil {
		l.LogInfo("Error while uploading ENV vars", logrus.Fields{
			"error": err,
		})
	}

	endpoint := os.Getenv("S3_ENDPOINT")
	region := os.Getenv("S3_REGION")
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						PartitionID:       "aws",
						URL:               endpoint,
						SigningRegion:     region,
						HostnameImmutable: true, // <- esto fuerza path-style
					}, nil
				},
			),
		),
	)
	if err != nil {
		return nil, err
	}

	return &S3Client{
		Client: s3.New(s3.Options{
			Credentials:      cfg.Credentials,
			Region:           cfg.Region,
			EndpointResolver: s3.EndpointResolverFromURL(endpoint),
			UsePathStyle:     true,
		}),
		Bucket: bucket,
		Endpoint: endpoint,
	}, nil
}

func (s *S3Client) Upload(foto FotoUpload) (string, error) {

	key := fmt.Sprintf("%s/%s", foto.Path, foto.Filename)


	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &key,
		Body:   bytes.NewReader(foto.Buffer),
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s/%s", s.Endpoint, s.Bucket, key)
	return url, nil
}
