package controller

// AppController is the top-level container holding all domain controllers.
// It is the single injection point used by the router.
// Each field is typed as an anonymous interface (structural typing) so the concrete
// struct satisfies it implicitly without a named interface declaration.
type AppController struct {
	Certificate interface {
		GetCertificates(c Context) error
		GetCertificate(c Context) error
		ImportCertificate(c Context) error
		DeleteCertificate(c Context) error
		RevokeCertificate(c Context) error
		GetExpiringCertificates(c Context) error
		ValidateCertificate(c Context) error
	}
	CSR interface {
		GetCSRs(c Context) error
		GetCSR(c Context) error
		SubmitCSR(c Context) error
		ApproveCSR(c Context) error
		RejectCSR(c Context) error
	}
	Chain interface {
		GetChains(c Context) error
		GetChain(c Context) error
		CreateChain(c Context) error
		DeleteChain(c Context) error
		ValidateChain(c Context) error
	}
}
