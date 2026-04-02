package controller

import (
	"net/http"
	"strconv"

	"github.com/your-org/x509-clean-architecture/domain/model"
	"github.com/your-org/x509-clean-architecture/usecase/interactor"
)

// csrController handles HTTP requests for CSR operations.
type csrController struct {
	csrInteractor interactor.CSRIteractor
}

// CSRIteractorController defines the HTTP contract for CSR endpoints.
type CSRIteractorController interface {
	GetCSRs(c Context) error
	GetCSR(c Context) error
	SubmitCSR(c Context) error
	ApproveCSR(c Context) error
	RejectCSR(c Context) error
}

// NewCSRIteractorController constructs a CSR controller wired to its interactor.
func NewCSRIteractorController(ci interactor.CSRIteractor) CSRIteractorController {
	return &csrController{ci}
}

// GetCSRs returns all CSRs, optionally filtered by status query param.
// GET /csrs?status=pending
func (sc *csrController) GetCSRs(c Context) error {
	status := c.QueryParam("status")

	var csrs []*model.CSR
	var err error

	if status == "pending" {
		csrs, err = sc.csrInteractor.ListPendingCSRs()
	} else {
		csrs, err = sc.csrInteractor.ListAllCSRs()
	}

	_ = csrs
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, csrs)
}

// GetCSR returns a single CSR by ID.
// GET /csrs/:id
func (sc *csrController) GetCSR(c Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	csr, err := sc.csrInteractor.GetCSR(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, csr)
}

// SubmitCSRRequest is the HTTP body for submitting a CSR.
type SubmitCSRRequest struct {
	Pem            string `json:"pem"`
	RequesterEmail string `json:"requester_email"`
}

// SubmitCSR parses and submits a PEM-encoded CSR.
// POST /csrs
func (sc *csrController) SubmitCSR(c Context) error {
	var req SubmitCSRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	csr, err := sc.csrInteractor.SubmitCSR(req.Pem, req.RequesterEmail)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, csr)
}

// ApproveCSRRequest is the HTTP body for approving a CSR.
type ApproveCSRRequest struct {
	ApproverID uint `json:"approver_id"`
}

// ApproveCSR transitions a CSR to approved status.
// POST /csrs/:id/approve
func (sc *csrController) ApproveCSR(c Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req ApproveCSRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	csr, err := sc.csrInteractor.ApproveCSR(uint(id), req.ApproverID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, csr)
}

// RejectCSRRequest is the HTTP body for rejecting a CSR.
type RejectCSRRequest struct {
	Notes string `json:"notes"`
}

// RejectCSR transitions a CSR to rejected status.
// POST /csrs/:id/reject
func (sc *csrController) RejectCSR(c Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	var req RejectCSRRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	csr, err := sc.csrInteractor.RejectCSR(uint(id), req.Notes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, csr)
}
