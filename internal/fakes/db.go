package fakes

import (
	"sort"

	"github.com/ophymx/s3d/internal/meta"
)

type Bucket struct {
	Meta    meta.BucketData
	Objects map[string]meta.ObjectData
}

func NewBucket(data meta.BucketData) *Bucket {
	return &Bucket{
		Meta:    data,
		Objects: make(map[string]meta.ObjectData),
	}
}

type lazy struct {
	meta.ObjectData
}

func (l lazy) Get() (meta.ObjectData, error) {
	return l.ObjectData, nil
}

type DB struct {
	Buckets map[string]*Bucket
}

func NewDB() *DB {
	return &DB{Buckets: make(map[string]*Bucket)}
}

func (db *DB) Get(target meta.Target) (data meta.ObjectData, err error) {
	if bucket, found := db.Buckets[target.Bucket()]; found {
		if obj, found := bucket.Objects[target.Key()]; found {
			return obj, nil
		}
		err = meta.ErrKeyNotFound
		return
	}
	err = meta.ErrBucketNotFound
	return
}

func (db *DB) Put(target meta.Target, data meta.ObjectData) (err error) {
	bucket, found := db.Buckets[target.Bucket()]
	if !found {
		return meta.ErrBucketNotFound
	}
	bucket.Objects[target.Key()] = data
	return
}

func (db *DB) Delete(target meta.Target) (err error) {
	if bucket, found := db.Buckets[target.Bucket()]; found {
		if _, found := bucket.Objects[target.Key()]; found {
			delete(bucket.Objects, target.Key())
			return
		}
		return meta.ErrKeyNotFound
	}
	return meta.ErrBucketNotFound
}

func (db *DB) CreateBucket(bucket string, data meta.BucketData) (err error) {
	if _, found := db.Buckets[bucket]; found {
		return
	}
	db.Buckets[bucket] = NewBucket(data)
	return
}

func (db *DB) DeleteBucket(bucket string) (err error) {
	if _, found := db.Buckets[bucket]; found {
		delete(db.Buckets, bucket)
		return
	}
	return meta.ErrBucketNotFound
}

func (db *DB) ListBuckets() (buckets []meta.Bucket, err error) {
	names := make([]string, 0, len(db.Buckets))
	for name := range db.Buckets {
		names = append(names, name)
	}
	sort.Strings(names)
	buckets = make([]meta.Bucket, 0, len(db.Buckets))
	for _, name := range names {
		buckets = append(buckets, meta.Bucket{
			Name:     name,
			Metadata: db.Buckets[name].Meta,
		})
	}
	return
}

func (db *DB) ForEachInBucket(name, seek string, forEach meta.ForEachFunc) (err error) {
	bucket, found := db.Buckets[name]
	if !found {
		return meta.ErrBucketNotFound
	}
	keys := make([]string, 0, len(bucket.Objects))
	for key := range bucket.Objects {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var next bool
	for _, key := range keys {
		if key < seek {
			continue
		}
		next, err = forEach(key, lazy{bucket.Objects[key]})
		if err != nil {
			return
		}
		if !next {
			break
		}
	}
	return
}

func (db *DB) Close() (err error) {
	return
}
