package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/cockroachdb/errors"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/walnuts1018/PRFExample/server/domain/entity"
	"github.com/walnuts1018/PRFExample/server/domain/model"
)

type WebAuthnSessionKey string

const (
	WebauthnCredentialCreationSessionKey  WebAuthnSessionKey = "webauthn-credential-creation-session"
	WebauthnCredentialAssertionSessionKey WebAuthnSessionKey = "webauthn-credential-assertion-session"
)

func (u *Usecase) SaveWebAuthnSession(
	ctx context.Context,
	sessionID model.SessionID,
	key WebAuthnSessionKey,
	session *webauthn.SessionData,
) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(session); err != nil {
		return fmt.Errorf("failed to encode session: %w", err)
	}

	expiresAt := synchro.In[tz.AsiaTokyo](session.Expires)
	ttl := expiresAt.Sub(u.clock.Now())
	if expiresAt.IsZero() {
		ttl = 1 * time.Hour
	}

	if err := u.sessionRepository.SaveSession(ctx, sessionID, string(key), buf, ttl); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

func (u *Usecase) GetWebAuthnSession(
	ctx context.Context,
	sessionID model.SessionID,
	key WebAuthnSessionKey,
) (*webauthn.SessionData, error) {
	sessEncoded, err := u.sessionRepository.GetSession(ctx, sessionID, string(key))
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	if sessEncoded == nil {
		return nil, errors.New("session not found")
	}

	session := &webauthn.SessionData{}
	if err := json.NewDecoder(sessEncoded).Decode(session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	if err := u.sessionRepository.DeleteSession(ctx, sessionID, string(key)); err != nil {
		return nil, fmt.Errorf("failed to delete session: %w", err)
	}

	return session, nil
}

const (
	UserIDInRegistrationSessionKey = "user-id-in-registration-session"
	UserIDInVerificationSessionKey = "user-id-in-verification-session"
	LoginUserIDSessionKey          = "login-user-id-session"
)

func (u *Usecase) SaveUserIDInRegistrationSession(ctx context.Context, sessionID model.SessionID, userID entity.UserID) error {
	if err := u.sessionRepository.SaveSession(ctx, sessionID, UserIDInRegistrationSessionKey, bytes.NewReader(userID[:]), 15*time.Minute); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

func (u *Usecase) GetUserIDInRegistrationSession(ctx context.Context, sessionID model.SessionID) (entity.UserID, error) {
	r, err := u.sessionRepository.GetSession(ctx, sessionID, UserIDInRegistrationSessionKey)
	if err != nil {
		return entity.UserID{}, fmt.Errorf("failed to get session: %w", err)
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return entity.UserID{}, fmt.Errorf("failed to read session: %w", err)
	}

	if err := u.sessionRepository.DeleteSession(ctx, sessionID, UserIDInRegistrationSessionKey); err != nil {
		return entity.UserID{}, fmt.Errorf("failed to delete session: %w", err)
	}

	return entity.UserID(buf), nil
}

func (u *Usecase) SaveUserIDInVerificationSession(ctx context.Context, sessionID model.SessionID, userID entity.UserID) error {
	if err := u.sessionRepository.SaveSession(ctx, sessionID, UserIDInVerificationSessionKey, bytes.NewReader(userID[:]), 15*time.Minute); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

func (u *Usecase) GetUserIDInVerificationSession(ctx context.Context, sessionID model.SessionID) (entity.UserID, error) {
	r, err := u.sessionRepository.GetSession(ctx, sessionID, UserIDInVerificationSessionKey)
	if err != nil {
		return entity.UserID{}, fmt.Errorf("failed to get session: %w", err)
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return entity.UserID{}, fmt.Errorf("failed to read session: %w", err)
	}

	if err := u.sessionRepository.DeleteSession(ctx, sessionID, UserIDInVerificationSessionKey); err != nil {
		return entity.UserID{}, fmt.Errorf("failed to delete session: %w", err)
	}

	return entity.UserID(buf), nil
}

func (u *Usecase) SaveLoginUserIDSession(ctx context.Context, sessionID model.SessionID, userID entity.UserID) error {
	if err := u.sessionRepository.SaveSession(ctx, sessionID, LoginUserIDSessionKey, bytes.NewReader(userID[:]), 24*time.Hour); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

func (u *Usecase) GetLoginUserIDSession(ctx context.Context, sessionID model.SessionID) (entity.UserID, error) {
	r, err := u.sessionRepository.GetSession(ctx, sessionID, LoginUserIDSessionKey)
	if err != nil {
		// TODO: Session切れを伝えるようにする
		return entity.UserID{}, fmt.Errorf("failed to get session: %w", err)
	}

	buf, err := io.ReadAll(r)
	if err != nil {
		return entity.UserID{}, fmt.Errorf("failed to read session: %w", err)
	}

	return entity.UserID(buf), nil
}

func (u *Usecase) DeleteLoginUserIDSession(ctx context.Context, sessionID model.SessionID) error {
	if err := u.sessionRepository.DeleteSession(ctx, sessionID, LoginUserIDSessionKey); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
