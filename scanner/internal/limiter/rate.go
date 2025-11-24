package limiter

import (
	"context"
	"maps"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type AccoutnLimiter struct {
	limiter  map[uint]*rate.Limiter
	accounts map[uint]Account
	mu       sync.Mutex
}

func NewAccountLimiter() *AccoutnLimiter {
	return &AccoutnLimiter{
		limiter:  make(map[uint]*rate.Limiter),
		accounts: make(map[uint]Account),
		mu:       sync.Mutex{},
	}
}
func (a *AccoutnLimiter) SetAccount(id uint, cookie string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.limiter[id] = rate.NewLimiter(rate.Every(time.Second*2), 4)
	a.accounts[id] = Account{
		ID:     id,
		Cookie: cookie,
	}
}

func (a *AccoutnLimiter) Wait(ctx context.Context, id uint) error {
	a.mu.Lock()
	lim := a.limiter[id]
	a.mu.Unlock()
	if lim == nil {
		return nil
	}
	return lim.Wait(ctx)
}

func (a *AccoutnLimiter) Snapshot() map[uint]Account {
	a.mu.Lock()
	defer a.mu.Unlock()
	snapshot := make(map[uint]Account)
	maps.Copy(snapshot, a.accounts)
	return snapshot
}

func (a *AccoutnLimiter) GetAccount(ctx context.Context) (uint, string) {
	snapshot := a.Snapshot()
	for k, v := range snapshot {
		if v.Available() {
			a.Wait(ctx, k)
			return k, v.Cookie
		}
	}
	return 0, ""
}

func (a *AccoutnLimiter) Penalize(id uint) {
	a.mu.Lock()
	defer a.mu.Unlock()
	account := a.accounts[id]
	account.Penalize()
}

func (a *AccoutnLimiter) Reward(id uint) {
	a.mu.Lock()
	defer a.mu.Unlock()
	account := a.accounts[id]
	account.Reward()
}
