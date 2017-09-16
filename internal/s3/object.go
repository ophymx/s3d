package s3

import (
	"io"
	"net/http"
	"strconv"
)

////////////////////////
// Responses
////////////////////////

// Object is the result of ObjectService#Get()
type Object struct {
	File          io.ReadCloser
	ContentLength int64
	ContentType   string
	LastModified  string
	CacheControl  string
	ETag          ETag
	UserDefined   map[string]string
	VersionID     string
}

// Send writes headers and file to writer.
func (resp Object) Send(writer http.ResponseWriter) (err error) {
	writer.Header().Add(HdrContentLength, strconv.FormatInt(resp.ContentLength, 10))
	if resp.ContentType != "" {
		writer.Header().Add(HdrContentType, resp.ContentType)
	}
	// .Add() canonicalizes ETag to 'Etag'
	writer.Header()[HdrETag] = resp.ETag.HeaderValue()
	writer.Header().Add(HdrLastModified, resp.LastModified)
	if resp.CacheControl != "" {
		writer.Header().Add(HdrCacheControl, resp.CacheControl)
	}
	for key, value := range resp.UserDefined {
		writer.Header().Add(AmzMetaPrefix+key, value)
	}
	if resp.VersionID != "" {
		writer.Header().Add(AmzVersionID, resp.VersionID)
	}

	if resp.File != nil {
		if _, err = io.Copy(writer, resp.File); err != nil {
			return
		}
		err = resp.File.Close()
	}
	return
}

func (resp Object) HTTPStatus() int {
	return http.StatusOK
}

type CopyObjectResult struct {
	XMLNS
	LastModified string
	ETag         ETag
}

func (results CopyObjectResult) Send(writer http.ResponseWriter) error {
	results.NS = NSS3
	return sendXML(writer, results)
}

func Created(etag ETag) SimpleResponse {
	return SimpleResponse{
		Status: http.StatusOK,
		Header: http.Header{
			HdrETag: []string{etag.String()},
		},
	}
}
