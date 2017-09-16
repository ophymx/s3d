package meta_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ophymx/s3d/internal/meta"
	"github.com/ophymx/s3d/internal/s3"
)

var (
	bucketDate = time.Date(2014, 5, 6, 3, 2, 1, 0, time.UTC)
)

var _ = Describe("DB", func() {
	var (
		dbpath string
		db     meta.DB
	)
	BeforeEach(func() { dbpath, db = setup() })
	AfterEach(func() { tearDown(dbpath, db) })

	Describe("ListBuckets", func() {
		Context("when no buckets exist", func() {
			It("returns an empty list of buckets", func() {
				Expect(db.ListBuckets()).To(Equal([]meta.Bucket{}))
			})
		})
	})

	Describe("CreateBucket", func() {
		Context("when bucket does not exist", func() {
			It("can be created", func() {
				Expect(db.CreateBucket("foo", meta.BucketData{CreationDate: bucketDate})).
					ToNot(HaveOccurred())
				Expect(db.ListBuckets()).To(Equal([]meta.Bucket{{
					Name:     "foo",
					Metadata: meta.BucketData{CreationDate: bucketDate},
				}}))
			})
		})
		Context("when bucket already exists", func() {
			BeforeEach(func() { must(db.CreateBucket("foo", meta.BucketData{})) })
			It("does not error", func() {
				Expect(db.CreateBucket("foo", meta.BucketData{CreationDate: bucketDate})).
					ToNot(HaveOccurred())
				Expect(db.ListBuckets()).To(Equal([]meta.Bucket{{
					Name:     "foo",
					Metadata: meta.BucketData{},
				}}))
			})
		})
	})

	Describe("DeleteBucket", func() {
		Context("when a bucket does not already exist", func() {
			It("returns a bucket not found error", func() {
				Expect(db.DeleteBucket("foo")).To(Equal(meta.ErrBucketNotFound))
			})
		})
	})

	Describe("Get", func() {
		Context("when bucket does not exist", func() {
			It("returns a bucket not found error", func() {
				_, err := db.Get(s3.NewResource("foo", "bar"))
				Expect(err).To(Equal(meta.ErrBucketNotFound))
			})
		})
		Context("when bucket exists but key does not", func() {
			BeforeEach(func() { must(db.CreateBucket("foo", meta.BucketData{})) })
			It("returns a key not found error", func() {
				_, err := db.Get(s3.NewResource("foo", "bar"))
				Expect(err).To(Equal(meta.ErrKeyNotFound))
			})
		})
	})

	Describe("Put", func() {
		Context("when bucket is missing", func() {
			It("returns a bucket not found error response", func() {
				Expect(db.Put(s3.NewResource("foo", "bar"), meta.ObjectData{})).
					To(Equal(meta.ErrBucketNotFound))
			})
		})
	})
})
