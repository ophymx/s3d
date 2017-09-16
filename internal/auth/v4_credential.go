package auth

import (
	"path"
	"strings"
)

type Credential struct {
	AccessKeyID string
	Date        string
	Region      string
	Service     string
}

var (
	errInvalidCredFormat         = AuthorizationQueryParametersError(`Error parsing the X-Amz-Credential parameter; the Credential is mal-formed; expecting "&lt;YOUR-AKID&gt;/YYYYMMDD/REGION/SERVICE/aws4_request".`)
	errInvalidCredTerminalFormat = `Error parsing the X-Amz-Credential parameter; incorrect terminal "%s". This endpoint uses "aws4_request".`
)

func ParseCredential(component string) (credential Credential, err error) {
	parts := strings.SplitN(component, "/", 5)
	if len(parts) != 5 {
		err = errInvalidCredFormat
		return
	}
	if parts[4] != "aws4_request" {
		err = AuthorizationQueryParametersErrorf(errInvalidCredTerminalFormat, parts[4])
	}
	credential = Credential{
		AccessKeyID: parts[0],
		Date:        parts[1],
		Region:      parts[2],
		Service:     parts[3],
	}
	return
}

func (c Credential) SigningKey(secretKey string) SigningKey {
	key := hmacSha256([]byte("AWS4"+secretKey), c.Date)
	key = hmacSha256(key, c.Region)
	key = hmacSha256(key, c.Service)
	key = hmacSha256(key, "aws4_request")
	return SigningKey(key)
}

func (c Credential) Scope() string {
	return path.Join(c.Date, c.Region, c.Service, "aws4_request")
}

func (c Credential) String() string {
	return path.Join(c.AccessKeyID, c.Scope())
}
