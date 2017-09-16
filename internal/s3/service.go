package s3

import (
	"net/http"
)

type ListAllMyBucketsResultBucket struct {
	Name         string
	CreationDate string
}

type ListAllMyBucketsResult struct {
	XMLNS
	Owner   OwnerResult
	Buckets []ListAllMyBucketsResultBucket `xml:">Bucket"`
}

func (results ListAllMyBucketsResult) Send(writer http.ResponseWriter) error {
	results.NS = NSS3
	return sendXMLHeader(writer, results)
}
