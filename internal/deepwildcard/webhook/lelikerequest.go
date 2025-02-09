package webhook

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"slices"

	stepWebhook "github.com/smallstep/certificates/webhook"
)

type LELikeRequest struct {
	*stepWebhook.RequestBody
}

func UnmarshalWebkookRequestJson(body []byte) (*LELikeRequest, error) {
	m := map[string]interface{}{}
	err := json.Unmarshal(body, &m)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling payload (1): %w", err)
	}

	if ok, field := checkNoUnknownFields(m, webhookPayloadRootKnownFields); !ok {
		return nil, fmt.Errorf("field \"%s\" is unexpected in payload", field)
	}

	// Check all x509Certificate fields are known
	if c, ok := m["x509Certificate"].(map[string]interface{}); !ok {
		return nil, fmt.Errorf("\"%s\" field in payload is not an object", "x509Certificate")
	} else {
		if ok, field := checkNoUnknownFields(c, x509CertificateRootKnownFields); !ok {
			return nil, fmt.Errorf("field \"%s\" is unexpected in x509Certificate", field)
		}
	}

	whr := &LELikeRequest{}
	err = json.Unmarshal(body, whr)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling payload (2): %w", err)
	}

	if ok, reason := checkCertIsLeLike(whr.X509Certificate); !ok {
		return nil, fmt.Errorf("disallowed certificate: %s", reason)
	}

	return whr, nil
}

func checkNoUnknownFields(m map[string]interface{}, fields []string) (bool, string) {
	for field := range m {
		if !slices.Contains(fields, field) {
			return false, field
		}
	}
	return true, ""
}

func checkCertIsLeLike(c *stepWebhook.X509Certificate) (bool, string) {
	if !c.Subject.IsEmpty() || len(c.RawSubject) != 0 {
		return false, "non-empty subjects not allowed"
	}

	if c.IPAddresses != nil && len(c.IPAddresses) != 0 {
		return false, "ip addresses are not allowed"
	}

	if c.EmailAddresses != nil && len(c.EmailAddresses) != 0 {
		return false, "email addresses are not allowed"
	}

	if c.URIs != nil && len(c.URIs) != 0 {
		return false, "uri addresses are not allowed"
	}

	if c.DNSNames != nil && len(c.DNSNames) != 0 {
		return false, "dns names in wrong place (should be in \"sans\")"
	}

	if len(c.Extensions) != 0 {
		return false, "dns names in wrong place (should be in \"sans\")"
	}

	if len(c.UnknownExtKeyUsage) != 0 {
		return false, "unknown extended key usage not allowed"
	}

	if c.KeyUsage != 1 { // digitalSignature
		return false, "key usage must be \"digitalSignature\""
	}

	if !slices.Contains(c.ExtKeyUsage, 1) { // serverAuth
		return false, "key usage must be have \"serverAuth\""
	}

	if slices.ContainsFunc(c.ExtKeyUsage,
		func(eku x509.ExtKeyUsage) bool {
			return eku != x509.ExtKeyUsageClientAuth &&
				eku != x509.ExtKeyUsageServerAuth
		}) {
		return false, "extended key usage must have \"serverAuth\" and optionally \"clientAuth\". No other allowed"
	}

	for _, s := range c.SANs {
		if s.Type != "dns" {
			return false, fmt.Sprintf("only \"dns\" SANs are allowed, got \"%s\"", s.Type)
		}
		if string(s.ASN1Value) != "" {
			return false, fmt.Sprintf("only \"dns\" SANs are allowed without ASN1 representation, got \"%s\"", string(s.ASN1Value))
		}
	}

	if c.BasicConstraints != nil {
		return false, "No basic constraints must be provided, found one"
	}
	if c.NameConstraints != nil {
		return false, "No name constraints must be provided, found one"
	}
	return true, ""
}
