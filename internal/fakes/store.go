package fakes

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"

	"github.com/ophymx/s3d/internal/blob"
)

var (
	ErrNoSuchBucket = errors.New("no such bucket")
	ErrNoSuchKey    = errors.New("no such key")
)

type Store struct {
	Buckets map[string]map[string]*bytes.Buffer
}

type bytesSize int64

func (b bytesSize) Size() int64 {
	return int64(b)
}

type NopWriteCloser struct {
	io.Writer
}

func (NopWriteCloser) Close() error {
	return nil
}

func NewStore() *Store {
	return &Store{Buckets: make(map[string]map[string]*bytes.Buffer)}
}

func (s *Store) Copy(src, dst blob.Resource) (err error) {
	if bucket, found := s.Buckets[src.Bucket()]; found {
		if b, found := bucket[src.Key()]; found {
			s.CreateBucket(dst.Bucket())
			s.Buckets[dst.Bucket()][dst.Key()] = bytes.NewBuffer(b.Bytes())
			return nil
		}
		return ErrNoSuchKey
	}
	return ErrNoSuchBucket
}

func (s *Store) Create(resource blob.Resource) (writer io.WriteCloser, err error) {
	s.CreateBucket(resource.Bucket())
	bucket := s.Buckets[resource.Bucket()]
	w := new(bytes.Buffer)
	bucket[resource.Key()] = w
	return NopWriteCloser{w}, nil
}

func (s *Store) Delete(resource blob.Resource) (err error) {
	if bucket, found := s.Buckets[resource.Bucket()]; found {
		if _, found := bucket[resource.Key()]; found {
			delete(bucket, resource.Key())
			return nil
		}
		return ErrNoSuchKey
	}
	return ErrNoSuchBucket
}

func (s *Store) CreateBucket(bucket string) (err error) {
	if _, found := s.Buckets[bucket]; !found {
		s.Buckets[bucket] = make(map[string]*bytes.Buffer)
	}
	return
}

func (s *Store) DeleteBucket(bucket string) (err error) {
	delete(s.Buckets, bucket)
	return
}

func (s *Store) Get(resource blob.Resource) (reader io.ReadCloser, err error) {
	if bucket, found := s.Buckets[resource.Bucket()]; found {
		if b, found := bucket[resource.Key()]; found {
			return ioutil.NopCloser(b), nil
		}
		return nil, ErrNoSuchKey
	}
	return nil, ErrNoSuchBucket
}

func (s *Store) Info(resource blob.Resource) (info blob.Info, err error) {
	if bucket, found := s.Buckets[resource.Bucket()]; found {
		if b, found := bucket[resource.Key()]; found {
			return bytesSize(b.Len()), nil
		}
		return nil, ErrNoSuchKey
	}
	return nil, ErrNoSuchBucket
}

func (s *Store) IsNoSuchKey(err error) bool {
	return err == ErrNoSuchBucket || err == ErrNoSuchKey
}

func (s *Store) MD5(resource blob.Resource) (chksum string, err error) {
	var reader io.Reader
	if reader, err = s.Get(resource); err != nil {
		return
	}

	digest := md5.New()
	_, err = io.Copy(digest, reader)
	if err != nil {
		return
	}

	return hex.EncodeToString(digest.Sum(nil)), nil
}
