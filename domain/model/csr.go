package model

import (
	"crypto/x509"
	"time"
)

// CSR represents a Certificate Signing Request in the domain layer.
type CSR struct {
	ID                 uint     `json:"id" gorm:"primaryKey"`
	Subject            string   `json:"subject"`             // CN requested by the requester
	Pem                string   `json:"pem"`                 // PEM-encoded CSR
	KeyAlgorithm       string   `json:"key_algorithm"`       // RSA, ECDSA, Ed25519
	SignatureAlgorithm string   `json:"signature_algorithm"` //EX:SHA256 and RSA
	DNSNames           []string `json:"dns_names"`
	IPAddresses        []string `json:"ip_addresses"`
	// RequesterEmail     string     `json:"requester_email"`
	Status     CSRStatus  `json:"status"` // pending, approved, rejected, issued
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	RejectedAt *time.Time `json:"rejected_at,omitempty"`
	ApproverID *uint      `json:"approver_id,omitempty"`
	// Notes              string     `json:"notes,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// CSRStatus is the lifecycle state of a CSR.
type CSRStatus string

const (
	CSRStatusPending  CSRStatus = "pending"
	CSRStatusApproved CSRStatus = "approved"
	CSRStatusRejected CSRStatus = "rejected"
	CSRStatusIssued   CSRStatus = "issued"
)

// TableName overrides GORM's default table name inference.
func (CSR) TableName() string {
	return "csrs"
}
