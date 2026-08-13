// Package httpx provides small HTTP helpers shared by the handler
// packages (user, wallet): JSON encoding/decoding and PostgreSQL
// constraint-violation detection, so each package does not have to
// reimplement them.
package httpx

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/uptrace/bun/driver/pgdriver"
)

// WriteJSON encodes v as JSON and writes it with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError writes a {"error": message} JSON body with the given status code.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// DecodeJSON decodes the request body into v.
func DecodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// IsConflict reports whether err is a PostgreSQL integrity constraint
// violation (e.g. a duplicate username), which callers should map to
// HTTP 409 Conflict instead of 500.
func IsConflict(err error) bool {
	var pgErr pgdriver.Error
	return errors.As(err, &pgErr) && pgErr.IntegrityViolation()
}
