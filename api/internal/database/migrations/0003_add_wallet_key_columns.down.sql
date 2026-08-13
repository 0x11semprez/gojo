ALTER TABLE wallet
  DROP COLUMN private_view_key,
  DROP COLUMN public_view_key,
  DROP COLUMN password_salt;
