package repository

import (
	"github.com/your-org/x509-clean-architecture/domain/model"
	"github.com/your-org/x509-clean-architecture/usecase/repository"
	"gorm.io/gorm"
)

// certificateRepository is the concrete driven adapter for CertificateRepository.
// It implements the use-case-layer interface using GORM.
type certificateRepository struct {
	db *gorm.DB
}

// NewCertificateRepository constructs the concrete repository.
// It returns the use-case-layer interface type, not the concrete type.
func NewCertificateRepository(db *gorm.DB) repository.CertificateRepository {
	return &certificateRepository{db: db}
}

func (r *certificateRepository) FindAll() ([]*model.Certificate, error) {
	var certs []*model.Certificate
	err := r.db.Find(&certs).Error
	return certs, err
}

func (r *certificateRepository) FindByID(id uint) (*model.Certificate, error) {
	var cert model.Certificate
	err := r.db.First(&cert, id).Error
	return &cert, err
}

func (r *certificateRepository) FindBySerial(serial string) (*model.Certificate, error) {
	var cert model.Certificate
	err := r.db.Where("serial = ?", serial).First(&cert).Error
	return &cert, err
}

func (r *certificateRepository) FindByFingerprint(fingerprint string) (*model.Certificate, error) {
	var cert model.Certificate
	err := r.db.Where("fingerprint = ?", fingerprint).First(&cert).Error
	return &cert, err
}

func (r *certificateRepository) FindBySubject(subject string) ([]*model.Certificate, error) {
	var certs []*model.Certificate
	err := r.db.Where("subject LIKE ?", "%"+subject+"%").Find(&certs).Error
	return certs, err
}

func (r *certificateRepository) FindByIssuer(issuer string) ([]*model.Certificate, error) {
	var certs []*model.Certificate
	err := r.db.Where("issuer LIKE ?", "%"+issuer+"%").Find(&certs).Error
	return certs, err
}

func (r *certificateRepository) FindByStatus(status model.CertStatus) ([]*model.Certificate, error) {
	var certs []*model.Certificate
	err := r.db.Where("status = ?", status).Find(&certs).Error
	return certs, err
}

func (r *certificateRepository) FindByProfile(profile string) ([]*model.Certificate, error) {
	var certs []*model.Certificate
	err := r.db.Where("profile = ?", profile).Find(&certs).Error
	return certs, err
}

func (r *certificateRepository) FindExpiring(withinDays int) ([]*model.Certificate, error) {
	var certs []*model.Certificate
	err := r.db.
		Where("not_after BETWEEN ? AND ?",
			gorm.Expr("NOW()"),
			gorm.Expr("DATE_ADD(NOW(), INTERVAL ? DAY)", withinDays),
		).
		Where("status = ?", model.CertStatusActive).
		Find(&certs).Error
	return certs, err
}

func (r *certificateRepository) FindRevoked() ([]*model.Certificate, error) {
	var certs []*model.Certificate
	err := r.db.Where("is_revoked = ?", true).Find(&certs).Error
	return certs, err
}

func (r *certificateRepository) Create(cert *model.Certificate) (*model.Certificate, error) {
	err := r.db.Create(cert).Error
	return cert, err
}

func (r *certificateRepository) Update(cert *model.Certificate) (*model.Certificate, error) {
	err := r.db.Save(cert).Error
	return cert, err
}

func (r *certificateRepository) Delete(id uint) error {
	return r.db.Delete(&model.Certificate{}, id).Error
}
