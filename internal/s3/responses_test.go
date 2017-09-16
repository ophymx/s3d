package s3_test

import (
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ophymx/s3d/internal/fakes"
	"github.com/ophymx/s3d/internal/fixtures"
	"github.com/ophymx/s3d/internal/s3"
)

func fixture(asset string) string {
	value, err := fixtures.Asset("xml/" + asset + ".xml")
	Expect(err).NotTo(HaveOccurred())
	return tidyXML(value)
}

var newlineIndentRegex = regexp.MustCompile("\\n\\s*")

func tidyXML(b []byte) string {
	return string(newlineIndentRegex.ReplaceAll(b, []byte{}))
}

func getBody(capture *fakes.ResponseWriter) string {
	return tidyXML(capture.Buffer.Bytes())
}

var _ = Describe("Responses", func() {
	var (
		resp    s3.Response
		capture *fakes.ResponseWriter
	)

	BeforeEach(func() {
		capture = fakes.NewResponseWriter()
	})

	Describe("AccessControlPolicy", func() {
		Specify("Send", func() {
			resp = s3.AccessControlPolicy{
				Owner: s3.OwnerResult{
					ID:          "75aa57f09aa0c8caeab4f8c24e99d10f8e7faeebf76c078efc7c6caea54ba06a",
					DisplayName: "CustomersName@amazon.com",
				},
				AccessControlList: []s3.Grant{
					s3.NewUserGrant(
						"75aa57f09aa0c8caeab4f8c24e99d10f8e7faeebf76c078efc7c6caea54ba06a",
						"CustomersName@amazon.com",
						s3.PermFull,
					),
				},
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("AccessControlPolicy")))
		})
	})

	Describe("BucketLoggingStatus", func() {
		Specify("Send", func() {
			resp = s3.BucketLoggingStatus{
				LoggingEnabled: &s3.LoggingEnabled{
					TargetBucket: "mybucketlogs",
					TargetPrefix: "mybucket-access_log-/",
					TargetGrants: []s3.Grant{
						s3.NewEmailGrant("user@company.com", s3.PermRead),
					},
				},
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("BucketLoggingStatus")))
		})

		Describe("disabled", func() {
			Specify("Send", func() {
				resp = s3.BucketLoggingStatus{}
				resp.Send(capture)
				Expect(getBody(capture)).To(Equal(fixture("BucketLoggingStatus-disabled")))
			})
		})
	})

	Describe("CompleteMultipartUploadResult", func() {
		Specify("Send", func() {
			resp = s3.CompleteMultipartUploadResult{
				Location: "http://Example-Bucket.s3.amazonaws.com/Example-Object",
				Bucket:   "Example-Bucket",
				Key:      "Example-Object",
				ETag:     s3.NewETag("3858f62230ac3c915f300c664312c11f-9"),
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("CompleteMultipartUploadResult")))
		})
	})

	Describe("CopyPartResult", func() {
		Specify("Send", func() {
			resp = s3.CopyPartResult{
				LastModified: "2009-10-28T22:32:00",
				ETag:         s3.NewETag("9b2cf535f27731c974343645a3985328"),
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("CopyPartResult")))
		})
	})

	Describe("CORSConfiguration", func() {
		Specify("Send", func() {
			resp = s3.CORSConfiguration{
				CORSRules: []s3.CORSRule{
					{
						AllowedOrigin: "http://www.example.com",
						AllowedMethod: "GET",
						MaxAgeSeconds: 3000,
						ExposeHeader:  "x-amz-server-side-encryption",
					},
				},
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("CORSConfiguration")))
		})
	})

	Describe("DeleteResult", func() {
		Specify("Send", func() {
			resp = s3.DeleteResult{
				Deleted: []s3.Deleted{
					{
						Key: "sample1.txt",
					},
				},
				Errors: []s3.DeletedError{
					{
						Key:     "sample2.txt",
						Code:    "AccessDenied",
						Message: "Access Denied",
					},
				},
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("DeleteResult")))
		})
	})

	Describe("InitiateMultipartUploadResult", func() {
		Specify("Send", func() {
			resp = s3.InitiateMultipartUploadResult{
				Bucket:   "example-bucket",
				Key:      "example-object",
				UploadId: "EXAMPLEJZ6e0YupT2h66iePQCc9IEbYbDUy4RTpMeoSMLPRp8Z5o1u8feSRonpvnWsKKG35tI2LB9VDPiCgTy.Gq2VxQLYjrue4Nq.NBdqI-",
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("InitiateMultipartUploadResult")))
		})
	})

	Describe("LifecycleConfiguration", func() {
		Specify("Send", func() {
			resp = s3.LifecycleConfiguration{
				Rules: []s3.Rule{
					{
						ID:     "Archive and then delete rule",
						Prefix: "projectdocs/",
						Status: "Enabled",
						Transition: &s3.Transition{
							Days:         365,
							StorageClass: "GLACIER",
						},
						Expiration: &s3.Expiration{
							Days: 3650,
						},
					},
				},
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("LifecycleConfiguration")))
		})
	})

	Describe("ListAllMyBucketsResult", func() {
		Specify("Send", func() {
			resp = s3.ListAllMyBucketsResult{
				Owner: s3.OwnerResult{
					ID:          "bcaf1ffd86f461ca5fb16fd081034f",
					DisplayName: "webfile",
				},
				Buckets: []s3.ListAllMyBucketsResultBucket{
					{
						Name:         "quotes",
						CreationDate: "2006-02-03T16:45:09.000Z",
					},
					{
						Name:         "samples",
						CreationDate: "2006-02-03T16:41:58.000Z",
					},
				},
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("ListAllMyBucketsResult")))
		})
	})

	Describe("ListBucketResult", func() {
		Specify("Send", func() {
			resp = s3.ListBucketResult{
				Name:           "bucket",
				MaxKeys:        1000,
				IsTruncated:    false,
				Delimiter:      "/",
				CommonPrefixes: s3.CommonPrefixes{"N/"},
				Contents: []s3.ContentResult{
					{
						Key:          "my-image.jpg",
						LastModified: "2009-10-12T17:50:30.000Z",
						ETag:         s3.NewETag("fba9dede5f27731c9771645a39863328"),
						Size:         434234,
						Owner: s3.OwnerResult{
							ID:          "75aa57f09aa0c8caeab4f8c24e99d10f8e7faeebf76c078efc7c6caea54ba06a",
							DisplayName: "mtd@amazon.com",
						},
						StorageClass: "STANDARD",
					},
					{
						Key:          "my-third-image.jpg",
						LastModified: "2009-10-12T17:50:30.000Z",
						ETag:         s3.NewETag("1b2cf535f27731c974343645a3985328"),
						Size:         64994,
						Owner: s3.OwnerResult{
							ID:          "75aa57f09aa0c8caeab4f8c24e99d10f8e7faeebf76c078efc7c6caea54ba06a",
							DisplayName: "mtd@amazon.com",
						},
						StorageClass: "STANDARD",
					},
				},
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("ListBucketResult")))
		})

		Describe("prefix", func() {
			Specify("Send", func() {
				resp = s3.ListBucketResult{
					Name:        "quotes",
					Prefix:      "N",
					Marker:      "Ned",
					MaxKeys:     40,
					IsTruncated: false,
					Contents: []s3.ContentResult{
						{
							Key:          "Nelson",
							LastModified: "2006-01-01T12:00:00.000Z",
							ETag:         s3.NewETag("828ef3fdfa96f00ad9f27c383fc9ac7f"),
							Size:         5,
							Owner: s3.OwnerResult{
								ID:          "bcaf161ca5fb16fd081034f",
								DisplayName: "webfile",
							},
							StorageClass: "STANDARD",
						},
						{
							Key:          "Neo",
							LastModified: "2006-01-01T12:00:00.000Z",
							ETag:         s3.NewETag("828ef3fdfa96f00ad9f27c383fc9ac7f"),
							Size:         4,
							Owner: s3.OwnerResult{
								ID:          "bcaf1ffd86a5fb16fd081034f",
								DisplayName: "webfile",
							},
							StorageClass: "STANDARD",
						},
					},
				}
				Expect(resp.Send(capture)).NotTo(HaveOccurred())
				Expect(getBody(capture)).To(Equal(fixture("ListBucketResult-prefix")))
			})
		})
	})

	Describe("LocationConstraint", func() {
		Specify("Send", func() {
			resp = s3.LocationConstraint{
				Location: "EU",
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("LocationConstraint")))
		})

		Describe("US-Classic", func() {
			Specify("Send", func() {
				resp = s3.LocationConstraint{}
				resp.Send(capture)
				Expect(getBody(capture)).To(Equal(fixture("LocationConstraint-US-Classic")))
			})
		})
	})

	Describe("NotificationConfiguration", func() {
		Specify("Send", func() {
			resp = s3.NotificationConfiguration{
				TopicConfiguration: s3.NewTopicConfiguration(
					"arn:aws:sns:us-east-1:123456789012:myTopic",
					"s3:ReducedRedundancyLostObject",
				),
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("NotificationConfiguration")))
		})
	})

	Describe("RequestPaymentConfiguration", func() {
		Specify("Send", func() {
			resp = s3.RequestPaymentConfiguration{
				Payer: "Requester",
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("RequestPaymentConfiguration")))
		})
	})

	Describe("Tagging", func() {
		Specify("Send", func() {
			resp = s3.Tagging{
				TagSet: []s3.Tag{
					{
						Key:   "Project",
						Value: "Project One",
					},
					{
						Key:   "User",
						Value: "jsmith",
					},
				},
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("Tagging")))
		})
	})

	Describe("VersioningConfiguration", func() {
		Specify("Send", func() {
			resp = s3.VersioningConfiguration{
				Status: "Enabled",
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("VersioningConfiguration")))
		})
	})
})
