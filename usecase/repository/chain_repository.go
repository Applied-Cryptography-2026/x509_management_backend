package repository

import (
	"github.com/your-org/x509-clean-architecture/domain/model"
)

// ChainRepository is the use-case-layer interface for certificate chain persistence.
type ChainRepository interface {
	FindAll() ([]*model.CertificateChain, error)
	FindByID(id uint) (*model.CertificateChain, error)
	FindByName(name string) (*model.CertificateChain, error)
	FindByType(chainType model.ChainType) ([]*model.CertificateChain, error)
	Create(chain *model.CertificateChain) (*model.CertificateChain, error)
	Update(chain *model.CertificateChain) (*model.CertificateChain, error)
	Delete(id uint) error
}
