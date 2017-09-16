package blob

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
)

const (
	dirMode = os.FileMode(0750)
)

type fsStore struct {
	root string
}

// NewFsStore provides filesystem backed Store at root.
func NewFsStore(root string) (store Store) {
	os.MkdirAll(root, dirMode)
	return fsStore{root: root}
}

func (fs fsStore) Get(resource Resource) (object io.ReadCloser, err error) {
	return os.Open(fs.path(resource))
}

func (fs fsStore) Copy(src, dst Resource) (err error) {
	if src.Bucket() == dst.Bucket() && src.Key() == dst.Key() {
		return
	}

	if err = fs.mkParent(dst); err != nil {
		return
	}

	// unlink before writting
	if err = fs.Delete(dst); err != nil && !fs.IsNoSuchKey(err) {
		return
	}

	return os.Link(fs.path(src), fs.path(dst))
}

func (fs fsStore) Create(resource Resource) (writer io.WriteCloser, err error) {
	if err = fs.mkParent(resource); err != nil {
		return
	}

	// unlink before writting
	if err = fs.Delete(resource); err != nil && !fs.IsNoSuchKey(err) {
		return
	}

	return os.Create(fs.path(resource))
}

func (fs fsStore) Delete(resource Resource) (err error) {
	return os.Remove(fs.path(resource))
}

func (fs fsStore) CreateBucket(bucket string) (err error) {
	return
}

func (fs fsStore) DeleteBucket(bucket string) (err error) {
	return os.RemoveAll(filepath.Join(fs.root, bucket))
}

func (fs fsStore) Info(resource Resource) (info Info, err error) {
	return os.Stat(fs.path(resource))
}

func (fs fsStore) IsNoSuchKey(err error) bool {
	return os.IsNotExist(err)
}

func (fs fsStore) MD5(resource Resource) (result string, err error) {
	file, err := os.Open(fs.path(resource))
	if err != nil {
		return
	}
	defer file.Close()

	digest := md5.New()
	_, err = io.Copy(digest, file)
	if err != nil {
		return
	}

	return hex.EncodeToString(digest.Sum(nil)), nil
}

func (fs fsStore) mkParent(resource Resource) (err error) {
	return os.MkdirAll(filepath.Dir(fs.path(resource)), dirMode)
}

func (fs fsStore) path(resource Resource) string {
	return filepath.Join(fs.root, resource.Bucket(), resource.Key())
}
