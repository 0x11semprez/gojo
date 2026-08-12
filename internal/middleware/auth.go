package middleware

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidCredentials is returned when the username does not exist or the
// password does not match. Both cases return the same error on purpose, so a
// caller (or an attacker) cannot use the error to guess valid usernames.
var ErrInvalidCredentials = errors.New("invalid username or password")

// HashPassword turns a plaintext password into a bcrypt hash that is safe to
// store in the database (the "secret" BYTEA column of the users table).
// bcrypt embeds a random salt in its output, so hashing the same password
// twice never produces the same bytes.
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// Login authenticates a username/password pair against the database and
// returns the matching user's id on success.
func Login(ctx context.Context, db *bun.DB, username, password string) (string, error) {
	var (
		id   string
		hash []byte
	)

	// One round trip: fetch the id and the stored bcrypt hash for that
	// username. We query the table directly instead of the user.User model
	// to avoid an import cycle (package user already imports middleware).
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

	// CompareHashAndPassword re-hashes the candidate password with the salt
	// extracted from the stored hash and compares the result in constant
	// time, so it is safe against timing attacks.
	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return id, nil
}
