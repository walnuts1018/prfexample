package scylladb

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/walnuts1018/PRFExample/server/domain/model"
)

func (db *ScyllaDB) SaveSession(ctx context.Context, sessionID model.SessionID, key string, data io.Reader, ttl time.Duration) error {
	b, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("failed to read session data: %w", err)
	}
	q := qb.Insert(sessionsTable.Name()).
		Columns(sessionsMeta.Columns...).
		TTL(ttl).
		QueryContext(ctx, db.sess).
		BindMap(qb.M{
			"session_id": sessionID,
			"key":        key,
			"data":       b,
		})
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}
	return nil
}

func (db *ScyllaDB) GetSession(ctx context.Context, sessionID model.SessionID, key string) (io.Reader, error) {
	var sess Session
	q := db.sess.Query(sessionsTable.Select()).BindMap(qb.M{"session_id": sessionID, "key": key})
	if err := q.GetRelease(&sess); err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return bytes.NewReader(sess.Data), nil
}

func (db *ScyllaDB) DeleteSession(ctx context.Context, sessionID model.SessionID, key string) error {
	q := db.sess.Query(sessionsTable.Delete()).BindMap(qb.M{"session_id": sessionID, "key": key})
	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
