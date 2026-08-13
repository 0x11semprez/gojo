package integrationtest

import (
	"context"
	"net/http"
	"testing"

	"gojo/internal/app"
	"gojo/internal/config"
	"gojo/internal/cryptography/keygen"
	"gojo/internal/cryptography/vault"
	"gojo/internal/database"
	"gojo/internal/server"
	"gojo/internal/user"
	"gojo/internal/wallet"
)

// walletOwnerPassword is the plaintext password used both to create
// the "wallet-owner" test account and to derive the inner
// envelope-encryption key when creating its wallets.
const walletOwnerPassword = "wallet-owner-integration-test-password"

// TestCreateWallet is an integration test that starts a real App
// against the test database and exercises wallet.CreateWallet for
// both supported networks: it checks that the stored row round-trips
// through decryption (outer layer with the server's
// WalletEncryptionKey, inner layer derived from the account's
// password) back to key material of the expected size, and that
// public data (public key, bitcoin address) is stored in clear.
func TestCreateWallet(t *testing.T) {
	cfg, err := config.NewConfigTest()
	if err != nil {
		t.Fatal(err)
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	srv, err := server.NewServe(cfg, mux)
	if err != nil {
		t.Fatal(err)
	}

	testApp := app.App{}.NewApp(srv, db, cfg)

	userID, err := user.CreateUser(testApp, "wallet-owner", walletOwnerPassword)
	if err != nil {
		t.Fatal(err)
	}

	outerKey, err := vault.KeyFromHex(cfg.WalletEncryptionKey)
	if err != nil {
		t.Fatalf("invalid WalletEncryptionKey in test config: %v", err)
	}

	t.Run("bitcoin", func(t *testing.T) {
		created, err := wallet.CreateWallet(testApp, userID, walletOwnerPassword, keygen.Bitcoin, "my btc wallet")
		if err != nil {
			t.Fatal(err)
		}
		if len(created.PublicKey) != 33 {
			t.Errorf("CreateWallet: len(PublicKey) = %d, want 33 (compressed secp256k1 public key)", len(created.PublicKey))
		}
		if created.Address == "" {
			t.Error("CreateWallet: Address is empty, want a bitcoin address")
		}

		w := fetchWallet(t, testApp, created.Id)

		if w.Network != "bitcoin" {
			t.Errorf("Network = %q, want %q", w.Network, "bitcoin")
		}
		if w.Address == "" {
			t.Error("Address is empty, want a bitcoin address")
		}
		if len(w.PublicKey) != 33 {
			t.Errorf("len(PublicKey) = %d, want 33 (compressed secp256k1 public key)", len(w.PublicKey))
		}
		if w.PrivateViewKey != nil || w.PublicViewKey != nil {
			t.Error("bitcoin wallet has non-nil view key columns, want nil")
		}

		privateKey := decryptEnvelope(t, outerKey, walletOwnerPassword, w.PasswordSalt, w.PrivateKey)
		if len(privateKey) != 32 {
			t.Errorf("decrypted private key length = %d, want 32", len(privateKey))
		}
	})

	t.Run("monero", func(t *testing.T) {
		created, err := wallet.CreateWallet(testApp, userID, walletOwnerPassword, keygen.Monero, "my xmr wallet")
		if err != nil {
			t.Fatal(err)
		}
		if len(created.PublicViewKey) != 32 {
			t.Errorf("CreateWallet: len(PublicViewKey) = %d, want 32 (view public key)", len(created.PublicViewKey))
		}

		w := fetchWallet(t, testApp, created.Id)

		if w.Network != "monero" {
			t.Errorf("Network = %q, want %q", w.Network, "monero")
		}
		if w.Address != "" {
			t.Errorf("Address = %q, want empty (monero address derivation is not implemented)", w.Address)
		}
		if len(w.PublicKey) != 32 {
			t.Errorf("len(PublicKey) = %d, want 32 (spend public key)", len(w.PublicKey))
		}
		if len(w.PublicViewKey) != 32 {
			t.Errorf("len(PublicViewKey) = %d, want 32 (view public key)", len(w.PublicViewKey))
		}

		spendKey := decryptEnvelope(t, outerKey, walletOwnerPassword, w.PasswordSalt, w.PrivateKey)
		if len(spendKey) != 32 {
			t.Errorf("decrypted private spend key length = %d, want 32", len(spendKey))
		}

		viewKey := decryptEnvelope(t, outerKey, walletOwnerPassword, w.PasswordSalt, w.PrivateViewKey)
		if len(viewKey) != 32 {
			t.Errorf("decrypted private view key length = %d, want 32", len(viewKey))
		}
	})

	t.Run("wrong password fails to decrypt", func(t *testing.T) {
		created, err := wallet.CreateWallet(testApp, userID, walletOwnerPassword, keygen.Bitcoin, "")
		if err != nil {
			t.Fatal(err)
		}

		w := fetchWallet(t, testApp, created.Id)

		outer, err := vault.Decrypt(outerKey, w.PrivateKey)
		if err != nil {
			t.Fatalf("outer decryption (server key) failed: %v", err)
		}

		wrongInnerKey := vault.DeriveKey("not-the-right-password", w.PasswordSalt)
		if _, err := vault.Decrypt(wrongInnerKey, outer); err == nil {
			t.Error("decrypting the inner layer with the wrong password succeeded, want an error")
		}
	})

	t.Run("unsupported network is rejected", func(t *testing.T) {
		if _, err := wallet.CreateWallet(testApp, userID, walletOwnerPassword, keygen.Network("dogecoin"), ""); err == nil {
			t.Error("CreateWallet() with an unsupported network succeeded, want an error")
		}
	})
}

// fetchWallet loads a wallet row by id, for assertions in
// TestCreateWallet.
func fetchWallet(t *testing.T, testApp *app.App, walletID string) *wallet.Wallet {
	t.Helper()

	w := new(wallet.Wallet)
	err := testApp.Database.NewSelect().
		Model(w).
		Where("id = ?", walletID).
		Scan(context.Background())
	if err != nil {
		t.Fatalf("fetch wallet %q: %v", walletID, err)
	}
	return w
}

// decryptEnvelope reverses wallet.CreateWallet's envelope encryption:
// outer layer first (server key), then inner layer (password-derived
// key), and fails the test if either layer does not decrypt.
func decryptEnvelope(t *testing.T, outerKey [vault.KeySize]byte, password string, salt, ciphertext []byte) []byte {
	t.Helper()

	inner, err := vault.Decrypt(outerKey, ciphertext)
	if err != nil {
		t.Fatalf("outer decryption (server key) failed: %v", err)
	}

	innerKey := vault.DeriveKey(password, salt)
	plaintext, err := vault.Decrypt(innerKey, inner)
	if err != nil {
		t.Fatalf("inner decryption (password-derived key) failed: %v", err)
	}
	return plaintext
}
