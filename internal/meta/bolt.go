package meta

import (
	"log"

	"github.com/boltdb/bolt"
)

const (
	bucketMetadataKey = "%%%%meta%%%%"
)

type boltDB struct {
	bdb      *bolt.DB
	encoding Encoding
}

// NewDB returns a DB backed by bolt.DB.
func NewDB(root string, encoding Encoding) (db DB, err error) {
	bdb, err := bolt.Open(root, 0640, nil)
	if err != nil {
		return
	}
	return boltDB{bdb: bdb, encoding: encoding}, nil
}

func (db boltDB) Get(target Target) (data ObjectData, err error) {
	bucket := target.Bucket()
	key := target.Key()
	err = db.bdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			log.Printf("Get: bucket not found: %s", bucket)
			return ErrBucketNotFound
		}

		objBytes := b.Get([]byte(key))
		if objBytes == nil {
			log.Printf("Get: key not found: (%s) %s", bucket, key)
			return ErrKeyNotFound
		}
		data, err = db.encoding.DecodeObject(objBytes)
		return err
	})
	return
}

func (db boltDB) Put(target Target, data ObjectData) error {
	bucket := target.Bucket()
	key := target.Key()
	return db.bdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			log.Printf("Put: bucket not found: %s", bucket)
			return ErrBucketNotFound
		}

		objBytes, err := db.encoding.EncodeObject(data)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), objBytes)
	})
}

func (db boltDB) Delete(target Target) error {
	bucket := target.Bucket()
	key := target.Key()
	return db.bdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			log.Printf("Delete: bucket not found: %s", bucket)
			return ErrBucketNotFound
		}
		return b.Delete([]byte(key))
	})
}

func (db boltDB) CreateBucket(bucket string, data BucketData) error {
	return db.bdb.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(bucket))
		if err == bolt.ErrBucketExists {
			return nil
		} else if err != nil {
			return err
		}

		dataBytes, err := db.encoding.EncodeBucket(data)
		if err != nil {
			return err
		}

		return b.Put([]byte(bucketMetadataKey), dataBytes)
	})
}

func (db boltDB) DeleteBucket(bucket string) error {
	return db.bdb.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucket))
		if err == bolt.ErrBucketNotFound {
			return ErrBucketNotFound
		}
		return err
	})
}

func (db boltDB) ListBuckets() (buckets []Bucket, err error) {
	err = db.bdb.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			dataBytes := b.Get([]byte(bucketMetadataKey))
			if len(dataBytes) == 0 {
				return ErrMissingBucketMetadata
			}

			data, encErr := db.encoding.DecodeBucket(dataBytes)
			if encErr != nil {
				return encErr
			}

			buckets = append(buckets, Bucket{
				Name:     string(name),
				Metadata: data,
			})

			return nil
		})
	})
	if buckets == nil {
		buckets = []Bucket{}
	}
	return
}

func (db boltDB) ForEachInBucket(bucket, seek string, fn ForEachFunc) error {
	return db.bdb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			log.Printf("Put: bucket not found: %s", bucket)
			return ErrBucketNotFound
		}

		c := b.Cursor()
		for k, v := c.Seek([]byte(seek)); k != nil; k, v = c.Next() {
			key := string(k)
			if key == bucketMetadataKey {
				continue
			}
			next, err := fn(key, lazyObject{data: v, encoding: db.encoding})
			if err != nil {
				return err
			}
			if !next {
				break
			}
		}
		return nil
	})
}

func (db boltDB) Close() error {
	return db.bdb.Close()
}

type lazyObject struct {
	data     []byte
	encoding Encoding
}

func (o lazyObject) Get() (ObjectData, error) {
	return o.encoding.DecodeObject(o.data)
}
