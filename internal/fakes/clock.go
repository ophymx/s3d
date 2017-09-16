package fakes

import (
	"time"
)

type Clock struct {
	Times []time.Time
}

func NewClock(times ...time.Time) *Clock {
	return &Clock{Times: times}
}

func (c *Clock) Add(now time.Time) {
	c.Times = append(c.Times, now)
}

func (c *Clock) Now() time.Time {
	switch len(c.Times) {
	case 0:
		return time.Now().UTC()
	case 1:
		return c.Times[0].UTC()
	default:
		now := c.Times[0]
		c.Times = c.Times[1:]
		return now.UTC()
	}
}
