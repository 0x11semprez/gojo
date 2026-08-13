```mermaid
flowchart LR

User[User]

API[Gojo API - Go]

Generator["Key Generator - Rust\n(bitcoin secp256k1, monero ed25519)"]

Vault["Envelope encryption\npassword-derived key + server key"]

DB[(PostgreSQL: accounts and wallets)]

User -->|HTTP requests| API

API -->|spawns subprocess, network arg| Generator

Generator -->|JSON on stdout: private + public key| API

API -->|private key| Vault

Vault -->|encrypted private key| DB

API -->|public key, address: stored in clear| DB
```

Public keys and addresses are meant to be shared, so they are stored in clear. Private keys never are: they go through two layers first, one derived from the account's password (see [api.md](api.md)), one keyed by a secret only the server holds — see [database-rbac.md](database-rbac.md) for who can read what in Postgres, and [entropy.md](entropy.md) for how the generator produces key material in the first place.
