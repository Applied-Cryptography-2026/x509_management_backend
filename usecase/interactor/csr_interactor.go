package interactor

import (
	"crypto/x509"
	"errors"
	"net"
	"time"

	"github.com/your-org/x509-clean-architecture/domain/errdomain"
	"github.com/your-org/x509-clean-architecture/domain/model"
	"github.com/your-org/x509-clean-architecture/usecase/repository"
)

// csrInteractor holds the CSR use-case business logic.
type csrInteractor struct {
	CSRRepository repository.CSRRepository
	DBRepository repository.DBRepository
}

// CSRIteractor defines the contract for CSR-related use cases.
type CSRIteractor interface {
	SubmitCSR(pem string, requesterEmail string) (*model.CSR, error)
	ApproveCSR(id uint, approverID uint) (*model.CSR, error)
	RejectCSR(id uint, notes string) (*model.CSR, error)
	GetCSR(id uint) (*model.CSR, error)
	ListPendingCSRs() ([]*model.CSR, error)
	ListAllCSRs() ([]*model.CSR, error)
}

// NewCSRIteractor wires the CSR interactor.
func NewCSRIteractor(
	csrRepo repository.CSRRepository,
	dbRepo repository.DBRepository,
) CSRIteractor {
	return &csrInteractor{
		CSRRepository: csrRepo,
		DBRepository:  dbRepo,
	}
}

// SubmitCSR parses a PEM CSR and persists it in pending state.
func (ci *csrInteractor) SubmitCSR(pem string, requesterEmail string) (*model.CSR, error) {
	der, err := decodeCSRPEM(pem)
	if err != nil {
		return nil, errdomain.ErrCSRInvalid
	}

	csrX509, err := x509.ParseCertificateRequest(der)
	if err != nil {
		return nil, errdomain.ErrCSRInvalid
	}

	now := nowFunc()
	domainCSR := &model.CSR{
		Subject:             csrX509.Subject.String(),
		Pem:                 pem,
		KeyAlgorithm:       csrX509.PublicKeyAlgorithm.String(),
		SignatureAlgorithm: csrX509.SignatureAlgorithm.String(),
		DNSNames:           csrX509.DNSNames,
		IPAddresses:        formatIPAddresses(csrX509.IPAddresses),
		RequesterEmail:     requesterEmail,
		Status:             model.CSRStatusPending,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return ci.CSRRepository.Create(domainCSR)
}

// ApproveCSR transitions a CSR from pending to approved.
func (ci *csrInteractor) ApproveCSR(id uint, approverID uint) (*model.CSR, error) {
	data, err := ci.DBRepository.Transaction(func(i interface{}) (interface{}, error) {
		csr, err := ci.CSRRepository.FindByID(id)
		if err != nil {
			return nil, err
		}
		if csr.Status != model.CSRStatusPending {
			return nil, errdomain.ErrCSRNotApproved
		}

		now := nowFunc()
		csr.Status = model.CSRStatusApproved
		csr.ApprovedAt = &now
		csr.ApproverID = &approverID
		csr.UpdatedAt = now
		return ci.CSRRepository.Update(csr)
	})

	if err != nil {
		return nil, err
	}
	csr, ok := data.(*model.CSR)
	if !ok {
		return nil, errors.New("interactor: cast error")
	}
	return csr, nil
}

// RejectCSR transitions a CSR from pending to rejected.
func (ci *csrInteractor) RejectCSR(id uint, notes string) (*model.CSR, error) {
	data, err := ci.DBRepository.Transaction(func(i interface{}) (interface{}, error) {
		csr, err := ci.CSRRepository.FindByID(id)
		if err != nil {
			return nil, err
		}
		if csr.Status != model.CSRStatusPending {
			return nil, errdomain.ErrCSRRejected
		}

		now := nowFunc()
		csr.Status = model.CSRStatusRejected
		csr.RejectedAt = &now
		csr.Notes = notes
		csr.UpdatedAt = now
		return ci.CSRRepository.Update(csr)
	})

	if err != nil {
		return nil, err
	}
	csr, ok := data.(*model.CSR)
	if !ok {
		return nil, errors.New("interactor: cast error")
	}
	return csr, nil
}

// GetCSR retrieves a single CSR by ID.
func (ci *csrInteractor) GetCSR(id uint) (*model.CSR, error) {
	return ci.CSRRepository.FindByID(id)
}

// ListPendingCSRs returns all CSRs awaiting approval.
func (ci *csrInteractor) ListPendingCSRs() ([]*model.CSR, error) {
	return ci.CSRRepository.FindPending()
}

// ListAllCSRs returns all CSRs.
func (ci *csrInteractor) ListAllCSRs() ([]*model.CSR, error) {
	return ci.CSRRepository.FindAll()
}

// ---------------------------------------------------------------------------
// Private helpers

func decodeCSRPEM(pemStr string) ([]byte, error) {
	// Stub — real impl: infrastructure/crypto/pem.go
	return nil, errors.New("not implemented")
}

func formatIPAddresses(ips []net.IP) []string {
	out := make([]string, len(ips))
	for i, ip := range ips {
		out[i] = ip.String()
	}
	return out
}

// nowFunc is injectable for testability.
var nowFunc = func() time.Time { return time.Now() }
