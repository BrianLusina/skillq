package minio

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/storage"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

// MinioStorageClient is a wrapper around minio that enables interactions with a Minio cluster
type MinioStorageClient struct {
	client *minio.Client
	log    logger.Logger
}

// NewClient creates a new Minio Storage Client
func NewClient(config Config, log logger.Logger) (storage.StorageClient, error) {
	creds := credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, config.Token)
	options := &minio.Options{
		Creds:  creds,
		Secure: config.UseSSL,
	}

	client, err := minio.New(config.Endpoint, options)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create minio storage client")
	}

	ctxCancel, err := client.HealthCheck(2 * time.Second)
	if err != nil {
		log.Warn("health check already started")
	}
	defer ctxCancel()

	if client.IsOnline() {
		log.Infof("Minio client is online & connected to %s", config.Endpoint)
	}

	return &MinioStorageClient{
		client: client,
		log:    log,
	}, nil
}

func (sc *MinioStorageClient) Upload(ctx context.Context, item storage.StorageItem) (string, error) {
	bucket := item.Bucket
	if err := sc.CreateBucket(ctx, bucket); err != nil {
		// check to see if we already own this bucket, which happens if this is run twice
		if exists, errBucketExists := sc.BucketExists(ctx, bucket); errBucketExists == nil && exists {
			sc.log.Warn(fmt.Sprintf("Bucket Already exists %s", bucket))
		} else {
			sc.log.Errorf("Failed to create bucket: %v", err)
			return "", err
		}
	}

	sc.log.Infof("Successfully created bucket %s", bucket)

	document, err := storage.GetDocumentData(item.Content)
	if err != nil {
		sc.log.Errorf("Failed to retrieve document data from item %v", item)
		return "", errors.Wrapf(err, "failed to retrieve document data")
	}

	// upload document
	reader := bytes.NewReader(document.Data)

	info, err := sc.client.PutObject(ctx, bucket, item.Name, reader, reader.Size(), minio.PutObjectOptions{
		ContentType:  item.ContentType,
		UserMetadata: item.Metadata,
		Progress:     reader,
	})
	if err != nil {
		sc.log.Errorf("Failed to upload document: %v", err)
		return "", errors.Wrapf(err, "failed to upload document")
	}

	if err := sc.client.SetBucketPolicy(ctx, bucket, ReadAndWritePolicyJson(bucket)); err != nil {
		sc.log.Errorf("Failed to set bucket policy for bucket %s with error %s", bucket, err)
	}

	sc.log.Infof("Successfully uploaded document to bucket %s at location %s", info.Bucket, info.Location)

	loc := fmt.Sprintf("%s/%s/%s.%s", sc.client.EndpointURL(), bucket, item.Name, document.FileExtension)

	return loc, nil
}

func (sc *MinioStorageClient) CreateBucket(ctx context.Context, bucket string) error {
	err := sc.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
		ObjectLocking: false,
		// TODO: set the region
		// Region: sc.region,
	})
	if err != nil {
		sc.log.Errorf("Failed to create bucket %s, with reason: %v:", bucket, err)
		return err
	}

	sc.log.Infof("Successfully created bucket %s", bucket)
	return nil
}

func (sc *MinioStorageClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	exists, err := sc.client.BucketExists(ctx, bucketName)
	if err == nil && exists {
		sc.log.Infof("Bucket exists %s", bucketName)
		return true, nil
	} else if err != nil {
		sc.log.Errorf("Failed to check existence of bucket:%s", bucketName)
		return false, err
	}
	return exists, nil
}
