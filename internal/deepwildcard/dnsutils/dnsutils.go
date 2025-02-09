package dnsutils

import (
	"regexp"
	"slices"
	"strings"
)

var correctTld = regexp.MustCompile(`(?i)^(?:[a-z]{2,63}|xn--[a-z0-9]{59})$`)
var correctDC = regexp.MustCompile(`(?i)^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)
var forbiddenTld = regexp.MustCompile(`(?i)^(?:arpa)$`)

func IsValidDnsDomain(domain string) bool {
	dcsRev := SplitDnsDomainToComponentsReversed(domain)
	if len(dcsRev) < 2 {
		return false
	}
	for idx, dc := range dcsRev {
		if (idx == 0 && !isValidDnsTld(dc)) || (idx != 0 && !isValidDnsDomainComponent(dc)) {
			return false
		}
	}
	return true
}

func SplitDnsDomainToComponents(domain string) []string {
	return strings.Split(domain, ".")
}

func SplitDnsDomainToComponentsReversed(domain string) []string {
	dcs := SplitDnsDomainToComponents(domain)
	slices.Reverse(dcs)
	return dcs
}

func isValidDnsDomainComponent(dc string) bool {
	return correctDC.MatchString(dc)
}

func isValidDnsTld(tld string) bool {
	return isValidDnsDomainComponent(tld) &&
		correctTld.MatchString(tld) &&
		!forbiddenTld.MatchString(tld)
}
