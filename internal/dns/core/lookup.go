package core

import (
	"context"
	"fmt"
	"github.com/jmwri/dnsdiff/internal/dns/domain"
	"github.com/jmwri/dnsdiff/internal/log"
	"github.com/miekg/dns"
	"go.uber.org/zap"
	"net"
	"sort"
	"time"
)

func Lookup(ctx context.Context, ns string, host string) (domain.Result, error) {
	log.WithFields(ctx, zap.String("ns", ns), zap.String("host", host))
	var result domain.Result
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		result, err = doLookup(ctx, ns, host)
		if err != nil {
			time.Sleep(time.Second * time.Duration(attempt))
		} else {
			break
		}
	}
	return result, err
}

func doLookup(ctx context.Context, ns string, host string) (domain.Result, error) {
	result := domain.NewResult()
	var err error

	c := new(dns.Client)
	c.SingleInflight = true

	// A
	if result.A, result.AErr, err = lookupA(ctx, c, ns, host); err != nil {
		return result, err
	}

	// AAAA
	if result.AAAA, result.AAAAErr, err = lookupAAAA(ctx, c, ns, host); err != nil {
		return result, err
	}

	// CNAME
	if result.CNAME, result.CNAMEErr, err = lookupCNAME(ctx, c, ns, host); err != nil {
		return result, err
	}

	// MX
	if result.MX, result.MXErr, err = lookupMX(ctx, c, ns, host); err != nil {
		return result, err
	}

	// NS
	if result.NS, result.NSErr, err = lookupNS(ctx, c, ns, host); err != nil {
		return result, err
	}

	// SRV
	if result.SRV, result.SRVErr, err = lookupSRV(ctx, c, ns, host); err != nil {
		return result, err
	}

	// TXT
	if result.TXT, result.TXTErr, err = lookupTXT(ctx, c, ns, host); err != nil {
		return result, err
	}

	// SOA
	if result.SOA, result.SOAErr, err = lookupSOA(ctx, c, ns, host); err != nil {
		return result, err
	}

	// CAA
	if result.CAA, result.CAAErr, err = lookupCAA(ctx, c, ns, host); err != nil {
		return result, err
	}

	// PTR
	if result.PTR, result.PTRErr, err = lookupPTR(ctx, c, ns, host); err != nil {
		return result, err
	}

	// SPF
	if result.SPF, result.SPFErr, err = lookupSPF(ctx, c, ns, host); err != nil {
		return result, err
	}

	return result, nil
}

func lookupA(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)
	seenRes := make(map[string]bool)
	var rCode string

	// Account for DNS load balancing by querying multiple times
	for i := 1; i <= 10; i++ {
		r, err := doRequest(ctx, c, ns, host, dns.TypeA)
		if err != nil {
			return res, "", fmt.Errorf("failed to lookup A: %w", err)
		}
		for _, a := range r.Answer {
			record, ok := a.(*dns.A)
			if !ok {
				continue
			}
			resStr := record.A.String()
			if seenRes[resStr] {
				continue
			}
			seenRes[resStr] = true
			res = append(res, resStr)
		}
		rCode = dns.RcodeToString[r.Rcode]
	}
	sort.Strings(res)
	return res, rCode, nil
}

func lookupAAAA(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)
	seenRes := make(map[string]bool)
	var rCode string

	// Account for DNS load balancing by querying multiple times
	for i := 1; i <= 10; i++ {
		r, err := doRequest(ctx, c, ns, host, dns.TypeAAAA)
		if err != nil {
			return res, "", fmt.Errorf("failed to lookup AAAA: %w", err)
		}
		for _, a := range r.Answer {
			record, ok := a.(*dns.AAAA)
			if !ok {
				continue
			}
			resStr := record.AAAA.String()
			if seenRes[resStr] {
				continue
			}
			seenRes[resStr] = true
			res = append(res, resStr)
		}
		rCode = dns.RcodeToString[r.Rcode]
	}
	sort.Strings(res)
	return res, rCode, nil
}

func lookupCNAME(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)

	r, err := doRequest(ctx, c, ns, host, dns.TypeCNAME)
	if err != nil {
		return res, "", fmt.Errorf("failed to lookup CNAME: %w", err)
	}
	for _, a := range r.Answer {
		record, ok := a.(*dns.CNAME)
		if !ok {
			continue
		}
		res = append(res, record.Target)
	}
	sort.Strings(res)
	return res, dns.RcodeToString[r.Rcode], nil
}

