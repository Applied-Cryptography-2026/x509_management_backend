package errdomain

import "errors"

// Domain-level sentinel errors.
// These are defined here so the domain layer never imports external packages.

var (
	// Certificate errors
	ErrCertNotFound       = errors.New("certificate not found")
	ErrCertAlreadyExists   = errors.New("certificate already exists")
	ErrCertExpired         = errors.New("certificate has expired")
	ErrCertNotYetValid     = errors.New("certificate is not yet valid")
	ErrCertRevoked         = errors.New("certificate has been revoked")
	ErrCertInvalid         = errors.New("certificate is invalid")
	ErrCertKeyMismatch     = errors.New("certificate does not match its private key")
	ErrCertChainBroken     = errors.New("certificate chain is broken")
	ErrCertParseFailed     = errors.New("failed to parse certificate")
	ErrCertUnsupportedAlgo = errors.New("unsupported signature algorithm")

	// CSR errors
	ErrCSRNotFound       = errors.New("CSR not found")
	ErrCSRAlreadyExists  = errors.New("CSR already exists")
	ErrCSRInvalid        = errors.New("CSR is invalid")
	ErrCSRRejected       = errors.New("CSR was rejected")
	ErrCSRNotApproved    = errors.New("CSR is not approved")

	// CA errors
	ErrCANotFound      = errors.New("CA certificate not found")
	ErrCANotCA        = errors.New("certificate is not a CA")
	ErrCANotSelfSigned = errors.New("CA certificate must be self-signed or signed by a parent CA")

	// Chain errors
	ErrChainIncomplete = errors.New("certificate chain is incomplete")
	ErrChainUntrusted  = errors.New("certificate chain is not trusted")

	// Key errors
	ErrKeyNotFound      = errors.New("private key not found")
	ErrKeyDecryptFailed = errors.New("failed to decrypt private key")
	ErrKeyAlgorithm     = errors.New("incompatible key algorithm")
)
