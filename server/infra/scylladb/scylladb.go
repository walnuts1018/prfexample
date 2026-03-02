package scylladb

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/migrate"
	"github.com/walnuts1018/PRFExample/server/config"
	migratefs "github.com/walnuts1018/PRFExample/server/infra/scylladb/migrate"
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
	cluster := gocql.NewCluster(cfg.Endpoint)
	cluster.Keyspace = cfg.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.User,
		Password: cfg.Password,
	}
	if cfg.CACertPath != "" && cfg.ClientCertPath != "" && cfg.ClientKeyPath != "" {
		cluster.SslOpts = &gocql.SslOptions{
			CaPath:                 cfg.CACertPath,
			CertPath:               cfg.ClientCertPath,
			KeyPath:                cfg.ClientKeyPath,
			EnableHostVerification: false, // TODO: Enable host verification in production
		}
	}

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

func (db *ScyllaDB) Migrate(ctx context.Context) error {
	if err := migrate.FromFS(ctx, db.sess, migratefs.Scripts); err != nil {
		return fmt.Errorf("failed to migrate ScyllaDB: %w", err)
	}
	return nil
}
