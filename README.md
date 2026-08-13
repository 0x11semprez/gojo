why gojo ? </br>
coldcard was hacked due to a conditional/build configuration that caused the secure RNG path not to be used as intended.</br>
</br>
<https://www.trmlabs.com/resources/blog/the-largest-hardware-wallet-exploit-of-2026-inside-the-usd-116-million-coldcard-hack> </br>
</br>
gojo is not a wallet. </br>
gojo is a custodial KMS (key management service) for bitcoin and monero keys, exposed through a golang api. </br>
gojo generates extremely secure keys by mixing several hardware-specific entropy sources in rust, instead of trusting a single RNG path. </br>
gojo is gojo </br>
</br>
gojo is designed to be used as a white-label infrastructure layer: it handles account management and key generation, while the integrator owns branding, marketing, the frontend experience, and wallet operations. </br>
</br>
architecture: </br>

- postgres: keeps the encrypted keys and the accounts </br>
- golang api: handles accounts and authentication, stores/serves keys, always encrypted </br>
- rust: generates the keys (bitcoin secp256k1, monero ed25519), collecting multiple hardware entropy sources before seeding </br>

## Quickstart

Prerequisites: Go 1.26+, Rust (2024 edition), PostgreSQL, `make`.

**1. Database (one-time).** As a postgres superuser:

```sql
CREATE ROLE gojo_api WITH LOGIN PASSWORD 'change-me';
CREATE ROLE gojo_migrations WITH LOGIN PASSWORD 'change-me';
CREATE DATABASE gojo OWNER gojo_migrations;
```

`gojo_migrations` owns the schema; `gojo_api` only gets what each migration grants it. See [database-rbac.md](api/docs/database-rbac.md).

**2. Clone, configure, run.**

```bash
git clone git@github.com:0x11semprez/gojo.git && cd gojo
cp .env.example .env
# edit .env: match the passwords from step 1, and set
# WALLET_ENCRYPTION_KEY to the output of `openssl rand -hex 32`
make run
```

**3. Create an account and a wallet.**

```bash
# the password is generated server-side, shown once -- save it
curl -s localhost:6969/users -X POST -d '{"username":"alice"}'
# {"id":"...","username":"alice","password":"..."}

curl -s localhost:6969/wallets -X POST -d '{
  "username": "alice",
  "password": "<password from above>",
  "network": "bitcoin"
}'
# {"id":"...","network":"bitcoin","public_key":"02...","address":"bc1..."}
```

## Docs

no need to have a long readme, just check all schemas in the docs files. we focused more on design system because coding skills are not as valuable in the eyes of others anymore... </br>
</br>

- [architecture-context.md](api/docs/architecture-context.md) -- the whole project, one schema
- [database-rbac.md](api/docs/database-rbac.md) -- the two Postgres roles and what each can do
- [entropy.md](api/docs/entropy.md) -- how the rust generator seeds its keys
- [api.md](api/docs/api.md) -- login, server-generated passwords, wallet-key encryption

we love create
