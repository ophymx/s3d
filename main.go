package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ophymx/s3d/internal/blob"
	"github.com/ophymx/s3d/internal/meta"
	"github.com/ophymx/s3d/internal/meta/dbenc"
	"github.com/ophymx/s3d/internal/s3"
	"github.com/ophymx/s3d/internal/server"
)

func main() {
	config := defaultConfig()
	var accessKey, secretKey, displayName, hosts string

	flag.StringVar(&config.DataRoot, "d", filepath.Join(os.TempDir(), "s3d"), "s3d data root")
	flag.IntVar(&config.Port, "p", 8080, "port")
	flag.StringVar(&accessKey, "a", "", "aws access key id")
	flag.StringVar(&secretKey, "s", "", "aws secret access key")
	flag.StringVar(&displayName, "n", "Example Account", "credential display name")
	flag.StringVar(&hosts, "h", "", "additional hosts to use when parsing bucket names")
	flag.Parse()

	if accessKey != "" && secretKey != "" {
		config.Credentials = append(config.Credentials, s3.Credential{
			AccessKeyID: accessKey,
			SecretKey:   secretKey,
			DisplayName: displayName,
		})
	}
	if hosts != "" {
		config.Hostnames = append(config.Hostnames, strings.Split(hosts, ",")...)
	}

	start(config)
}

func defaultConfig() config {
	c := config{
		S3: s3.Config{
			Region: "us-east-1",
			HostID: "====host/id====",
		},
		Hostnames: []string{
			"s3.amazonaws.com",
		},
	}
	return c
}

func start(config config) {
	bucketParser := s3.NewBucketParser(config.Hostnames)
	credentials := config.getCredentialsMap()
	store := blob.NewFsStore(filepath.Join(config.DataRoot, "buckets"))
	db, err := meta.NewDB(filepath.Join(config.DataRoot, "meta.db"), dbenc.MsgPack)
	if err != nil {
		log.Fatal(err)
	}

	handler := server.NewHandler(db, store, bucketParser, credentials, config.S3)
	log.Printf("Start server, port: %d, data: %s, hostID: %s", config.Port, config.DataRoot, config.S3.HostID)
	http.ListenAndServe(config.listenAddr(), handler)
}
