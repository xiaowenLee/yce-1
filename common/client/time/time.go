package time

import (
	"time"
)

type LocalTime struct {
}

func NewLocalTime() *LocalTime {
	return &LocalTime{}
}

func (local *LocalTime) String() string {
	t := time.Now()
	return t.Format(time.RFC3339)
}
