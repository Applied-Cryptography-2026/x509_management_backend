package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// Validator performs cryptographic validation of certificates and chains.
type Validator struct{}

// NewValidator creates a new Validator instance.
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateCertificate performs structural and crypto validation of a single certificate.
func (v *Validator) ValidateCertificate(cert *x509.Certificate) error {
	if cert == nil {
		return errors.New("crypto/validator: certificate is nil")
	}
	if cert.NotAfter.Before(cert.NotBefore) {
		return errors.New("crypto/validator: NotAfter is before NotBefore")
	}
	if cert.Signature == nil {
		return errors.New("crypto/validator: signature is empty")
	}
	return nil
}

// ValidateChain validates a certificate chain against a trusted root pool.
func (v *Validator) ValidateChain(chain []*x509.Certificate, roots *x509.CertPool) error {
	if len(chain) == 0 {
		return errors.New("crypto/validator: chain is empty")
	}

	opts := x509.VerifyOptions{
		DNSName:       chain[0].DNSNames[0], // leaf DNS name
		Intermediates: x509.NewCertPool(),
		Roots:         roots,
	}

	// Add all but the last cert (root) to intermediates pool
	for _, cert := range chain[:len(chain)-1] {
		opts.Intermediates.AddCert(cert)
	}

	_, err := chain[0].Verify(opts)
	if err != nil {
		return fmt.Errorf("crypto/validator: chain verification failed: %w", err)
	}

	return nil
}

// ParseChainPEM parses a concatenated PEM chain into individual certificates.
func (v *Validator) ParseChainPEM(pemChain string) ([]*x509.Certificate, error) {
	var certs []*x509.Certificate
	remaining := pemChain

	for {
		block, rest := pem.Decode([]byte(remaining))
		if block == nil {
			break
		}
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("crypto/validator: failed to parse cert: %w", err)
			}
			certs = append(certs, cert)
		}
		remaining = string(rest)
	}

	return certs, nil
}

// VerifySignature verifies that cert was signed by issuer.
func (v *Validator) VerifySignature(cert, issuer *x509.Certificate) error {
	return cert.CheckSignatureFrom(issuer)
}

// KeyUsageFromString converts a string key usage name to pkix.KeyUsage.
func KeyUsageFromString(name string) (x509.KeyUsage, bool) {
	switch name {
	case "digitalSignature":
		return x509.KeyUsageDigitalSignature, true
	case "keyEncipherment":
		return x509.KeyUsageKeyEncipherment, true
	case "dataEncipherment":
		return x509.KeyUsageDataEncipherment, true
	case "keyCertSign":
		return x509.KeyUsageCertSign, true
	case "cRLSign":
		return x509.KeyUsageCRLSign, true
	case "keyAgreement":
		return x509.KeyUsageKeyAgreement, true
	case "encipherOnly":
		return x509.KeyUsageEncipherOnly, true
	case "decipherOnly":
		return x509.KeyUsageDecipherOnly, true
	default:
		return 0, false
	}
}
