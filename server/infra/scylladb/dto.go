package scylladb

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/scylladb/gocql "
	"github.com/walnuts1018/PRFExample/server/domain/entity"
)

type User struct {
	Id          gocql.UUID
	PrfSalt     string
	IsTemporary bool
}

func (u User) ToEntity() entity.User {
	return entity.User{
		ID:          entity.UserID(u.Id),
		PRFSalt:     u.PrfSalt,
		IsTemporary: u.IsTemporary,
	}
}

func UserFromEntity(entity entity.User) User {
	return User{
		Id:          gocql.UUID(entity.ID),
		PrfSalt:     entity.PRFSalt,
		IsTemporary: entity.IsTemporary,
	}
}

type WebAuthnCredentialID string

func (id WebAuthnCredentialID) ToEntity() (entity.WebAuthnCredentialID, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(id))
	if err != nil {
		return entity.WebAuthnCredentialID{}, err
	}
	return entity.WebAuthnCredentialID(decoded), nil
}

func WebAuthnCredentialIDFromEntity(id entity.WebAuthnCredentialID) WebAuthnCredentialID {
	return WebAuthnCredentialID(base64.StdEncoding.EncodeToString([]byte(id)))
}

type WebAuthnCredential struct {
	Id         WebAuthnCredentialID
	UserID     gocql.UUID
	Credential []byte
	CreatedAt  time.Time
}

func (c WebAuthnCredential) ToEntity() (entity.WebAuthnCredential, error) {
	id, err := c.Id.ToEntity()
	if err != nil {
		return entity.WebAuthnCredential{}, err
	}

	var cred webauthn.Credential
	if err := json.Unmarshal(c.Credential, &cred); err != nil {
		return entity.WebAuthnCredential{}, err
	}

	return entity.WebAuthnCredential{
		ID:         id,
		UserID:     entity.UserID(c.UserID),
		Credential: &cred,
		CreatedAt:  synchro.In[tz.AsiaTokyo](c.CreatedAt),
	}, nil
}

func WebAuthnCredentialFromEntity(entity entity.WebAuthnCredential) (WebAuthnCredential, error) {
	cred, err := json.Marshal(entity.Credential)
	if err != nil {
		return WebAuthnCredential{}, err
	}

	return WebAuthnCredential{
		Id:         WebAuthnCredentialIDFromEntity(entity.ID),
		UserID:     gocql.UUID(entity.UserID),
		Credential: cred,
		CreatedAt:  entity.CreatedAt.StdTime(),
	}, nil
}

type EncryptedData struct {
	Id        gocql.UUID
	UserId    gocql.UUID
	Data      []byte
	Iv        []byte
	UpdatedAt time.Time
}

func (e EncryptedData) ToEntity() entity.EncryptedData {
	return entity.EncryptedData{
		ID:        entity.EncryptedDataID(e.Id),
		UserID:    entity.UserID(e.UserId),
		Data:      e.Data,
		IV:        e.Iv,
		UpdatedAt: synchro.In[tz.AsiaTokyo](e.UpdatedAt),
	}
}

func EncryptedDataFromEntity(entity entity.EncryptedData) EncryptedData {
	return EncryptedData{
		Id:        gocql.UUID(entity.ID),
		UserId:    gocql.UUID(entity.UserID),
		Data:      entity.Data,
		Iv:        entity.IV,
		UpdatedAt: entity.UpdatedAt.StdTime(),
	}
}

type Session struct {
	SessionId string
	Key       string
	Data      []byte
}
