package utils

import "strconv"

func StringConstraint(value string, min, max, defaultValue int64) int64 {
	if res, err := strconv.Atoi(value); err == nil {
		defaultValue = int64(res)
	}
	if defaultValue < min {
		defaultValue = min
	}
	if defaultValue > max {
		defaultValue = max
	}
	return defaultValue
}

func Max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
