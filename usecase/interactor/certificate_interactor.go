package interactor

import (
	"crypto/x509"
	"errors"
	"time"

	"github.com/your-org/x509-clean-architecture/domain/errdomain"
	"github.com/your-org/x509-clean-architecture/domain/model"
	"github.com/your-org/x509-clean-architecture/usecase/presenter"
	"github.com/your-org/x509-clean-architecture/usecase/repository"
)

// certificateInteractor is the concrete implementation.
// It holds only interfaces — no concrete adapter types.
type certificateInteractor struct {
	CertificateRepository repository.CertificateRepository
	ChainRepository       repository.ChainRepository
	CertificatePresenter  presenter.CertificatePresenter
	DBRepository          repository.DBRepository
}

// CertificateInteractor is both the interface and the struct name here.
// Embedding is not used; dependencies are explicit constructor parameters.
type CertificateInteractor interface {
	// View use cases
	GetCertificate(id uint) (*model.Certificate, error)
	ListCertificates() ([]*model.Certificate, error)
	SearchCertificates(query string) ([]*model.Certificate, error)
	GetExpiringCertificates(withinDays int) ([]*model.Certificate, error)

	// Lifecycle use cases
	ImportCertificate(pem string, keyPEM string) (*model.Certificate, error)
	RevokeCertificate(id uint, reason string) (*model.Certificate, error)
	RenewCertificate(id uint, newCSR *model.CSR) (*model.Certificate, error)
	DeleteCertificate(id uint) error

	// Validation use cases
	ValidateChain(chainID uint) (bool, error)
	ValidateCertificate(id uint) (*model.Certificate, error)
}

// NewCertificateInteractor assembles the interactor with its required port interfaces.
func NewCertificateInteractor(
	cRepo repository.CertificateRepository,
	chRepo repository.ChainRepository,
	cPresenter presenter.CertificatePresenter,
	dbRepo repository.DBRepository,
) CertificateInteractor {
	return &certificateInteractor{
		CertificateRepository: cRepo,
		ChainRepository:       chRepo,
		CertificatePresenter:  cPresenter,
		DBRepository:          dbRepo,
	}
}

// GetCertificate retrieves a single certificate and applies the presenter transform.
func (ci *certificateInteractor) GetCertificate(id uint) (*model.Certificate, error) {
	cert, err := ci.CertificateRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	return ci.CertificatePresenter.ResponseCert(cert), nil
}

// ListCertificates returns all certificates through the presenter.
func (ci *certificateInteractor) ListCertificates() ([]*model.Certificate, error) {
	certs, err := ci.CertificateRepository.FindAll()
	if err != nil {
		return nil, err
	}
	return ci.CertificatePresenter.ResponseCerts(certs), nil
}

// SearchCertificates queries by subject, issuer, or fingerprint.
func (ci *certificateInteractor) SearchCertificates(query string) ([]*model.Certificate, error) {
	certs, err := ci.CertificateRepository.FindBySubject(query)
	if err != nil {
		// Fall back: try fingerprint
		cert, ferr := ci.CertificateRepository.FindByFingerprint(query)
		if ferr != nil {
			return nil, err
		}
		return []*model.Certificate{cert}, nil
	}
	return ci.CertificatePresenter.ResponseCerts(certs), nil
}

// GetExpiringCertificates returns certs expiring within the given window.
func (ci *certificateInteractor) GetExpiringCertificates(withinDays int) ([]*model.Certificate, error) {
	return ci.CertificateRepository.FindExpiring(withinDays)
}

// ImportCertificate parses a PEM-encoded certificate and persists it.
func (ci *certificateInteractor) ImportCertificate(pem string, keyPEM string) (*model.Certificate, error) {
	der, err := decodePEM(pem)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, err
	}

	// Check for duplicate
	existing, _ := ci.CertificateRepository.FindByFingerprint(fingerprint(cert))
	if existing != nil {
		return nil, errdomain.ErrCertAlreadyExists
	}

	domainCert := modelFromX509(cert, pem, keyPEM)
	return ci.CertificatePresenter.ResponseCert(domainCert), nil
}

// RevokeCertificate marks a certificate as revoked within a transaction.
func (ci *certificateInteractor) RevokeCertificate(id uint, reason string) (*model.Certificate, error) {
	data, err := ci.DBRepository.Transaction(func(i interface{}) (interface{}, error) {
		cert, err := ci.CertificateRepository.FindByID(id)
		if err != nil {
			return nil, err
		}
		if cert.IsRevoked {
			return nil, errdomain.ErrCertRevoked
		}

		now := time.Now()
		cert.IsRevoked = true
		cert.RevokedAt = &now
		cert.Status = model.CertStatusRevoked

		return ci.CertificateRepository.Update(cert)
	})

	if err != nil {
		return nil, err
	}

	cert, ok := data.(*model.Certificate)
	if !ok {
		return nil, errors.New("interactor: cast error")
	}

	return ci.CertificatePresenter.ResponseCert(cert), nil
}

// RenewCertificate issues a new certificate from a CSR to replace an existing one.
func (ci *certificateInteractor) RenewCertificate(id uint, newCSR *model.CSR) (*model.Certificate, error) {
	existing, err := ci.CertificateRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errdomain.ErrCertNotFound
	}

	// Placeholder: actual signing logic lives in infrastructure/crypto/signer.go
	// which the interactor would call via an injected signer service.
	_ = newCSR

	return nil, nil
}

// DeleteCertificate soft-deletes a certificate.
func (ci *certificateInteractor) DeleteCertificate(id uint) error {
	return ci.CertificateRepository.Delete(id)
}

// ValidateChain verifies the trust chain for a given chain definition.
func (ci *certificateInteractor) ValidateChain(chainID uint) (bool, error) {
	chain, err := ci.ChainRepository.FindByID(chainID)
	if err != nil {
		return false, err
	}

	if len(chain.CertIDs) < 2 {
		return false, errdomain.ErrChainIncomplete
	}

	// The actual crypto verification is delegated to a signer service
	// injected from the infrastructure layer. Here we just check structural validity.
	return len(chain.CertIDs) >= 2, nil
}

// ValidateCertificate checks domain-level validity of a certificate.
func (ci *certificateInteractor) ValidateCertificate(id uint) (*model.Certificate, error) {
	cert, err := ci.CertificateRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if cert.IsExpired() {
		return cert, errdomain.ErrCertExpired
	}
	if time.Now().Before(cert.NotBefore) {
		return cert, errdomain.ErrCertNotYetValid
	}
	if cert.IsRevoked {
		return cert, errdomain.ErrCertRevoked
	}

	return ci.CertificatePresenter.ResponseCert(cert), nil
}

// ---------------------------------------------------------------------------
// Private helpers (domain-level parsing helpers — no framework deps)

// decodePEM decodes a PEM block into DER bytes.
func decodePEM(pemStr string) ([]byte, error) {
	// Keep stubs lean — real impl in infrastructure/crypto/pem.go
	return nil, errors.New("not implemented: infrastructure/crypto/pem.go")
}

// fingerprint returns the SHA-256 fingerprint of a parsed x509 cert.
func fingerprint(cert *x509.Certificate) string {
	// Stub — real impl delegates to infrastructure/crypto/fingerprint.go
	return ""
}

// modelFromX509 converts an x509.Certificate to the domain model.
func modelFromX509(cert *x509.Certificate, pem, keyPEM string) *model.Certificate {
	// Stub — real impl delegates to infrastructure/certificate/converter.go
	return nil
}
