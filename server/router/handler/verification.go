package handler

import (
	"errors"
	"log/slog"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/labstack/echo/v5"
	"github.com/walnuts1018/PRFExample/server/domain/entity"
	"github.com/walnuts1018/PRFExample/server/usecase"
)

func (h *Handler) GetWebAuthnCredentialAssertion(c *echo.Context) error {
	ctx := c.Request().Context()
	sessionID := mustGetSessionID(c)

	var params struct {
		UserID string `query:"user_id"`
	}

	if err := c.Bind(&params); err != nil {
		slog.WarnContext(ctx, "failed to bind query parameters", slog.Any("error", err))
		return c.JSON(400, map[string]string{"error": "invalid request"})
	}

	userID, err := entity.ParseUserID(params.UserID)
	if err != nil {
		slog.WarnContext(ctx, "failed to parse user ID", slog.Any("error", err))
		return c.JSON(400, map[string]string{"error": "invalid user ID"})
	}

	assertion, session, err := h.u.BeginWebAuthnLogin(ctx, userID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to begin webauthn login", slog.Any("error", err))
		return errors.New("failed to begin webauthn login")
	}

	if err := h.u.SaveWebAuthnSession(ctx, sessionID, usecase.WebauthnCredentialAssertionSessionKey, session); err != nil {
		slog.ErrorContext(ctx, "failed to save webauthn credential assertion session", slog.Any("error", err))
		return errors.New("failed to save webauthn credential assertion session")
	}

	if err := h.u.SaveUserIDInVerificationSession(ctx, sessionID, userID); err != nil {
		slog.ErrorContext(ctx, "failed to save user ID in verification session", slog.Any("error", err))
		return errors.New("failed to save user ID in verification session")
	}

	return c.JSON(200, assertion)
}

func (h *Handler) VerifyWebAuthnCredentialAssertion(c *echo.Context) error {
	ctx := c.Request().Context()
	sessionID := mustGetSessionID(c)

	session, err := h.u.GetWebAuthnSession(ctx, sessionID, usecase.WebauthnCredentialAssertionSessionKey)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get webauthn credential assertion session", slog.Any("error", err))
		return errors.New("failed to get webauthn credential assertion session")
	}

	userID, err := h.u.GetUserIDInVerificationSession(ctx, sessionID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get user ID in verification session", slog.Any("error", err))
		return errors.New("failed to get user ID in verification session")
	}

	assertion, err := protocol.ParseCredentialRequestResponseBody(c.Request().Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse credential assertion", slog.Any("error", err))
		return errors.New("failed to parse credential assertion")
	}

	user, wc, err := h.u.FinishWebAuthnLogin(ctx, userID, *session, assertion)
	if err != nil {
		slog.ErrorContext(ctx, "failed to finish webauthn login", slog.Any("error", err))
		return errors.New("failed to finish webauthn login")
	}

	if err := h.u.SaveLoginUserIDSession(ctx, sessionID, user.ID); err != nil {
		slog.ErrorContext(ctx, "failed to save login user ID session", slog.Any("error", err))
		return errors.New("failed to save login user ID session")
	}

	type response struct {
		UserID     string                    `json:"user_id"`
		Credential entity.WebAuthnCredential `json:"credential"`
	}

	return c.JSON(200, response{
		UserID:     user.ID.String(),
		Credential: wc,
	})
}
