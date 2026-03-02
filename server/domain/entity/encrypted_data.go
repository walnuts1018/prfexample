package entity

import (
	"fmt"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
	"github.com/google/uuid"
)

type EncryptedDataID uuid.UUID

func (id EncryptedDataID) String() string {
	return uuid.UUID(id).String()
}

func ParseEncryptedDataID(s string) (EncryptedDataID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return EncryptedDataID(uuid.Nil), fmt.Errorf("failed to parse EncryptedDataID: %w", err)
	}
	return EncryptedDataID(id), nil
}

func NewEncryptedDataID() (EncryptedDataID, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return EncryptedDataID(uuid.Nil), fmt.Errorf("failed to generate EncryptedDataID: %w", err)
	}
	return EncryptedDataID(id), nil
}

type EncryptedData struct {
	ID        EncryptedDataID
	UserID    UserID
	Data      []byte
	IV        []byte
	UpdatedAt synchro.Time[tz.AsiaTokyo]
}
