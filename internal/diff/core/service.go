package core

import (
	"context"
	"fmt"
	"github.com/jmwri/dnsdiff/internal/diff/domain"
	dnsCore "github.com/jmwri/dnsdiff/internal/dns/core"
	"net"
)

func NewService(nsA, nsB string) *Service {
	return &Service{
		resA: dnsCore.ResolverForNS(nsA),
		resB: dnsCore.ResolverForNS(nsB),
	}
}

type Service struct {
	resA, resB *net.Resolver
}

func (s *Service) HostDiff(ctx context.Context, host string) (domain.HostDiff, error) {
	var diff domain.HostDiff
	resultA, err := dnsCore.Lookup(ctx, s.resA, host)
	if err != nil {
		return diff, fmt.Errorf("failed A lookup: %w", err)
	}
	resultB, err := dnsCore.Lookup(ctx, s.resB, host)
	if err != nil {
		return diff, fmt.Errorf("failed B lookup: %w", err)
	}
	diff, err = resultDiff(resultA, resultB)
	if err != nil {
		return diff, fmt.Errorf("failed to generate diff: %w", err)
	}
	if len(diff) == 0 && resultA.IsEmpty() {
		return []string{"No records found in NS A"}, nil
	}
	return diff, nil
}
