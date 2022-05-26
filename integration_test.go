//nolint:lll,paralleltest
package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	is := is.New(t)

	conf := newConfig()
	srv := newServer(conf)

	for _, tc := range []struct {
		key   string
		value string
	}{
		{key: "foo", value: "bar"},
	} {
		t.Run(tc.key+"/"+tc.value, func(t *testing.T) {
			routeKey := "/v1/" + tc.key
			routeKeyValue := routeKey + "/" + tc.value

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
				{name: "get", method: http.MethodGet, path: routeKey, expectedStatus: http.StatusOK, expectedResp: tc.value, auth: true},
				{name: "del", method: http.MethodDelete, path: routeKey, expectedStatus: http.StatusOK, auth: true},
				{name: "get after del", method: http.MethodGet, path: routeKey, expectedStatus: http.StatusNotFound, auth: true},
			} {
				t.Run(step.name, func(t *testing.T) {
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
}
