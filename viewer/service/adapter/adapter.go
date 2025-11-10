package adapter

import "github.com/Yoak3n/troll/viewer/config"

type Adapter struct {
	Account map[string]float32
	Proxy   map[string]bool
}

func NewAdapter() {
	conf := config.GetConfiguration()
	a := &Adapter{
		Account: make(map[string]float32),
		Proxy:   make(map[string]bool),
	}
	for _, cookie := range conf.Cookies {
		a.Account[cookie.Data] = 0
	}
	for _, proxy := range conf.Proxies {
		a.Proxy[proxy.Data] = true
	}
}

func (a *Adapter) FindAvaliableAccountCookie() string {
	leastPayloadAccount := ""
	for cookie, payload := range a.Account {
		if payload > a.Account[leastPayloadAccount] {
			leastPayloadAccount = cookie
		}
	}
	a.Account[leastPayloadAccount] += 0.2
	return leastPayloadAccount
}

func (a *Adapter) ReleaseAccount(cookie string) float32 {
	a.Account[cookie] -= 0.2
	return a.Account[cookie]
}

func (a *Adapter) FindAvaliableProxy() string {
	for p, ok := range a.Proxy {
		if ok {
			return p
		}
	}
	return ""
}
