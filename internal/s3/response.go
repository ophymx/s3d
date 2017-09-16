package s3

import "net/http"

type Response interface {
	HTTPStatus() int
	Send(http.ResponseWriter) error
}
