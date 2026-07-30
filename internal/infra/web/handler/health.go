package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

const RouteHealth string = "/health"

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// @Summary      Health check
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "ok"
// @Router       /health [get]
func (hh *HealthHandler) Handle(c *echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
