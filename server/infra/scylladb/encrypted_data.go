package scylladb

import (
	"context"
	"fmt"
	"iter"

	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/walnuts1018/PRFExample/server/domain/entity"
)

func (db *ScyllaDB) StoreEncryptedData(ctx context.Context, encryptedData entity.EncryptedData) error {
	q := qb.Insert(encryptedDataTable.Name()).
		Columns(encryptedDataMeta.Columns...).
		Unique().
		QueryContext(ctx, db.sess).
		BindStruct(EncryptedDataFromEntity(encryptedData))
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to store encrypted data: %w", err)
	}
	return nil
}

func (db *ScyllaDB) UpdateEncryptedData(ctx context.Context, encryptedData entity.EncryptedData) error {
	q := db.sess.Query(encryptedDataTable.Update("data", "iv", "updated_at")).BindStruct(EncryptedDataFromEntity(encryptedData))
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to update encrypted data: %w", err)
	}
	return nil
}

func (db *ScyllaDB) GetEncryptedData(ctx context.Context, userID entity.UserID, id entity.EncryptedDataID) (entity.EncryptedData, error) {
	var encryptedData EncryptedData
	q := db.sess.Query(encryptedDataTable.Select()).BindMap(qb.M{"id": uuid.UUID(id), "user_id": uuid.UUID(userID)})
	if err := q.GetRelease(&encryptedData); err != nil {
		return entity.EncryptedData{}, fmt.Errorf("failed to get encrypted data: %w", err)
	}
	return encryptedData.ToEntity(), nil
}

func (db *ScyllaDB) ListEncryptedData(ctx context.Context, userID entity.UserID) (iter.Seq2[entity.EncryptedDataID, entity.EncryptedData], error) {
	var encryptedDataList []EncryptedData
	q := db.sess.Query(encryptedDataTable.Select()).BindMap(qb.M{"user_id": uuid.UUID(userID)})
	if err := q.SelectRelease(&encryptedDataList); err != nil {
		return nil, fmt.Errorf("failed to list encrypted data: %w", err)
	}

	return func(yield func(entity.EncryptedDataID, entity.EncryptedData) bool) {
		for _, encryptedData := range encryptedDataList {
			if !yield(entity.EncryptedDataID(encryptedData.Id), encryptedData.ToEntity()) {
				return
			}
		}
	}, nil
}

func (db *ScyllaDB) DeleteEncryptedData(ctx context.Context, userID entity.UserID, id entity.EncryptedDataID) error {
	q := db.sess.Query(encryptedDataTable.Delete()).BindMap(qb.M{"id": uuid.UUID(id), "user_id": uuid.UUID(userID)})
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to delete encrypted data: %w", err)
	}
	return nil
}
