package entity

import (
	"fmt"

	"github.com/walnuts1018/PRFExample/server/util/random"
)

// TODO: modelにうつす
type PRFSalt = string

func NewPRFSalt(r random.Random) (PRFSalt, error) {
	salt, err := r.SecureString(32, random.Alphanumeric)
	if err != nil {
		return "", fmt.Errorf("failed to generate PRF salt: %w", err)
	}
	return salt, nil
}
