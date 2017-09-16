package main

import (
	"strconv"

	"github.com/ophymx/s3d/internal/s3"
)

// S3dConfig configuration for s3d.
type config struct {
	DataRoot    string
	Port        int
	Hostnames   []string
	S3          s3.Config
	Credentials []s3.Credential
}

func (c config) listenAddr() string {
	return ":" + strconv.Itoa(c.Port)
}

func (c config) getCredentialsMap() map[string]s3.Credential {
	lookup := make(map[string]s3.Credential)
	for _, cred := range c.Credentials {
		lookup[cred.AccessKeyID] = cred
	}
	return lookup
}
