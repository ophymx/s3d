package s3

import (
	"sort"
	"strings"
)

type BucketParser struct {
	suffixes []string
}

func NewBucketParser(hostnames []string) BucketParser {
	m := map[string]bool{}
	s := []string{}
	for _, host := range hostnames {
		if host == "" {
			continue
		}
		if host[0] != '.' {
			host = "." + host
		}
		if !m[host] {
			m[host] = true
			s = append(s, host)
		}
	}
	sort.Sort(hostSuffixes(s))
	return BucketParser{s}
}

func (p BucketParser) Parse(host string) (bucket string) {
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[0:idx]
	}

	for _, suffix := range p.suffixes {
		if len(host) > len(suffix) && strings.HasSuffix(host, suffix) {
			bucket = host[0 : len(host)-len(suffix)]
			break
		}
	}
	return
}

// Sort hostnames in descending length order
type hostSuffixes []string

func (s hostSuffixes) Len() int {
	return len(s)
}

func (s hostSuffixes) Less(i, j int) bool {
	if len(s[i]) == len(s[j]) {
		return s[i] >= s[j]
	}
	return len(s[i]) > len(s[j])
}

func (s hostSuffixes) Swap(i, j int) {
	a := s[j]
	s[j] = s[i]
	s[i] = a
}
