package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/crunchydeer30/lets-go/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	defer ts.Close()

	code, _, body := ts.get(t, "/ping")
	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	fmt.Println(app.templateCache)

	defer ts.Close()

	tests := []struct {
		name     string
		url      string
		wantCode int
		wantBody string
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, "An old silent pond..."},
		// {"Non-existent ID", "/snippet/2", http.StatusNotFound, "No snippet with ID 2"},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, "Not Found\n"},
		{"String ID", "/snippet/abc", http.StatusNotFound, "Not Found\n"},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, "Not Found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.url)
			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
