package scylladb

import (
	"context"
	"embed"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/migrate"
	"github.com/walnuts1018/PRFExample/server/config"
	"github.com/walnuts1018/PRFExample/server/usecase"
)

type ScyllaDB struct {
	sess gocqlx.Session
}

var _ usecase.UserRepository = (*ScyllaDB)(nil)
var _ usecase.WebAuthnCredentialRepository = (*ScyllaDB)(nil)
var _ usecase.EncryptedDataRepository = (*ScyllaDB)(nil)
var _ usecase.SessionRepository = (*ScyllaDB)(nil)

func NewScyllaDB(cfg config.ScyllaDB) (*ScyllaDB, error) {
	cluster := gocql.NewCluster(cfg.Host + ":" + fmt.Sprint(cfg.Port))
	cluster.Keyspace = cfg.Keyspace

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return nil, fmt.Errorf("failed to create ScyllaDB session: %w", err)
	}

	return &ScyllaDB{
		sess: session,
	}, nil
}

func (db *ScyllaDB) Close() {
	db.sess.Close()
}

//go:embed migrate/*.cql
var migrateFS embed.FS

func (db *ScyllaDB) Migrate(ctx context.Context) error {
	if err := migrate.FromFS(ctx, db.sess, migrateFS); err != nil {
		return fmt.Errorf("failed to migrate ScyllaDB: %w", err)
	}
	return nil
}
