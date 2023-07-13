package util

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	uuid2 "github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type UploadClient struct {
	client     *s3.Client
	presign    *s3.PresignClient
	BucketName string
}

func NewUploadClient(bucketName string, accountId string, keyId string, secret string) (*UploadClient, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(keyId, secret, "")),
	)
	if err != nil {
		log.Err(err).Msg("error while uploading")
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	presign := s3.NewPresignClient(client)
	return &UploadClient{client: client, BucketName: bucketName, presign: presign}, nil
}

func (u *UploadClient) Upload(content []byte, ctx context.Context) (*s3.PutObjectOutput, error) {
	uuid, _ := uuid2.NewUUID()
	uuidStr := uuid.String()
	uuidPtr := &uuidStr
	res, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(content),
		Key:    uuidPtr,
		Bucket: &u.BucketName,
	})
	if err != nil {
		log.Err(err).Msg("error while uploading")
	}
	return res, nil
}

func (u *UploadClient) Presign(ctx context.Context) (*v4.PresignedHTTPRequest, string, error) {
	uuid, _ := uuid2.NewUUID()
	uuidStr := uuid.String()
	uuidPtr := &uuidStr

	presigned, err := u.presign.PresignPutObject(ctx, &s3.PutObjectInput{Bucket: &u.BucketName, Key: uuidPtr})
	if err != nil {
		log.Err(err).Msg("error while uploading")
		return nil, uuidStr, err
	}
	return presigned, uuidStr, nil
}
