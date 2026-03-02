CREATE TABLE users (
   id UUID,
   prfsalt TEXT,
   is_temporary BOOLEAN,
   PRIMARY KEY (id)
);

CREATE TABLE webauthn_credentials (
   id TEXT,
   user_id UUID,
   credential BLOB,
   created_at TIMESTAMP,
   PRIMARY KEY (user_id, id)
) WITH CLUSTERING ORDER BY (id ASC);

CREATE TABLE encrypted_data (
   id UUID,
   user_id UUID,
   data BLOB,
   iv BLOB,
   updated_at TIMESTAMP,
   PRIMARY KEY (id, user_id)
) WITH CLUSTERING ORDER BY (user_id ASC);


CREATE TABLE sessions (
   session_id TEXT,
   key TEXT,
   data BLOB,
   PRIMARY KEY ((session_id, key))
);

