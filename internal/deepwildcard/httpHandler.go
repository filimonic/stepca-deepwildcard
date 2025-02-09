package deepwildcard

import (
	"deepwildcard/internal/deepwildcard/webhook"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (dw *dhServer) httpX509AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		dw.writeDenied(w, "MNA", "method not allowed")
		return
	}

	if r.Header.Get("X-Request-Id") == "" ||
		r.Header.Get("X-Smallstep-Signature") == "" ||
		r.Header.Get("X-Smallstep-Webhook-Id") == "" {
		dw.writeDenied(w, "NO_STEPCA", "requestee is not step-ca")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		dw.writeDenied(w, "NO_READ", fmt.Sprintf("failed to read body: %s", err))
		return
	}

	whr, err := webhook.UnmarshalWebkookRequestJson(body)
	if err != nil {
		dw.writeDenied(w, "NO_UNMARSH", fmt.Sprintf("failed to unmarshal body: %s", err))
		return
	}

	if len(whr.X509Certificate.SANs) <= 0 {
		dw.writeDenied(w, "NO_SANS", "no SANs found in certificate")
		return
	}

	allowedSansToLog := []string{}
	for _, san := range whr.X509Certificate.SANs {
		if result := dw.validator.IsDomainAllowed(san.Value); !result.Allowed {
			dw.writeDenied(w, result.Reason.Code, result.Reason.Message.Error())
			return
		}
		allowedSansToLog = append(allowedSansToLog, san.Value)
	}

	dw.writeAllowed(w, "ALL_VALID", fmt.Sprintf("All allowed: %s", strings.Join(allowedSansToLog, ", ")))
	return
}

func (dw *dhServer) writeDenied(w http.ResponseWriter, code string, message string) {
	dw.LogF("DENY  :%s (%s)", message, code)
	if err := webhook.CreateDenied(code, message).Write(w); err != nil {
		errText := fmt.Sprintf("Failed to write \"%s\" response: %s", code, err.Error())
		dw.LogF("FAIL  :%s", errText)
		http.Error(w, errText, http.StatusInternalServerError)
	}
}

func (dw *dhServer) writeAllowed(w http.ResponseWriter, code string, message string) {
	dw.LogF("ALLOW :%s (%s)", message, code)
	if err := webhook.CreateAllowed().Write(w); err != nil {
		errText := fmt.Sprintf("Failed to write allow response: %s", err.Error())
		dw.LogF("FAIL  :%s", errText)
		http.Error(w, errText, http.StatusInternalServerError)
	}
}
