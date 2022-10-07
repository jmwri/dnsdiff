package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmwri/dnsdiff/internal/dns/domain"
	"net"
	"sort"
	"strings"
	"time"
)

func Lookup(ctx context.Context, res *net.Resolver, host string) (domain.Result, error) {
	var result domain.Result
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		result, err = doLookup(ctx, res, host)
		if err != nil {
			if isTempDNSErr(err) {
				time.Sleep(time.Second * time.Duration(attempt))
				continue
			}
			return result, fmt.Errorf("failed lookup: %w", err)
		}
		break
	}
	return result, nil
}

func doLookup(ctx context.Context, res *net.Resolver, host string) (domain.Result, error) {
	var result domain.Result
	var err error

	if result.Addresses, err = res.LookupHost(ctx, host); err != nil && !canIgnoreDNSErr(err) {
		return result, fmt.Errorf("failed to lookup host: %w", err)
	}

	if result.CNAME, err = res.LookupCNAME(ctx, host); err != nil && !canIgnoreDNSErr(err) {
		return result, fmt.Errorf("failed to lookup cname: %w", err)
	}

	if result.MX, err = res.LookupMX(ctx, host); err != nil && !canIgnoreDNSErr(err) {
		return result, fmt.Errorf("failed to lookup mx: %w", err)
	}

	if result.NS, err = res.LookupNS(ctx, host); err != nil && !canIgnoreDNSErr(err) {
		return result, fmt.Errorf("failed to lookup ns: %w", err)
	}

	srvService, srvProto, srvHost, srvErr := getSrvParams(host)
	if srvErr == nil {
		if _, result.SRV, err = res.LookupSRV(ctx, srvService, srvProto, srvHost); err != nil && !canIgnoreDNSErr(err) {
			return result, fmt.Errorf("failed to lookup srv: %w", err)
		}
	}

	if result.TXT, err = res.LookupTXT(ctx, host); err != nil && !canIgnoreDNSErr(err) {
		return result, fmt.Errorf("failed to lookup txt: %w", err)
	}

	return sortResult(result), nil
}

func isTempDNSErr(err error) bool {
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return dnsErr.Temporary()
	}
	return false
}

func canIgnoreDNSErr(err error) bool {
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return dnsErr.IsNotFound || dnsErr.Err == "lame referral"
	}
	return false
}

var ErrInvalidSrvHost = errors.New("invalid host for srv")
var ErrInvalidSrvSvc = errors.New("invalid service for srv")
var ErrInvalidSrvProto = errors.New("invalid service for srv")

func getSrvParams(host string) (service string, proto string, shortHost string, err error) {
	splitHost := strings.Split(host, ".")
	if len(splitHost) < 3 {
		err = ErrInvalidSrvHost
		return
	}
	service = splitHost[0]
	if len(service) < 2 || service[0:1] != "_" {
		err = ErrInvalidSrvSvc
		return
	}
	service = service[1:]

	proto = splitHost[1]
	if len(proto) < 2 || proto[0:1] != "_" {
		err = ErrInvalidSrvProto
		return
	}
	proto = proto[1:]
	shortHost = strings.Join(splitHost[2:], ".")
	return
}

func sortResult(r domain.Result) domain.Result {
	sort.Strings(r.Addresses)
	sort.Slice(r.MX, func(i, j int) bool {
		return r.MX[i].Host < r.MX[j].Host
	})
	sort.Slice(r.NS, func(i, j int) bool {
		return r.NS[i].Host < r.NS[j].Host
	})
	sort.Slice(r.SRV, func(i, j int) bool {
		return r.SRV[i].Priority < r.SRV[j].Priority
	})
	sort.Strings(r.TXT)
	return r
}
