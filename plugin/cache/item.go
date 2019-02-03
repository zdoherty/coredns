package cache

import (
	"time"

	"github.com/coredns/coredns/plugin/cache/freq"

	"github.com/miekg/dns"
)

type cacher interface {
	fromMsg(*dns.Msg, time.Time, time.Duration)
	toMsg(*dns.Msg, time.Time) *dns.Msg
	ttl(time.Time) int
	m() meta
}

type meta struct {
	origTTL uint32
	stored  time.Time
	*freq.Freq
}

func (m meta) ttl(now time.Time) int {
	ttl := int(m.origTTL) - int(now.UTC().Sub(m.stored).Seconds())
	return ttl
}

func copyExtra(extra []dns.RR) []dns.RR {
	ex := make([]dns.RR, len(extra))
	// Don't copy OPT records as these are hop-by-hop.
	j := 0
	for _, e := range extra {
		if e.Header().Rrtype == dns.TypeOPT {
			continue
		}
		ex[j] = e
		j++
	}
	return ex[:j]
}
