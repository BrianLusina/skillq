package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/storage"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	// "github.com/aws/aws-sdk-go-v2/feature/s3/ma"
	"github.com/pkg/errors"
)

// S3StorageClient is a structure that is a wrapper around an AWS S3 client
type S3StorageClient struct {
	// s3Client is the AWS S3 s3Client
	s3Client *awsS3.Client

	// region is the region this client operates in
	region string

	// log is an application logger
	log logger.Logger
}

// NewClient creates a new S3 storage client
func NewClient(config Config, log logger.Logger) (storage.StorageClient, error) {
	awsCredentials := credentials.NewStaticCredentialsProvider(config.KeyID, config.SecretKey, config.SessionToken)
	s3Config := aws.Config{
		Region:      *aws.String(config.Region),
		Credentials: awsCredentials,
	}

	client := awsS3.NewFromConfig(s3Config)

	return &S3StorageClient{s3Client: client, log: log, region: config.Region}, nil
}

func (sc *S3StorageClient) Upload(ctx context.Context, item storage.StorageItem) (string, error) {
	bucket := item.Bucket
	if ok, err := sc.BucketExists(ctx, bucket); err != nil || !ok {
		err := sc.CreateBucket(ctx, bucket)
		if err != nil {
			return "", errors.Wrapf(err, "failed to create bucket")
		}
	}

	sc.log.Infof("Uploading storage item %v", item)

	document, err := storage.GetDocumentData(item.Content)
	if err != nil {
		sc.log.Errorf("Failed to retrieve document data from item %v", item)
		return "", errors.Wrapf(err, "failed to retrieve document data")
	}

	// upload document
	bufferedReader := bytes.NewReader(document.Data)
	key := fmt.Sprintf("%s.%s", item.Name, document.FileExtension)

	_, err = sc.s3Client.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		Body:                 bufferedReader,
		ContentType:          aws.String(document.MimeType),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
		StorageClass:         types.StorageClassIntelligentTiering,
		ACL:                  types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		sc.log.Errorf("Failed to upload item %v, Err: %v", item, err)
		return "", errors.Wrapf(err, "failed to upload storage item %v", item)
	}

	location := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", bucket, sc.region, key)

	return location, nil
}

// CreateBucket creates a bucket with the specified name in the specified Region.
func (sc *S3StorageClient) CreateBucket(ctx context.Context, name string) error {
	_, err := sc.s3Client.CreateBucket(ctx, &awsS3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(sc.region),
		},
	})
	if err != nil {
		sc.log.Errorf("Couldn't create bucket %v in Region %v. Here's why: %v\n", name, sc.region, err)
		return err
	}

	return nil
}

func (sc *S3StorageClient) Download(ctx context.Context, bucket, key string) (*storage.StorageDownloadItem, error) {
	sc.log.Infof("Downloading storage item from bucket %v with key: %s", bucket, key)

	result, err := sc.s3Client.GetObject(ctx, &awsS3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		sc.log.Errorf("Failed to download item from bucket %v with key: %s, Err: %v", bucket, key, err)
		return nil, errors.Wrapf(err, "failed to upload storage item %v", bucket)
	}

	body, err := io.ReadAll(result.Body)
	if err != nil {
		sc.log.Errorf("Failed to read result body of download from bucket %v with key: %s, Err: %v", bucket, key, err)
		return nil, errors.Wrapf(err, "failed to read result body of download")
	}

	return &storage.StorageDownloadItem{
		Content:      body,
		ContentType:  *result.ContentType,
		Bucket:       bucket,
		Name:         key,
		LastModified: result.LastModified,
	}, nil
}

// BucketExists checks whether a bucket exists in the current account.
func (sc *S3StorageClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := sc.s3Client.HeadBucket(ctx, &awsS3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				sc.log.Errorf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				sc.log.Errorf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		sc.log.Infof("Bucket %v exists and you already own it.", bucketName)
	}

	return exists, err
}
