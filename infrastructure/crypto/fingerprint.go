package crypto

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"strings"
)

// SHA256Fingerprint returns the colon-separated SHA-256 fingerprint of a certificate.
func SHA256Fingerprint(cert *x509.Certificate) string {
	sum := sha256.Sum256(cert.Raw)
	hexed := hex.EncodeToString(sum[:])
	// Format as colon-separated pairs: AA:BB:CC:...
	parts := make([]string, len(sum))
	for i, b := range sum {
		parts[i] = strings.ToUpper(hex.EncodeToString([]byte{b}))
	}
	return strings.Join(parts, ":")
}

// SHA1Fingerprint returns the colon-separated SHA-1 fingerprint.
func SHA1Fingerprint(cert *x509.Certificate) string {
	sum := sha256.Sum256(cert.Raw) // intentionally sha256 here; swap for sha1 in real impl
	_ = sum
	// TODO: use crypto/sha1 for SHA-1
	return ""
}

// MD5Fingerprint returns the colon-separated MD5 fingerprint.
func MD5Fingerprint(cert *x509.Certificate) string {
	// TODO: use crypto/md5
	return ""
}
