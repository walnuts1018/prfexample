package usecase

import (
	"context"
	"io"
	"iter"
	"time"

	"github.com/Code-Hex/synchro/tz"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/walnuts1018/PRFExample/server/domain/entity"
	"github.com/walnuts1018/PRFExample/server/domain/model"
	"github.com/walnuts1018/PRFExample/server/util/clock"
	"github.com/walnuts1018/PRFExample/server/util/random"
)

type UserRepository interface {
	CreateTemporaryUser(ctx context.Context, user entity.User) error
	PromoteTemporaryUser(ctx context.Context, userID entity.UserID) error
	GetUserByID(ctx context.Context, id entity.UserID) (entity.User, error)
}

type EncryptedDataRepository interface {
	StoreEncryptedData(ctx context.Context, encryptedData entity.EncryptedData) error
	UpdateEncryptedData(ctx context.Context, encryptedData entity.EncryptedData) error
	GetEncryptedData(ctx context.Context, userID entity.UserID, id entity.EncryptedDataID) (entity.EncryptedData, error)
	ListEncryptedData(ctx context.Context, userID entity.UserID) (iter.Seq2[entity.EncryptedDataID, entity.EncryptedData], error)
	DeleteEncryptedData(ctx context.Context, userID entity.UserID, id entity.EncryptedDataID) error
}

type WebAuthnCredentialRepository interface {
	StoreWebAuthnCredential(ctx context.Context, webAuthnCredential entity.WebAuthnCredential) error
	ListWebAuthnCredentialsByUserID(ctx context.Context, userID entity.UserID) (iter.Seq[entity.WebAuthnCredential], error)
	GetWebAuthnCredential(ctx context.Context, id entity.WebAuthnCredentialID, userID entity.UserID) (entity.WebAuthnCredential, error)
	UpdateWebAuthnCredentialOnLogin(ctx context.Context, id entity.WebAuthnCredentialID, userID entity.UserID, cred *webauthn.Credential) error
	DeleteWebAuthnCredential(ctx context.Context, userID entity.UserID, id entity.WebAuthnCredentialID) error
}

type SessionRepository interface {
	SaveSession(ctx context.Context, sessionID model.SessionID, key string, data io.Reader, ttl time.Duration) error
	GetSession(ctx context.Context, sessionID model.SessionID, key string) (io.Reader, error)
	DeleteSession(ctx context.Context, sessionID model.SessionID, key string) error
}

type Usecase struct {
	encryptedDataRepository      EncryptedDataRepository
	webAuthnCredentialRepository WebAuthnCredentialRepository
	userRepository               UserRepository
	sessionRepository            SessionRepository

	// Service
	webAuthn *webauthn.WebAuthn

	random random.Random
	clock  clock.Clock[tz.AsiaTokyo]
}

func NewUsecase(
	encryptedDataRepository EncryptedDataRepository,
	webAuthnCredentialRepository WebAuthnCredentialRepository,
	userRepository UserRepository,
	sessionRepository SessionRepository,

	webAuthn *webauthn.WebAuthn,

	random random.Random,
	clock clock.Clock[tz.AsiaTokyo],
) *Usecase {
	return &Usecase{
		encryptedDataRepository:      encryptedDataRepository,
		webAuthnCredentialRepository: webAuthnCredentialRepository,
		userRepository:               userRepository,
		sessionRepository:            sessionRepository,

		webAuthn: webAuthn,

		random: random,
		clock:  clock,
	}
}
