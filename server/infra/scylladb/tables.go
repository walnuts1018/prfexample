package scylladb

import "github.com/scylladb/gocqlx/v3/table"

var usersMeta = table.Metadata{
	Name:    "users",
	Columns: []string{"id", "prf_salt", "is_temporary"},
	PartKey: []string{"id"},
}
var usersTable = table.New(usersMeta)

var webAuthnCredentialsMeta = table.Metadata{
	Name:    "webauthn_credentials",
	Columns: []string{"id", "user_id", "credential", "created_at"},
	PartKey: []string{"user_id"},
	SortKey: []string{"id"},
}
var webAuthnCredentialsTable = table.New(webAuthnCredentialsMeta)

var encryptedDataMeta = table.Metadata{
	Name:    "encrypted_data",
	Columns: []string{"id", "user_id", "data", "iv", "updated_at"},
	PartKey: []string{"id"},
	SortKey: []string{"user_id"},
}
var encryptedDataTable = table.New(encryptedDataMeta)

var sessionsMeta = table.Metadata{
	Name:    "sessions",
	Columns: []string{"session_id", "key", "data"},
	PartKey: []string{"session_id", "key"},
}
var sessionsTable = table.New(sessionsMeta)
