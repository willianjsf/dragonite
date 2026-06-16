package httputil

import (
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	RequestTimeout time.Duration = 2 * time.Second
)

// UnimplementedHandler representa um handler não implementado ainda.
func UnimplementedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

// GetClientIP extracts the real IP address of the client from the HTTP request.
func GetClientIP(r *http.Request) string {
	// 1. Try the X-Real-IP header first (common in Nginx setups)
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// 2. Try the X-Forwarded-For header.
	// This header can contain a comma-separated list of IPs if the request
	// passed through multiple proxies. The first IP is the original client.
	ips := r.Header.Get("X-Forwarded-For")
	if ips != "" {
		splitIps := strings.Split(ips, ",")
		if len(splitIps) > 0 {
			return strings.TrimSpace(splitIps[0])
		}
	}

	// 3. Fallback to standard RemoteAddr
	// RemoteAddr usually comes in the format "IP:PORT" (e.g., "192.168.1.5:43212")
	// We use net.SplitHostPort to strip the port and keep only the IP.
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If it fails (e.g., there is no port), just return the raw string
		return r.RemoteAddr
	}

	return ip
}

// Extrai o token the autenticacao do usuário
func ExtractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// The Authorization header is expected to be in the format "Bearer <token>"
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {
		return ""
	}

	return splitAuth[1]
}
