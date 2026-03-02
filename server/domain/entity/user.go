package entity

import (
	"fmt"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

type UserID uuid.UUID

func NewUserID() (UserID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return UserID(uuid.Nil), fmt.Errorf("failed to generate UserID: %w", err)
	}
	return UserID(id), nil
}

func ParseUserID(s string) (UserID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return UserID(uuid.Nil), fmt.Errorf("failed to parse UserID: %w", err)
	}
	return UserID(id), nil
}

func (id UserID) String() string {
	return uuid.UUID(id).String()
}

func (id UserID) MarshalText() ([]byte, error) {
	return []byte(uuid.UUID(id).String()), nil
}

func (id *UserID) UnmarshalText(data []byte) error {
	u, err := uuid.Parse(string(data))
	if err != nil {
		return err
	}
	*id = UserID(u)

	return nil
}

type User struct {
	ID          UserID
	PRFSalt     PRFSalt
	IsTemporary bool
}

type WebAuthnUser struct {
	User
	Credentials []webauthn.Credential
}

func (u *WebAuthnUser) WebAuthnID() []byte {
	return u.ID[:]
}

func (u *WebAuthnUser) WebAuthnName() string {
	return u.ID.String()
}

func (u *WebAuthnUser) WebAuthnDisplayName() string {
	return u.ID.String()
}

func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

var _ webauthn.User = new(WebAuthnUser)
