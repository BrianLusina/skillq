package minio

// Config is the configuration for setting up a minio client
type Config struct {
	// PublicUrl is the public URL that is used as a base URL for public documents/images etc
	PublicUrl string

	// Endpoint is the Minio Host
	Endpoint string

	// AccessKeyID is the access key ID
	AccessKeyID string

	// SecretAccessKey is the secret access key
	SecretAccessKey string

	// UseSSL is a boolean defining whether to use a SSL connection or not
	UseSSL bool

	// Token is the token to use for the connection, can be a blank string
	Token string
}
