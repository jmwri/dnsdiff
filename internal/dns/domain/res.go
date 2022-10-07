package domain

import "net"

type Result struct {
	Addresses []string
	CNAME     string
	MX        []*net.MX
	NS        []*net.NS
	SRV       []*net.SRV
	TXT       []string
}

func (r Result) IsEmpty() bool {
	if len(r.Addresses) > 0 {
		return false
	}
	if r.CNAME != "" {
		return false
	}
	if len(r.MX) > 0 {
		return false
	}
	if len(r.SRV) > 0 {
		return false
	}
	if len(r.TXT) > 0 {
		return false
	}
	return true
}
