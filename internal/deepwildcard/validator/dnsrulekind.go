package validator

type DnsRuleKind string

const DNSRULEKIND_GRAND = DnsRuleKind(`+.`)
const DNSRULEKIND_CHILD = DnsRuleKind(`*.`)
const DNSRULEKIND_EXACT = DnsRuleKind(``)

func (k DnsRuleKind) String() string {
	return string(k)
}
