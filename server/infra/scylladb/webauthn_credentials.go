package scylladb

import (
	"context"
	"fmt"
	"iter"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/walnuts1018/PRFExample/server/domain/entity"
)

func (db *ScyllaDB) StoreWebAuthnCredential(ctx context.Context, webAuthnCredential entity.WebAuthnCredential) error {
	e, err := WebAuthnCredentialFromEntity(webAuthnCredential)
	if err != nil {
		return fmt.Errorf("failed to convert WebAuthnCredential entity to DTO: %w", err)
	}

	q := db.sess.Query(webAuthnCredentialsTable.Insert()).BindStruct(e)
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to store WebAuthnCredential: %w", err)
	}
	return nil
}

func (db *ScyllaDB) ListWebAuthnCredentialsByUserID(ctx context.Context, userID entity.UserID) (iter.Seq[entity.WebAuthnCredential], error) {
	var webAuthnCredentialList []WebAuthnCredential
	q := db.sess.Query(webAuthnCredentialsTable.Select()).BindMap(qb.M{"user_id": gocql.UUID(userID)})
	if err := q.SelectRelease(&webAuthnCredentialList); err != nil {
		return nil, fmt.Errorf("failed to list WebAuthnCredentials: %w", err)
	}

	return func(yield func(entity.WebAuthnCredential) bool) {
		for _, webAuthnCredential := range webAuthnCredentialList {
			cred, err := webAuthnCredential.ToEntity()
			if err != nil {
				continue
			}
			if !yield(cred) {
				return
			}
		}
	}, nil
}

func (db *ScyllaDB) GetWebAuthnCredential(ctx context.Context, id entity.WebAuthnCredentialID) (entity.WebAuthnCredential, error) {
	var webAuthnCredential WebAuthnCredential
	q := db.sess.Query(webAuthnCredentialsTable.Select()).BindMap(qb.M{"id": WebAuthnCredentialIDFromEntity(id)})
	if err := q.GetRelease(&webAuthnCredential); err != nil {
		return entity.WebAuthnCredential{}, fmt.Errorf("failed to get WebAuthnCredential by ID: %w", err)
	}
	return webAuthnCredential.ToEntity()
}

func (db *ScyllaDB) UpdateWebAuthnCredentialOnLogin(ctx context.Context, id entity.WebAuthnCredentialID, cred *webauthn.Credential) error {
	q := qb.Update(webAuthnCredentialsTable.Name()).
		Set("credential").
		Where(qb.Eq("id")).
		QueryContext(ctx, db.sess).
		BindMap(qb.M{"id": WebAuthnCredentialIDFromEntity(id), "credential": cred})
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to update WebAuthnCredential on login: %w", err)
	}
	return nil
}

func (db *ScyllaDB) DeleteWebAuthnCredential(ctx context.Context, userID entity.UserID, id entity.WebAuthnCredentialID) error {
	q := db.sess.Query(webAuthnCredentialsTable.Delete()).BindMap(qb.M{"id": WebAuthnCredentialIDFromEntity(id), "user_id": gocql.UUID(userID)})
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to delete WebAuthnCredential: %w", err)
	}
	return nil
}
