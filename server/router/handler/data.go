package handler

import (
	"errors"
	"log/slog"
	"maps"

	"github.com/labstack/echo/v5"
	"github.com/walnuts1018/PRFExample/server/domain/entity"
)

func (h Handler) ListEncryptedData(c *echo.Context) error {
	ctx := c.Request().Context()
	sessionID := mustGetSessionID(c)

	userID, err := h.u.GetLoginUserIDSession(ctx, sessionID)
	if err != nil {
		slog.WarnContext(ctx, "unauthorized access to list encrypted data", "session_id", sessionID, "error", err)
		return c.JSON(401, map[string]string{"error": "unauthorized"})
	}

	data, err := h.u.ListEncryptedData(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list encrypted data", "error", err)
		return errors.New("failed to list encrypted data")
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
		slog.WarnContext(ctx, "unauthorized access to get encrypted data", "session_id", sessionID, "data_id", params.ID, "error", err)
		return c.JSON(401, map[string]string{"error": "unauthorized"})
	}

	id, err := entity.ParseEncryptedDataID(params.ID)
	if err != nil {
		slog.WarnContext(ctx, "invalid encrypted data ID format", "data_id", params.ID, "error", err)
		return c.JSON(400, map[string]string{"error": "invalid data ID"})
	}

	data, err := h.u.GetEncryptedData(ctx, userID, id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get encrypted data", "data_id", params.ID, "error", err)
		return errors.New("failed to get encrypted data")
	}

	return c.JSON(200, data)
}

func (h Handler) SaveEncryptedData(c *echo.Context) error {
	ctx := c.Request().Context()
	sessionID := mustGetSessionID(c)

	userID, err := h.u.GetLoginUserIDSession(ctx, sessionID)
	if err != nil {
		slog.WarnContext(ctx, "unauthorized access to save encrypted data", "session_id", sessionID, "error", err)
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
		slog.ErrorContext(ctx, "failed to save encrypted data", "user_id", userID, "error", err)
		return errors.New("failed to save encrypted data")
	}

	return c.JSON(200, data)
}
