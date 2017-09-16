package s3

type CompleteMultipartUpload struct {
	Parts []Part `xml:"Part"`
}
