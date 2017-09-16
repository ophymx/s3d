package ops

import (
	"time"

	"github.com/ophymx/s3d/internal/meta"
	"github.com/ophymx/s3d/internal/s3"
)

type serviceOps struct {
	db meta.DB
}

func NewService(db meta.DB) ServiceOperations {
	return serviceOps{db: db}
}

func (srv serviceOps) ListBuckets() s3.Response {
	result := s3.ListAllMyBucketsResult{
		Owner: s3.OwnerResult{},
	}
	buckets, err := srv.db.ListBuckets()
	if err != nil {
		return s3.InternalError(err)
	}
	for _, bucket := range buckets {
		result.Buckets = append(result.Buckets, s3.ListAllMyBucketsResultBucket{
			Name:         bucket.Name,
			CreationDate: bucket.Metadata.CreationDate.Format(time.RFC3339),
		})
	}
	return result
}
