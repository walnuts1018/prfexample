package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/labstack/echo/v5"
	"github.com/walnuts1018/PRFExample/server/domain/entity"
	"github.com/walnuts1018/PRFExample/server/usecase"
)

func (h *Handler) GetWebAuthnCredentialCreation(c *echo.Context) error {
	ctx := c.Request().Context()
	sessionID := mustGetSessionID(c)

	userID, creation, sess, err := h.u.BeginWebAuthnRegistration(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to begin webauthn registration", slog.Any("error", err))
		return errors.New("failed to begin webauthn registration")
	}

	if err := h.u.SaveWebAuthnSession(ctx, sessionID, usecase.WebauthnCredentialCreationSessionKey, sess); err != nil {
		slog.ErrorContext(ctx, "failed to save webauthn session", slog.Any("error", err))
		return errors.New("failed to save webauthn session")
	}

	if err := h.u.SaveUserIDInRegistrationSession(ctx, sessionID, userID); err != nil {
		slog.ErrorContext(ctx, "failed to save user ID in registration session", slog.Any("error", err))
		return errors.New("failed to save user ID in registration session")
	}

	return c.JSON(200, creation)
}

func (h *Handler) CreateWebAuthnCredential(c *echo.Context) error {
	ctx := c.Request().Context()
	sessionID := mustGetSessionID(c)

	creationData, err := protocol.ParseCredentialCreationResponseBody(c.Request().Body)
	if err != nil {
		slog.WarnContext(ctx, "failed to parse credential creation response body", slog.Any("error", err))
		return c.JSON(400, map[string]string{"error": "invalid request body"})
	}

	session, err := h.u.GetWebAuthnSession(ctx, sessionID, usecase.WebauthnCredentialCreationSessionKey)
	if err != nil {
		slog.WarnContext(ctx, "failed to get webauthn session", slog.Any("error", err))
		return c.JSON(400, map[string]string{"error": "invalid session"})
	}

	userID, err := h.u.GetUserIDInRegistrationSession(ctx, sessionID)
	if err != nil {
		slog.WarnContext(ctx, "failed to get user ID in registration session", slog.Any("error", err))
		return c.JSON(400, map[string]string{"error": "invalid session"})
	}

	user, wc, err := h.u.FinishWebAuthnRegistration(context.WithoutCancel(ctx), userID, creationData, session)
	if err != nil {
		slog.ErrorContext(ctx, "failed to finish webauthn registration", slog.Any("error", err))
		return errors.New("failed to finish webauthn registration")
	}

	if err := h.u.SaveLoginUserIDSession(ctx, sessionID, user.ID); err != nil {
		slog.ErrorContext(ctx, "failed to save login user ID session", slog.Any("error", err))
		return errors.New("failed to save login user ID session")
	}

	type response struct {
		UserID     entity.UserID             `json:"user_id"`
		Credential entity.WebAuthnCredential `json:"credential"`
	}

	return c.JSON(200, response{
		UserID:     user.ID,
		Credential: wc,
	})
}
