package valueutil

import "time"

func TimePointer(t time.Time) *time.Time {
	return &t
}

func TimeValue(t *time.Time) time.Time {
	return *t
}
