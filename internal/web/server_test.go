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
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleIndex(w, r, true)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expectedLinks := []string{"10m", "30m", "1h", "3h", "6h", "12h", "24h", "Max"}
	for _, link := range expectedLinks {
		if !strings.Contains(rr.Body.String(), link) {
			t.Errorf("handler returned body does not contain link for %v", link)
		}
	}
}
