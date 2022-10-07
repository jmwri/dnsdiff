package core

import (
	"fmt"
	dnsDomain "github.com/jmwri/dnsdiff/internal/dns/domain"
	"github.com/kr/pretty"
	"reflect"
	"strings"
)

func resultDiff(a, b dnsDomain.Result) ([]string, error) {
	// Ignore any diff if NS are the same.
	if reflect.DeepEqual(a.NS, b.NS) {
		return []string{}, nil
	}
	diffs := pretty.Diff(a, b)
	for i, diff := range diffs {
		split := strings.Split(diff, ":")
		switch split[0] {
		case "Addresses":
			diffs[i] = fmt.Sprintf("%s: %v != %v", split[0], a.Addresses, b.Addresses)
		case "TXT":
			diffs[i] = fmt.Sprintf("%s: %v != %v", split[0], a.TXT, b.TXT)
		}
	}
	return diffs, nil
}
