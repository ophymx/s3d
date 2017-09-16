package s3

import "net/http"

const (
	NSS3  = "http://s3.amazonaws.com/doc/2006-03-01/"
	NSXsi = "http://www.w3.org/2001/XMLSchema-instance"
)

type StatucOKResponse struct {
}

func (x StatucOKResponse) HTTPStatus() int {
	return http.StatusOK
}

type XMLNS struct {
	NS string `xml:"xmlns,attr"`
	StatucOKResponse
}

// ACL
type Permission string

const (
	PermRead     Permission = "READ"
	PermWrite               = "WRITE"
	PermReadACP             = "READ_ACP"
	PermWriteACP            = "WRITE_ACP"
	PermFull                = "FULL_CONTROL"
)

type GrantType string

const (
	GrantUser  GrantType = "CanonicalUser"
	GrantGroup           = "Group"
	GrantEmail           = "AmazonCustomerByEmail"
)

const (
	LogDeliveryGroup        = "http://acs.amazonaws.com/groups/s3/LogDelivery"
	AllUsersGroup           = "http://acs.amazonaws.com/groups/global/AllUsers"
	AuthenticatedUsersGroup = "http://acs.amazonaws.com/groups/global/AuthenticatedUsers"
)

type Grantee struct {
	NSXsi        string    `xml:"xmlns:xsi,attr"`
	Type         GrantType `xml:"xsi:type,attr"`
	URI          string    `xml:",omitempty"`
	ID           string    `xml:",omitempty"`
	DisplayName  string    `xml:",omitempty"`
	EmailAddress string    `xml:",omitempty"`
}

type Grant struct {
	Grantee    Grantee
	Permission Permission
}

func NewGroupGrant(uri string, permission Permission) (grant Grant) {
	return Grant{
		Grantee: Grantee{
			NSXsi: NSXsi,
			Type:  GrantGroup,
			URI:   uri,
		},
		Permission: permission,
	}
}

func NewUserGrant(id, displayName string, permission Permission) (grant Grant) {
	return Grant{
		Grantee: Grantee{
			NSXsi:       NSXsi,
			Type:        GrantUser,
			ID:          id,
			DisplayName: displayName,
		},
		Permission: permission,
	}
}

func NewEmailGrant(emailAddress string, permission Permission) (grant Grant) {
	return Grant{
		Grantee: Grantee{
			NSXsi:        NSXsi,
			Type:         GrantEmail,
			EmailAddress: emailAddress,
		},
		Permission: permission,
	}
}

type AccessControlPolicy struct {
	Owner             OwnerResult
	AccessControlList []Grant `xml:"AccessControlList>Grant"`
	StatucOKResponse
}

func (results AccessControlPolicy) Send(writer http.ResponseWriter) error {
	return sendXML(writer, results)
}

// LifeCycle
type Expiration struct {
	Days int    `xml:",omitempty"`
	Date string `xml:",omitempty"`
}

type Transition struct {
	Days         int    `xml:",omitempty"`
	Date         string `xml:",omitempty"`
	StorageClass string
}

type NoncurrentVersionExpiration struct {
	NoncurrentDays int
}

type NoncurrentVersionTransition struct {
	NoncurrentDays int
	StorageClass   string
}

type Rule struct {
	ID                          string
	Prefix                      string
	Status                      string
	Transition                  *Transition                  `xml:",omitempty"`
	Expiration                  *Expiration                  `xml:",omitempty"`
	NoncurrentVersionTransition *NoncurrentVersionTransition `xml:",omitempty"`
	NoncurrentVersionExpiration *NoncurrentVersionExpiration `xml:",omitempty"`
}

type LifecycleConfiguration struct {
	XMLNS
	Rules []Rule `xml:"Rule"`
}

func (results LifecycleConfiguration) Send(writer http.ResponseWriter) error {
	results.NS = NSS3
	return sendXMLHeader(writer, results)
}

// Location
type LocationConstraint struct {
	XMLNS
	Location string `xml:",chardata"`
}

func (results LocationConstraint) Send(writer http.ResponseWriter) error {
	results.NS = NSS3
	return sendXMLHeader(writer, results)
}

// Bucket Logging
type LoggingEnabled struct {
	TargetBucket string
	TargetPrefix string
	TargetGrants []Grant `xml:"TargetGrants>Grant"`
}

