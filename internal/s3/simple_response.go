package s3

import "net/http"

// SimpleResponse responses with no body.
type SimpleResponse struct {
	Status int
	Header http.Header
	Body   []byte
}

func (resp SimpleResponse) Send(writer http.ResponseWriter) (err error) {
	writerHeader := writer.Header()

	for name, values := range resp.Header {
		for _, value := range values {
			writerHeader[name] = append(writerHeader[name], value)
		}
	}
	writer.WriteHeader(resp.Status)
	if len(resp.Body) > 0 {
		_, err = writer.Write(resp.Body)
	}
	return nil
}

func (resp SimpleResponse) HTTPStatus() int {
	return resp.Status
}

func NoContent() SimpleResponse {
	return SimpleResponse{
		Status: http.StatusNoContent,
	}
}
