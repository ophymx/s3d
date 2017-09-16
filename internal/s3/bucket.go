package s3

import (
	"encoding/xml"
	"net/http"
	"net/url"
)

type CreateBucketConfiguration struct {
	XMLName            xml.Name `xml:"CreateBucketConfiguration"`
	LocationConstraint string
}

type OwnerResult struct {
	ID          string
	DisplayName string
}

type ContentResult struct {
	Key          string
	LastModified string
	ETag         ETag
	Size         int64
	Owner        OwnerResult
	StorageClass string
}

type CommonPrefixes []string

func (prefixes CommonPrefixes) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	if len(prefixes) == 0 {
		return nil
	}

	e.EncodeToken(start)

	for _, prefix := range prefixes {
		marshalString(e, "Prefix", prefix)
	}

	return e.EncodeToken(start.End())
}

type ListBucketResult struct {
	XMLNS
	Name           string
	Prefix         string
	Marker         string
	NextMarker     string
	MaxKeys        int
	Delimiter      string
	EncodingType   string `xml:",omitempty"`
	IsTruncated    bool
	CommonPrefixes CommonPrefixes
	Contents       []ContentResult
}

func (result *ListBucketResult) AppendPrefix(prefix string) bool {
	n := len(result.CommonPrefixes)
	if result.EncodingType == "url" {
		prefix = (&url.URL{Path: prefix}).EscapedPath()
	}
	if n == 0 || result.CommonPrefixes[n-1] != prefix {
		if result.IsFull() {
			return false
		}
		result.CommonPrefixes = append(result.CommonPrefixes, prefix)
	}
	return true
}

func (result *ListBucketResult) IsFull() bool {
	return len(result.CommonPrefixes)+len(result.Contents) >= result.MaxKeys
}

func (result ListBucketResult) Send(writer http.ResponseWriter) error {
	result.NS = NSS3
	return sendXMLHeader(writer, result)
}
