package cache

import (
	"time"

	"github.com/coredns/coredns/plugin/cache/freq"

	"github.com/miekg/dns"
)

type delegationItem struct {
	Ns    []dns.RR
	Extra []dns.RR

	meta
}

func (i *delegationItem) ttl(now time.Time) int { return i.meta.ttl(now) }
func (i *delegationItem) m() meta               { return i.meta }

func (i *delegationItem) fromMsg(m *dns.Msg, now time.Time, d time.Duration) {
	i.Ns = m.Ns
	i.Extra = copyExtra(m.Extra)

	i.origTTL = uint32(d.Seconds())
	i.stored = now.UTC()
	i.Freq = new(freq.Freq)
}

func (i *delegationItem) toMsg(r *dns.Msg, now time.Time) *dns.Msg {
	m := new(dns.Msg)
	m.SetReply(r)
	m.RecursionAvailable = true
	m.Rcode = dns.RcodeSuccess
	m.Ns = make([]dns.RR, len(i.Ns))
	m.Extra = make([]dns.RR, len(i.Extra))

	ttl := uint32(i.meta.ttl(now))

	for j, r := range i.Ns {
		m.Ns[j] = dns.Copy(r)
		m.Ns[j].Header().Ttl = ttl
	}
	for j, r := range i.Extra {
		m.Extra[j] = dns.Copy(r)
		m.Extra[j].Header().Ttl = ttl
	}
	return m
}
