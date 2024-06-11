package grpcserver

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/w-h-a/pkg/utils/errorutils"
	"google.golang.org/grpc/codes"
)

func ToControllerHandler(grpcFormattedMethod string) (controller string, handler string, err error) {
	parts := strings.Split(grpcFormattedMethod, "/")

	if len(parts) != 3 || len(parts[1]) == 0 || len(parts[2]) == 0 {
		return controller, handler, fmt.Errorf("malformed method name: %v", grpcFormattedMethod)
	}

	controller = parts[1]

	handler = parts[2]

	return controller, handler, nil
}

func ToErrorCode(err error) codes.Code {
	e, ok := err.(*errorutils.Error)
	if !ok {
		return codes.Unknown
	}

	switch e.Code {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	case http.StatusNotFound:
		return codes.NotFound
	case http.StatusRequestTimeout:
		return codes.DeadlineExceeded
	case http.StatusInternalServerError:
		return codes.Internal
	}

	return codes.Unknown
}
