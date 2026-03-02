package usecase

import (
	"context"
	"fmt"

	"github.com/walnuts1018/PRFExample/server/domain/entity"
)

func (u *Usecase) createTemporaryUser(ctx context.Context) (entity.User, error) {
	userID, err := entity.NewUserID()
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to generate UserID: %w", err)
	}

	salt, err := entity.NewPRFSalt(u.random)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to generate PRF salt: %w", err)
	}

	user := entity.User{
		ID:          userID,
		PRFSalt:     salt,
		IsTemporary: true,
	}
	if err := u.userRepository.CreateTemporaryUser(ctx, user); err != nil {
		return entity.User{}, fmt.Errorf("failed to create temporary user: %w", err)
	}
	return user, nil
}

func (u *Usecase) getUser(ctx context.Context, id entity.UserID, allowTemporary ...bool) (entity.User, error) {
	user, err := u.userRepository.GetUserByID(ctx, id)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to get user by ID: %w", err)
	}
	if user.IsTemporary && (len(allowTemporary) == 0 || !allowTemporary[0]) {
		return entity.User{}, fmt.Errorf("user is temporary")
	}

	return user, nil
}
