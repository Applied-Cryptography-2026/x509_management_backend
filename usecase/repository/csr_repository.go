package repository

import (
	"github.com/your-org/x509-clean-architecture/domain/model"
)

// CSRRepository is the use-case-layer interface for CSR persistence (input port).
type CSRRepository interface {
	FindAll() ([]*model.CSR, error)
	FindByID(id uint) (*model.CSR, error)
	FindBySubject(subject string) ([]*model.CSR, error)
	FindByStatus(status model.CSRStatus) ([]*model.CSR, error)
	FindPending() ([]*model.CSR, error)
	Create(csr *model.CSR) (*model.CSR, error)
	Update(csr *model.CSR) (*model.CSR, error)
	Delete(id uint) error
}
