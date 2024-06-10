package minio

import "fmt"

// ReadOnlyPolicyJson returns the policy JSON setting for a bucket
func ReadOnlyPolicyJson(bucket string) string {
	resource := fmt.Sprintf(`["arn:aws:s3:::%s/*"]`, bucket)

	return fmt.Sprintf(`{
		"Version": "2012-10-17", 
		"Statement": [
			{
				"Action": ["s3:GetObject"], 
				"Effect": "Allow",
				"Sid": "",
				"Principal": {"AWS": ["*"]},
				"Resource": %s,
			},
		]
	}`, resource)
}
