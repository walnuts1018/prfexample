package handler

import (
	"github.com/labstack/echo/v5"
	"github.com/walnuts1018/PRFExample/server/domain/model"
)

const (
	sessionIDHeaderKey  = "X-Session-ID"
	sessionIDContextKey = "session_id"
)

func (h Handler) SessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		sessionID := model.SessionID(c.Request().Header.Get(sessionIDHeaderKey))
		if sessionID == "" {
			var err error
			sessionID, err = model.NewSessionID(h.rand)
			if err != nil {
				return c.JSON(500, map[string]string{"error": "failed to generate session ID"})
			}
			c.Response().Header().Set(sessionIDHeaderKey, string(sessionID))
		}
		c.Set(sessionIDContextKey, sessionID)
		return next(c)
	}
}

func mustGetSessionID(c *echo.Context) model.SessionID {
	sessionID, ok := c.Get(sessionIDContextKey).(model.SessionID)
	if !ok {
		panic("session ID not found in context")
	}
	return sessionID
}
