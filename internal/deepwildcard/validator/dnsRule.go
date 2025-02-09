package validator

import (
	"deepwildcard/internal/deepwildcard/dnsutils"
	"fmt"
	"strings"
)

type DnsRule string

func (r DnsRule) String() string {
	return string(r)
}

func (r *DnsRule) IsValid() bool {
	return dnsutils.IsValidDnsDomain(r.trimAllKinds())
}

func (r *DnsRule) Kind() DnsRuleKind {
	if strings.HasPrefix(r.String(), DNSRULEKIND_CHILD.String()) {
		return DNSRULEKIND_CHILD
	}
	if strings.HasPrefix(r.String(), DNSRULEKIND_GRAND.String()) {
		return DNSRULEKIND_GRAND
	}
	return DNSRULEKIND_EXACT
}

func (r *DnsRule) IsKind(kind DnsRuleKind) bool {
	return r.Kind() == kind
}

func (r *DnsRule) trimAllKinds() string {
	d := r.String()
	for _, k := range []string{
		DNSRULEKIND_GRAND.String(),
		DNSRULEKIND_CHILD.String(),
	} {
		d = strings.TrimPrefix(d, k)
	}
	return d
}

func (r *DnsRule) IsMatchDomain(domain string) bool {
	kind := r.Kind()
	ruleDomain := r.trimAllKinds()
	ruleDCs := dnsutils.SplitDnsDomainToComponents(ruleDomain)
	ruleDCCount := len(ruleDCs)
	DCs := dnsutils.SplitDnsDomainToComponents(domain)
	DCCount := len(DCs)
	isSameSuffix := strings.HasSuffix(domain, fmt.Sprintf(".%s", r.trimAllKinds()))
	switch kind {
	case DNSRULEKIND_EXACT:
		return strings.EqualFold(domain, r.String())
	case DNSRULEKIND_CHILD:
		return isSameSuffix && (DCCount == (ruleDCCount + 1))
	case DNSRULEKIND_GRAND:
		return isSameSuffix && (DCCount > (ruleDCCount + 1))
	default:
		return false
	}
}
