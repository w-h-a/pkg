package grpcclient

import (
	"fmt"
	"strings"
)

func ToGRPCMethod(method string) string {
	parts := strings.Split(method, ".")

	if len(parts) != 2 {
		return method
	}

	return fmt.Sprintf("/%s/%s", parts[0], parts[1])
}
