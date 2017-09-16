package s3

import (
	"fmt"
	"net/http"
	"strings"
)

func AccessDenied(message string) ErrorResponse {
	return NewErrorResponse("AccessDenied", http.StatusForbidden, message)
}

func AccountProblem(message string) ErrorResponse {
	return NewErrorResponse("AccountProblem", http.StatusForbidden, message)
}

func AmbiguousGrantByEmailAddress(message string) ErrorResponse {
	return NewErrorResponse("AmbiguousGrantByEmailAddress", http.StatusBadRequest, message)
}

func AuthorizationQueryParametersError(message string) ErrorResponse {
	return NewErrorResponse("AuthorizationQueryParametersError", http.StatusBadRequest, message)
}

func BadDigest(message string) ErrorResponse {
	return NewErrorResponse("BadDigest", http.StatusBadRequest, message)
}

func BucketAlreadyExists(message string) ErrorResponse {
	return NewErrorResponse("BucketAlreadyExists", http.StatusConflict, message)
}

func BucketAlreadyOwnedByYou(message string) ErrorResponse {
	return NewErrorResponse("BucketAlreadyOwnedByYou", http.StatusConflict, message)
}

func BucketNotEmpty(message string) ErrorResponse {
	return NewErrorResponse("BucketNotEmpty", http.StatusConflict, message)
}

func CredentialsNotSupported(message string) ErrorResponse {
	return NewErrorResponse("CredentialsNotSupported", http.StatusBadRequest, message)
}

func CrossLocationLoggingProhibited(message string) ErrorResponse {
	return NewErrorResponse("CrossLocationLoggingProhibited", http.StatusForbidden, message)
}

func EntityTooSmall(message string) ErrorResponse {
	return NewErrorResponse("EntityTooSmall", http.StatusBadRequest, message)
}

func EntityTooLarge(message string) ErrorResponse {
	return NewErrorResponse("EntityTooLarge", http.StatusBadRequest, message)
}

func ExpiredToken(message string) ErrorResponse {
	return NewErrorResponse("ExpiredToken", http.StatusBadRequest, message)
}

func IllegalVersioningConfigurationException(message string) ErrorResponse {
	return NewErrorResponse("IllegalVersioningConfigurationException", http.StatusBadRequest, message)
}

func IncompleteBody(message string) ErrorResponse {
	return NewErrorResponse("IncompleteBody", http.StatusBadRequest, message)
}

func IncorrectNumberOfFilesInPostRequest(message string) ErrorResponse {
	return NewErrorResponse("IncorrectNumberOfFilesInPostRequest", http.StatusBadRequest, message)
}

func InlineDataTooLarge(message string) ErrorResponse {
	return NewErrorResponse("InlineDataTooLarge", http.StatusBadRequest, message)
}

func InternalError(err error) ErrorResponse {
	return NewErrorResponse("InternalError", http.StatusInternalServerError, err.Error())
}

func InternalErrorf(format string, args ...interface{}) ErrorResponse {
	return NewErrorResponse("InternalError", http.StatusInternalServerError, fmt.Sprintf(format, args...))
}

func InvalidAccessKeyID(key string) ErrorResponse {
	return ErrorResponse{
		Code:    "InvalidAccessKeyId",
		Status:  http.StatusForbidden,
		Message: "The AWS Access Key Id you provided does not exist in our records.",
		Params:  map[string]string{"AWSAccessKeyId": key},
	}
}

func InvalidArgument(message, name, value string) ErrorResponse {
	return ErrorResponse{
		Code:    "InvalidArgument",
		Status:  http.StatusBadRequest,
		Message: message,
		Params: map[string]string{
			"ArgumentName":  name,
			"ArgumentValue": value,
		},
	}
}

func InvalidBucketName(message string) ErrorResponse {
	return NewErrorResponse("InvalidBucketName", http.StatusBadRequest, message)
}

func InvalidBucketState(message string) ErrorResponse {
	return NewErrorResponse("InvalidBucketState", http.StatusConflict, message)
}

func InvalidDigest(digest string) ErrorResponse {
	return ErrorResponse{
		Code:    "InvalidDigest",
		Status:  http.StatusBadRequest,
		Message: "The Content-MD5 you specified was invalid.",
		Params:  map[string]string{"Content-MD5": digest},
	}
}

func InvalidEncryptionAlgorithmError(message string) ErrorResponse {
	return NewErrorResponse("InvalidEncryptionAlgorithmError", http.StatusBadRequest, message)
}

