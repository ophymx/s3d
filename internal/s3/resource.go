package s3

import (
	"path"
	"strings"
)

type Resource struct {
	bucket string
	key    string
}

func (res Resource) Key() string {
	return res.key
}

func (res Resource) Bucket() string {
	return res.bucket
}

func (res Resource) String() string {
	return path.Join(res.bucket, res.key)
}

func (res Resource) IsEmpty() bool {
	return res.bucket == "" && res.key == ""
}

func NewResource(bucket, key string) Resource {
	return Resource{
		bucket: bucket,
		key:    key,
	}
}

func ParseResource(res string) (resource Resource) {
	segments := strings.SplitN(strings.TrimLeft(res, "/"), "/", 2)
	switch len(segments) {
	case 0:
	case 1:
		resource.bucket = segments[0]
	default:
		resource.bucket = segments[0]
		resource.key = segments[1]
	}
	return
}
