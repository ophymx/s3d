package fakes

import (
	"bytes"
	"net/http"
)

type ResponseWriter struct {
	header http.Header
	Status int
	*bytes.Buffer
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		header: make(http.Header),
		Buffer: &bytes.Buffer{},
	}
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) WriteHeader(status int) {
	if w.Status != 0 {
		panic("already wrote header")
	}
	w.Status = status
}
