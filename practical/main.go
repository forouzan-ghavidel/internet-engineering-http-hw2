package main

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret, err := readJWTSecret()
	if err != nil {
		log.Fatal(err)
	}

	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://backend:8081"
	}

	target, err := url.Parse(backendURL)
	if err != nil || target.Hostname() == "" || target.User != nil ||
		(target.Scheme != "http" && target.Scheme != "https") {
		log.Fatal("BACKEND_URL must be an HTTP or HTTPS URL without credentials")
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = 10 * time.Second
	proxy.Transport = transport
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		// Do not log request URLs, headers, or bodies, which may contain secrets.
		log.Printf("backend request failed (%T)", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	server := &http.Server{
		Addr:              ":80",
		Handler:           jwtAuth(proxy, secret),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Print("HTTP reverse proxy listening on :80")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func readJWTSecret() ([]byte, error) {
	// The local .env stores at least 32 cryptographically random bytes as base64.
	secret, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_SECRET"))
	unique := make(map[byte]bool)
	for _, b := range secret {
		unique[b] = true
	}
	if err != nil || len(secret) < 32 || len(unique) < 16 {
		return nil, errors.New("JWT_SECRET must be base64-encoded random bytes (at least 32 bytes); missing or weak key")
	}
	return secret, nil
}

func jwtAuth(next http.Handler, secret []byte) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header.Values("Authorization")
		parts := strings.Fields(r.Header.Get("Authorization"))
		if len(headers) != 1 || len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (any, error) {
			return secret, nil
		}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())
		if err != nil || !token.Valid {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r.Header.Del("Authorization")
		next.ServeHTTP(w, r)
	})
}
