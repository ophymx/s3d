package ops

import (
	"errors"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ophymx/s3d/internal/blob"
	"github.com/ophymx/s3d/internal/clock"
	"github.com/ophymx/s3d/internal/meta"
	"github.com/ophymx/s3d/internal/s3"
)

const (
	encodingTypeURL = "url"
)

var (
	errInvalidMaxKeys = errors.New("invalid max-keys")
)

type bucketOps struct {
	db    meta.DB
	store blob.Store
	clock clock.Clock
}

func NewBucket(db meta.DB, store blob.Store, clock clock.Clock) BucketOperations {
	return bucketOps{db: db, store: store, clock: clock}
}

func (srv bucketOps) Create(bucket string) s3.Response {
	if response := validateBucket(bucket); response != nil {
		return response
	}
	bktMeta := meta.BucketData{CreationDate: srv.clock.Now()}
	if err := srv.db.CreateBucket(bucket, bktMeta); err != nil {
		return s3.InternalError(err)
	}
	if err := srv.store.CreateBucket(bucket); err != nil {
		return s3.InternalError(err)
	}
	return s3.NoContent()
}

func (srv bucketOps) Delete(bucket string) s3.Response {
	var err error
	if err = srv.db.DeleteBucket(bucket); err == meta.ErrBucketNotFound {
		return s3.NoSuchBucket(bucket)
	} else if err != nil {
		return s3.InternalError(err)
	}
	if err = srv.store.DeleteBucket(bucket); err != nil {
		return s3.InternalError(err)
	}
	return s3.NoContent()
}

func (srv bucketOps) ListBucket(bucket string, query url.Values) s3.Response {
	maxKeys, err := parseMaxKeys(query.Get("max-keys"))
	if err != nil {
		return s3.InvalidArgument("Argument maxKeys must be an integer between 0 and 2147483647", "maxKeys", query.Get("max-keys"))
	}
	encodingType := query.Get("encoding-type")
	if encodingType != "" && encodingType != encodingTypeURL {
		return s3.InvalidArgument("Invalid Encoding Method specified in Request", "encoding-type", encodingType)
	}

	result := s3.ListBucketResult{
		Name:         bucket,
		Marker:       query.Get("marker"),
		Prefix:       query.Get("prefix"),
		Delimiter:    query.Get("delimiter"),
		EncodingType: encodingType,
		MaxKeys:      maxKeys,
	}

	err = srv.db.ForEachInBucket(bucket, result.Marker, func(key string, lazy meta.LazyObject) (bool, error) {
		if !strings.HasPrefix(key, result.Prefix) {
			return false, nil
		}

		if result.Delimiter != "" {
			pl := len(result.Prefix)
			if idx := strings.Index(key[pl:len(key)], result.Delimiter); idx != -1 {
				key = key[0 : pl+idx+1]
				if result.EncodingType == encodingTypeURL {
					key = urlEncodePath(key)
				}
				if key != result.Marker {
					if !result.AppendPrefix(key) {
						result.IsTruncated = true
						return false, nil
					}
					result.NextMarker = key
				}
				return true, nil
			}
		}

		if key == result.Marker {
			return true, nil
		}
		if result.IsFull() {
			result.IsTruncated = true
			return false, nil
		}

		object, innerErr := lazy.Get()
		if innerErr != nil {
			return false, innerErr
		}

		innerErr = srv.checkStore(object, s3.NewResource(bucket, key))
		if innerErr != nil {
			return false, innerErr
		}
		if result.EncodingType == encodingTypeURL {
			key = urlEncodePath(key)
		}
		result.NextMarker = ""

		result.Contents = append(result.Contents, s3.ContentResult{
			Key:          key,
			LastModified: object.LastModified.Format(time.RFC3339),
			ETag:         s3.NewETag(object.ContentMD5),
			Size:         object.Size,
			Owner:        s3.OwnerResult{},
			StorageClass: "STANDARD",
		})
		return true, nil
	})
	if err != nil {
		if err == meta.ErrBucketNotFound {
			return s3.NoSuchBucket(bucket)
		}
		return s3.InternalError(err)
	}
	if !result.IsTruncated {
		result.NextMarker = ""
	}

	return result
}

func (srv bucketOps) checkStore(data meta.ObjectData, resource s3.Resource) (err error) {
	info, err := srv.store.Info(resource)
	if err != nil {
		return
	}
	if info.Size() != data.Size {
		log.Printf("size mismatch: %s, db(%v), fs(%v)", resource, data.Size, info.Size())
		data.Size = info.Size()
		if data.ContentMD5, err = srv.store.MD5(resource); err != nil {
			return
		}
		if err = srv.db.Put(resource, data); err != nil {
			return
		}
	}
	return
}

var labelRegex = regexp.MustCompile("^[a-z0-9]([a-z0-9-]*[a-z0-9])?$")

func validateBucket(bucket string) (response s3.Response) {
	switch {
	case len(bucket) < 3:
		return s3.InvalidBucketName("BucketName too short")
	case len(bucket) > 63:
		return s3.InvalidBucketName("BucketName too long")
	}

	for _, label := range strings.Split(bucket, ".") {
		if !labelRegex.Match([]byte(label)) {
			return s3.InvalidBucketName("BucketName not formated correctly")
		}
	}

	return
}

const maxListBucketKeys int = 1000

func parseMaxKeys(value string) (maxKeys int, err error) {
	if value == "" {
		return maxListBucketKeys, nil
	}
	if maxKeys, err = strconv.Atoi(value); err != nil {
		return
	}
	switch {
	case maxKeys < 0:
		return 0, errInvalidMaxKeys
	case maxKeys > 2147483647:
		return 0, errInvalidMaxKeys
	case maxKeys > maxListBucketKeys:
		return maxListBucketKeys, nil
	default:
		return
	}
}

func urlEncodePath(path string) string {
	return (&url.URL{Path: path}).EscapedPath()
}
