package s3

import (
	"encoding/xml"
	"net/http"
)

func sendXMLHeader(writer http.ResponseWriter, results interface{}) (err error) {
	writer.Header().Add("Content-Type", "application/xml")
	xmlBytes, err := xml.Marshal(results)
	if err != nil {
		return
	}
	writer.Write([]byte(xml.Header))
	writer.Write(xmlBytes)
	writer.Write([]byte{'\n'})
	return nil
}

func sendXML(writer http.ResponseWriter, results interface{}) (err error) {
	writer.Header().Add("Content-Type", "application/xml")
	xmlBytes, err := xml.Marshal(results)
	if err != nil {
		return
	}
	writer.Write(xmlBytes)
	writer.Write([]byte{'\n'})
	return nil
}

func marshalString(e *xml.Encoder, name, value string) (err error) {
	if len(value) == 0 {
		return nil
	}

	start := xml.StartElement{Name: xml.Name{Local: name}}
	if err = e.EncodeToken(start); err != nil {
		return
	}
	if err = e.EncodeToken(xml.CharData(value)); err != nil {
		return
	}
	return e.EncodeToken(start.End())
}
