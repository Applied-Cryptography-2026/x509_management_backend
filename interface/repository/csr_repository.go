package repository

import (
	"github.com/your-org/x509-clean-architecture/domain/model"
	"github.com/your-org/x509-clean-architecture/usecase/repository"
	"gorm.io/gorm"
)

// csrRepository is the concrete driven adapter for CSRRepository.
type csrRepository struct {
	db *gorm.DB
}

// NewCSRRepository constructs the concrete CSR repository.
func NewCSRRepository(db *gorm.DB) repository.CSRRepository {
	return &csrRepository{db: db}
}

func (r *csrRepository) FindAll() ([]*model.CSR, error) {
	var csrs []*model.CSR
	err := r.db.Find(&csrs).Error
	return csrs, err
}

func (r *csrRepository) FindByID(id uint) (*model.CSR, error) {
	var csr model.CSR
	err := r.db.First(&csr, id).Error
	return &csr, err
}

func (r *csrRepository) FindBySubject(subject string) ([]*model.CSR, error) {
	var csrs []*model.CSR
	err := r.db.Where("subject LIKE ?", "%"+subject+"%").Find(&csrs).Error
	return csrs, err
}

func (r *csrRepository) FindByStatus(status model.CSRStatus) ([]*model.CSR, error) {
	var csrs []*model.CSR
	err := r.db.Where("status = ?", status).Find(&csrs).Error
	return csrs, err
}

func (r *csrRepository) FindPending() ([]*model.CSR, error) {
	var csrs []*model.CSR
	err := r.db.Where("status = ?", model.CSRStatusPending).Find(&csrs).Error
	return csrs, err
}

func (r *csrRepository) Create(csr *model.CSR) (*model.CSR, error) {
	err := r.db.Create(csr).Error
	return csr, err
}

func (r *csrRepository) Update(csr *model.CSR) (*model.CSR, error) {
	err := r.db.Save(csr).Error
	return csr, err
}

func (r *csrRepository) Delete(id uint) error {
	return r.db.Delete(&model.CSR{}, id).Error
}
