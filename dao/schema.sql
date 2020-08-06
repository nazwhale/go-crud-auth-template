CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  password VARCHAR(255),
  created_at TIMESTAMP
);

/* Lifeline example */
/* https://github.com/nazwhale/lifeline-api/blob/master/db/schema.sql */
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  created_at TIMESTAMP, /* local type is DATE */
  last_login_at TIMESTAMP, /* local type is DATE */
  login_type VARCHAR(255),
  external_login_id VARCHAR(255),
  deleted_at TIMESTAMP
);