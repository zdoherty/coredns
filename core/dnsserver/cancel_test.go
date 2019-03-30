package dnsserver

import (
	"context"
	"testing"

	"github.com/miekg/dns"
)

// We can't include the actual cancel plugin here, because that create a cyclic dependency. Just make
// a plugin with the same name because that's what triggers the check.
type cancel struct{}

func (c cancel) Name() string { return "cancel" }
func (c cancel) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return 0, nil
}

func TestCancelEnabled(t *testing.T) {
	c := testConfig("dns", cancel{})
	s, err := NewServer("127.0.0.1:53", []*Config{c})

	if err != nil {
		t.Fatalf("Expected no error for NewServer, got %s", err)
	}
	if !s.cancel {
		t.Errorf("Expected cancel to be true, got false")
	}
}

func TestCancelDisabled(t *testing.T) {
	c := testConfig("dns", testPlugin{})
	s, err := NewServer("127.0.0.1:53", []*Config{c})

	if err != nil {
		t.Fatalf("Expected no error for NewServer, got %s", err)
	}
	if s.cancel {
		t.Errorf("Expected cancel to be false, got true")
	}
}
