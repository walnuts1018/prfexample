package handler

import (
	"maps"

	"github.com/labstack/echo/v5"
	"github.com/walnuts1018/PRFExample/server/domain/entity"
)

func (h Handler) ListEncryptedData(c *echo.Context) error {
	ctx := c.Request().Context()
	sessionID := mustGetSessionID(c)

	userID, err := h.u.GetLoginUserIDSession(ctx, sessionID)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "unauthorized"})
	}

	data, err := h.u.ListEncryptedData(ctx, userID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "failed to list encrypted data"})
	}

	return c.JSON(200, maps.Collect(data))
}

func (h Handler) GetEncryptedData(c *echo.Context) error {
	ctx := c.Request().Context()
	sessionID := mustGetSessionID(c)

	var params struct {
		ID string `path:"id"`
	}
	if err := c.Bind(&params); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}

	userID, err := h.u.GetLoginUserIDSession(ctx, sessionID)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "unauthorized"})
	}

	id, err := entity.ParseEncryptedDataID(params.ID)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "invalid data ID"})
	}

	data, err := h.u.GetEncryptedData(ctx, userID, id)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "failed to get encrypted data"})
	}

	return c.JSON(200, data)
}

func (h Handler) SaveEncryptedData(c *echo.Context) error {
	ctx := c.Request().Context()
	sessionID := mustGetSessionID(c)

	userID, err := h.u.GetLoginUserIDSession(ctx, sessionID)
	if err != nil {
		return c.JSON(401, map[string]string{"error": "unauthorized"})
	}

	var params struct {
		Data []byte `json:"data"`
		IV   []byte `json:"iv"`
	}
	if err := c.Bind(&params); err != nil {
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}

	user := entity.User{ID: userID}
	data, err := h.u.SaveEncryptedData(ctx, user, params.Data, params.IV)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "failed to save encrypted data"})
	}

	return c.JSON(200, data)
}
