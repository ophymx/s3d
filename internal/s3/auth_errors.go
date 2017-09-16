package s3

import "github.com/ophymx/s3d/internal/auth"

func AuthError(err error) ErrorResponse {
	if authErr, ok := err.(auth.AuthError); ok {
		switch authErr.Code {
		case auth.ErrCodeInvalidArgument:
			return InvalidArgument(authErr.Message, authErr.Params["name"], authErr.Params["value"])
		case auth.ErrCodeMissingSecurityHeader:
			return MissingSecurityHeader(authErr.Message)
		case auth.ErrCodeMissingSecurityElement:
			return MissingSecurityElement(authErr.Message)
		case auth.ErrAuthorizationQueryParametersError:
			return AuthorizationQueryParametersError(authErr.Message)
		case auth.ErrAccessDenied:
			return AccessDenied(authErr.Message)
		case auth.ErrInvalidAccessKeyID:
			return InvalidAccessKeyID(authErr.Params["key"])
		case auth.ErrSignatureDoesNotMatch:
			return SignatureDoesNotMatch(authErr.Message, authErr.Params["key"], authErr.Params["sts"], authErr.Params["signature"])
		}
	}
	return InternalError(err)
}
