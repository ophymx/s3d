package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	qryAwsAccessKeyID = "AWSAccessKeyId"
	qrySignature      = "Signature"
	qryExpires        = "Expires"
)

const (
	errInvalidDateExpires = "Invalid date (should be seconds since epoch): "
)

var (
	errQueryV2Missing = AccessDenied("Query-string authentication requires the Signature, Expires and AWSAccessKeyId parameters")
)

type AuthorizationV2 struct {
	AccessKeyID string
	Signature   string
	Date        time.Time
	Expires     int64
	Resource    Resource
}

func ParseAuthorizationV2(resource Resource, param string, date time.Time, values Values) (auth AuthorizationV2, err error) {
	parts := strings.SplitN(param, ":", 2)
	if len(parts) != 2 {
		err = InvalidArgument(errInvalidFormat, hdrAuthorization, param)
		return
	}
	auth = AuthorizationV2{
		AccessKeyID: parts[0],
		Signature:   parts[1],
		Date:        date,
		Resource:    resource,
	}
	return
}

func ParseAuthorizationQueryV2(resource Resource, values Values) (auth AuthorizationV2, err error) {
	auth = AuthorizationV2{
		AccessKeyID: values.Get(qryAwsAccessKeyID),
		Signature:   values.Get(qrySignature),
		Resource:    resource,
	}

	expires := values.Get(qryExpires)
	if auth.AccessKeyID == "" || auth.Signature == "" || expires == "" {
		return AuthorizationV2{}, errQueryV2Missing
	}
	if auth.Expires, err = strconv.ParseInt(expires, 10, 64); err != nil {
		err = AccessDenied(errInvalidDateExpires + expires)
	}
	return
}

func (auth AuthorizationV2) GetAccessKeyID() string {
	return auth.AccessKeyID
}

func (auth AuthorizationV2) Verify(secretKey string, req *http.Request) error {
	canonicalReq := CanonicalRequestV2{
		Method:      req.Method,
		ContentMD5:  req.Header.Get(hdrContentMD5),
		ContentType: req.Header.Get(hdrContentType),
		Date:        auth.Date.Format(""),
		Headers:     map[string]string{},
		Resource:    auth.Resource,
	}
	sts := canonicalReq.StringToSign()
	if SigningKeyV2(secretKey).Verify(sts, auth.Signature) {
		return nil
	}

	return nil
}

func hmacSha1(key []byte, value string) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(value))
	return mac.Sum(nil)
}
