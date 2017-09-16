package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	errQueryV4Missing     = AuthorizationQueryParametersError("Query-string authentication version 4 requires the X-Amz-Algorithm, X-Amz-Credential, X-Amz-Signature, X-Amz-Date, X-Amz-SignedHeaders, and X-Amz-Expires parameters.")
	errQueryV4Unsupported = AuthorizationQueryParametersError(`X-Amz-Algorithm only supports "AWS4-HMAC-SHA256"`)
	errExpiresFormat      = AuthorizationQueryParametersError("X-Amz-Expires should be a number")
)

type AuthorizationV4 struct {
	Credential    Credential
	Date          time.Time
	Expires       int
	SignedHeaders []string
	Signature     string
}

func (auth AuthorizationV4) GetAccessKeyID() string {
	return auth.Credential.AccessKeyID
}

func (auth AuthorizationV4) Verify(secretKey string, req *http.Request) error {
	canonicalReq := CanonicalRequest{
		Method:        req.Method,
		URI:           req.URL.Path,
		Query:         map[string]string{},
		Headers:       map[string]string{},
		SignedHeaders: auth.SignedHeaders,
		Date:          auth.Date.Format(ISO8601),
		PayloadHash:   req.Header.Get(AmzContentSHA256),
	}
	canonicalReq.Headers[hdrHost] = req.Host
	for _, header := range auth.SignedHeaders {
		if header != "host" {
			canonicalReq.Headers[header] = req.Header.Get(header)
		}
	}

	query := req.URL.Query()
	for key := range query {
		if strings.ToLower(key) != strings.ToLower(AmzSignature) {
			canonicalReq.Query[key] = query.Get(key)
		}
	}
	sts := canonicalReq.StringToSign(auth.Credential)

	log.Printf("canonical-request: \n%v\n", canonicalReq)
	log.Printf("string-to-sign: \n%s\n", sts)

	if auth.Credential.SigningKey(secretKey).Verify(sts, auth.Signature) {
		return nil
	}
	return SignatureDoesNotMatch(auth.GetAccessKeyID(), sts, auth.Signature)
}

func ParseAuthorizationV4(param string, date time.Time, values Values) (auth AuthorizationV4, err error) {
	auth.Date = date
	for _, component := range strings.Split(param, ",") {
		entry := strings.SplitN(component, "=", 2)
		if len(entry) != 2 {
			err = InvalidArgument("needs '='", "", "")
			break
		}
		key, value := strings.ToLower(strings.Trim(entry[0], " ")), strings.Trim(entry[1], " ")

		switch key {
		case "credential":
			auth.Credential, err = ParseCredential(value)
		case "signedheaders":
			auth.SignedHeaders = strings.Split(strings.ToLower(value), ";")
		case "signature":
			auth.Signature = value
		case "":
			continue
		default:
			err = InvalidArgument("invalid component", key, value)
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		return
	}
	if auth.Credential == (Credential{}) || auth.SignedHeaders == nil || auth.Signature == "" {
		err = MissingSecurityElement("missing component(s)")
	}
	return
}

func ParseAuthorizationQueryV4(values Values) (auth AuthorizationV4, err error) {
	algorithm := values.Get(AmzAlgorithm)
	if algorithm == "" {
		return
	}
	if algorithm != Aws4HmacSha256 {
		return AuthorizationV4{}, errQueryV4Unsupported
	}

	if cred := values.Get(AmzCredential); cred == "" {
		return AuthorizationV4{}, errQueryV4Missing
	} else if auth.Credential, err = ParseCredential(cred); err != nil {
		return
	}

	if date := values.Get(AmzDate); date == "" {
		return AuthorizationV4{}, errQueryV4Missing
	} else if auth.Date, err = parseDate(date); err != nil {
		return
	}

	if expires := values.Get(AmzExpires); expires == "" {
		return AuthorizationV4{}, errQueryV4Missing
	} else if auth.Expires, err = parseExpires(expires); err != nil {
		return
	}

	if auth.SignedHeaders = strings.Split(values.Get(AmzSignedHeaders), ";"); len(auth.SignedHeaders) == 0 {
		return AuthorizationV4{}, errQueryV4Missing
	}

	if auth.Signature = values.Get(AmzSignature); auth.Signature == "" {
		return AuthorizationV4{}, errQueryV4Missing
	}
	return
}

func parseExpires(value string) (expires int, err error) {
	if expires, err = strconv.Atoi(value); err != nil {
		err = errExpiresFormat
	}
	return
}

func hmacSha256(key []byte, value string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(value))
	return mac.Sum(nil)
}
