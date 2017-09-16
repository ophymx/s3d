package ops_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ophymx/s3d/internal/fakes"
	"github.com/ophymx/s3d/internal/fixtures"
	"github.com/ophymx/s3d/internal/meta"
	"github.com/ophymx/s3d/internal/ops"
	"github.com/ophymx/s3d/internal/s3"
)

var _ = Describe("ServiceService", func() {
	var (
		db  *fakes.DB
		srv ops.ServiceOperations
	)
	BeforeEach(func() {
		db = fakes.NewDB()
		srv = ops.NewService(db)
	})

	Describe("ListBuckets", func() {
		Context("when there are no buckets", func() {
			It("returns an empty list of buckets", func() {
				Expect(srv.ListBuckets()).To(Equal(s3.ListAllMyBucketsResult{}))
			})
		})
		Context("when buckets exist", func() {
			BeforeEach(func() {
				db.Buckets = map[string]*fakes.Bucket{
					"foo": fakes.NewBucket(meta.BucketData{CreationDate: fixtures.Time1}),
					"bar": fakes.NewBucket(meta.BucketData{CreationDate: fixtures.Time2}),
				}
			})
			It("returns a list of buckets", func() {
				Expect(srv.ListBuckets()).To(Equal(s3.ListAllMyBucketsResult{
					Buckets: []s3.ListAllMyBucketsResultBucket{
						{
							Name:         "bar",
							CreationDate: "2015-06-07T04:03:02Z",
						},
						{
							Name:         "foo",
							CreationDate: "2014-05-06T03:02:01Z",
						},
					},
				}))
			})
		})
	})
})
