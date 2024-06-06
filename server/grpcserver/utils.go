package grpcserver

import (
	"fmt"
	"strings"
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
