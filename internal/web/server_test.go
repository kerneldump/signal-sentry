package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleIndex(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleIndex)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedLinks := []string{"5m", "15m", "30m", "45m", "1h", "24h"}
	for _, link := range expectedLinks {
		if !strings.Contains(rr.Body.String(), link) {
			t.Errorf("handler returned body does not contain link for %v", link)
		}
	}
}
