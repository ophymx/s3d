package auth

import (
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SignatureScheme string

const (
	Aws4HmacSha256 = "AWS4-HMAC-SHA256"
	Aws2           = "AWS"
)

type Authorization interface {
	GetAccessKeyID() string
	Verify(secretKey string, req *http.Request) error
}

type Resource interface {
	Key() string
	Bucket() string
}

type Values interface {
	Get(key string) string
	Del(key string)
}

type parser interface {
	parseHeaders(resource Resource, param string, date time.Time, headers Values)
	parseQuery(resource Resource, query Values)
}

func GetAuth(resource Resource, headers Values, query Values) (auth Authorization, err error) {
	if authHeader := headers.Get(hdrAuthorization); len(authHeader) > 0 {
		var date time.Time
		if date, err = getDate(headers); err != nil {
			return
		}

		var scheme, param string
		if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 {
			scheme, param = parts[0], parts[1]
		}

		switch scheme {
		case Aws4HmacSha256:
			return ParseAuthorizationV4(param, date, headers)
		case Aws2:
			return ParseAuthorizationV2(resource, param, date, headers)
		case "":
			return nil, InvalidArgument(errHeaderSpacing, hdrAuthorization, authHeader)
		default:
			return nil, InvalidArgument(errUnsupportedType, hdrAuthorization, authHeader)
		}
	}

	switch {
	case query.Get(qryAwsAccessKeyID) != "":
		return ParseAuthorizationQueryV2(resource, query)
	case query.Get(AmzAlgorithm) != "":
		return ParseAuthorizationQueryV4(query)
	default:
		return
	}
}

func urlPathEscape(s string) string {
	u := &url.URL{Path: s}
	return u.String()
}

func hexEncode(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}
