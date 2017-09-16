package auth

import "fmt"

// ErrCode auth error codes to help match up to S3 errors codes.
type ErrCode int

const (
	ErrCodeInvalidArgument               ErrCode = iota
	ErrCodeMissingSecurityHeader                 = iota
	ErrCodeMissingSecurityElement                = iota
	ErrAuthorizationQueryParametersError         = iota
	ErrAccessDenied                              = iota
	ErrInvalidAccessKeyID                        = iota
	ErrSignatureDoesNotMatch                     = iota
)

const (
	errInvalidFormat     = "AWS authorization header is invalid.  Expected AwsAccessKeyId:signature"
	errHeaderSpacing     = "Authorization header is invalid -- one and only one ' ' (space) required"
	errUnsupportedType   = "Unsupported Authorization Type"
	errAccessKeyNotFound = "The AWS Access Key Id you provided does not exist in our records."
	errSignature         = "The request signature we calculated does not match the signature you provided. Check your key and signing method."
)

type AuthError struct {
	Code    ErrCode
	Message string
	Params  map[string]string
}

func (err AuthError) Error() string {
	return err.Message
}

func MissingSecurityHeader(msg string) AuthError {
	return AuthError{
		Code:    ErrCodeMissingSecurityHeader,
		Message: msg,
	}
}

func InvalidArgument(msg, name, value string) AuthError {
	return AuthError{
		Code:    ErrCodeInvalidArgument,
		Message: msg,
		Params: map[string]string{
			"name":  name,
			"value": value,
		},
	}
}

func MissingSecurityElement(msg string) AuthError {
	return AuthError{
		Code:    ErrCodeMissingSecurityElement,
		Message: msg,
	}
}

func AccessDenied(msg string) AuthError {
	return AuthError{
		Code:    ErrAccessDenied,
		Message: msg,
	}
}

func AuthorizationQueryParametersError(msg string) AuthError {
	return AuthError{
		Code:    ErrAuthorizationQueryParametersError,
		Message: msg,
	}
}

func AuthorizationQueryParametersErrorf(msg string, args ...interface{}) AuthError {
	return AuthError{
		Code:    ErrAuthorizationQueryParametersError,
		Message: fmt.Sprintf(msg, args...),
	}
}

func InvalidAccessKeyId(msg string) AuthError {
	return AuthError{
		Code:    ErrInvalidAccessKeyID,
		Message: msg,
	}
}

func SignatureDoesNotMatch(accessKeyID, sts, signature string) AuthError {
	return AuthError{
		Code:    ErrSignatureDoesNotMatch,
		Message: errSignature,
		Params: map[string]string{
			"key":       accessKeyID,
			"sts":       sts,
			"signature": signature,
		},
	}
}
