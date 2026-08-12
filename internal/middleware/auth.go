// Package middleware provides authentication-related functions:
// password hashing and credential verification (login).
package middleware

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials is returned when the username does not exist
// or the password does not match. Both cases deliberately return the
// same error, so that a caller (or an attacker) cannot use it to
// guess which usernames exist.
var ErrInvalidCredentials = errors.New("invalid username or password")

// HashPassword turns a plaintext password into a bcrypt hash, which
// can be safely stored in the database (the "secret" BYTEA column in
// the users table). bcrypt embeds a random salt in its result, so
// hashing the same password twice never produces the same bytes.
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Login verifies a username/password pair against the database and
// returns the matching user's id on success.
func Login(ctx context.Context, db *bun.DB, username, password string) (string, error) {
	var (
		id   string
		hash []byte
	)

	// A single database round trip: fetch the id and the bcrypt hash
	// stored for this username. We query the table directly rather
	// than the user.User model to avoid an import cycle (the user
	// package already imports middleware).
	err := db.NewSelect().
		Table("users").
		Column("id", "secret").
		Where("username = ?", username).
		Where("deleted_at IS NULL").
		Scan(ctx, &id, &hash)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	// CompareHashAndPassword re-hashes the provided password with the
	// salt extracted from the stored hash and compares the result in
	// constant time, which protects against timing attacks.
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return id, nil
}
