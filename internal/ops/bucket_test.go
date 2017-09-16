package ops_test

import (
	"bytes"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ophymx/s3d/internal/fakes"
	"github.com/ophymx/s3d/internal/fixtures"
	"github.com/ophymx/s3d/internal/meta"
	"github.com/ophymx/s3d/internal/ops"
	"github.com/ophymx/s3d/internal/s3"
)

var (
	t1 = time.Date(2014, 5, 6, 3, 2, 1, 0, time.UTC)
	t2 = time.Date(2015, 6, 7, 4, 3, 2, 0, time.UTC)
)

var _ = Describe("BucketService", func() {
	var (
		db    *fakes.DB
		store *fakes.Store
		clock *fakes.Clock
		srv   ops.BucketOperations
	)
	BeforeEach(func() {
		db = fakes.NewDB()
		store = fakes.NewStore()
		clock = fakes.NewClock(t1, t2)
		srv = ops.NewBucket(db, store, clock)
	})

	Describe("Create", func() {
		Context("when a bucket named 'foo' does not exist", func() {
			It("can create a bucket named 'foo'", func() {
				Expect(srv.Create("foo")).To(Equal(s3.NoContent()))
				Expect(db.Buckets).To(Equal(map[string]*fakes.Bucket{
					"foo": fakes.NewBucket(meta.BucketData{CreationDate: t1}),
				}))
				Expect(store.Buckets).To(Equal(map[string]map[string]*bytes.Buffer{
					"foo": map[string]*bytes.Buffer{},
				}))
			})
		})
		Context("when a bucket named 'foo' does exist", func() {
			BeforeEach(func() { srv.Create("foo") })
			It("does not error trying to create a bucket named 'foo'", func() {
				Expect(srv.Create("foo")).To(Equal(s3.NoContent()))
				Expect(db.Buckets).To(Equal(map[string]*fakes.Bucket{
					"foo": fakes.NewBucket(meta.BucketData{CreationDate: t1}),
				}))
				Expect(store.Buckets).To(Equal(map[string]map[string]*bytes.Buffer{
					"foo": map[string]*bytes.Buffer{},
				}))
			})
		})
		Context("when bucket name is invalid", func() {
			It("errors with InvalidBucketName", func() {
				Expect(srv.Create("fo")).To(Equal(s3.InvalidBucketName("BucketName too short")))
				Expect(srv.Create("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")).
					To(Equal(s3.InvalidBucketName("BucketName too long")))
				Expect(srv.Create("-abcd")).
					To(Equal(s3.InvalidBucketName("BucketName not formated correctly")))
				Expect(srv.Create("foo.-abcd")).
					To(Equal(s3.InvalidBucketName("BucketName not formated correctly")))
				Expect(srv.Create("foo.abcd-")).
					To(Equal(s3.InvalidBucketName("BucketName not formated correctly")))
				Expect(db.Buckets).To(BeEmpty())
				Expect(store.Buckets).To(BeEmpty())
			})
		})
	})
	Describe("Delete", func() {
		Context("when a bucket named 'foo' does exist", func() {
			BeforeEach(func() { srv.Create("foo") })
			It("can delete the bucket", func() {
				Expect(srv.Delete("foo")).To(Equal(s3.NoContent()))
				Expect(db.Buckets).To(BeEmpty())
				Expect(store.Buckets).To(BeEmpty())
			})
		})
		Context("when a bucket named 'foo' does not exist", func() {
			It("errors with NoSuchBucket", func() {
				Expect(srv.Delete("foo")).To(Equal(s3.NoSuchBucket("foo")))
			})
		})
	})
	Describe("ListBucket", func() {
		Context("when the bucket does not exist", func() {
			It("erros with NoSuchBucket", func() {
				Expect(srv.ListBucket("foo", url.Values{})).To(Equal(s3.NoSuchBucket("foo")))
			})
		})
		Context("when the bucket exists", func() {
			BeforeEach(func() { srv.Create("foo") })
			Context("and it is empty", func() {
				It("returns an empty list result", func() {
					Expect(srv.ListBucket("foo", url.Values{})).To(Equal(s3.ListBucketResult{
						Name:    "foo",
						MaxKeys: 1000,
					}))
				})
			})

			Context("and it has 2 object", func() {
				BeforeEach(func() {
					db.Buckets["foo"].Objects["bar/Example file.txt"] = meta.ObjectData{
						ContentMD5:   "content-md5[baz]",
						Size:         3,
						ContentType:  "plain/text",
						LastModified: fixtures.Time1,
					}
					db.Buckets["foo"].Objects["example.jpeg"] = meta.ObjectData{
						ContentMD5:   "content-md5[]",
						Size:         0,
						ContentType:  "image/jpeg",
						LastModified: fixtures.Time2,
					}
					store.Buckets["foo"]["bar/Example file.txt"] = bytes.NewBufferString("baz")
					store.Buckets["foo"]["example.jpeg"] = bytes.NewBuffer([]byte{})
				})
				Context("and not sending any parameters", func() {
					It("returns a list with full key", func() {
						Expect(srv.ListBucket("foo", url.Values{})).To(Equal(s3.ListBucketResult{
							Name:    "foo",
							MaxKeys: 1000,
							Contents: []s3.ContentResult{
								{
									Key:          "bar/Example file.txt",
									LastModified: "2014-05-06T03:02:01Z",
									ETag:         s3.ETag("content-md5[baz]"),
									Size:         3,
									StorageClass: "STANDARD",
								},
								{
									Key:          "example.jpeg",
									LastModified: "2015-06-07T04:03:02Z",
									ETag:         s3.ETag("content-md5[]"),
									Size:         0,
									StorageClass: "STANDARD",
								},
							},
						}))
					})
				})

				Context("and setting delimiter to '/'", func() {
					It("returns a list with full key", func() {
						Expect(srv.ListBucket("foo", url.Values{"delimiter": []string{"/"}})).To(Equal(s3.ListBucketResult{
							Name:           "foo",
							MaxKeys:        1000,
							Delimiter:      "/",
							CommonPrefixes: s3.CommonPrefixes{"bar/"},
							Contents: []s3.ContentResult{
								{
									Key:          "example.jpeg",
									LastModified: "2015-06-07T04:03:02Z",
									ETag:         s3.ETag("content-md5[]"),
									Size:         0,
									StorageClass: "STANDARD",
								},
							},
						}))
					})
				})

				Context("and setting encoding type to 'url'", func() {
					It("returns a list with url encoded keys", func() {
						Expect(srv.ListBucket("foo", url.Values{"encoding-type": []string{"url"}})).To(Equal(s3.ListBucketResult{
							Name:         "foo",
							MaxKeys:      1000,
							EncodingType: "url",
							Contents: []s3.ContentResult{
								{
									Key:          "bar/Example%20file.txt",
									LastModified: "2014-05-06T03:02:01Z",
									ETag:         s3.ETag("content-md5[baz]"),
									Size:         3,
									StorageClass: "STANDARD",
								},
								{
									Key:          "example.jpeg",
									LastModified: "2015-06-07T04:03:02Z",
									ETag:         s3.ETag("content-md5[]"),
									Size:         0,
									StorageClass: "STANDARD",
								},
							},
						}))
					})
				})

				Context("and setting prefix to 'bar/'", func() {
					It("returns a list of onyl object who's key have prefix", func() {
						Expect(srv.ListBucket("foo", url.Values{"prefix": []string{"bar/"}})).To(Equal(s3.ListBucketResult{
							Name:    "foo",
							MaxKeys: 1000,
							Prefix:  "bar/",
							Contents: []s3.ContentResult{
								{
									Key:          "bar/Example file.txt",
									LastModified: "2014-05-06T03:02:01Z",
									ETag:         s3.ETag("content-md5[baz]"),
									Size:         3,
									StorageClass: "STANDARD",
								},
							},
						}))
					})
				})

				Context("and setting max-keys to 1", func() {
					It("returns a list of onyl object who's key have prefix", func() {
						Expect(srv.ListBucket("foo", url.Values{"max-keys": []string{"1"}})).To(Equal(s3.ListBucketResult{
							Name:        "foo",
							MaxKeys:     1,
							IsTruncated: true,
							Contents: []s3.ContentResult{
								{
									Key:          "bar/Example file.txt",
									LastModified: "2014-05-06T03:02:01Z",
									ETag:         s3.ETag("content-md5[baz]"),
									Size:         3,
									StorageClass: "STANDARD",
								},
							},
						}))
					})
				})

				Context("and setting max-keys to invalid value", func() {
					It("returns an InvalidArgument response", func() {
						Expect(srv.ListBucket("foo", url.Values{"max-keys": []string{"foo"}})).
							To(Equal(s3.InvalidArgument("Argument maxKeys must be an integer between 0 and 2147483647", "maxKeys", "foo")))
						Expect(srv.ListBucket("foo", url.Values{"max-keys": []string{"-1"}})).
							To(Equal(s3.InvalidArgument("Argument maxKeys must be an integer between 0 and 2147483647", "maxKeys", "-1")))
						Expect(srv.ListBucket("foo", url.Values{"max-keys": []string{"2147483649"}})).
							To(Equal(s3.InvalidArgument("Argument maxKeys must be an integer between 0 and 2147483647", "maxKeys", "2147483649")))
					})
				})

				Context("and setting max-keys to 1001", func() {
					It("returns a list with max-keys set to 1000", func() {
						Expect(srv.ListBucket("foo", url.Values{"max-keys": []string{"1001"}})).To(Equal(s3.ListBucketResult{
							Name:    "foo",
							MaxKeys: 1000,
							Contents: []s3.ContentResult{
								{
									Key:          "bar/Example file.txt",
									LastModified: "2014-05-06T03:02:01Z",
									ETag:         s3.ETag("content-md5[baz]"),
									Size:         3,
									StorageClass: "STANDARD",
								},
								{
									Key:          "example.jpeg",
									LastModified: "2015-06-07T04:03:02Z",
									ETag:         s3.ETag("content-md5[]"),
									Size:         0,
									StorageClass: "STANDARD",
								},
							},
						}))
					})
				})

				Context("and setting marker 'example.jpeg'", func() {
					It("returns a list with max-keys set to 1000", func() {
						Expect(srv.ListBucket("foo", url.Values{"marker": []string{"bar/Example file.txt"}})).To(Equal(s3.ListBucketResult{
							Name:    "foo",
							MaxKeys: 1000,
							Marker:  "bar/Example file.txt",
							Contents: []s3.ContentResult{
								{
									Key:          "example.jpeg",
									LastModified: "2015-06-07T04:03:02Z",
									ETag:         s3.ETag("content-md5[]"),
									Size:         0,
									StorageClass: "STANDARD",
								},
							},
						}))
					})

					Context("and setting delimiter to '/' and max-keys to 1", func() {
						It("returns the common prefix and a next marker of 'bar/'", func() {
							values := url.Values{}
							values.Set("delimiter", "/")
							values.Set("max-keys", "1")
							Expect(srv.ListBucket("foo", values)).To(Equal(s3.ListBucketResult{
								Name:           "foo",
								MaxKeys:        1,
								NextMarker:     "bar/",
								IsTruncated:    true,
								Delimiter:      "/",
								CommonPrefixes: s3.CommonPrefixes{"bar/"},
							}))
						})
					})
				})
			})
		})
	})
})
