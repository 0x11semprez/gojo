```mermaid
sequenceDiagram
    participant C as Client
    participant A as API (Go)
    participant DB as Postgres

    C->>A: POST /users {username}
    A->>A: generate a random 256-bit password (crypto/rand)
    A->>A: bcrypt-hash it
    A->>DB: store username + bcrypt hash
    A-->>C: {id, username, password}
    note over A,C: password shown once here, never stored in clear again

    C->>A: POST /login {username, password}
    A->>DB: fetch the bcrypt hash
    A->>A: bcrypt compare (constant time)
    A-->>C: {id}

    C->>A: POST /wallets {username, password, network}
    A->>A: login (as above)
    A->>A: derive a key from password + a fresh random salt
    A->>A: ask the rust generator for a keypair
    A->>A: encrypt the private key: password-derived layer, then server-key layer
    A->>DB: store encrypted private key + public key in clear
    A-->>C: {id, public_key, address}
```

The password is never chosen by a human: it is generated server-side, 256 bits of `crypto/rand`, so brute-forcing it is not a realistic attack. Only its bcrypt hash is stored — one-way, so the server itself cannot recompute the wallet-encryption key without the client sending the password again on each request. See [entropy.md](entropy.md) for how the key material itself is generated, and [database-rbac.md](database-rbac.md) for who can read what in Postgres.
