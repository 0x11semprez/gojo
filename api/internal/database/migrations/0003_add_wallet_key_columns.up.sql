ALTER TABLE wallet
  ADD COLUMN private_view_key BYTEA,
  ADD COLUMN public_view_key BYTEA,
  ADD COLUMN password_salt BYTEA NOT NULL;

-- private_view_key/public_view_key are only populated for monero
-- wallets (bitcoin has a single keypair); password_salt is always
-- set, it seeds the password-derived key used as the inner layer of
-- the wallet's envelope encryption.
