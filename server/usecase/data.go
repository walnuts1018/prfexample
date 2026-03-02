package usecase

import (
	"context"
	"fmt"
	"iter"

	"github.com/walnuts1018/PRFExample/server/domain/entity"
)

func (u *Usecase) ListEncryptedData(ctx context.Context, userID entity.UserID) (iter.Seq2[entity.EncryptedDataID, entity.EncryptedData], error) {
	return u.encryptedDataRepository.ListEncryptedData(ctx, userID)
}

func (u *Usecase) GetEncryptedData(ctx context.Context, userID entity.UserID, id entity.EncryptedDataID) (entity.EncryptedData, error) {
	return u.encryptedDataRepository.GetEncryptedData(ctx, userID, id)
}

func (u *Usecase) SaveEncryptedData(ctx context.Context, user entity.User, data []byte, iv []byte) (entity.EncryptedData, error) {
	id, err := entity.NewEncryptedDataID()
	if err != nil {
		return entity.EncryptedData{}, fmt.Errorf("failed to generate EncryptedDataID: %w", err)
	}

	e := entity.EncryptedData{
		ID:        id,
		UserID:    user.ID,
		Data:      data,
		IV:        iv,
		UpdatedAt: u.clock.Now(),
	}

	if err := u.encryptedDataRepository.StoreEncryptedData(ctx, e); err != nil {
		return entity.EncryptedData{}, fmt.Errorf("failed to store encrypted data: %w", err)
	}

	return e, nil
}
