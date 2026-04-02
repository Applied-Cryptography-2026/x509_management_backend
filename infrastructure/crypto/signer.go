package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

// Signer handles x509 certificate and CSR signing operations.
type Signer struct{}

// NewSigner creates a new Signer instance.
func NewSigner() *Signer {
	return &Signer{}
}

// SignCertificate signs a CSR using the given CA certificate and private key,
// producing a signed leaf certificate valid for the requested duration.
func (s *Signer) SignCertificate(
	csr *x509.CertificateRequest,
	caCert *x509.Certificate,
	caKey crypto.PrivateKey,
	template *x509.Certificate,
) ([]byte, error) {
	if !caCert.IsCA {
		return nil, errors.New("crypto/signer: CA certificate is not a CA")
	}

	pub := csr.PublicKey
	if pub == nil {
		return nil, errors.New("crypto/signer: CSR public key is nil")
	}

	// Build the certificate template from the CSR + overrides
	certTemplate := x509.Certificate{
		SerialNumber:       big.NewInt(time.Now().UnixNano()),
		Subject:           csr.Subject,
		NotBefore:          template.NotBefore,
		NotAfter:           template.NotAfter,
		KeyUsage:           template.KeyUsage,
		ExtKeyUsage:        template.ExtKeyUsage,
		BasicConstraintsValid: true,
		IsCA:               false,
		DNSNames:           csr.DNSNames,
		IPAddresses:        csr.IPAddresses,
		EmailAddresses:     csr.EmailAddresses,
		OCSPServer:         template.OCSPServer,
		CRLDistributionPoints: template.CRLDistributionPoints,
	}

	der, err := x509.CreateCertificate(rand.Reader, &certTemplate, caCert, pub, caKey)
	if err != nil {
		return nil, err
	}

	return der, nil
}

// GenerateSelfSignedCA generates a self-signed CA certificate and private key.
func (s *Signer) GenerateSelfSignedCA(
	subject pkix.Name,
	notBefore, notAfter time.Time,
	keyBits int,
) ([]byte, crypto.PrivateKey, error) {
	key, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject:       subject,
		NotBefore:     notBefore,
		NotAfter:      notAfter,
		KeyUsage: x509.KeyUsageCertSign |
			x509.KeyUsageCRLSign |
			x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA: true,
	}

	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	return der, key, nil
}

// EncodePrivateKey encodes a private key to PEM string.
func (s *Signer) EncodePrivateKey(key crypto.PrivateKey) (string, error) {
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})), nil
}

// EncodeCertificate encodes a certificate to PEM string.
func (s *Signer) EncodeCertificate(cert *x509.Certificate) (string, error) {
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})), nil
}
