package middleware

import (
	"mail/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkHeaders   bool
	}{
		{
			name:           "OPTIONS request",
			method:         http.MethodOptions,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.HTTPServer.AllowedIPsByCORS = []string{"http://localhost:3000"}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(tt.method, "/test", nil)
			rr := httptest.NewRecorder()

			handler := CORS(nextHandler, cfg)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.checkHeaders {
				assert.Equal(t, cfg.HTTPServer.AllowedIPsByCORS[0], rr.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "GET, POST, PUT, OPTIONS, DELETE", rr.Header().Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "Content-Type", rr.Header().Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
			}
		})
	}
}
