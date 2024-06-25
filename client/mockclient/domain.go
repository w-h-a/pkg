package mockclient

import "github.com/w-h-a/pkg/utils/errorutils"

type Response struct {
	Response interface{}       `json:"response,omitempty"`
	Err      *errorutils.Error `json:"err,omitempty"`
}
