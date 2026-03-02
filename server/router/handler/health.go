package handler

import "github.com/labstack/echo/v5"

func (h Handler) Readiness(c *echo.Context) error {
	return c.JSON(200, map[string]string{"status": "ok"})
}

func (h Handler) Liveness(c *echo.Context) error {
	return c.JSON(200, map[string]string{"status": "ok"})
}
