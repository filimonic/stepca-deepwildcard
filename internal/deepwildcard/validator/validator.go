package validator

import (
	"deepwildcard/internal/deepwildcard/dnsutils"
	"fmt"
	"strings"
)

type Validator struct {
	config *Config
}

func New(options ...ValidatorOption) (*Validator, error) {
	v := &Validator{}
	for _, o := range options {
		err := o(v)
		if err != nil {
			return nil, err
		}
	}

	for _, rule := range v.config.Dns.Allow {
		if !dnsutils.IsValidDnsDomain(rule.trimAllKinds()) {
			return nil, fmt.Errorf("\"%s\" is not a valid dns \"allow\" rule", rule)
		}
		if strings.ToLower(string(rule)) != string(rule) {
			return nil, fmt.Errorf("\"%s\" dns \"allow\" mult be lowercase", rule)
		}
	}

	for _, rule := range v.config.Dns.Deny {
		if !dnsutils.IsValidDnsDomain(rule.trimAllKinds()) {
			return nil, fmt.Errorf("\"%s\" is not a valid dns \"deny\" rule", rule)
		}
		if strings.ToLower(string(rule)) != string(rule) {
			return nil, fmt.Errorf("\"%s\" dns \"deny\" mult be lowercase", rule)
		}
	}

	return v, nil
}

func (v *Validator) IsDomainAllowed(domain string) *Result {
	if len(v.config.Dns.Allow) < 0 {
		return ResultDenied("VD_ALLOW_EMPTY", "allow list is empty. everything will be denied")
	}

	if !dnsutils.IsValidDnsDomain(domain) {
		return ResultDenied("VD_INVALID", fmt.Sprintf("\"%s\" is not valid domain name", domain))
	}

	// Check against allow list
	domain = strings.ToLower(domain)
	allowedMatched := false
	for _, rule := range v.config.Dns.Allow {
		if !allowedMatched && rule.IsMatchDomain(domain) {
			allowedMatched = true
			break
		}
	}

	if !allowedMatched {
		return ResultDenied("VD_ALLOW_NOT_MATCH", fmt.Sprintf("no \"allow\" rules matched \"%s\" domain", domain))
	}

	deniedMatched := false
	var denyMatchedRule DnsRule
	for _, rule := range v.config.Dns.Deny {
		if !deniedMatched && rule.IsMatchDomain(domain) {
			deniedMatched = true
			denyMatchedRule = rule
			break
		}
	}

	if allowedMatched && deniedMatched {
		return ResultDenied("VD_DENY_MATCH", fmt.Sprintf("\"deny\" rule \"%s\" matched \"%s\" domain", denyMatchedRule.String(), domain))
	}

	if allowedMatched && !deniedMatched {
		return ResultAllowed()
	}

	return ResultDenied("VD_UNEXPECTED", "unexpected behavoir")
}
