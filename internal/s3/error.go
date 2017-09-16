package s3

import (
	"encoding/xml"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Status    int
	Code      string
	Message   string
	RequestID string
	HostID    string

	Params map[string]string
}

func NewErrorResponse(code string, status int, message string) ErrorResponse {
	return ErrorResponse{
		Code:    code,
		Status:  status,
		Message: message,
	}
}

func (err ErrorResponse) Error() string {
	return err.Code + ": " + err.Message
}

func (err ErrorResponse) HTTPStatus() int {
	return err.Status
}

func (err ErrorResponse) Send(writer http.ResponseWriter) error {
	header := writer.Header()
	header.Add(HdrContentType, "application/xml")
	err.RequestID = header.Get(AmzRequestID)
	err.HostID = header.Get(AmzHostID)

	log.Printf("[%s] Error %s: %s", err.RequestID, err.Code, err.Message)

	writer.WriteHeader(err.Status)
	return sendXMLHeader(writer, err)
}

func (resp ErrorResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	start.Name.Local = "Error"
	if err = e.EncodeToken(start); err != nil {
		return
	}
	if err = marshalString(e, "Code", resp.Code); err != nil {
		return
	}
	if err = marshalString(e, "Message", resp.Message); err != nil {
		return
	}
	for name, value := range resp.Params {
		if err = marshalString(e, name, value); err != nil {
			return
		}
	}
	if err = marshalString(e, "RequestId", resp.RequestID); err != nil {
		return
	}
	if err = marshalString(e, "HostId", resp.HostID); err != nil {
		return
	}
	return e.EncodeToken(start.End())
}
