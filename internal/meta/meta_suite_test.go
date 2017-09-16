package meta_test

import (
	"io/ioutil"
	"os"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/ophymx/s3d/internal/meta"
	"github.com/ophymx/s3d/internal/meta/dbenc"

	"testing"
)

func TestMeta(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Meta Suite")
}

var _ = ginkgo.Describe("nothing", func() {
	ginkgo.It("can help me jump to ginkgo definitions", func() {
		gomega.Expect("to also be able to jump to gomega definitions").
			To(gomega.ContainSubstring("gomega"))
	})
})

func setup() (dbpath string, db meta.DB) {
	f, err := ioutil.TempFile("", "s3-db-test-")
	if err != nil {
		panic(err)
	}
	dbpath = f.Name()
	if err = f.Close(); err != nil {
		panic(err)
	}
	if db, err = meta.NewDB(dbpath, dbenc.MsgPack); err != nil {
		panic(err)
	}
	return
}

func tearDown(dbpath string, db meta.DB) {
	if err := db.Close(); err != nil {
		panic(err)
	}
	if err := os.RemoveAll(dbpath); err != nil {
		panic(err)
	}
}

func must(err error) {
	if err != nil {
		panic(must)
	}
}