func lookupMX(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)

	r, err := doRequest(ctx, c, ns, host, dns.TypeMX)
	if err != nil {
		return res, "", fmt.Errorf("failed to lookup MX: %w", err)
	}
	for _, a := range r.Answer {
		record, ok := a.(*dns.MX)
		if !ok {
			continue
		}
		res = append(res, record.Mx)
	}
	sort.Strings(res)
	return res, dns.RcodeToString[r.Rcode], nil
}

func lookupNS(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)

	r, err := doRequest(ctx, c, ns, host, dns.TypeNS)
	if err != nil {
		return res, "", fmt.Errorf("failed to lookup NS: %w", err)
	}
	for _, a := range r.Ns {
		record, ok := a.(*dns.NS)
		if !ok {
			continue
		}
		res = append(res, record.Ns)
	}
	sort.Strings(res)
	return res, dns.RcodeToString[r.Rcode], nil
}

func lookupSRV(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)

	r, err := doRequest(ctx, c, ns, host, dns.TypeSRV)
	if err != nil {
		return res, "", fmt.Errorf("failed to lookup SRV: %w", err)
	}
	for _, a := range r.Answer {
		record, ok := a.(*dns.SRV)
		if !ok {
			continue
		}
		res = append(res, fmt.Sprintf("%d %d %d %s", record.Priority, record.Weight, record.Port, record.Target))
	}
	sort.Strings(res)
	return res, dns.RcodeToString[r.Rcode], nil
}

func lookupTXT(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)

	r, err := doRequest(ctx, c, ns, host, dns.TypeTXT)
	if err != nil {
		return res, "", fmt.Errorf("failed to lookup TXT: %w", err)
	}
	for _, a := range r.Answer {
		record, ok := a.(*dns.TXT)
		if !ok {
			continue
		}
		for _, v := range record.Txt {
			res = append(res, v)
		}
	}
	sort.Strings(res)
	return res, dns.RcodeToString[r.Rcode], nil
}

func lookupSOA(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)

	r, err := doRequest(ctx, c, ns, host, dns.TypeSOA)
	if err != nil {
		return res, "", fmt.Errorf("failed to lookup SOA: %w", err)
	}
	for _, a := range r.Answer {
		record, ok := a.(*dns.SOA)
		if !ok {
			continue
		}
		res = append(res, fmt.Sprintf("%s %s %d %d %d %d %d", record.Ns, record.Mbox, record.Serial, record.Refresh, record.Retry, record.Expire, record.Minttl))
	}
	sort.Strings(res)
	return res, dns.RcodeToString[r.Rcode], nil
}

func lookupCAA(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)

	r, err := doRequest(ctx, c, ns, host, dns.TypeCAA)
	if err != nil {
		return res, "", fmt.Errorf("failed to lookup CAA: %w", err)
	}
	for _, a := range r.Answer {
		record, ok := a.(*dns.CAA)
		if !ok {
			continue
		}
		res = append(res, fmt.Sprintf("%d %s %s", record.Flag, record.Tag, record.Value))
	}
	sort.Strings(res)
	return res, dns.RcodeToString[r.Rcode], nil
}

func lookupPTR(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)

	r, err := doRequest(ctx, c, ns, host, dns.TypePTR)
	if err != nil {
		return res, "", fmt.Errorf("failed to lookup PTR: %w", err)
	}
	for _, a := range r.Answer {
		record, ok := a.(*dns.PTR)
		if !ok {
			continue
		}
		res = append(res, record.Ptr)
	}
	sort.Strings(res)
	return res, dns.RcodeToString[r.Rcode], nil
}

func lookupSPF(ctx context.Context, c *dns.Client, ns, host string) ([]string, string, error) {
	res := make([]string, 0)

	r, err := doRequest(ctx, c, ns, host, dns.TypeSPF)
	if err != nil {
		return res, "", fmt.Errorf("failed to lookup SPF: %w", err)
	}
	for _, a := range r.Answer {
		record, ok := a.(*dns.SPF)
		if !ok {
			continue
		}
		for _, v := range record.Txt {
			res = append(res, v)
		}
	}
	sort.Strings(res)
	return res, dns.RcodeToString[r.Rcode], nil
}

func doRequest(ctx context.Context, c *dns.Client, ns, host string, recordType uint16) (*dns.Msg, error) {
	for {
		m := new(dns.Msg)
		m.RecursionDesired = true
		m.SetQuestion(host, recordType)
		r, _, err := c.ExchangeContext(ctx, m, net.JoinHostPort(ns, "53"))
		if err != nil {
			return r, fmt.Errorf("failed to lookup: %w", err)
		}

		if r.MsgHdr.Authoritative {
			return r, nil
		}

		if len(r.Ns) == 0 {
			return r, fmt.Errorf("unable to recurse")
		}

		ns = r.Ns[0].(*dns.NS).Ns
	}
}
