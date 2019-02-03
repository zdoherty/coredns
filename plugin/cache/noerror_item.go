package cache

import (
	"time"

	"github.com/coredns/coredns/plugin/cache/freq"

	"github.com/miekg/dns"
)

type noErrorItem struct {
	AuthenticatedData bool
	Answer            []dns.RR
	Extra             []dns.RR

	meta
}

func (i *noErrorItem) ttl(now time.Time) int { return i.meta.ttl(now) }
func (i *noErrorItem) m() meta               { return i.meta }

func (i *noErrorItem) fromMsg(m *dns.Msg, now time.Time, d time.Duration) {
	i.AuthenticatedData = m.AuthenticatedData
	i.Answer = m.Answer
	i.Extra = copyExtra(m.Extra)

	i.origTTL = uint32(d.Seconds())
	i.stored = now.UTC()
	i.Freq = new(freq.Freq)
}

func (i *noErrorItem) toMsg(r *dns.Msg, now time.Time) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(r)
	m.AuthenticatedData = i.AuthenticatedData
	m.RecursionAvailable = true
	m.Rcode = dns.RcodeSuccess

	m.Answer = make([]dns.RR, len(i.Answer))
	m.Extra = make([]dns.RR, len(i.Extra))

	ttl := uint32(i.meta.ttl(now))

	for j, r := range i.Answer {
		m.Answer[j] = dns.Copy(r)
		m.Answer[j].Header().Ttl = ttl
	}
	for j, r := range i.Extra {
		m.Extra[j] = dns.Copy(r)
		m.Extra[j].Header().Ttl = ttl
	}
	return m
}
