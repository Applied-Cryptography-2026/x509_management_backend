package certificate

import (
	"crypto/x509"
	"net"
	"strings"
	"time"

	"github.com/your-org/x509-clean-architecture/domain/model"
	"github.com/your-org/x509-clean-architecture/infrastructure/crypto"
)

// Converter converts between x509.Certificate and domain model.
type Converter struct{}

// NewConverter creates a new Converter.
func NewConverter() *Converter {
	return &Converter{}
}

// ToDomain converts an x509.Certificate to the domain Certificate model.
func (c *Converter) ToDomain(cert *x509.Certificate, pemStr, keyPEM string) *model.Certificate {
	status := model.CertStatusActive
	if time.Now().Before(cert.NotBefore) {
		status = model.CertStatusPending
	} else if time.Now().After(cert.NotAfter) {
		status = model.CertStatusExpired
	}

	return &model.Certificate{
		Subject:     cert.Subject.String(),
		Issuer:      cert.Issuer.String(),
		Serial:      cert.SerialNumber.Text(16),
		Fingerprint: crypto.SHA256Fingerprint(cert),
		NotBefore:   cert.NotBefore,
		NotAfter:    cert.NotAfter,
		KeyUsage:    keyUsageStrings(cert.KeyUsage),
		ExtKeyUsage: extKeyUsageStrings(cert.ExtKeyUsage),
		DNSNames:    cert.DNSNames,
		IPAddresses: ipStrings(cert.IPAddresses),
		IsCA:        cert.IsCA,
		CRLURL:      strings.Join(cert.CRLDistributionPoints, ","),
		OCSPURL:     strings.Join(cert.OCSPServer, ","),
		CertPEM:     pemStr,
		KeyPEM:      keyPEM,
		Status:      status,
	}
}

// ToX509 converts a domain Certificate PEM string back to an x509.Certificate.
func (c *Converter) ToX509(cert *model.Certificate) (*x509.Certificate, error) {
	return crypto.ParseCertificatePEM(cert.CertPEM)
}

// ---------------------------------------------------------------------------

func keyUsageStrings(ku x509.KeyUsage) []string {
	mapping := []struct {
		bit  x509.KeyUsage
		name string
	}{
		{x509.KeyUsageDigitalSignature, "digitalSignature"},
		{x509.KeyUsageContentCommitment, "contentCommitment"}, // Non-Repudiation
		{x509.KeyUsageKeyEncipherment, "keyEncipherment"},
		{x509.KeyUsageDataEncipherment, "dataEncipherment"},
		{x509.KeyUsageKeyAgreement, "keyAgreement"},
		{x509.KeyUsageCertSign, "keyCertSign"},
		{x509.KeyUsageCRLSign, "cRLSign"},
		{x509.KeyUsageEncipherOnly, "encipherOnly"},
		{x509.KeyUsageDecipherOnly, "decipherOnly"},
	}

	var out []string
	for _, m := range mapping {
		if ku&m.bit != 0 {
			out = append(out, m.name)
		}
	}
	return out
}

func extKeyUsageStrings(eku []x509.ExtKeyUsage) []string {
	mapping := map[x509.ExtKeyUsage]string{
		x509.ExtKeyUsageServerAuth:      "serverAuth",
		x509.ExtKeyUsageClientAuth:      "clientAuth",
		x509.ExtKeyUsageCodeSigning:     "codeSigning",
		x509.ExtKeyUsageEmailProtection: "emailProtection",
		x509.ExtKeyUsageTimeStamping:    "timeStamping",
		x509.ExtKeyUsageOCSPSigning:     "ocspSigning",
	}
	var out []string
	for _, u := range eku {
		if name, ok := mapping[u]; ok {
			out = append(out, name)
		}
	}
	return out
}

func ipStrings(ips []net.IP) []string {
	out := make([]string, len(ips))
	for i, ip := range ips {
		out[i] = ip.String()
	}
	return out
}
