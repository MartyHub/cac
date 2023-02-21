package internal

import "time"

type clock interface {
	now() time.Time
}

type utcClock struct{}

func (c utcClock) now() time.Time {
	return time.Now().UTC()
}

type fixedClock struct {
	t time.Time
}

func (c fixedClock) now() time.Time {
	return c.t
}
