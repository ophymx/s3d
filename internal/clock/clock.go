package clock

import "time"

type Clock interface {
	// Always returns a time.Time in UTC
	Now() time.Time
}

type realClock struct {
}

var Real = realClock{}

func New() Clock {
	return realClock{}
}

func (realClock) Now() time.Time {
	return time.Now().UTC()
}
