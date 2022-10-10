package core

import (
	"context"
	"fmt"
	"github.com/jmwri/dnsdiff/internal/diff/domain"
	dnsCore "github.com/jmwri/dnsdiff/internal/dns/core"
)

func NewService(nsA, nsB string) *Service {
	return &Service{
		nsA: nsA,
		nsB: nsB,
	}
}

type Service struct {
	nsA, nsB string
}

func (s *Service) HostDiff(ctx context.Context, host string) (domain.HostDiff, error) {
	var diff domain.HostDiff
	resultA, err := dnsCore.Lookup(ctx, s.nsA, host)
	if err != nil {
		return diff, fmt.Errorf("failed A lookup: %w", err)
	}
	resultB, err := dnsCore.Lookup(ctx, s.nsB, host)
	if err != nil {
		return diff, fmt.Errorf("failed B lookup: %w", err)
	}
	diff, err = resultDiff(resultA, resultB)
	if err != nil {
		return diff, fmt.Errorf("failed to generate diff: %w", err)
	}
	return diff, nil
}
