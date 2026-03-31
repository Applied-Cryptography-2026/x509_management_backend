package repository

import (
	"github.com/your-org/x509-clean-architecture/domain/model"
	"github.com/your-org/x509-clean-architecture/usecase/repository"
	"gorm.io/gorm"
)

// chainRepository is the concrete driven adapter for ChainRepository.
type chainRepository struct {
	db *gorm.DB
}

// NewChainRepository constructs the concrete chain repository.
func NewChainRepository(db *gorm.DB) repository.ChainRepository {
	return &chainRepository{db: db}
}

func (r *chainRepository) FindAll() ([]*model.CertificateChain, error) {
	var chains []*model.CertificateChain
	err := r.db.Find(&chains).Error
	return chains, err
}

func (r *chainRepository) FindByID(id uint) (*model.CertificateChain, error) {
	var chain model.CertificateChain
	err := r.db.First(&chain, id).Error
	return &chain, err
}

func (r *chainRepository) FindByName(name string) (*model.CertificateChain, error) {
	var chain model.CertificateChain
	err := r.db.Where("name = ?", name).First(&chain).Error
	return &chain, err
}

func (r *chainRepository) FindByType(chainType model.ChainType) ([]*model.CertificateChain, error) {
	var chains []*model.CertificateChain
	err := r.db.Where("chain_type = ?", chainType).Find(&chains).Error
	return chains, err
}

func (r *chainRepository) Create(chain *model.CertificateChain) (*model.CertificateChain, error) {
	err := r.db.Create(chain).Error
	return chain, err
}

func (r *chainRepository) Update(chain *model.CertificateChain) (*model.CertificateChain, error) {
	err := r.db.Save(chain).Error
	return chain, err
}

func (r *chainRepository) Delete(id uint) error {
	return r.db.Delete(&model.CertificateChain{}, id).Error
}
