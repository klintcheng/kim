package wire

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsGrpcError check err type with a grpc status
func IsGrpcError(err error, code codes.Code) bool {
	if err == nil {
		return false
	}
	if st, ok := status.FromError(err); ok {
		return st.Code() == code
	}
	return false
}
