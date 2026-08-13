// Package wallet contains the business logic for generating and
// storing bitcoin/monero wallets: it asks package cryptography/keygen
// for freshly generated key material, envelope-encrypts the private
// key(s) with package cryptography/vault, and persists the result.
package wallet

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"

	"gojo/internal/app"
	"gojo/internal/cryptography/keygen"
	"gojo/internal/cryptography/vault"
	"gojo/internal/httpx"
	"gojo/internal/middleware"
)

// CreateWallet generates a new keypair for the given network (via the
// rust generator), envelope-encrypts the private key material with
// two layers -- an inner layer derived from the user's plaintext
// password, an outer layer keyed by the server's WalletEncryptionKey
// -- and stores the resulting wallet. It returns the stored wallet,
// so callers can read back the public data (public key, address)
// without a second database round trip.
//
// The plaintext password is only ever used in memory to derive the
// inner key; it is never stored.
func CreateWallet(a *app.App, userID, password string, network keygen.Network, name string) (*Wallet, error) {
	salt, err := vault.NewSalt()
	if err != nil {
		return nil, fmt.Errorf("cannot generate password salt: %w", err)
	}
	innerKey := vault.DeriveKey(password, salt)

	outerKey, err := vault.KeyFromHex(a.Config.WalletEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("invalid wallet encryption key: %w", err)
	}

	ctx := context.Background()

	w := &Wallet{
		UserId:       userID,
		Name:         name,
		Network:      string(network),
		PasswordSalt: salt,
	}

	switch network {
	case keygen.Bitcoin:
		keys, err := keygen.GenerateBitcoinKeys(ctx, a.Config)
		if err != nil {
			return nil, fmt.Errorf("cannot generate bitcoin keys: %w", err)
		}

		privateKey, err := envelopeEncrypt(innerKey, outerKey, keys.PrivateKey)
		if err != nil {
			return nil, err
		}

		w.PrivateKey = privateKey
		w.PublicKey = keys.PublicKey
		w.Address = keys.Address

	case keygen.Monero:
		keys, err := keygen.GenerateMoneroKeys(ctx, a.Config)
		if err != nil {
			return nil, fmt.Errorf("cannot generate monero keys: %w", err)
		}

		privateSpendKey, err := envelopeEncrypt(innerKey, outerKey, keys.PrivateSpendKey)
		if err != nil {
			return nil, err
		}
		privateViewKey, err := envelopeEncrypt(innerKey, outerKey, keys.PrivateViewKey)
		if err != nil {
			return nil, err
		}

		w.PrivateKey = privateSpendKey
		w.PublicKey = keys.PublicSpendKey
		w.PrivateViewKey = privateViewKey
		w.PublicViewKey = keys.PublicViewKey

	default:
		return nil, fmt.Errorf("unsupported network: %q", network)
	}

	if _, err := a.Database.NewInsert().Model(w).Exec(ctx); err != nil {
		return nil, err
	}

	return w, nil
}

// createWalletRequest is the JSON body expected by
// CreateWalletHandler. Network must be "bitcoin" or "monero"; Name is
// an optional user-facing label. Username/Password authenticate the
// caller (see middleware.Login) -- the password is also the account's
// plaintext password used to derive the wallet's inner encryption key,
// and is never stored.
type createWalletRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Network  string `json:"network"`
	Name     string `json:"name"`
}

// createWalletResponse is the JSON body returned by
// CreateWalletHandler. Only public data is returned: PublicKey (and,
// for monero, PublicViewKey) never need to stay secret, unlike the
// private key material, which is never sent back over the wire.
// Address is only set for bitcoin (see Wallet.Address).
type createWalletResponse struct {
	Id            string `json:"id"`
	Network       string `json:"network"`
	PublicKey     string `json:"public_key"`
	Address       string `json:"address,omitempty"`
	PublicViewKey string `json:"public_view_key,omitempty"`
}

// CreateWalletHandler handles "POST /wallets": it authenticates the
// caller (see middleware.Login), then generates a fresh keypair for
// the requested network (bitcoin or monero, see package
// cryptography/keygen) and stores it envelope-encrypted under the
// account's password (see CreateWallet). A wallet can only be created
// by proving ownership of the account with a valid username/password
// pair -- knowing the account's id alone is not enough.
func CreateWalletHandler(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createWalletRequest
		if err := httpx.DecodeJSON(r, &req); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if req.Username == "" || req.Password == "" || req.Network == "" {
			httpx.WriteError(w, http.StatusBadRequest, "username, password and network are required")
			return
		}

		network := keygen.Network(req.Network)
		if network != keygen.Bitcoin && network != keygen.Monero {
			httpx.WriteError(w, http.StatusBadRequest,
				fmt.Sprintf("unsupported network %q: must be %q or %q", req.Network, keygen.Bitcoin, keygen.Monero))
			return
		}

		userID, err := middleware.Login(r.Context(), a.Database, req.Username, req.Password)
		if err != nil {
			if errors.Is(err, middleware.ErrInvalidCredentials) {
				httpx.WriteError(w, http.StatusUnauthorized, "invalid username or password")
				return
			}
			log.Printf("create wallet: login: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, "cannot create wallet")
			return
		}

		wal, err := CreateWallet(a, userID, req.Password, network, req.Name)
		if err != nil {
			log.Printf("create wallet: %v", err)
			httpx.WriteError(w, http.StatusInternalServerError, "cannot create wallet")
			return
		}

		httpx.WriteJSON(w, http.StatusCreated, createWalletResponse{
			Id:            wal.Id,
			Network:       wal.Network,
			PublicKey:     hex.EncodeToString(wal.PublicKey),
			Address:       wal.Address,
			PublicViewKey: hex.EncodeToString(wal.PublicViewKey),
		})
	}
}

// envelopeEncrypt encrypts plaintext twice: first under innerKey (the
// password-derived key, tied to the user), then under outerKey (the
// server's wallet encryption key). Recovering the private key
// requires reversing both layers, in the opposite order.
func envelopeEncrypt(innerKey, outerKey [vault.KeySize]byte, plaintext []byte) ([]byte, error) {
	inner, err := vault.Encrypt(innerKey, plaintext)
	if err != nil {
		return nil, fmt.Errorf("cannot apply inner encryption: %w", err)
	}

	outer, err := vault.Encrypt(outerKey, inner)
	if err != nil {
		return nil, fmt.Errorf("cannot apply outer encryption: %w", err)
	}

	return outer, nil
}
