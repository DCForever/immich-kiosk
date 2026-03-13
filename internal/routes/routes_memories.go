package routes

import (
	"context"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/labstack/echo/v5"

	"github.com/damongolding/immich-kiosk/internal/common"
	"github.com/damongolding/immich-kiosk/internal/config"
	"github.com/damongolding/immich-kiosk/internal/i18n"
	"github.com/damongolding/immich-kiosk/internal/kiosk"
	memoriesComponent "github.com/damongolding/immich-kiosk/internal/templates/components/memories"
)

// Memories returns an echo.HandlerFunc for the on-demand memories view.
// GET/POST /memories renders the collage or "No memories for this day" / error state.
// Includes back button to return to slideshow via POST /asset/new.
func Memories(baseConfig *config.Config, com *common.Common) echo.HandlerFunc {
	return func(c *echo.Context) error {
		requestData, err := InitializeRequestData(c, baseConfig)
		if err != nil {
			return err
		}
		if requestData == nil {
			return nil
		}

		requestConfig := requestData.RequestConfig
		requestID := requestData.RequestID
		deviceID := requestData.DeviceID

		if !requestConfig.Memories {
			t := i18n.T()
			return RenderMessage(c, t("memories_error_load"), t("memories_no_photos"))
		}

		requestCtx := common.CopyContext(c)
		provider := getProvider(context.Background(), requestConfig)
		queries := requestCtx.URL.Query()

		viewData, err := ProcessMemoriesCollage(provider, requestConfig, requestID, deviceID, queries)
		if err != nil {
			log.Debug(requestID, "Memories on-demand failed", "error", err)
			return RenderMemoriesEmpty(c, requestConfig, queries, com.Secret())
		}

		viewData.ShowMemoriesBackButton = true
		return RenderMemoriesView(c, viewData, com.Secret())
	}
}

// RenderMemoriesView renders the memories collage with back button.
func RenderMemoriesView(c *echo.Context, viewData common.ViewData, secret string) error {
	return Render(c, http.StatusOK, memoriesComponent.MemoriesView(viewData, secret))
}

// RenderMemoriesEmpty renders "No memories for this day" with back button for 45s.
func RenderMemoriesEmpty(c *echo.Context, requestConfig config.Config, queries map[string][]string, secret string) error {
	t := i18n.T()
	requestConfig.Layout = kiosk.LayoutCollage
	requestConfig.Duration = memoriesDurationSeconds
	viewData := common.ViewData{
		RequestID:              c.Response().Header().Get(echo.HeaderXRequestID),
		DeviceID:               c.Request().Header.Get("kiosk-device-id"),
		Config:                 requestConfig,
		MemoryCaption:          t("memories_no_photos"),
		Queries:                queries,
		ShowMemoriesBackButton: true,
	}
	return Render(c, http.StatusOK, memoriesComponent.MemoriesEmpty(viewData, secret))
}
