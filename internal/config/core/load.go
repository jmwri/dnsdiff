package core

import (
	"context"
	"fmt"
	"github.com/jmwri/dnsdiff/internal/config/domain"
	"github.com/jmwri/dnsdiff/internal/config/port"
)

func Load(ctx context.Context, loader port.ConfigLoader) (domain.Config, error) {
	cfg, err := loader(ctx)
	if err != nil {
		return cfg, fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.NSA == "" {
		return cfg, fmt.Errorf("invalid NS A")
	}
	if cfg.NSB == "" {
		return cfg, fmt.Errorf("invalid NS B")
	}
	return cfg, nil
}
