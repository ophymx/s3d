package meta

import (
	"errors"
	"time"
)

var (
	ErrBucketNotFound        = errors.New("metadata bucket not found")
	ErrKeyNotFound           = errors.New("metadata key not found")
	ErrMissingBucketMetadata = errors.New("bucket metadata not found")
)

type DB interface {
	Get(target Target) (data ObjectData, err error)
	Put(target Target, data ObjectData) error
	Delete(target Target) error
	CreateBucket(bucket string, data BucketData) error
	DeleteBucket(bucket string) error
	ListBuckets() (buckets []Bucket, err error)
	ForEachInBucket(bucket, seek string, forEach ForEachFunc) error
	Close() error
}

type Target interface {
	Bucket() string
	Key() string
}

type LazyObject interface {
	Get() (obj ObjectData, err error)
}

type ForEachFunc func(key string, obj LazyObject) (next bool, err error)

type Bucket struct {
	Name     string
	Metadata BucketData
}

type BucketData struct {
	CreationDate time.Time
}

type ObjectData struct {
	ContentMD5   string
	Size         int64
	CacheControl string
	LastModified time.Time
	ContentType  string
	VersionID    string
	UserDefined  map[string]string
}

type Encoding interface {
	EncodeBucket(data BucketData) ([]byte, error)
	DecodeBucket(b []byte) (BucketData, error)
	EncodeObject(data ObjectData) ([]byte, error)
	DecodeObject(b []byte) (ObjectData, error)
}
