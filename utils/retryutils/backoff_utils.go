package retryutils

import (
	"math"
	"time"
)

func ExponentialBackoff(attempts int) time.Duration {
	if attempts > 13 {
		return 2 * time.Minute
	}
	return time.Duration(math.Pow(float64(attempts), math.E)) * 100 * time.Millisecond
}
