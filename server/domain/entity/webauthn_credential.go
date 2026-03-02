package entity

import (
	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebAuthnCredentialID []byte

type WebAuthnCredential struct {
	ID         WebAuthnCredentialID
	UserID     UserID
	Credential *webauthn.Credential
	CreatedAt  synchro.Time[tz.AsiaTokyo]
}

func NewCredential(user User, cred *webauthn.Credential) (WebAuthnCredential, error) {
	return WebAuthnCredential{
		ID:         cred.ID,
		UserID:     user.ID,
		Credential: cred,
		CreatedAt:  synchro.Now[tz.AsiaTokyo](),
	}, nil
}
