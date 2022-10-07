package core

import (
	"context"
	"fmt"
	"net"
	"time"
)

func ResolverForNS(ns string) *net.Resolver {
	return &net.Resolver{
		PreferGo:     true,
		StrictErrors: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(5000),
			}
			return d.DialContext(ctx, network, fmt.Sprintf("%s:%d", ns, 53))
		},
	}
}
