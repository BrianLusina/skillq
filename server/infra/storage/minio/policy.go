package minio

import "fmt"

func ReadAndWritePolicyJson(bucket string) string {
	resource := fmt.Sprintf(`["arn:aws:s3:::%s/*"]`, bucket)

	return fmt.Sprintf(`{
		"Version": "2012-10-17", 
		"Statement": [
			{
				"Sid": "ListObjectsInBucket",
				"Effect": "Allow",
				"Action": ["s3:ListBucket"], 
				"Principal": {"AWS": ["*"]},
				"Resource": %s,
			},
            {
            	"Sid": "AllObjectActions",
            	"Effect": "Allow",
            	"Action": "s3:*Object",
            	"Resource": %s
			}
		]
	}`, resource, resource)
}