func InvalidLocationConstraint(message string) ErrorResponse {
	return NewErrorResponse("InvalidLocationConstraint", http.StatusBadRequest, message)
}

func InvalidObjectState(message string) ErrorResponse {
	return NewErrorResponse("InvalidObjectState", http.StatusForbidden, message)
}

func InvalidPart(message string) ErrorResponse {
	return NewErrorResponse("InvalidPart", http.StatusBadRequest, message)
}

func InvalidPartOrder(message string) ErrorResponse {
	return NewErrorResponse("InvalidPartOrder", http.StatusBadRequest, message)
}

func InvalidPayer(message string) ErrorResponse {
	return NewErrorResponse("InvalidPayer", http.StatusForbidden, message)
}

func InvalidPolicyDocument(message string) ErrorResponse {
	return NewErrorResponse("InvalidPolicyDocument", http.StatusBadRequest, message)
}

func InvalidRange(message string) ErrorResponse {
	return NewErrorResponse("InvalidRange", http.StatusRequestedRangeNotSatisfiable, message)
}

func InvalidRequest(message string) ErrorResponse {
	return NewErrorResponse("InvalidRequest", http.StatusBadRequest, message)
}

func InvalidSecurity(message string) ErrorResponse {
	return NewErrorResponse("InvalidSecurity", http.StatusForbidden, message)
}

func InvalidSOAPRequest(message string) ErrorResponse {
	return NewErrorResponse("InvalidSOAPRequest", http.StatusBadRequest, message)
}

func InvalidStorageClass(message string) ErrorResponse {
	return NewErrorResponse("InvalidStorageClass", http.StatusBadRequest, message)
}

func InvalidTargetBucketForLogging(message string) ErrorResponse {
	return NewErrorResponse("InvalidTargetBucketForLogging", http.StatusBadRequest, message)
}

func InvalidToken(message string) ErrorResponse {
	return NewErrorResponse("InvalidToken", http.StatusBadRequest, message)
}

func InvalidURI(message string) ErrorResponse {
	return NewErrorResponse("InvalidURI", http.StatusBadRequest, message)
}

func KeyTooLong(message string) ErrorResponse {
	return NewErrorResponse("KeyTooLong", http.StatusBadRequest, message)
}

func MalformedACLError(message string) ErrorResponse {
	return NewErrorResponse("MalformedACLError", http.StatusBadRequest, message)
}

func MalformedPOSTRequest(message string) ErrorResponse {
	return NewErrorResponse("MalformedPOSTRequest", http.StatusBadRequest, message)
}

func MalformedXML(message string) ErrorResponse {
	return NewErrorResponse("MalformedXML", http.StatusBadRequest, message)
}

func MaxMessageLengthExceeded(message string) ErrorResponse {
	return NewErrorResponse("MaxMessageLengthExceeded", http.StatusBadRequest, message)
}

func MaxPostPreDataLengthExceededError(message string) ErrorResponse {
	return NewErrorResponse("MaxPostPreDataLengthExceededError", http.StatusBadRequest, message)
}

func MetadataTooLarge(message string) ErrorResponse {
	return NewErrorResponse("MetadataTooLarge", http.StatusBadRequest, message)
}

func MethodNotAllowed(message string) ErrorResponse {
	return NewErrorResponse("MethodNotAllowed", http.StatusMethodNotAllowed, message)
}

func MissingContentLength(message string) ErrorResponse {
	return NewErrorResponse("MissingContentLength", http.StatusLengthRequired, message)
}

func MissingRequestBodyError(message string) ErrorResponse {
	return NewErrorResponse("MissingRequestBodyError", http.StatusBadRequest, message)
}

func MissingSecurityElement(message string) ErrorResponse {
	return NewErrorResponse("MissingSecurityElement", http.StatusBadRequest, message)
}

func MissingSecurityHeader(message string) ErrorResponse {
	return NewErrorResponse("MissingSecurityHeader", http.StatusBadRequest, message)
}

func NoLoggingStatusForKey(message string) ErrorResponse {
	return NewErrorResponse("NoLoggingStatusForKey", http.StatusBadRequest, message)
}

func NoSuchBucket(bucket string) ErrorResponse {
	return ErrorResponse{
		Code:    "NoSuchBucket",
		Status:  http.StatusNotFound,
		Message: "The specified bucket does not exist",
		Params:  map[string]string{"BucketName": bucket},
	}
}

