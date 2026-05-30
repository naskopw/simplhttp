package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandleConfig(t *testing.T) {
	e := echo.New()
	s := &Server{
		echo: e,
		config: Config{
			MaxSizeBytes: 500,
			ReadOnly:     true,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/simpl/api/config", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := s.handleConfig(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var config ServerConfig
	err = json.Unmarshal(rec.Body.Bytes(), &config)
	assert.NoError(t, err)

	assert.Equal(t, int64(500), config.MaxSizeBytes)
	assert.True(t, config.ReadOnly)
}

func TestSPAFallback(t *testing.T) {
	s := NewServer(Config{
		BaseRoute: "/my-app",
		Dir:       ".",
	})

	// Test route within BaseRoute should return index.html (SPA fallback)
	req := httptest.NewRequest(http.MethodGet, "/my-app/some-page", nil)
	rec := httptest.NewRecorder()
	
	s.echo.ServeHTTP(rec, req)
	
	// We expect 200 OK because of the SPA fallback serving index.html
	// Note: in tests ui.Assets dist/index.html might not exist or be empty if not built
	// but the handler should at least try to serve it.
	// Since we are using NewServer, it sets up the HTTPErrorHandler.
	assert.Equal(t, http.StatusOK, rec.Code)

	// Test route within /simpl/api should NOT return SPA fallback (should be 404 if not found)
	req = httptest.NewRequest(http.MethodGet, "/simpl/api/not-found", nil)
	rec = httptest.NewRecorder()
	s.echo.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
