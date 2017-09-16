package dbenc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ophymx/s3d/internal/fixtures"
	"github.com/ophymx/s3d/internal/meta/dbenc"
)

var _ = Describe("MsgPack", func() {
	Describe("Bucket", func() {
		It("can encode and decode bucket meta data", func() {
			b, err := dbenc.MsgPack.EncodeBucket(fixtures.BucketMetadata())
			Expect(err).ToNot(HaveOccurred())
			Expect(b).ToNot(BeEmpty())

			Expect(dbenc.MsgPack.DecodeBucket(b)).To(Equal(fixtures.BucketMetadata()))
		})
	})
	Describe("Object", func() {
		It("can encode and decode object meta data", func() {
			b, err := dbenc.MsgPack.EncodeObject(fixtures.ObjectMetadata())
			Expect(err).ToNot(HaveOccurred())
			Expect(b).ToNot(BeEmpty())

			Expect(dbenc.MsgPack.DecodeObject(b)).To(Equal(fixtures.ObjectMetadata()))
		})
	})
})
