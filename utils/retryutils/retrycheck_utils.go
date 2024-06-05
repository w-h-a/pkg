package retryutils

import "github.com/w-h-a/pkg/utils/errorutils"

func RetryOnError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}

	e := errorutils.ParseError(err.Error())
	if e == nil {
		return false, nil
	}

	switch e.Code {
	// retry on timeout or internal server error
	case 408, 500:
		return true, nil
	default:
		return false, nil
	}
}
