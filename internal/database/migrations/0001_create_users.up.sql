CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

  username TEXT NOT NULL UNIQUE,
  secret BYTEA NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);

-- Migrations run as gojo_migrations, but the app connects as gojo_api.
-- Without this grant, gojo_api gets "permission denied" on a freshly
-- created table even though the table exists.
GRANT SELECT, INSERT, UPDATE, DELETE ON users TO gojo_api;
