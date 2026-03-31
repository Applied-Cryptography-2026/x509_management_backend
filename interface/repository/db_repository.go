package repository

import (
	"log"

	"github.com/your-org/x509-clean-architecture/usecase/repository"
	"gorm.io/gorm"
)

// dbRepository is the concrete driven adapter for DBRepository.
// It wraps GORM transactions with a closure-based API matching the use-case port.
type dbRepository struct {
	db *gorm.DB
}

// NewDBRepository constructs the transaction-aware DB repository.
func NewDBRepository(db *gorm.DB) repository.DBRepository {
	return &dbRepository{db: db}
}

// Transaction executes fn inside a GORM database transaction.
// On any error returned by fn, the transaction is rolled back.
// On panic from fn, the transaction is also rolled back.
// On success, the transaction is committed.
func (r *dbRepository) Transaction(
	fn func(interface{}) (interface{}, error),
) (data interface{}, err error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if p := recover(); p != nil {
			log.Printf("db_repository: recovered panic during transaction: %v", p)
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Printf("db_repository: rolling back transaction: %v", err)
			tx.Rollback()
		} else {
			if commitErr := tx.Commit().Error; commitErr != nil {
				err = commitErr
			}
		}
	}()

	data, err = fn(tx)
	return data, err
}
