```mermaid
flowchart TD

CPU["CPU vendor + brand + frequency\n(machine_info)"]
RNG["OS CSPRNG\n32 random bytes (OsRng)"]

Seed["seed_material\n(bytes, concatenated)"]

BTC["SHA-256, re-hashed with a nonce\nuntil it's a valid secp256k1 scalar"]
XMR["hash_to_scalar\n(ed25519 scalar field)"]

PrivBTC[(bitcoin private key)]
PrivSpend[(monero private spend key)]
PrivView[(monero private view key)]

CPU --> Seed
RNG --> Seed

Seed --> BTC --> PrivBTC
Seed --> XMR --> PrivSpend
PrivSpend -->|hash_to_scalar of the spend key| PrivView
```

Two independent sources feed every key: a hardware fingerprint (CPU vendor/brand/frequency) and the OS's own secure random generator. Neither is trusted alone — the point is to not depend on a single RNG path, [as coldcard learned the hard way](../../README.md).

Diagnostics (cpu status, memory, cores) go to stderr; only the generated key material is printed as JSON on stdout (see `generator/src/main.rs`).
