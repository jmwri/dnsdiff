package diff

import (
	"context"
	"github.com/jmwri/dnsdiff/internal/diff/domain"
)

type Service interface {
	HostDiff(ctx context.Context, host string) (domain.HostDiff, error)
}
