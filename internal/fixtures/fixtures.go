package fixtures

import "time"
import "github.com/ophymx/s3d/internal/meta"

//go:generate go-bindata --pkg fixtures -o xml.go --debug xml

var (
	Time1              = time.Date(2014, 5, 6, 3, 2, 1, 0, time.UTC)
	Time2              = time.Date(2015, 6, 7, 4, 3, 2, 2000000, time.UTC)
	BucketCreationData = time.Date(2017, 8, 1, 2, 9, 8, 0, time.UTC)
)

func BucketMetadata() meta.BucketData {
	return meta.BucketData{
		CreationDate: BucketCreationData,
	}
}

const (
	ObjectContentMD5         = "4672ce371fb3c1170a9e71bc4b2810b9"
	ObjectCacheControl       = "max-age=31536000"
	ObjectContentType        = "application/x-iso9660-image"
	ObjectVersionID          = "222222"
	ObjectSize         int64 = 718274560
)

var (
	ObjectLastModified = time.Date(2017, 2, 11, 5, 1, 0, 1000000, time.UTC)
)

func ObjectMetadata() meta.ObjectData {
	return meta.ObjectData{
		ContentMD5:   ObjectContentMD5,
		Size:         ObjectSize,
		CacheControl: ObjectCacheControl,
		LastModified: ObjectLastModified,
		ContentType:  ObjectContentType,
		VersionID:    ObjectVersionID,
		UserDefined:  ObjectUserDefined(),
	}
}

func ObjectUserDefined() map[string]string {
	return map[string]string{
		"name":    "zesty",
		"flavor":  "server",
		"arch":    "amd64",
		"version": "17.04",
	}
}
