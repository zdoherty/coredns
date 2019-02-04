package cache

import (
	"context"
	"testing"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

func BenchmarkCacheMemory(b *testing.B) {
	c := New()
	c.Next = BackendMemoryHandler()
	ctx := context.TODO()

	names := []string{"example1", "example2", "a", "b", "c", "d", "e", "f", "g", "h", "i", "A", "B", "C", "D", "E", "F", "G", "H"}

	reqs := make([]*dns.Msg, len(names))
	for i, q := range names {
		reqs[i] = new(dns.Msg)
		reqs[i].SetQuestion(q+".example.org.", dns.TypeA)
	}

	b.ReportAllocs()
	b.StartTimer()
	j := 0
	for i := 0; i < b.N; i++ {
		req := reqs[j]
		c.ServeDNS(ctx, &test.ResponseWriter{}, req)
		j = (j + 1) % len(names)
	}
}

func BackendMemoryHandler() plugin.Handler {
	return plugin.HandlerFunc(func(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Response = true
		m.RecursionAvailable = true

		owner := m.Question[0].Name
		// uppercase letter return NXDOMAIN
		if byte(owner[0]) >= 65 && byte(owner[0]) < 90 {
			m.Ns = []dns.RR{test.SOA("example.org IN SOA 1 2 3 4 5 5")}
			m.Rcode = dns.RcodeNameError
			w.WriteMsg(m)
			return dns.RcodeSuccess, nil
		}
		m.Answer = []dns.RR{test.A(owner + " 303 IN A 127.0.0.53")}
		w.WriteMsg(m)
		return dns.RcodeSuccess, nil
	})
}
