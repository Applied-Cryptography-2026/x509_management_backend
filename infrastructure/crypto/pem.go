package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// DecodePEMBlock extracts the first matching PEM block and returns its DER bytes.
func DecodePEMBlock(pemStr string, blockType string) ([]byte, error) {
	remaining := pemStr
	for {
		block, rest := pem.Decode([]byte(remaining))
		if block == nil {
			return nil, errors.New("crypto/pem: no PEM block found of type " + blockType)
		}
		if block.Type == blockType {
			return block.Bytes, nil
		}
		remaining = string(rest)
	}
}

// EncodeToPEM DER-encodes a value and wraps it in a PEM block.
func EncodeToPEM(der []byte, blockType string) string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  blockType,
		Bytes: der,
	}))
}

// ParseCertificatePEM parses a PEM-encoded certificate.
func ParseCertificatePEM(pemStr string) (*x509.Certificate, error) {
	der, err := DecodePEMBlock(pemStr, "CERTIFICATE")
	if err != nil {
		return nil, fmt.Errorf("crypto/pem: ParseCertificatePEM: %w", err)
	}
	return x509.ParseCertificate(der)
}

// ParsePrivateKeyPEM parses a PEM-encoded private key (supports RSA, ECDSA, Ed25519).
func ParsePrivateKeyPEM(pemStr string) (any, error) {
	der, err := DecodePEMBlock(pemStr, "PRIVATE KEY")
	if err != nil {
		// Try PKCS8
		der, err = DecodePEMBlock(pemStr, "EC PRIVATE KEY")
		if err != nil {
			der, err = DecodePEMBlock(pemStr, "RSA PRIVATE KEY")
		}
	}
	if err != nil {
		return nil, fmt.Errorf("crypto/pem: ParsePrivateKeyPEM: %w", err)
	}
	return parsePrivateKey(der)
}

func parsePrivateKey(der []byte) (any, error) {
	// Try PKCS8 first (preferred)
	key, err := x509.ParsePKCS8PrivateKey(der)
	if err == nil {
		return key, nil
	}
	// Fall back: try PKCS1 RSA
	_, err = x509.ParsePKCS1PrivateKey(der)
	if err == nil {
		return x509.ParsePKCS1PrivateKey(der)
	}
	// Fall back: try EC
	_, err = x509.ParseECPrivateKey(der)
	if err == nil {
		return x509.ParseECPrivateKey(der)
	}
	return nil, errors.New("crypto/pem: unable to parse private key format")
}

// ParseCSRPEM parses a PEM-encoded CSR.
func ParseCSRPEM(pemStr string) (*x509.CertificateRequest, error) {
	der, err := DecodePEMBlock(pemStr, "CERTIFICATE REQUEST")
	if err != nil {
		return nil, fmt.Errorf("crypto/pem: ParseCSRPEM: %w", err)
	}
	return x509.ParseCertificateRequest(der)
}
