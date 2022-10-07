package main

import (
	"context"
	"flag"
	"fmt"
	configCore "github.com/jmwri/dnsdiff/internal/config/core"
	configDomain "github.com/jmwri/dnsdiff/internal/config/domain"
	"github.com/jmwri/dnsdiff/internal/diff"
	diffCore "github.com/jmwri/dnsdiff/internal/diff/core"
	"github.com/jmwri/dnsdiff/internal/log"
	"github.com/jmwri/dnsdiff/internal/util"
	"go.uber.org/zap"
	"io"
	"os"
)

var ctx context.Context
var cancel context.CancelFunc

func init() {
	ctx, cancel = context.WithCancel(context.Background())
}

var cfg configDomain.Config

func init() {
	var err error
	cfg, err = configCore.Load(ctx, func(ctx context.Context) (configDomain.Config, error) {
		nsA := flag.String("a", "", "First DNS server")
		nsB := flag.String("b", "", "Second DNS server")
		hostsPath := flag.String("hosts", "", "Hosts path")
		outPath := flag.String("out", "stdout", "Output file path. 'stdout' for stdout.")
		parent := flag.String("p", "", "Parent host appended to all hosts")
		flag.Parse()

		cfg := configDomain.Config{
			NSA:       *nsA,
			NSB:       *nsB,
			HostsPath: *hostsPath,
			OutPath:   *outPath,
			Parent:    *parent,
		}

		return cfg, nil
	})
	if err != nil {
		log.Fatal(ctx, "invalid args", zap.Error(err))
	}
}

func main() {
	var diffService diff.Service = diffCore.NewService(cfg.NSA, cfg.NSB)

	var outputWriter io.WriteCloser = os.Stdout
	if cfg.OutPath != "stdout" {
		var err error
		outputWriter, err = os.Create(cfg.OutPath)
		if err != nil {
			log.Fatal(ctx, "failed to open output file", zap.String("path", cfg.OutPath), zap.Error(err))
			return
		}
		defer func(outputWriter io.WriteCloser) {
			_ = outputWriter.Close()
		}(outputWriter)
	}

	hosts, err := util.ReadLines(cfg.HostsPath)
	if err != nil {
		log.Fatal(ctx, "failed to read hosts", zap.Error(err))
	}

	adjustedHosts := make([]string, 0)

	if cfg.Parent != "" {
		for _, host := range hosts {
			if host == "" {
				continue
			}
			if host == "@" {
				adjustedHosts = append(adjustedHosts, cfg.Parent)
			} else {
				adjustedHosts = append(adjustedHosts, fmt.Sprintf("%s.%s", host, cfg.Parent))
			}
		}
	}

	for _, host := range adjustedHosts {
		if host == "" {
			continue
		}
		diffLines, err := diffService.HostDiff(ctx, host)
		if err != nil {
			log.Fatal(ctx, "failed to get diff for hosts", zap.Error(err))
		}
		if len(diffLines) == 0 {
			_, err = fmt.Fprintf(outputWriter, "%s same or delegated\n", host)
			if err != nil {
				log.Fatal(ctx, "failed to write output", zap.String("path", cfg.OutPath), zap.Error(err))
			}
			continue
		}
		_, err = fmt.Fprintln(outputWriter, host)
		if err != nil {
			log.Fatal(ctx, "failed to write output", zap.String("path", cfg.OutPath), zap.Error(err))
		}
		for _, diffLine := range diffLines {
			_, err = fmt.Fprintf(outputWriter, "    %s\n", diffLine)
			if err != nil {
				log.Fatal(ctx, "failed to write output", zap.String("path", cfg.OutPath), zap.Error(err))
			}
		}
	}
}
