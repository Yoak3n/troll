package limiter

import (
	"time"
)

type Account struct {
	ID            uint
	Cookie        string
	FailedCount   uint
	NextAvaliable time.Time
}

func (a *Account) Available() bool {
	return a.NextAvaliable.Before(time.Now())
}

func (a *Account) Penalize() {
	a.FailedCount++
	a.NextAvaliable = time.Now().Add(time.Minute * 5)
}

func (a *Account) Reward() {
	a.FailedCount = 0
	a.NextAvaliable = time.Time{}
}
