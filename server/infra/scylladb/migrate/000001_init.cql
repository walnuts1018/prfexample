CREATE TABLE IF NOT EXISTS users (
   id UUID,
   prfsalt TEXT,
   is_temporary BOOLEAN,
   PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS webauthn_credentials (
   id TEXT,
   user_id UUID,
   credential BLOB,
   created_at TIMESTAMP,
   PRIMARY KEY (user_id, id)
) WITH CLUSTERING ORDER BY (id ASC);

CREATE TABLE IF NOT EXISTS encrypted_data (
   id UUID,
   user_id UUID,
   data BLOB,
   iv BLOB,
   updated_at TIMESTAMP,
   PRIMARY KEY (user_id, id)
) WITH CLUSTERING ORDER BY (id ASC);

CREATE TABLE IF NOT EXISTS sessions (
   session_id TEXT,
   key TEXT,
   data BLOB,
   PRIMARY KEY ((session_id, key))
);

