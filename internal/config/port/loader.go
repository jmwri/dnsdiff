package port

import (
	"context"
	"github.com/jmwri/dnsdiff/internal/config/domain"
)

type ConfigLoader func(ctx context.Context) (domain.Config, error)
