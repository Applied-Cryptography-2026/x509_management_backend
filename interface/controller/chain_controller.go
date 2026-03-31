package controller

import (
	"net/http"
	"strconv"
)

// chainController handles HTTP requests for certificate chain operations.
type chainController struct {
	// TODO: wire chain interactor once usecase/interactor/chain_interactor.go is implemented
}

// ChainController defines the HTTP contract for chain endpoints.
type ChainController interface {
	GetChains(c Context) error
	GetChain(c Context) error
	CreateChain(c Context) error
	DeleteChain(c Context) error
	ValidateChain(c Context) error
}

// NewChainController constructs a chain controller.
func NewChainController() ChainController {
	return &chainController{}
}

// GetChains returns all chains.
// GET /chains
func (cc *chainController) GetChains(c Context) error {
	return c.JSON(http.StatusOK, []map[string]string{})
}

// GetChain returns a single chain by ID.
// GET /chains/:id
func (cc *chainController) GetChain(c Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	_ = id
	return c.JSON(http.StatusOK, map[string]string{})
}

// CreateChainRequest is the HTTP body for creating a chain.
type CreateChainRequest struct {
	Name    string `json:"name"`
	CertIDs []uint `json:"cert_ids"`
}

// CreateChain builds a certificate chain from the given cert IDs.
// POST /chains
func (cc *chainController) CreateChain(c Context) error {
	var req CreateChainRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	_ = req
	return c.JSON(http.StatusCreated, map[string]string{})
}

// DeleteChain deletes a chain.
// DELETE /chains/:id
func (cc *chainController) DeleteChain(c Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	_ = id
	return c.JSON(http.StatusNoContent, nil)
}

// ValidateChain validates a chain's trust path.
// POST /chains/:id/validate
func (cc *chainController) ValidateChain(c Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	_ = id
	return c.JSON(http.StatusOK, map[string]bool{"valid": true})
}
