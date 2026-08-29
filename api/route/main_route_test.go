package route

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRouterRejectsMalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/suggest", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	NewRouter().ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", resp.Code, http.StatusBadRequest)
	}
	if !strings.Contains(resp.Body.String(), "A-1-10001") {
		t.Fatalf("response = %s, want client validation code", resp.Body.String())
	}
}
