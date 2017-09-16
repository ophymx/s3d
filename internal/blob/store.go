package blob

import (
	"io"
)

// Store represents an object store.
// Only implementation so far is a filesystem backed store.
// Initialize with NewFsStore.
type Store interface {

	// Get fetches object from store.
	Get(resource Resource) (object io.ReadCloser, err error)

	// Copy copies object in store.
	Copy(src, dst Resource) (err error)

	// Create creates object in store.
	Create(resource Resource) (writer io.WriteCloser, err error)

	// Delete deletes object in store.
	Delete(resource Resource) (err error)

	// CreateBucket creates bucket in store.
	CreateBucket(bucket string) (err error)

	// DeleteBucket deletes all objects in bucket and bucket in store.
	DeleteBucket(bucket string) (err error)

	// Info fetches object metadata from store.
	Info(resource Resource) (info Info, err error)

	// IsNoSuchKey tests if error is a result of a missing object.
	IsNoSuchKey(err error) bool

	// Compute MD5 of object in store.
	MD5(resource Resource) (md5 string, err error)
}

// Info interface to read metadata of an object in store.
type Info interface {
	Size() int64
}

type Resource interface {
	Bucket() string
	Key() string
}
