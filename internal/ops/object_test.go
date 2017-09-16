package ops_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ophymx/s3d/internal/fakes"
	"github.com/ophymx/s3d/internal/fixtures"
	"github.com/ophymx/s3d/internal/meta"
	"github.com/ophymx/s3d/internal/ops"
	"github.com/ophymx/s3d/internal/s3"
)

func stringBody(content string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(content))
}

var _ = Describe("ObjectService", func() {
	var (
		db    *fakes.DB
		store *fakes.Store
		clock *fakes.Clock
		srv   ops.ObjectOperations
	)
	BeforeEach(func() {
		db = fakes.NewDB()
		store = fakes.NewStore()
		clock = fakes.NewClock()
		srv = ops.NewObject(db, store, clock)
	})

	Describe("Get", func() {
		Context("when bucket does not exist", func() {
			It("returns a 404 not found response", func() {
				Expect(srv.Get(s3.NewResource("foo", "bar.txt"))).To(Equal(s3.NoSuchBucket("foo")))
			})
		})
		Context("when bucket exists", func() {
			BeforeEach(func() {
				db.CreateBucket("foo", meta.BucketData{CreationDate: fixtures.Time1})
				store.CreateBucket("foo")
			})
			Context("but key does not", func() {
				It("returns a 404 not found response", func() {
					Expect(srv.Get(s3.NewResource("foo", "bar.txt"))).To(Equal(s3.NoSuchKey("bar.txt")))
				})
			})
			Context("and key exists", func() {
				BeforeEach(func() {
					db.Put(s3.NewResource("foo", "bar.txt"), meta.ObjectData{
						ContentMD5:   "content-md5[baz]",
						LastModified: fixtures.Time1,
					})
					store.Buckets["foo"]["bar.txt"] = bytes.NewBufferString("baz")
				})
				It("returns an object response", func() {
					Expect(srv.Get(s3.NewResource("foo", "bar.txt"))).To(Equal(s3.Object{
						ContentLength: 3,
						LastModified:  "2014-05-06T03:02:01Z",
						ETag:          "content-md5[baz]",
						File:          ioutil.NopCloser(bytes.NewBufferString("baz")),
					}))
				})
			})
		})
	})

	Describe("Head", func() {
		Context("when bucket does not exist", func() {
			It("returns a 404 not found response", func() {
				Expect(srv.Head(s3.NewResource("foo", "bar.txt"))).To(Equal(s3.NoSuchBucket("foo")))
			})
		})
		Context("when bucket exists", func() {
			BeforeEach(func() {
				db.CreateBucket("foo", meta.BucketData{CreationDate: fixtures.Time1})
				store.CreateBucket("foo")
			})
			Context("but key does not", func() {
				It("returns a 404 not found response", func() {
					Expect(srv.Head(s3.NewResource("foo", "bar.txt"))).To(Equal(s3.NoSuchKey("bar.txt")))
				})
			})
			Context("and key exists", func() {
				BeforeEach(func() {
					db.Put(s3.NewResource("foo", "bar.txt"), meta.ObjectData{
						ContentMD5:   "content-md5[baz]",
						LastModified: fixtures.Time1,
					})
					store.Buckets["foo"]["bar.txt"] = bytes.NewBufferString("baz")
				})
				It("returns an object response", func() {
					Expect(srv.Head(s3.NewResource("foo", "bar.txt"))).To(Equal(s3.Object{
						ContentLength: 3,
						LastModified:  "2014-05-06T03:02:01Z",
						ETag:          "content-md5[baz]",
					}))
				})
			})
		})
	})

	Describe("Put", func() {
		Context("when bucket does not exist", func() {
			It("returns a 404 not found response", func() {
				Expect(srv.Put(s3.NewResource("foo", "bar.txt"), "plain/text", stringBody("baz"))).
					To(Equal(s3.NoSuchBucket("foo")))
			})
		})

		Context("when object does not exist", func() {
			BeforeEach(func() {
				db.CreateBucket("foo", meta.BucketData{CreationDate: fixtures.Time1})
				store.CreateBucket("foo")
			})

			It("creates the object", func() {
				Expect(srv.Put(s3.NewResource("foo", "bar.txt"), "plain/text", stringBody("baz"))).
					To(Equal(s3.Created(s3.NewETag("73feffa4b7f6bb68e44cf984c85f6e88"))))
			})
		})
	})

	Describe("Copy", func() {
		Context("when source bucket does not exist", func() {
			It("returns a 404 not found", func() {
				Expect(srv.Copy(s3.NewResource("foo1", "bar1.txt"), s3.NewResource("foo2", "bar2.txt"))).
					To(Equal(s3.NoSuchBucket("foo1")))
			})
		})
		Context("when source bucket exists", func() {
			BeforeEach(func() {
				db.CreateBucket("foo1", meta.BucketData{CreationDate: fixtures.Time1})
				store.CreateBucket("foo1")
			})
			Context("but source key does not", func() {
				It("returns a 404 not found", func() {
					Expect(srv.Copy(s3.NewResource("foo1", "bar1.txt"), s3.NewResource("foo2", "bar2.txt"))).
						To(Equal(s3.NoSuchKey("bar1.txt")))
				})
			})
			Context("and source key exists", func() {
				BeforeEach(func() {
					db.Put(s3.NewResource("foo1", "bar1.txt"), meta.ObjectData{
						ContentMD5:   "content-md5[baz]",
						LastModified: fixtures.Time1,
					})
					store.Buckets["foo1"]["bar1.txt"] = bytes.NewBufferString("baz")
				})
				Context("but destination bucket does not", func() {
					It("returns a 404 not found", func() {
						Expect(srv.Copy(s3.NewResource("foo1", "bar1.txt"), s3.NewResource("foo2", "bar2.txt"))).
							To(Equal(s3.NoSuchBucket("foo2")))
					})
				})
				Context("and destination bucket exists", func() {
					BeforeEach(func() {
						db.CreateBucket("foo2", meta.BucketData{CreationDate: fixtures.Time1})
						store.CreateBucket("foo2")
					})
					It("copies the file", func() {
						clock.Add(fixtures.Time2)
						Expect(srv.Copy(s3.NewResource("foo1", "bar1.txt"), s3.NewResource("foo2", "bar2.txt"))).
							To(Equal(s3.CopyObjectResult{
								LastModified: "2015-06-07T04:03:02Z",
								ETag:         "content-md5[baz]",
							}))
						Expect(db.Buckets["foo2"].Objects["bar2.txt"]).To(Equal(meta.ObjectData{
							LastModified: fixtures.Time2,
							ContentMD5:   "content-md5[baz]",
						}))
						Expect(store.Buckets["foo2"]["bar2.txt"]).To(Equal(store.Buckets["foo1"]["bar1.txt"]))
					})
				})
			})
		})
	})

	Describe("Delete", func() {
		Context("when bucket does not exist", func() {
			It("returns a 404 error", func() {
				Expect(srv.Delete(s3.NewResource("foo", "bar.txt"))).
					To(Equal(s3.NoSuchBucket("foo")))
			})
		})
		Context("when bucket exists", func() {
			BeforeEach(func() {
				db.CreateBucket("foo", meta.BucketData{CreationDate: fixtures.Time1})
				store.CreateBucket("foo")
			})
			Context("but object does not", func() {
				It("does not return an error", func() {
					Expect(srv.Delete(s3.NewResource("foo", "bar.txt"))).
						To(Equal(s3.NoContent()))
				})
			})
			Context("and object exists", func() {
				BeforeEach(func() {
					db.Put(s3.NewResource("foo", "bar.txt"), meta.ObjectData{
						ContentMD5:   "content-md5[baz]",
						LastModified: fixtures.Time1,
					})
					store.Buckets["foo"]["bar.txt"] = bytes.NewBufferString("baz")
				})
				It("deletes the object", func() {
					Expect(srv.Delete(s3.NewResource("foo", "bar.txt"))).
						To(Equal(s3.NoContent()))
					Expect(db.Buckets["foo"].Objects).To(BeEmpty())
					Expect(store.Buckets["foo"]).To(BeEmpty())
				})
			})
		})
	})
})
