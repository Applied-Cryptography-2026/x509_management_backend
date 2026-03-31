package presenter

import (
	"github.com/your-org/x509-clean-architecture/domain/model"
)

// CertificatePresenter is a use-case-layer interface (output port).
// Defines how the interactor formats its output for the controller.
// The concrete implementation lives in interface/presenter/ (driven adapter).
type CertificatePresenter interface {
	ResponseCert(cert *model.Certificate) *model.Certificate
	ResponseCerts(certs []*model.Certificate) []*model.Certificate
	ResponseChain(chain *model.CertificateChain) *model.CertificateChain
	ErrorResponse(err error) map[string]interface{}
}