func NoSuchKey(resource string) ErrorResponse {
	return ErrorResponse{
		Code:    "NoSuchKey",
		Status:  http.StatusNotFound,
		Message: "No such key",
		Params:  map[string]string{"Resource": resource},
	}
}

func NoSuchLifecycleConfiguration(message string) ErrorResponse {
	return NewErrorResponse("NoSuchLifecycleConfiguration", http.StatusNotFound, message)
}

func NoSuchUpload(message string) ErrorResponse {
	return NewErrorResponse("NoSuchUpload", http.StatusNotFound, message)
}

func NoSuchVersion(message string) ErrorResponse {
	return NewErrorResponse("NoSuchVersion", http.StatusNotFound, message)
}

func NotImplemented(message string) ErrorResponse {
	return NewErrorResponse("NotImplemented", http.StatusNotImplemented, message)
}

func NotSignedUp(message string) ErrorResponse {
	return NewErrorResponse("NotSignedUp", http.StatusForbidden, message)
}

func NotSuchBucketPolicy(message string) ErrorResponse {
	return NewErrorResponse("NotSuchBucketPolicy", http.StatusNotFound, message)
}

func OperationAborted(message string) ErrorResponse {
	return NewErrorResponse("OperationAborted", http.StatusConflict, message)
}

func PermanentRedirect(message string) ErrorResponse {
	return NewErrorResponse("PermanentRedirect", http.StatusMovedPermanently, message)
}

func PreconditionFailed(message string) ErrorResponse {
	return NewErrorResponse("PreconditionFailed", http.StatusPreconditionFailed, message)
}

func Redirect(message string) ErrorResponse {
	return NewErrorResponse("Redirect", http.StatusTemporaryRedirect, message)
}

func RestoreAlreadyInProgress(message string) ErrorResponse {
	return NewErrorResponse("RestoreAlreadyInProgress", http.StatusConflict, message)
}

func RequestIsNotMultiPartContent(message string) ErrorResponse {
	return NewErrorResponse("RequestIsNotMultiPartContent", http.StatusBadRequest, message)
}

func RequestTimeout(message string) ErrorResponse {
	return NewErrorResponse("RequestTimeout", http.StatusBadRequest, message)
}

func RequestTimeTooSkewed(message string) ErrorResponse {
	return NewErrorResponse("RequestTimeTooSkewed", http.StatusForbidden, message)
}

func RequestTorrentOfBucketError(message string) ErrorResponse {
	return NewErrorResponse("RequestTorrentOfBucketError", http.StatusBadRequest, message)
}

func SignatureDoesNotMatch(message, accessKey, sts, signature string) ErrorResponse {
	return ErrorResponse{
		Code:    "SignatureDoesNotMatch",
		Status:  http.StatusForbidden,
		Message: message,
		Params: map[string]string{
			"AWSAccessKeyId":    accessKey,
			"StringToSign":      sts,
			"SignatureProvided": signature,
			"StringToSignBytes": strings.Trim(fmt.Sprintf("%v", []byte(sts)), "[]"),
		},
	}
}

func ServiceUnavailable(message string) ErrorResponse {
	return NewErrorResponse("ServiceUnavailable", http.StatusServiceUnavailable, message)
}

func SlowDown(message string) ErrorResponse {
	return NewErrorResponse("SlowDown", http.StatusServiceUnavailable, message)
}

func TemporaryRedirect(message string) ErrorResponse {
	return NewErrorResponse("TemporaryRedirect", http.StatusTemporaryRedirect, message)
}

func TokenRefreshRequired(message string) ErrorResponse {
	return NewErrorResponse("TokenRefreshRequired", http.StatusBadRequest, message)
}

func TooManyBuckets(message string) ErrorResponse {
	return NewErrorResponse("TooManyBuckets", http.StatusBadRequest, message)
}

func UnexpectedContent(message string) ErrorResponse {
	return NewErrorResponse("UnexpectedContent", http.StatusBadRequest, message)
}

func UnresolvableGrantByEmailAddress(message string) ErrorResponse {
	return NewErrorResponse("UnresolvableGrantByEmailAddress", http.StatusBadRequest, message)
}

func UserKeyMustBeSpecified(message string) ErrorResponse {
	return NewErrorResponse("UserKeyMustBeSpecified", http.StatusBadRequest, message)
}
