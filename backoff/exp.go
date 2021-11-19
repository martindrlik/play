package backoff

import (
	"math"
	"time"
)

func Exp(retry int) time.Duration {
	return time.Duration(math.Exp(float64(retry))) * time.Second
}
