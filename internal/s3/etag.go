package s3

import (
	"encoding/xml"
	"strings"
)

type ETag string

func NewETag(hash string) ETag {
	hash = strings.Trim(hash, `"`)
	if strings.Contains(hash, "&;") {
		hash = strings.TrimPrefix(hash, "&quot;")
		hash = strings.TrimSuffix(hash, "&quot;")
		hash = strings.TrimPrefix(hash, "&#34;")
		hash = strings.TrimSuffix(hash, "&#34;")
	}
	return ETag(hash)
}

func (etag ETag) MD5() string {
	return string(etag)
}

func (etag ETag) HeaderValue() []string {
	return []string{string(etag)}
}

func (etag ETag) String() string {
	return `"` + string(etag) + `"`
}

func (etag ETag) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	if err = e.EncodeToken(start); err != nil {
		return
	}
	if err = e.EncodeToken(xml.CharData(etag.String())); err != nil {
		return
	}
	return e.EncodeToken(start.End())
}
