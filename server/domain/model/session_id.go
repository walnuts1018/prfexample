package model

import (
	"fmt"

	"github.com/walnuts1018/PRFExample/server/util/random"
)

type SessionID string

func NewSessionID(rand random.Random) (SessionID, error) {
	s, err := rand.SecureString(64, random.Alphanumeric)
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}
	return SessionID(s), nil
}
