package domain

func NewResult() Result {
	return Result{
		A:     make([]string, 0),
		AAAA:  make([]string, 0),
		CNAME: make([]string, 0),
		MX:    make([]string, 0),
		NS:    make([]string, 0),
		SRV:   make([]string, 0),
		TXT:   make([]string, 0),
		SOA:   make([]string, 0),
		CAA:   make([]string, 0),
		PTR:   make([]string, 0),
	}
}

type Result struct {
	A        []string
	AErr     string
	AAAA     []string
	AAAAErr  string
	CNAME    []string
	CNAMEErr string
	MX       []string
	MXErr    string
	NS       []string
	NSErr    string
	SRV      []string
	SRVErr   string
	TXT      []string
	TXTErr   string
	SOA      []string
	SOAErr   string
	CAA      []string
	CAAErr   string
	PTR      []string
	PTRErr   string
}
