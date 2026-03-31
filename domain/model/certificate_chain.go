package model

import "time"

// CertificateChain represents a chain of trust
// end-user-cert → intermediate → root.
type CertificateChain struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name"`       // Human-readable name for the chain
	ChainType ChainType  `json:"chain_type"` // tls-server, tls-client, etc.
	CertIDs   []uint     `json:"cert_ids"`   // Ordered list of certificate IDs (leaf first)
	PemChain  string     `json:"pem_chain"`  // Concatenated PEM chain
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// ChainType describes the purpose of a certificate chain.
type ChainType string

const (
	ChainTypeTLSServer   ChainType = "tls-server"
	ChainTypeTLSClient   ChainType = "tls-client"
	ChainTypeCodeSigning ChainType = "code-signing"
	ChainTypeSMIME       ChainType = "smime"
	ChainTypeCustom      ChainType = "custom"
)

func (CertificateChain) TableName() string {
	return "certificate_chains"
}
