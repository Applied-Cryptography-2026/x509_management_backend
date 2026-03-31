package registry

import (
	"github.com/your-org/x509-clean-architecture/interface/controller"
	"gorm.io/gorm"
)

// registry is the composition root — all concrete types are instantiated here.
// It follows the manual DI pattern from golang-clean-architecture by manakuro.
type registry struct {
	db *gorm.DB
}

// Registry defines the top-level interface for the composition root.
type Registry interface {
	NewAppController() controller.AppController
}

// NewRegistry creates a new registry wired to the given database instance.
func NewRegistry(db *gorm.DB) Registry {
	return &registry{db: db}
}

// NewAppController assembles and returns the top-level AppController.
// All dependencies are assembled bottom-up here:
//   1. Infrastructure adapters (repositories) are created first
//   2. Passed into interactors (use cases)
//   3. Passed into controllers (HTTP handlers)
//   4. Collected into AppController
func (r *registry) NewAppController() controller.AppController {
	return controller.AppController{
		Certificate: r.NewCertificateController(),
		CSR:        r.NewCSRController(),
		Chain:      r.NewChainController(),
	}
}
