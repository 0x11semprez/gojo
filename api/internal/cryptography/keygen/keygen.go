// Package keygen bridges the golang API to the rust key generator
// (see the "generator" crate at the repository root): it runs the
// compiled binary as a subprocess for a given network and decodes the
// raw key material it prints as JSON on stdout. Nothing is written to
// disk by the generator; the key material only ever exists in memory
// on the Go side, which is responsible for encrypting it before
// storage (see package wallet).
package keygen

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"gojo/internal/config"
)

// Network identifies which key material the rust generator should produce.
type Network string

const (
	Bitcoin Network = "bitcoin"
	Monero  Network = "monero"
)

// BitcoinKeyMaterial holds the raw (unencrypted) key material for a
// freshly generated bitcoin keypair.
type BitcoinKeyMaterial struct {
	PrivateKey []byte
	PublicKey  []byte
	Address    string
}

// MoneroKeyMaterial holds the raw (unencrypted) key material for a
// freshly generated monero keypair (spend key and view key).
type MoneroKeyMaterial struct {
	PrivateSpendKey []byte
	PublicSpendKey  []byte
	PrivateViewKey  []byte
	PublicViewKey   []byte
}

// rawOutput mirrors the single-line JSON object the rust generator
// prints on stdout (see generator/src/main.rs). Fields are hex
// strings; fields that do not apply to the requested network are
// simply empty.
type rawOutput struct {
	Network         string `json:"network"`
	PrivateKey      string `json:"private_key"`
	PublicKey       string `json:"public_key"`
	Address         string `json:"address"`
	PrivateSpendKey string `json:"private_spend_key"`
	PublicSpendKey  string `json:"public_spend_key"`
	PrivateViewKey  string `json:"private_view_key"`
	PublicViewKey   string `json:"public_view_key"`
}

// binaryPath resolves the compiled rust generator binary:
// cfg.GeneratorBinPath when set, otherwise
// generator/target/release/generator relative to the repository root.
// The root is computed from this source file's location (like
// database.StartMigrations does for the migrations directory),
// so it does not depend on the process's working directory.
func binaryPath(cfg *config.Config) string {
	if cfg != nil && cfg.GeneratorBinPath != "" {
		return cfg.GeneratorBinPath
	}
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..")
	return filepath.Join(repoRoot, "generator", "target", "release", "generator")
}

// run executes the rust generator for the given network and parses
// its JSON stdout.
func run(ctx context.Context, cfg *config.Config, network Network) (rawOutput, error) {
	cmd := exec.CommandContext(ctx, binaryPath(cfg), string(network))

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return rawOutput{}, fmt.Errorf("generator %s failed: %w (stderr: %s)", network, err, stderr.String())
	}

	var out rawOutput
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &out); err != nil {
		return rawOutput{}, fmt.Errorf("cannot parse generator output: %w", err)
	}
	if out.Network != string(network) {
		return rawOutput{}, fmt.Errorf("generator returned network %q, expected %q", out.Network, network)
	}
	return out, nil
}

func decodeHex(field, value string) ([]byte, error) {
	b, err := hex.DecodeString(value)
	if err != nil {
		return nil, fmt.Errorf("invalid hex for %s: %w", field, err)
	}
	return b, nil
}

// GenerateBitcoinKeys runs the rust generator for the bitcoin network
// and decodes its output into raw key material.
func GenerateBitcoinKeys(ctx context.Context, cfg *config.Config) (*BitcoinKeyMaterial, error) {
	out, err := run(ctx, cfg, Bitcoin)
	if err != nil {
		return nil, err
	}

	privateKey, err := decodeHex("private_key", out.PrivateKey)
	if err != nil {
		return nil, err
	}
	publicKey, err := decodeHex("public_key", out.PublicKey)
	if err != nil {
		return nil, err
	}
	if out.Address == "" {
		return nil, errors.New("generator returned an empty bitcoin address")
	}

	return &BitcoinKeyMaterial{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    out.Address,
	}, nil
}

// GenerateMoneroKeys runs the rust generator for the monero network
// and decodes its output into raw key material.
func GenerateMoneroKeys(ctx context.Context, cfg *config.Config) (*MoneroKeyMaterial, error) {
	out, err := run(ctx, cfg, Monero)
	if err != nil {
		return nil, err
	}

	privateSpendKey, err := decodeHex("private_spend_key", out.PrivateSpendKey)
	if err != nil {
		return nil, err
	}
	publicSpendKey, err := decodeHex("public_spend_key", out.PublicSpendKey)
	if err != nil {
		return nil, err
	}
	privateViewKey, err := decodeHex("private_view_key", out.PrivateViewKey)
	if err != nil {
		return nil, err
	}
	publicViewKey, err := decodeHex("public_view_key", out.PublicViewKey)
	if err != nil {
		return nil, err
	}

	return &MoneroKeyMaterial{
		PrivateSpendKey: privateSpendKey,
		PublicSpendKey:  publicSpendKey,
		PrivateViewKey:  privateViewKey,
		PublicViewKey:   publicViewKey,
	}, nil
}
