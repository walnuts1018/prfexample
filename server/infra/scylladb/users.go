package scylladb

import (
	"context"
	"fmt"
	"time"

	"github.com/scylladb/gocql "
	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/walnuts1018/PRFExample/server/domain/entity"
)

const (
	temporaryUserTTL = 1 * time.Hour
)

func (db *ScyllaDB) CreateTemporaryUser(ctx context.Context, user entity.User) error {
	q := qb.Insert(usersTable.Name()).
		Columns(usersMeta.Columns...).
		TTL(temporaryUserTTL).
		Unique().
		QueryContext(ctx, db.sess).
		BindStruct(UserFromEntity(user))
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (db *ScyllaDB) PromoteTemporaryUser(ctx context.Context, userID entity.UserID) error {
	q := qb.Update(usersTable.Name()).
		Set("is_temporary").
		Where(qb.Eq("id")).
		TTL(0).
		QueryContext(ctx, db.sess).
		BindMap(qb.M{"id": gocql.UUID(userID), "is_temporary": false})
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to promote temporary user: %w", err)
	}
	return nil
}

func (db *ScyllaDB) GetUserByID(ctx context.Context, id entity.UserID) (entity.User, error) {
	var user User
	q := db.sess.Query(usersTable.Select()).BindMap(qb.M{"id": gocql.UUID(id)})
	if err := q.GetRelease(&user); err != nil {
		return entity.User{}, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return user.ToEntity(), nil
}
