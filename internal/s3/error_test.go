package s3_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ophymx/s3d/internal/fakes"
	"github.com/ophymx/s3d/internal/s3"
)

var _ = Describe("Responses", func() {
	var (
		resp    s3.Response
		capture *fakes.ResponseWriter
	)

	BeforeEach(func() {
		capture = fakes.NewResponseWriter()
		capture.Header().Add(s3.AmzRequestID, "42E981960003A41E")
		capture.Header().Add(s3.AmzHostID, "====host/id====")
	})

	Describe("ErrorResponse", func() {
		Specify("Send", func() {
			resp = s3.ErrorResponse{
				Status:  500,
				Code:    "InvalidFooBar",
				Message: "Everything broke",
				Params: map[string]string{
					"BucketName": "my-bucket",
				},
			}
			Expect(resp.Send(capture)).NotTo(HaveOccurred())
			Expect(getBody(capture)).To(Equal(fixture("Error")))
		})
	})
})
