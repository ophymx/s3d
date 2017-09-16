package auth

import (
	"bytes"
	"regexp"
	"sort"
	"strings"
)

type CanonicalRequestV2 struct {
	Method      string
	ContentMD5  string
	ContentType string
	Date        string
	Headers     map[string]string
	Resource    Resource
	Query       map[string]string
}

func (req CanonicalRequestV2) StringToSign() string {
	return req.String()
}

func (req CanonicalRequestV2) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(req.Method)
	buffer.WriteByte('\n')

	buffer.WriteString(req.ContentMD5)
	buffer.WriteByte('\n')

	buffer.WriteString(req.ContentType)
	buffer.WriteByte('\n')

	buffer.WriteString(req.Date)
	buffer.WriteByte('\n')

	req.writeCanonicalizedAmzHeaders(buffer)
	req.writeCanonicalizedResource(buffer)

	return buffer.String()
}

var foldingRegex = regexp.MustCompile("\\s*\\n\\s*")

func (req CanonicalRequestV2) writeCanonicalizedAmzHeaders(buffer bytes.Buffer) {
	keys := []string{}
	headers := map[string]string{}
	for key, value := range req.Headers {
		newKey := strings.ToLower(key)
		keys = append(keys, newKey)

		value = foldingRegex.ReplaceAllString(strings.TrimLeft(value, " "), " ")
		if v, exists := headers[newKey]; exists {
			headers[newKey] = v + "," + value
		} else {
			headers[newKey] = value
		}
	}
	sort.Strings(keys)

	for _, key := range keys {
		buffer.WriteString(key)
		buffer.WriteByte(':')
		buffer.WriteString(headers[key])
		buffer.WriteByte('\n')
	}
}

func (req CanonicalRequestV2) writeCanonicalizedResource(buffer bytes.Buffer) {
	buffer.WriteByte('/')
	buffer.WriteString(req.Resource.Bucket())

	if req.Resource.Key() != "" {
		buffer.WriteByte('/')
		buffer.WriteString(req.Resource.Key())
	}

	if len(req.Query) > 0 {

	}
}
