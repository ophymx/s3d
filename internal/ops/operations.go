package ops

import (
	"io"
	"net/url"

	"github.com/ophymx/s3d/internal/s3"
)

type ServiceOperations interface {
	ListBuckets() s3.Response
}

type BucketOperations interface {
	Create(bucket string) s3.Response
	Delete(bucket string) s3.Response
	ListBucket(bucket string, query url.Values) s3.Response
}

type ObjectOperations interface {
	Get(resource s3.Resource) s3.Response
	Head(resource s3.Resource) s3.Response
	Put(resource s3.Resource, contentType string, body io.ReadCloser) s3.Response
	Copy(src, dst s3.Resource) s3.Response
	Delete(resource s3.Resource) s3.Response
}
