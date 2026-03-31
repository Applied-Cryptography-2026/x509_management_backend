package registry

import (
	"github.com/your-org/x509-clean-architecture/interface/controller"
	ip "github.com/your-org/x509-clean-architecture/interface/presenter"
	ir "github.com/your-org/x509-clean-architecture/interface/repository"
	"github.com/your-org/x509-clean-architecture/usecase/interactor"
	"gorm.io/gorm"
)

// ---- Certificate wire ------------------------------------------------

// newCertificateInteractor assembles the certificate interactor.
func (r *registry) newCertificateInteractor() interactor.CertificateInteractor {
	return interactor.NewCertificateInteractor(
		ir.NewCertificateRepository(r.db),
		ir.NewChainRepository(r.db),
		ip.NewCertificatePresenter(),
		ir.NewDBRepository(r.db),
	)
}

// newCertificateController assembles the certificate controller.
func (r *registry) newCertificateController() controller.CertificateController {
	return controller.NewCertificateController(
		r.newCertificateInteractor(),
	)
}

// NewCertificateController is the public accessor used by NewAppController.
func (r *registry) NewCertificateController() controller.CertificateController {
	return r.newCertificateController()
}

// ---- CSR wire --------------------------------------------------------

// newCSRIteractor assembles the CSR interactor.
func (r *registry) newCSRIteractor() interactor.CSRIteractor {
	return interactor.NewCSRIteractor(
		ir.NewCSRRepository(r.db),
		ir.NewDBRepository(r.db),
	)
}

// newCSRController assembles the CSR controller.
func (r *registry) newCSRController() controller.CSRIteractorController {
	return controller.NewCSRIteractorController(
		r.newCSRIteractor(),
	)
}

// NewCSRController is the public accessor used by NewAppController.
func (r *registry) NewCSRController() controller.CSRIteractorController {
	return r.newCSRController()
}

// ---- Chain wire ------------------------------------------------------

// newChainController assembles the chain controller.
func (r *registry) newChainController() controller.ChainController {
	return controller.NewChainController()
}

// NewChainController is the public accessor used by NewAppController.
func (r *registry) NewChainController() controller.ChainController {
	return r.newChainController()
}
