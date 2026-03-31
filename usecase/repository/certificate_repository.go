package repository

import (
	"github.com/your-org/x509-clean-architecture/domain/model"
)

// CertificateRepository is a use-case-layer interface (input port).
// It defines what the interactor needs from the persistence layer.
// The concrete implementation lives in interface/repository/ (driven adapter).
type CertificateRepository interface {
	// FindAll returns all certificates.
	FindAll() ([]*model.Certificate, error)

	// FindByID returns a certificate by its ID.
	FindByID(id uint) (*model.Certificate, error)

	// FindBySerial returns a certificate by its serial number.
	FindBySerial(serial string) (*model.Certificate, error)

	// FindByFingerprint returns a certificate by its SHA-256 fingerprint.
	FindByFingerprint(fingerprint string) (*model.Certificate, error)

	// FindBySubject returns all certificates matching a subject string.
	FindBySubject(subject string) ([]*model.Certificate, error)

	// FindByIssuer returns all certificates issued by a given issuer.
	FindByIssuer(issuer string) ([]*model.Certificate, error)

	// FindByStatus returns all certificates with a given status.
	FindByStatus(status model.CertStatus) ([]*model.Certificate, error)

	// FindByProfile returns all certificates matching a profile (e.g. "tls-server").
	FindByProfile(profile string) ([]*model.Certificate, error)

	// FindExpiring returns all certificates expiring within the given duration.
	FindExpiring(withinDays int) ([]*model.Certificate, error)

	// FindRevoked returns all revoked certificates.
	FindRevoked() ([]*model.Certificate, error)

	// Create persists a new certificate and returns the saved entity.
	Create(cert *model.Certificate) (*model.Certificate, error)

	// Update persists changes to an existing certificate.
	Update(cert *model.Certificate) (*model.Certificate, error)

	// Delete soft-deletes a certificate (sets DeletedAt).
	Delete(id uint) error
}
