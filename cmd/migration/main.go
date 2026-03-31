package main

import (
	"log"

	"github.com/your-org/x509-clean-architecture/config"
	"github.com/your-org/x509-clean-architecture/domain/model"
	"github.com/your-org/x509-clean-architecture/infrastructure/datastore"
)

func main() {
	config.ReadConfig()

	db, err := datastore.NewDB()
	if err != nil {
		log.Fatalln("migration: failed to connect to database:", err)
	}
	defer func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}()

	// AutoMigrate runs GORM auto-migration.
	// For production use, prefer explicit SQL migration files (e.g., goose).
	err = db.AutoMigrate(
		&model.Certificate{},
		&model.CSR{},
		&model.CertificateChain{},
	)
	if err != nil {
		log.Fatalln("migration: AutoMigrate failed:", err)
	}

	log.Println("migration: all tables migrated successfully")
}
