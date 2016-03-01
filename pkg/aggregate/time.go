package aggregate

import (
	"time"
)

type Newest struct {
	Time  time.Time
	Valid bool
}

func NewNewest() *Newest {
	return &Newest{time.Time{}, false}
}

func (a *Newest) Check(t time.Time) {
	if a.Time.Before(t) {
		a.Time = t
		a.Valid = true
	}
}

func (a *Newest) Combine(t *Newest) {
	a.Check(t.Time)
}

type Oldest struct {
	Time  time.Time
	Valid bool
}

func NewOldest() *Oldest {
	return &Oldest{time.Now(), false}
}

func (a *Oldest) Check(t time.Time) {
	if a.Time.After(t) {
		a.Time = t
		a.Valid = true
	}
}

func (a *Oldest) Combine(t *Oldest) {
	a.Check(t.Time)
}
