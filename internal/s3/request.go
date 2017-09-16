package s3

import (
	"net/http"
	"net/url"
	"time"

	"github.com/ophymx/s3d/internal/auth"
)

const (
	emptySha256Hash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	emptySha1Hash   = "da39a3ee5e6b4b0d3255bfef95601890afd80709"
	emptyMd5Hash    = "d41d8cd98f00b204e9800998ecf8427e"
)

type Request struct {
	ID         string
	Method     string
	Resource   Resource
	Auth       auth.Authorization
	Credential Credential
	Query      url.Values
	Time       time.Time
	Host       string
	RawReq     *http.Request
}

func (req Request) getHeader(key string) string {
	return req.RawReq.Header.Get(key)
}
