package ops

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"time"

	"github.com/ophymx/s3d/internal/blob"
	"github.com/ophymx/s3d/internal/clock"
	"github.com/ophymx/s3d/internal/meta"
	"github.com/ophymx/s3d/internal/s3"
)

type objectOps struct {
	db    meta.DB
	store blob.Store
	clock clock.Clock
}

func NewObject(db meta.DB, store blob.Store, clock clock.Clock) ObjectOperations {
	return objectOps{db: db, store: store, clock: clock}
}

func (srv objectOps) Put(resource s3.Resource, contentType string, body io.ReadCloser) s3.Response {
	_, err := srv.db.Get(resource)
	if err == meta.ErrBucketNotFound {
		return s3.NoSuchBucket(resource.Bucket())
	}

	writer, err := srv.store.Create(resource)
	if err != nil {
		return s3.InternalError(err)
	}
	defer writer.Close()

	digest := md5.New()

	size, err := io.Copy(io.MultiWriter(writer, digest), body)
	if err != nil {
		return s3.InternalError(err)
	}

	contentMD5 := hex.EncodeToString(digest.Sum(nil))
	err = srv.db.Put(resource, meta.ObjectData{
		ContentMD5:   contentMD5,
		Size:         size,
		LastModified: srv.clock.Now(),
		ContentType:  contentType,
	})
	if err != nil {
		defer srv.store.Delete(resource)
		if err == meta.ErrBucketNotFound {
			return s3.NoSuchBucket(resource.Bucket())
		}
		return s3.InternalError(err)
	}

	return s3.Created(s3.NewETag(contentMD5))
}

func (srv objectOps) Copy(src, dst s3.Resource) s3.Response {
	log.Printf("Copy: %s -> %s", src, dst)

	objMeta, err := srv.db.Get(src)
	if err != nil {
		switch err {
		case meta.ErrBucketNotFound:
			return s3.NoSuchBucket(src.Bucket())
		case meta.ErrKeyNotFound:
			return s3.NoSuchKey(src.Key())
		default:
			return s3.InternalError(err)
		}
	}

	err = srv.store.Copy(src, dst)
	if srv.store.IsNoSuchKey(err) {
		return s3.NoSuchKey(src.Key())
	}
	if err != nil {
		return s3.InternalError(err)
	}

	objMeta.LastModified = srv.clock.Now()
	err = srv.db.Put(dst, objMeta)
	if err != nil {
		defer srv.store.Delete(dst)
		if err == meta.ErrBucketNotFound {
			return s3.NoSuchBucket(dst.Bucket())
		}
		return s3.InternalError(err)
	}

	return s3.CopyObjectResult{
		ETag:         s3.NewETag(objMeta.ContentMD5),
		LastModified: objMeta.LastModified.Format(time.RFC3339),
	}
}

// Get fetches the object and metadata.
func (srv objectOps) Get(resource s3.Resource) s3.Response {
	return srv.get(resource, false)
}

// Head fetches object metadata but not the object itself.
func (srv objectOps) Head(resource s3.Resource) s3.Response {
	return srv.get(resource, true)
}

func (srv objectOps) get(resource s3.Resource, head bool) s3.Response {
	objMeta, err := srv.db.Get(resource)
	if err != nil {
		switch err {
		case meta.ErrBucketNotFound:
			return s3.NoSuchBucket(resource.Bucket())
		case meta.ErrKeyNotFound:
			return s3.NoSuchKey(resource.Key())
		default:
			return s3.InternalError(err)
		}
	}

	info, err := srv.store.Info(resource)
	if srv.store.IsNoSuchKey(err) {
		return s3.NoSuchKey(resource.Key())
	}
	if err != nil {
		return s3.InternalError(err)
	}

	var file io.ReadCloser
	if !head {
		file, err = srv.store.Get(resource)
		if srv.store.IsNoSuchKey(err) {
			return s3.NoSuchKey(resource.Key())
		}
		if err != nil {
			return s3.InternalError(err)
		}
	}

	return s3.Object{
		File:          file,
		ContentLength: info.Size(),
		ETag:          s3.NewETag(objMeta.ContentMD5),
		ContentType:   objMeta.ContentType,
		LastModified:  objMeta.LastModified.Format(time.RFC3339),
		CacheControl:  objMeta.CacheControl,
		UserDefined:   objMeta.UserDefined,
		VersionID:     objMeta.VersionID,
	}
}

// Delete deletes object at bucket/key
func (srv objectOps) Delete(resource s3.Resource) s3.Response {
	_, err := srv.db.Get(resource)
	if err == meta.ErrBucketNotFound {
		return s3.NoSuchBucket(resource.Bucket())
	}

	err = srv.store.Delete(resource)
	if err != nil {
		if !srv.store.IsNoSuchKey(err) {
			return s3.InternalError(err)
		}
	}

	srv.db.Delete(resource)

	return s3.NoContent()
}
