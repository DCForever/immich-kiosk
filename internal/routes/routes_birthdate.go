// Package routes: birthdate edit API handlers for click-to-edit person birthdates.
// Only active when source is PhotoPrism and birthdate_file_path and birthdate_mapping_path are set.

package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/damongolding/immich-kiosk/internal/birthdate"
	"github.com/damongolding/immich-kiosk/internal/config"
	"github.com/damongolding/immich-kiosk/internal/i18n"
	"github.com/labstack/echo/v5"
)

// SetBirthdateRequest matches the birthdate-edit-api contract (personId, birthDate).
type SetBirthdateRequest struct {
	PersonID  string `json:"personId"`
	BirthDate string `json:"birthDate"`
}

// SetBirthdateResponse matches the birthdate-edit-api contract (status, message).
type SetBirthdateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// DeleteBirthdateRequest matches the birthdate-edit-api contract (personId).
type DeleteBirthdateRequest struct {
	PersonID string `json:"personId"`
}

// DeleteBirthdateResponse has the same shape as SetBirthdateResponse.
type DeleteBirthdateResponse struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// SetBirthdate handles POST /api/person-birthdate to add or update a person's birthdate.
func SetBirthdate(baseConfig *config.Config) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if baseConfig.Source != config.SourcePhotoPrism || strings.TrimSpace(baseConfig.BirthdateFilePath) == "" || strings.TrimSpace(baseConfig.BirthdateMappingPath) == "" {
			t := i18n.T()
			return c.JSON(http.StatusBadRequest, SetBirthdateResponse{
				Status:  "error",
				Message: t("birthdate_not_configured"),
			})
		}

		var req SetBirthdateRequest
		if err := c.Bind(&req); err != nil {
			t := i18n.T()
			return c.JSON(http.StatusBadRequest, SetBirthdateResponse{
				Status:  "error",
				Message: t("birthdate_invalid_request"),
			})
		}
		personID := strings.TrimSpace(req.PersonID)
		dob := strings.TrimSpace(req.BirthDate)
		if personID == "" {
			t := i18n.T()
			return c.JSON(http.StatusBadRequest, SetBirthdateResponse{
				Status:  "error",
				Message: t("birthdate_person_required"),
			})
		}
		if dob == "" {
			t := i18n.T()
			return c.JSON(http.StatusBadRequest, SetBirthdateResponse{
				Status:  "error",
				Message: t("birthdate_date_required"),
			})
		}
		if _, err := birthdate.ParseDOB(dob); err != nil {
			t := i18n.T()
			return c.JSON(http.StatusBadRequest, SetBirthdateResponse{
				Status:  "error",
				Message: t("birthdate_invalid_format"),
			})
		}
		if err := birthdate.ValidateDOBForSave(dob, time.Now()); err != nil {
			t := i18n.T()
			msg := t("birthdate_invalid_date")
			if err != nil {
				msg = err.Error()
			}
			return c.JSON(http.StatusBadRequest, SetBirthdateResponse{
				Status:  "error",
				Message: msg,
			})
		}

		err := birthdate.SetBirthdate(baseConfig.BirthdateFilePath, baseConfig.BirthdateMappingPath, personID, personID, dob)
		if err != nil {
			log.Warn("SetBirthdate failed", "personId", personID, "err", err)
			t := i18n.T()
			return c.JSON(http.StatusInternalServerError, SetBirthdateResponse{
				Status:  "error",
				Message: t("birthdate_save_failed"),
			})
		}
		t := i18n.T()
		return c.JSON(http.StatusOK, SetBirthdateResponse{
			Status:  "ok",
			Message: t("birthdate_saved"),
		})
	}
}

// DeleteBirthdate handles DELETE /api/person-birthdate to remove a person's birthdate.
func DeleteBirthdate(baseConfig *config.Config) echo.HandlerFunc {
	return func(c *echo.Context) error {
		if baseConfig.Source != config.SourcePhotoPrism || strings.TrimSpace(baseConfig.BirthdateFilePath) == "" || strings.TrimSpace(baseConfig.BirthdateMappingPath) == "" {
			t := i18n.T()
			return c.JSON(http.StatusBadRequest, DeleteBirthdateResponse{
				Status:  "error",
				Message: t("birthdate_not_configured"),
			})
		}

		var req DeleteBirthdateRequest
		if err := c.Bind(&req); err != nil {
			t := i18n.T()
			return c.JSON(http.StatusBadRequest, DeleteBirthdateResponse{
				Status:  "error",
				Message: t("birthdate_invalid_request"),
			})
		}
		personID := strings.TrimSpace(req.PersonID)
		if personID == "" {
			t := i18n.T()
			return c.JSON(http.StatusBadRequest, DeleteBirthdateResponse{
				Status:  "error",
				Message: t("birthdate_person_required"),
			})
		}

		err := birthdate.RemoveBirthdate(baseConfig.BirthdateFilePath, baseConfig.BirthdateMappingPath, personID, personID)
		if err != nil {
			log.Warn("DeleteBirthdate failed", "personId", personID, "err", err)
			t := i18n.T()
			return c.JSON(http.StatusInternalServerError, DeleteBirthdateResponse{
				Status:  "error",
				Message: t("birthdate_remove_failed"),
			})
		}
		t := i18n.T()
		return c.JSON(http.StatusOK, DeleteBirthdateResponse{
			Status:  "ok",
			Message: t("birthdate_removed"),
		})
	}
}
