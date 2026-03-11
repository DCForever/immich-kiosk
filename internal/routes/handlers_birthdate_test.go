package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/damongolding/immich-kiosk/internal/config"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetBirthdate_and_DeleteBirthdate(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	require.NoError(t, os.WriteFile(birthPath, []byte(`{}`), 0644))
	require.NoError(t, os.WriteFile(mapPath, []byte(`{}`), 0644))

	cfg := config.New()
	cfg.Source = config.SourcePhotoPrism
	cfg.BirthdateFilePath = birthPath
	cfg.BirthdateMappingPath = mapPath

	e := echo.New()

	// POST set birthdate
	body := []byte(`{"personId":"alice","birthDate":"1990-05-15"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/person-birthdate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	setHandler := SetBirthdate(cfg)
	require.NoError(t, setHandler(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	var setResp SetBirthdateResponse
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&setResp))
	assert.Equal(t, "ok", setResp.Status)

	// DELETE birthdate
	delBody := []byte(`{"personId":"alice"}`)
	req2 := httptest.NewRequest(http.MethodDelete, "/api/person-birthdate", bytes.NewReader(delBody))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)
	delHandler := DeleteBirthdate(cfg)
	require.NoError(t, delHandler(c2))
	assert.Equal(t, http.StatusOK, rec2.Code)
	var delResp DeleteBirthdateResponse
	require.NoError(t, json.NewDecoder(rec2.Body).Decode(&delResp))
	assert.Equal(t, "ok", delResp.Status)
}

func TestSetBirthdate_notConfigured(t *testing.T) {
	cfg := config.New()
	cfg.Source = config.SourcePhotoPrism
	cfg.BirthdateFilePath = ""
	cfg.BirthdateMappingPath = ""

	e := echo.New()
	body := []byte(`{"personId":"alice","birthDate":"1990-05-15"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/person-birthdate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	setHandler := SetBirthdate(cfg)
	require.NoError(t, setHandler(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSetBirthdate_invalidDate(t *testing.T) {
	dir := t.TempDir()
	birthPath := filepath.Join(dir, "birthdates.json")
	mapPath := filepath.Join(dir, "mapping.json")
	require.NoError(t, os.WriteFile(birthPath, []byte(`{}`), 0644))
	require.NoError(t, os.WriteFile(mapPath, []byte(`{}`), 0644))

	cfg := config.New()
	cfg.Source = config.SourcePhotoPrism
	cfg.BirthdateFilePath = birthPath
	cfg.BirthdateMappingPath = mapPath

	e := echo.New()
	body := []byte(`{"personId":"alice","birthDate":"2030-01-01"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/person-birthdate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	setHandler := SetBirthdate(cfg)
	require.NoError(t, setHandler(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
