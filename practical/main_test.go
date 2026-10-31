package main

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTAuth(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef") // Test key only.
	future := time.Now().Add(time.Hour).Unix()
	past := time.Now().Add(-time.Hour).Unix()
	sign := func(method jwt.SigningMethod, claims jwt.MapClaims, key any) string {
		t.Helper()
		token, err := jwt.NewWithClaims(method, claims).SignedString(key)
		if err != nil {
			t.Fatal(err)
		}
		return token
	}
	valid := sign(jwt.SigningMethodHS256, jwt.MapClaims{"exp": future}, secret)
	cases := []struct {
		name   string
		header string
		want   int
	}{
		{"missing header", "", 401},
		{"malformed scheme", "Basic " + valid, 401},
		{"missing token", "Bearer", 401},
		{"extra field", "Bearer " + valid + " extra", 401},
		{"malformed token", "Bearer not-a-jwt", 401},
		{"wrong signature", "Bearer " + sign(jwt.SigningMethodHS256, jwt.MapClaims{"exp": future}, []byte("wrong test key")), 401},
		{"expired", "Bearer " + sign(jwt.SigningMethodHS256, jwt.MapClaims{"exp": past}, secret), 401},
		{"missing exp", "Bearer " + sign(jwt.SigningMethodHS256, jwt.MapClaims{}, secret), 401},
		{"invalid exp", "Bearer " + sign(jwt.SigningMethodHS256, jwt.MapClaims{"exp": "invalid"}, secret), 401},
		{"HS384", "Bearer " + sign(jwt.SigningMethodHS384, jwt.MapClaims{"exp": future}, secret), 401},
		{"unsigned", "Bearer " + sign(jwt.SigningMethodNone, jwt.MapClaims{"exp": future}, jwt.UnsafeAllowNoneSignatureType), 401},
		{"future nbf", "Bearer " + sign(jwt.SigningMethodHS256, jwt.MapClaims{"exp": future, "nbf": future}, secret), 401},
		{"valid", "Bearer " + valid, 200},
		{"case insensitive scheme", "bearer " + valid, 200},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			backend := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Header.Get("Authorization") != "" {
					t.Error("Authorization reached backend")
				}
				body, _ := io.ReadAll(r.Body)
				if r.Method != "POST" || r.URL.Path != "/post" || r.URL.Query().Get("item") != "value" || string(body) != "hello" {
					t.Error("request changed")
				}
				w.WriteHeader(http.StatusOK)
			})
			// Even a valid query token must not authenticate a missing/invalid header.
			r := httptest.NewRequest("POST", "/post?item=value&access_token="+valid, strings.NewReader("hello"))
			if tc.header != "" {
				r.Header.Set("Authorization", tc.header)
			}
			w := httptest.NewRecorder()
			jwtAuth(backend, secret).ServeHTTP(w, r)
			if w.Code != tc.want {
				t.Fatalf("status = %d, want %d", w.Code, tc.want)
			}
			if (calls == 1) != (tc.want == 200) {
				t.Fatal("unauthorized request reached backend or valid request was blocked")
			}
		})
	}
}

func TestReadJWTSecret(t *testing.T) {
	cases := []struct {
		name  string
		value string
		valid bool
	}{
		{"missing", "", false},
		{"not base64", "invalid!", false},
		{"short", base64.StdEncoding.EncodeToString([]byte("short")), false},
		{"repeated bytes", base64.StdEncoding.EncodeToString(make([]byte, 32)), false},
		{"valid format", base64.StdEncoding.EncodeToString([]byte("0123456789abcdef0123456789abcdef")), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("JWT_SECRET", tc.value)
			_, err := readJWTSecret()
			if (err == nil) != tc.valid {
				t.Fatal("unexpected secret validation result")
			}
		})
	}
}
