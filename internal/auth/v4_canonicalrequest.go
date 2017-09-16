package auth

import (
	"bytes"
	"crypto/sha256"
	"net/url"
	"sort"
	"strings"
)

const (
	V4PayloadUnsigned            = "UNSIGNED-PAYLOAD"
	V4PayloadStreamingHmacSha256 = "STREAMING-AWS4-HMAC-SHA256-PAYLOAD"
)

type CanonicalRequest struct {
	Method        string
	URI           string
	Query         map[string]string
	Headers       map[string]string
	SignedHeaders []string
	Date          string
	PayloadHash   string
}

func (req CanonicalRequest) StringToSign(cred Credential) string {
	return strings.Join([]string{Aws4HmacSha256, req.Date, cred.Scope(), req.Digest()}, "\n")
}

func (req CanonicalRequest) String() string {
	var buffer bytes.Buffer
	req.write(&buffer)
	return buffer.String()
}

func (req CanonicalRequest) Digest() string {
	digest := sha256.New()
	req.write(WrapWriter(digest))
	return string(hexEncode(digest.Sum(nil)))
}

func (req CanonicalRequest) write(writer Writer) {
	writer.WriteString(strings.ToUpper(req.Method))
	writer.WriteByte('\n')
	writer.WriteString(urlPathEscape(req.URI))
	writer.WriteByte('\n')

	req.writeQuery(writer)
	writer.WriteByte('\n')

	req.writeHeaders(writer)
	writer.WriteByte('\n')

	req.writeSignedHeaders(writer)
	writer.WriteByte('\n')

	if len(req.PayloadHash) == 0 {
		writer.WriteString(V4PayloadUnsigned)
	} else {
		writer.WriteString(req.PayloadHash)
	}
}

func (req CanonicalRequest) writeQuery(writer Writer) {
	if len(req.Query) == 0 {
		return
	}
	keys := []string{}
	query := map[string]string{}
	for key := range req.Query {
		newKey := url.QueryEscape(key)
		keys = append(keys, newKey)
		query[newKey] = url.QueryEscape(req.Query[key])
	}
	sort.Strings(keys)
	writer.WriteString(keys[0])
	writer.WriteByte('=')
	writer.WriteString(query[keys[0]])
	for _, key := range keys[1:len(keys)] {
		writer.WriteByte('&')
		writer.WriteString(key)
		writer.WriteByte('=')
		writer.WriteString(query[key])
	}
}

func (req CanonicalRequest) writeHeaders(writer Writer) {
	if len(req.Headers) == 0 {
		return
	}
	keys := []string{}
	headers := map[string]string{}
	for key := range req.Headers {
		newKey := strings.ToLower(key)
		keys = append(keys, newKey)
		headers[newKey] = strings.Trim(req.Headers[key], " ")
	}
	sort.Strings(keys)
	for _, key := range keys {
		writer.WriteString(key)
		writer.WriteByte(':')
		writer.WriteString(headers[key])
		writer.WriteByte('\n')
	}
}

func (req CanonicalRequest) writeSignedHeaders(writer Writer) {
	headers := []string{}
	for _, key := range req.SignedHeaders {
		headers = append(headers, strings.ToLower(key))
	}
	sort.Strings(headers)
	writer.WriteString(strings.Join(headers, ";"))
}
