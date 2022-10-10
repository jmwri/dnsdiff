package core

import (
	dnsDomain "github.com/jmwri/dnsdiff/internal/dns/domain"
	"github.com/kr/pretty"
)

func resultDiff(a, b dnsDomain.Result) ([]string, error) {
	diffs := pretty.Diff(a, b)
	return diffs, nil
}
