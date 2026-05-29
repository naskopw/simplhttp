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

	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
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
