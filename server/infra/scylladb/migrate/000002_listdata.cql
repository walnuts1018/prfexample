DROP TABLE IF EXISTS encrypted_data;
CREATE TABLE IF NOT EXISTS encrypted_data (
   id UUID,
   user_id UUID,
   data BLOB,
   iv BLOB,
   updated_at TIMESTAMP,
   PRIMARY KEY (user_id, id)
) WITH CLUSTERING ORDER BY (id ASC);
