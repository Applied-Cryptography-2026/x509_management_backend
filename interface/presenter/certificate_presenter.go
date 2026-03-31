package presenter

import (
	"github.com/your-org/x509-clean-architecture/domain/model"
	"github.com/your-org/x509-clean-architecture/usecase/presenter"
)

// certificatePresenter is the concrete driven adapter for the presenter output port.
type certificatePresenter struct{}

// NewCertificatePresenter constructs the concrete presenter.
func NewCertificatePresenter() presenter.CertificatePresenter {
	return &certificatePresenter{}
}

// ResponseCert applies any presentation-layer enrichment to a single certificate.
func (cp *certificatePresenter) ResponseCert(cert *model.Certificate) *model.Certificate {
	// Example enrichment: mask the private key in responses
	// cert.KeyPEM = ""
	return cert
}

// ResponseCerts applies enrichment to a list of certificates.
func (cp *certificatePresenter) ResponseCerts(certs []*model.Certificate) []*model.Certificate {
	out := make([]*model.Certificate, len(certs))
	for i, cert := range certs {
		out[i] = cp.ResponseCert(cert)
	}
	return out
}

// ResponseChain returns the chain as-is (no enrichment needed).
func (cp *certificatePresenter) ResponseChain(chain *model.CertificateChain) *model.CertificateChain {
	return chain
}

// ErrorResponse formats an error into a structured HTTP response body.
func (cp *certificatePresenter) ErrorResponse(err error) map[string]interface{} {
	return map[string]interface{}{
		"error":   err.Error(),
		"type":    "certificate_error",
	}
}
