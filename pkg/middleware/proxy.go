package middleware

import (
	"net"
	"net/http"
	"regexp"
	"strings"
)

var (
	// De-facto standard header keys.
	xForwardedFor    = http.CanonicalHeaderKey("X-Forwarded-For")
	xForwardedHost   = http.CanonicalHeaderKey("X-Forwarded-Host")
	xForwardedProto  = http.CanonicalHeaderKey("X-Forwarded-Proto")
	xForwardedScheme = http.CanonicalHeaderKey("X-Forwarded-Scheme")
	xRealIP          = http.CanonicalHeaderKey("X-Real-IP")
)

var (
	forwarded  = http.CanonicalHeaderKey("Forwarded")
	forRegex   = regexp.MustCompile(`(?i)(?:for=)([^(;|,| )]+)`)
	protoRegex = regexp.MustCompile(`(?i)(?:proto=)(https|http)`)
)

func ProxyHeadersMiddleware(h http.Handler, trustedSubnets []*net.IPNet) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		trusted := false
		ip := net.ParseIP(r.RemoteAddr)
		for _, cidr := range trustedSubnets {
			if cidr.Contains(ip) {
				trusted = true
				break
			}
		}

		if !trusted {
			h.ServeHTTP(w, r)
			return
		}

		// Set the remote IP with the value passed from the proxy.
		if fwd := getIP(r); fwd != "" {
			r.RemoteAddr = fwd
		}

		// Set the scheme (proto) with the value passed from the proxy.
		if scheme := getScheme(r); scheme != "" {
			r.URL.Scheme = scheme
		}
		// Set the host with the value passed by the proxy
		if r.Header.Get(xForwardedHost) != "" {
			r.Host = r.Header.Get(xForwardedHost)
		}
		// Call the next handler in the chain.
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

// getIP retrieves the IP from the X-Forwarded-For, X-Real-IP and RFC7239
// Forwarded headers (in that order).
func getIP(r *http.Request) string {
	var addr string

	if fwd := r.Header.Get(xForwardedFor); fwd != "" {
		// Only grab the first (client) address. Note that '192.168.0.1,
		// 10.1.1.1' is a valid key for X-Forwarded-For where addresses after
		// the first may represent forwarding proxies earlier in the chain.
		s := strings.Index(fwd, ", ")
		if s == -1 {
			s = len(fwd)
		}
		addr = fwd[:s]
	} else if fwd := r.Header.Get(xRealIP); fwd != "" {
		// X-Real-IP should only contain one IP address (the client making the
		// request).
		addr = fwd
	} else if fwd := r.Header.Get(forwarded); fwd != "" {
		// match should contain at least two elements if the protocol was
		// specified in the Forwarded header. The first element will always be
		// the 'for=' capture, which we ignore. In the case of multiple IP
		// addresses (for=8.8.8.8, 8.8.4.4,172.16.1.20 is valid) we only
		// extract the first, which should be the client IP.
		if match := forRegex.FindStringSubmatch(fwd); len(match) > 1 {
			// IPv6 addresses in Forwarded headers are quoted-strings. We strip
			// these quotes.
			addr = strings.Trim(match[1], `"`)
		}
	}

	return addr
}

// getScheme retrieves the scheme from the X-Forwarded-Proto and RFC7239
// Forwarded headers (in that order).
func getScheme(r *http.Request) string {
	var scheme string

	// Retrieve the scheme from X-Forwarded-Proto.
	if proto := r.Header.Get(xForwardedProto); proto != "" {
		scheme = strings.ToLower(proto)
	} else if proto = r.Header.Get(xForwardedScheme); proto != "" {
		scheme = strings.ToLower(proto)
	} else if proto = r.Header.Get(forwarded); proto != "" {
		// match should contain at least two elements if the protocol was
		// specified in the Forwarded header. The first element will always be
		// the 'proto=' capture, which we ignore. In the case of multiple proto
		// parameters (invalid) we only extract the first.
		if match := protoRegex.FindStringSubmatch(proto); len(match) > 1 {
			scheme = strings.ToLower(match[1])
		}
	}

	return scheme
}
