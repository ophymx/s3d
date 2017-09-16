package s3

const (
	AmzRequestID  = "x-amz-request-id"
	AmzHostID     = "x-amz-id-2"
	AmzCopySource = "x-amz-copy-source"
	AmzVersionID  = "x-amz-version-id"
	AmzMetaPrefix = "x-amz-meta-"

	// Common headers
	HdrContentMD5    = "Content-MD5"
	HdrContentLength = "Content-Length"
	HdrContentType   = "Content-Type"
	HdrLastModified  = "Last-Modified"
	HdrCacheControl  = "Cache-Control"
	HdrETag          = "ETag"
)