type BucketLoggingStatus struct {
	XMLNS
	LoggingEnabled *LoggingEnabled `xml:",omitempty"`
}

func (results BucketLoggingStatus) Send(writer http.ResponseWriter) error {
	results.NS = NSS3
	return sendXMLHeader(writer, results)
}

// Multipart upload
type InitiateMultipartUploadResult struct {
	XMLNS
	Bucket   string
	Key      string
	UploadId string
}

func (results InitiateMultipartUploadResult) Send(writer http.ResponseWriter) error {
	results.NS = NSS3
	return sendXMLHeader(writer, results)
}

type CopyPartResult struct {
	LastModified string
	ETag         ETag
	StatucOKResponse
}

func (results CopyPartResult) Send(writer http.ResponseWriter) error {
	return sendXML(writer, results)
}

type CompleteMultipartUploadResult struct {
	XMLNS
	Location string
	Bucket   string
	Key      string
	ETag     ETag
}

func (results CompleteMultipartUploadResult) Send(writer http.ResponseWriter) error {
	results.NS = NSS3
	return sendXMLHeader(writer, results)
}

type Part struct {
	PartNumber   int
	LastModified string `xml:",omitempty"`
	ETag         ETag
	Size         int64 `xml:",omitempty"`
}

type ListPartsResult struct {
	XMLNS
	Bucket               string
	Key                  string
	UploadId             string
	Initiator            OwnerResult
	Owner                OwnerResult
	StorageClass         string
	PartNumberMarker     int
	NextPartNumberMarker int
	MaxParts             int
	IsTruncated          bool
	Parts                []Part `xml:"Part"`
}

type CORSRule struct {
	AllowedOrigin string
	AllowedMethod string
	MaxAgeSeconds int
	ExposeHeader  string
}

type CORSConfiguration struct {
	CORSRules []CORSRule `xml:"CORSRule"`
	StatucOKResponse
}

func (results CORSConfiguration) Send(writer http.ResponseWriter) error {
	return sendXML(writer, results)
}

type Deleted struct {
	Key                   string
	VersionId             string `xml:",omitempty"`
	DeleteMarker          string `xml:",omitempty"`
	DeleteMarkerVersionId string `xml:",omitempty"`
}

type DeletedError struct {
	Key       string
	Code      string
	Message   string
	VersionId string `xml:",omitempty"`
}

type DeleteResult struct {
	XMLNS
	Deleted []Deleted
	Errors  []DeletedError `xml:"Error"`
}

func (results DeleteResult) Send(writer http.ResponseWriter) error {
	results.NS = NSS3
	return sendXMLHeader(writer, results)
}

type TopicConfiguration struct {
	Topic string
	Event string
}

func NewTopicConfiguration(topic, event string) *TopicConfiguration {
	return &TopicConfiguration{
		Topic: topic,
		Event: event,
	}
}

type NotificationConfiguration struct {
	TopicConfiguration *TopicConfiguration `xml:",omitempty"`
	StatucOKResponse
}

func (results NotificationConfiguration) Send(writer http.ResponseWriter) error {
	return sendXML(writer, results)
}

type RequestPaymentConfiguration struct {
	XMLNS
	Payer string
}

func (results RequestPaymentConfiguration) Send(writer http.ResponseWriter) error {
	results.NS = NSS3
	return sendXMLHeader(writer, results)
}

type Tag struct {
	Key   string
	Value string
}

type Tagging struct {
	TagSet []Tag `xml:"TagSet>Tag"`
	StatucOKResponse
}

func (results Tagging) Send(writer http.ResponseWriter) error {
	return sendXML(writer, results)
}

type VersioningConfiguration struct {
	XMLNS
	Status    string `xml:",omitempty"`
	MfaDelete string `xml:",omitempty"`
}

func (results VersioningConfiguration) Send(writer http.ResponseWriter) error {
	results.NS = NSS3
	return sendXMLHeader(writer, results)
}

type IndexDocument struct {
	Suffix string
}

type ErrorDocument struct {
}

type RedirectAllRequestsTo struct {
}

type RoutingRules struct {
}

type WebsiteConfiguration struct {
	XMLNS
	RedirectAllRequestsTo RedirectAllRequestsTo
	IndexDocument         IndexDocument
	ErrorDocument         ErrorDocument
	RoutingRules          RoutingRules
}
