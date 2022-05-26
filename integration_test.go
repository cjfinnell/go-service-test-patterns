//nolint:lll
package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func FuzzIntegration(f *testing.F) {
	if testing.Short() {
		f.Skip("skipping integration test")
	}

	conf := newConfig()
	srv := newServer(conf)

	f.Add("key", "value")
	f.Add("e4e810c8c1ebc58a29df86222e68", "a57cfc8cad97fc576e01dfb3924e80102199")

	f.Fuzz(func(t *testing.T, fuzzKey, fuzzValue string) {
		if strings.TrimSpace(fuzzKey) == "" || strings.TrimSpace(fuzzValue) == "" {
			t.Skip()
		}
		if url.PathEscape(fuzzKey+fuzzValue) != fuzzKey+fuzzValue {
			t.Skip()
		}

		routeKey := "/v1/" + fuzzKey
		routeKeyValue := routeKey + "/" + fuzzValue

		for _, step := range []struct {
			name           string
			method         string
			path           string
			auth           bool
			expectedStatus int
			expectedResp   string
		}{
			{name: "status", method: http.MethodGet, path: "/status", expectedStatus: http.StatusOK},
			{name: "no auth get", method: http.MethodGet, path: routeKey, expectedStatus: http.StatusUnauthorized},
			{name: "no auth set", method: http.MethodPost, path: routeKeyValue, expectedStatus: http.StatusUnauthorized},
			{name: "no auth del", method: http.MethodDelete, path: routeKey, expectedStatus: http.StatusUnauthorized},
			{name: "get before set", method: http.MethodGet, path: routeKey, expectedStatus: http.StatusNotFound, auth: true},
			{name: "set", method: http.MethodPost, path: routeKeyValue, expectedStatus: http.StatusCreated, auth: true},
			{name: "get", method: http.MethodGet, path: routeKey, expectedStatus: http.StatusOK, expectedResp: fuzzValue, auth: true},
			{name: "del", method: http.MethodDelete, path: routeKey, expectedStatus: http.StatusOK, auth: true},
			{name: "get after del", method: http.MethodGet, path: routeKey, expectedStatus: http.StatusNotFound, auth: true},
		} {
			t.Run(step.name, func(t *testing.T) {
				is := is.New(t)
				w := httptest.NewRecorder()
				req := httptest.NewRequest(step.method, step.path, nil)
				if step.auth {
					req.Header.Set("Authorization", "test-auth")
				}
				srv.router.ServeHTTP(w, req)
				is.Equal(step.expectedStatus, w.Code)
				if step.expectedResp != "" {
					is.Equal(step.expectedResp, strings.TrimSpace(w.Body.String()))
				}
			})
		}
	})
}
