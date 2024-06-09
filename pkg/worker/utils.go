package worker

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func ParseDuration(value string, defaultUnit string) (time.Duration, error) {
	// Provided value is only a numeric string
	if v, err := strconv.Atoi(value); err == nil {
		value = fmt.Sprintf("%d"+defaultUnit, v)
	}

	d, err := time.ParseDuration(value)
	if err != nil {
		return time.Duration(0), fmt.Errorf("invalid time duration '%s' : %s", value, err.Error())
	}
	return d, nil
}

func timeoutToSeconds(t string) (uint32, error) {
	if t == "" {
		return 0, nil
	}
	d, err := ParseDuration(t, "s")
	if err != nil {
		return 0, err
	}
	s := d.Truncate(time.Second).Seconds()
	if s < 0 {
		return 0, fmt.Errorf("timeout cannot be negative")
	}
	i := int64(s)
	if i > math.MaxUint32 {
		return 0, fmt.Errorf("timeout cannot be over %d", math.MaxUint32)
	}
	return uint32(i), nil
}
