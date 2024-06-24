package google

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/storage"
	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	googleStorage "cloud.google.com/go/storage"
)

// GoogleStorageClient is a structure for handling google cloud storage
type GoogleStorageClient struct {
	log       logger.Logger
	client    *googleStorage.Client
	projectID string
}

// NewClient creates a new storage client with google cloud
func NewClient(config Config, log logger.Logger) (storage.StorageClient, error) {
	ctx := context.Background()

	client, err := googleStorage.NewClient(ctx,
		option.WithCredentials(&google.Credentials{
			ProjectID: config.ProjectID,
		}),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create google cloud storage client")
	}

	return &GoogleStorageClient{
		log:       log,
		client:    client,
		projectID: config.ProjectID,
	}, nil
}

// Upload uploads a new storage item and returns the URL to the stored item
func (sc *GoogleStorageClient) Upload(ctx context.Context, item storage.StorageItem) (string, error) {
	bucket := item.Bucket

	// does the bucket already exist? If it does, ignore and proceed to upload document, if it does not, create the bucket item first
	if ok, err := sc.BucketExists(ctx, bucket); !ok && err != nil {
		// create bucket
		sc.log.Infof("Bucket %s does not exist. Error info: %v. Creating bucket...", bucket, err)

		err := sc.CreateBucket(ctx, bucket)
		if err != nil {
			sc.log.Errorf("Failed to create bucket %s with error: %v", bucket, err)
			return "", errors.Wrapf(err, "failed to create bucket %s when uploading item: %v", bucket, item)
		}
	}

	// obtain the handle and set to upload the
	objectHandle := sc.client.Bucket(bucket).Object(item.Name)

	attrs, err := objectHandle.Attrs(ctx)
	if err != nil {
		sc.log.Errorf("Failed to obtain object handle attributes: %v", err)
		return "", errors.Wrapf(err, "object.Attrs: %v", err)
	}

	objectHandle = objectHandle.If(googleStorage.Conditions{DoesNotExist: true})
	objectHandle = objectHandle.If(googleStorage.Conditions{GenerationMatch: attrs.Generation})

	// upload an object with the storage writer
	objectWriter := objectHandle.NewWriter(ctx)

	// get the document data from the content, this is used to create a buffered reader
	document, err := storage.GetDocumentData(item.Content)
	if err != nil {
		sc.log.Errorf("Failed to retrieve document data from item %v", item)
		return "", errors.Wrapf(err, "failed to retrieve document data")
	}

	bufferedReader := bytes.NewReader(document.Data)

	if _, err := io.Copy(objectWriter, bufferedReader); err != nil {
		sc.log.Errorf("Failed to upload document data: %v with err: %v", item, err)
		return "", errors.Wrapf(err, "failed to upload document %v", item)
	}

	if err := objectWriter.Close(); err != nil {
		sc.log.Errorf("Failed to close object writer with err: %v", err)
		return "", errors.Wrapf(err, "failed to close object writer: %v", err)
	}

	return fmt.Sprintf("https://storage.cloud.google.com/%s/%s", bucket, item.Name), nil
}

// CreateBucket creates a bucket
func (sc *GoogleStorageClient) CreateBucket(ctx context.Context, bucketName string) error {
	if ok, err := sc.BucketExists(ctx, bucketName); err == nil && ok {
		// bucket already exists
		sc.log.Warn(fmt.Sprintf("bucket %s already exists", bucketName))
		return fmt.Errorf("bucket %s already exists", bucketName)
	}

	// Creates a Bucket instance.
	bucketHandle := sc.client.Bucket(bucketName)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Creates the new bucket.
	if err := bucketHandle.Create(ctx, sc.projectID, nil); err != nil {
		sc.log.Errorf("Failed to create bucket %s with error %v", bucketName, err)
		return errors.Wrapf(err, "failed to create bucket %s", bucketName)
	}

	return nil
}

// BucketExists checks if a bucket exists
func (sc *GoogleStorageClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	bucketHandle := sc.client.Bucket(bucketName)
	_, err := bucketHandle.Attrs(ctx)
	if err != nil {
		sc.log.Errorf("Failed to check existence of bucket: %s. Err: %v", bucketName, err)
		return false, errors.Wrapf(err, "bucket %s does not exist", bucketName)
	}

	return true, nil
}
