package plugin

import "context"

// Done is a non-blocking function that returns true if the context has been canceled.
// If a context is canceled it is highly likely the client that sent us this request has
// stopped caring. This means no reply whould be sent back.
//
// Typical use case for returning from a plugin when the context has been canceled:
// (p is a plugin.Handler)
//
// if plugin.Done(ctx) {
// 	return dns.RcodeSuccess, plugin.Error(p.Name(), ctx.Err())
// }
//
func Done(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
	return false
}
