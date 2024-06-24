package s3

// Config is the configuration for setting up AWS S3 client
type Config struct {
	// Region is the region that the bucket will be in
	Region string

	// KeyID is the AWS Access Key ID
	KeyID string

	// SecretKey is the AWS Secret Key
	SecretKey string

	// SessionToken is th AWS Session token
	SessionToken string
}
