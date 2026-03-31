package controller

import (
	"net/http"
	"strconv"

	"github.com/your-org/x509-clean-architecture/domain/model"
	"github.com/your-org/x509-clean-architecture/usecase/interactor"
)

// certificateController handles HTTP requests for certificate operations.
// It depends only on the use-case-layer interactor interface.
type certificateController struct {
	certificateInteractor interactor.CertificateInteractor
}

// CertificateController defines the HTTP contract for certificate endpoints.
type CertificateController interface {
	GetCertificates(c Context) error
	GetCertificate(c Context) error
	ImportCertificate(c Context) error
	DeleteCertificate(c Context) error
	RevokeCertificate(c Context) error
	GetExpiringCertificates(c Context) error
	ValidateCertificate(c Context) error
}

// NewCertificateController constructs a certificate controller wired to its interactor.
func NewCertificateController(ci interactor.CertificateInteractor) CertificateController {
	return &certificateController{ci}
}

// GetCertificates returns all certificates.
// GET /certificates
func (cc *certificateController) GetCertificates(c Context) error {
	certs, err := cc.certificateInteractor.ListCertificates()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, certs)
}

// GetCertificate returns a single certificate by ID.
// GET /certificates/:id
func (cc *certificateController) GetCertificate(c Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	cert, err := cc.certificateInteractor.GetCertificate(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cert)
}

// ImportCertificateRequest is the HTTP body for importing a certificate.
type ImportCertificateRequest struct {
	CertPEM string `json:"cert_pem"`
	KeyPEM  string `json:"key_pem,omitempty"`
}

// ImportCertificate parses and imports a PEM-encoded certificate.
// POST /certificates
func (cc *certificateController) ImportCertificate(c Context) error {
	var req ImportCertificateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	cert, err := cc.certificateInteractor.ImportCertificate(req.CertPEM, req.KeyPEM)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, cert)
}

// DeleteCertificate soft-deletes a certificate.
// DELETE /certificates/:id
func (cc *certificateController) DeleteCertificate(c Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	if err := cc.certificateInteractor.DeleteCertificate(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusNoContent, nil)
}

// RevokeCertificateRequest is the HTTP body for revoking a certificate.
type RevokeCertificateRequest struct {
	Reason string `json:"reason"`
}

// RevokeCertificate marks a certificate as revoked.
// POST /certificates/:id/revoke
func (cc *certificateController) RevokeCertificate(c Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req RevokeCertificateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	cert, err := cc.certificateInteractor.RevokeCertificate(uint(id), req.Reason)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cert)
}

// GetExpiringCertificates returns certificates expiring within a query-param window.
// GET /certificates/expiring?days=30
func (cc *certificateController) GetExpiringCertificates(c Context) error {
	days := 30 // default
	if d := c.QueryParam("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil {
			days = parsed
		}
	}

	certs, err := cc.certificateInteractor.GetExpiringCertificates(days)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, certs)
}

// ValidateCertificateRequest is the HTTP body for validating a certificate.
type ValidateCertificateRequest struct {
	CertPEM string `json:"cert_pem,omitempty"`
}

// ValidateCertificate validates a certificate by ID or inline PEM.
// POST /certificates/validate
func (cc *certificateController) ValidateCertificate(c Context) error {
	var req ValidateCertificateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	idStr := c.Param("id")
	if idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
		}
		cert, err := cc.certificateInteractor.ValidateCertificate(uint(id))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, cert)
	}

	return c.JSON(http.StatusBadRequest, map[string]string{"error": "id or cert_pem required"})
}
