package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

func main() {
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
		Handler:           proxy,
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
