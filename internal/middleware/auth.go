package middleware

import "net/http"

func Login(secret string) (string, *http.Response, error) {}
