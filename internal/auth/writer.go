package auth

import "io"

type Writer interface {
	Write(p []byte) (int, error)
	WriteByte(c byte) error
	WriteString(s string) (int, error)
}

func WrapWriter(writer io.Writer) Writer {
	return wWriter{writer}
}

type wWriter struct {
	io.Writer
}

func (w wWriter) WriteByte(c byte) error {
	_, err := w.Write([]byte{c})
	return err
}

func (w wWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}
