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
	statusCode := codes.Unknown

	if e, ok := err.(*errorutils.Error); ok {
		switch e.Code {
		case http.StatusBadRequest:
			statusCode = codes.InvalidArgument
		case http.StatusUnauthorized:
			statusCode = codes.Unauthenticated
		case http.StatusForbidden:
			statusCode = codes.PermissionDenied
		case http.StatusNotFound:
			statusCode = codes.NotFound
		case http.StatusRequestTimeout:
			statusCode = codes.DeadlineExceeded
		case http.StatusInternalServerError:
			statusCode = codes.Internal
		}

		return statusCode
	}

	return statusCode
}
