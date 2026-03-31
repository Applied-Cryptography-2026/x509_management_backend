# x509 Clean Architecture

An opinionated **Go (Golang) project scaffold** for managing the full lifecycle of x509 certificates — issuance, import, renewal, revocation, chain validation, and CRL/OCSP — built on the [golang-clean-architecture](https://github.com/manakuro/golang-clean-architecture) pattern by manakuro.

It is a **framework**, not a running application. Every file declares *what belongs there* and *why*; implement the `TODO` stubs to make it run.

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Dependency Rule](#2-dependency-rule)
3. [Directory Structure](#3-directory-structure)
4. [Domain Layer](#4-domain-layer)
5. [Use Case Layer](#5-use-case-layer)
6. [Interface Layer](#6-interface-layer)
7. [Infrastructure Layer](#7-infrastructure-layer)
8. [Registry — Composition Root](#8-registry--composition-root)
9. [Config](#9-config)
10. [Entry Points](#10-entry-points)
11. [Database Migrations](#11-database-migrations)
12. [CLI Tool](#12-cli-tool)
13. [Domain Entities](#13-domain-entities)
14. [HTTP API Reference](#14-http-api-reference)
15. [Development Workflow](#15-development-workflow)
16. [Next Steps — What to Implement](#16-next-steps--what-to-implement)

---

## 1. Architecture Overview

This project follows **Clean Architecture** (a.k.a. Ports & Adapters, or Hexagonal Architecture) adapted for Go. The codebase is divided into four concentric layers:

```
┌──────────────────────────────────────────────────────────────────────┐
│  infrastructure/           outermost — framework & drivers            │
│  datastore/                MySQL + GORM connection                    │
│  crypto/                   PEM, signing, fingerprinting, validation    │
│  router/                   Echo HTTP wiring + middleware               │
├──────────────────────────────────────────────────────────────────────┤
│  interface/                adapters — driving (controllers) & driven  │
│  controller/               HTTP handlers → interactors               │
│  presenter/                response transformers → controllers        │
│  repository/               GORM implementations → interactors        │
├──────────────────────────────────────────────────────────────────────┤
│  usecase/                  innermost application layer                │
│  interactor/               business logic orchestrators               │
│  repository/               repository interfaces (input ports)       │
│  presenter/                presenter interfaces (output ports)        │
├──────────────────────────────────────────────────────────────────────┤
│  domain/                   innermost — enterprise business rules      │
│  model/                    pure domain entities, no external imports   │
│  errdomain/                sentinel errors, no external imports        │
└──────────────────────────────────────────────────────────────────────┘
```

**Dependency rule:** Outer layers may import inner layers. Inner layers **never** import outer layers. This keeps business logic testable and independent of any framework.

---

## 2. Dependency Rule

```
domain/model/          ← pure Go, zero external imports
    ↓  (interface only)
usecase/repository/     ← imports domain → defines data access ports
usecase/interactor/     ← imports domain + usecase ports → business logic
    ↓
interface/repository/  ← imports usecase ports → GORM/MySQL implementations
interface/controller/  ← imports usecase ports → Echo HTTP handlers
    ↓
infrastructure/        ← imports domain + usecase (GORM, crypto, Echo)
    ↓
registry/              ← composition root: wires everything together
    ↓
cmd/                   ← application entry point
```

The **key principle:** when you import packages, arrows always point inward. `interface/repository` never imports `infrastructure` — it is `infrastructure`.

---

## 3. Directory Structure

```
x509-clean-architecture/
│
├── domain/                          ← innermost layer
│   ├── model/                       ← domain entities
│   │   ├── certificate.go
│   │   ├── csr.go
│   │   └── certificate_chain.go
│   └── errdomain/
│       └── errors.go               ← sentinel errors
│
├── usecase/                        ← application layer
│   ├── repository/                 ← port interfaces
│   │   ├── certificate_repository.go
│   │   ├── csr_repository.go
│   │   ├── chain_repository.go
│   │   └── db_repository.go
│   ├── presenter/                  ← output port interfaces
│   │   └── certificate_presenter.go
│   └── interactor/                 ← concrete business logic
│       ├── certificate_interactor.go
│       └── csr_interactor.go
│
├── interface/                      ← adapters layer
│   ├── controller/                ← driving adapters (HTTP handlers)
│   │   ├── context.go             ← Echo → abstract Context
│   │   ├── app_controller.go     ← top-level container
│   │   ├── certificate_controller.go
│   │   ├── csr_controller.go
│   │   └── chain_controller.go
│   ├── presenter/                  ← driven adapters (response transformers)
│   │   └── certificate_presenter.go
│   └── repository/                 ← driven adapters (GORM implementations)
│       ├── certificate_repository.go
│       ├── csr_repository.go
│       ├── chain_repository.go
│       └── db_repository.go
│
├── infrastructure/                 ← framework & drivers layer
│   ├── datastore/
│   │   └── db.go                  ← MySQL + GORM connection factory
│   ├── router/
│   │   └── router.go              ← Echo routing + middleware
│   ├── crypto/
│   │   ├── pem.go                 ← PEM encode/decode
│   │   ├── fingerprint.go         ← SHA-256/SHA-1/MD5 fingerprint
│   │   ├── signer.go              ← x509 signing, self-signed CA
│   │   └── validator.go           ← chain validation, signature verification
│   └── certificate/
│       └── converter.go           ← x509.Certificate ↔ domain model
│
├── registry/                       ← composition root (DI wiring)
│   ├── registry.go                ← Registry interface + AppController
│   └── wire.go                    ← all NewXxxController() wiring
│
├── config/                         ← configuration
│   ├── config.go                  ← Viper YAML loader
│   └── config.yml                 ← dev configuration
│
├── cmd/                           ← entry points
│   ├── app/main.go                ← HTTP server
│   ├── migration/main.go          ← GORM AutoMigrate
│   └── cli/main.go               ← CLI tool
│
├── db/migrations/                 ← goose SQL migrations
│   ├── 001_setup.sql              ← certificates table
│   ├── 002_csrs.sql               ← csrs table
│   └── 003_chains.sql             ← certificate_chains table
│
├── docker/
│   └── docker-compose.yml         ← MySQL 8 for local dev
│
├── testutil/
│   └── db.go                      ← sqlmock + GORM test helper
│
├── main.go                         ← root entry point (legacy)
├── Makefile                        ← build/run/migrate/test targets
└── .air.toml                       ← hot-reload config
```

---

## 4. Domain Layer

**Purpose:** Define enterprise business rules without any dependency on frameworks, databases, or HTTP.

### `domain/model/`

Contains **pure Go structs** annotated with GORM tags (for persistence) and JSON tags (for API responses). The tags live here so the domain model can be serialized, but the *logic* inside each model has no framework awareness.

#### `certificate.go` — Certificate Entity

```go
type Certificate struct {
    ID           uint       `json:"id" gorm:"primaryKey"`
    Subject      string     `json:"subject"`     // Distinguished Name
    Issuer       string     `json:"issuer"`       // Issuer DN
    Serial       string     `json:"serial"`       // Hex serial number
    Fingerprint  string     `json:"fingerprint"` // SHA-256 fingerprint (dedup key)
    NotBefore    time.Time  `json:"not_before"`
    NotAfter     time.Time  `json:"not_after"`
    KeyUsage     []string   `json:"key_usage"`     // e.g. digitalSignature, keyCertSign
    ExtKeyUsage  []string   `json:"ext_key_usage"` // e.g. serverAuth, clientAuth
    DNSNames     []string   `json:"dns_names"`
    IPAddresses  []string   `json:"ip_addresses"`
    IsCA         bool       `json:"is_ca"`
    IsRevoked    bool       `json:"is_revoked"`
    RevokedAt    *time.Time `json:"revoked_at,omitempty"`
    CRLURL       string     `json:"crl_url"`
    OCSPURL      string     `json:"ocsp_url"`
    CertPEM      string     `json:"cert_pem"`
    KeyPEM       string     `json:"key_pem,omitempty"` // encrypted at rest
    ParentID     *uint      `json:"parent_id,omitempty"` // FK to issuer CA
    Profile      string     `json:"profile"` // tls-server, tls-client, ca, etc.
    Status       CertStatus `json:"status"`  // active, expired, revoked, pending, hold
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
    DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
```

Domain methods:
- `IsExpired()` — `time.Now().After(NotAfter)`
- `IsValid()` — `Status == active && !Expired && !Revoked`

#### `csr.go` — CSR Entity

Represents a Certificate Signing Request submitted by an end-entity. Its lifecycle is a state machine:

```
pending → approved → issued
      ↘ rejected
```

#### `certificate_chain.go` — Chain Entity

Represents a full trust chain (leaf → intermediate(s) → root). Stored as an ordered list of certificate IDs and a concatenated PEM chain string.

### `domain/errdomain/errors.go`

Sentinel errors defined **inside the domain layer** (no external imports). This ensures that `usecase` packages can return domain-specific errors without importing any infrastructure package.

```go
var ErrCertNotFound   = errors.New("certificate not found")
var ErrCertExpired     = errors.New("certificate has expired")
var ErrCertRevoked     = errors.New("certificate has been revoked")
var ErrCSRInvalid      = errors.New("CSR is invalid")
var ErrChainUntrusted  = errors.New("certificate chain is not trusted")
// ... etc.
```

---

## 5. Use Case Layer

**Purpose:** Express application-specific business rules. This layer knows nothing about HTTP, databases, or crypto libraries. It orchestrates domain entities using port interfaces.

### `usecase/repository/` — Input Ports (Interfaces)

These interfaces live in the **use case layer** and are implemented by the **interface layer**. This is the core of the dependency inversion in Clean Architecture.

#### `certificate_repository.go`

```go
type CertificateRepository interface {
    FindAll() ([]*model.Certificate, error)
    FindByID(id uint) (*model.Certificate, error)
    FindBySerial(serial string) (*model.Certificate, error)
    FindByFingerprint(fingerprint string) (*model.Certificate, error)
    FindBySubject(subject string) ([]*model.Certificate, error)
    FindByIssuer(issuer string) ([]*model.Certificate, error)
    FindByStatus(status model.CertStatus) ([]*model.Certificate, error)
    FindByProfile(profile string) ([]*model.Certificate, error)
    FindExpiring(withinDays int) ([]*model.Certificate, error)
    FindRevoked() ([]*model.Certificate, error)
    Create(cert *model.Certificate) (*model.Certificate, error)
    Update(cert *model.Certificate) (*model.Certificate, error)
    Delete(id uint) error
}
```

#### `csr_repository.go`

```go
type CSRRepository interface {
    FindAll() ([]*model.CSR, error)
    FindByID(id uint) (*model.CSR, error)
    FindBySubject(subject string) ([]*model.CSR, error)
    FindByStatus(status model.CSRStatus) ([]*model.CSR, error)
    FindPending() ([]*model.CSR, error)
    Create(csr *model.CSR) (*model.CSR, error)
    Update(csr *model.CSR) (*model.CSR, error)
    Delete(id uint) error
}
```

#### `chain_repository.go`

```go
type ChainRepository interface {
    FindAll() ([]*model.CertificateChain, error)
    FindByID(id uint) (*model.CertificateChain, error)
    FindByName(name string) (*model.CertificateChain, error)
    FindByType(chainType model.ChainType) ([]*model.CertificateChain, error)
    Create(chain *model.CertificateChain) (*model.CertificateChain, error)
    Update(chain *model.CertificateChain) (*model.CertificateChain, error)
    Delete(id uint) error
}
```

#### `db_repository.go`

Abstraction over database transactions. Allows interactors to run multiple repository calls atomically.

```go
type DBRepository interface {
    Transaction(func(interface{}) (interface{}, error)) (interface{}, error)
}
```

### `usecase/presenter/` — Output Ports (Interfaces)

Defines how the interactor formats its results. The concrete transformer lives in `interface/presenter/`.

```go
type CertificatePresenter interface {
    ResponseCert(cert *model.Certificate) *model.Certificate
    ResponseCerts(certs []*model.Certificate) []*model.Certificate
    ResponseChain(chain *model.CertificateChain) *model.CertificateChain
    ErrorResponse(err error) map[string]interface{}
}
```

### `usecase/interactor/` — Business Logic

The interactor is the **use case**. It receives requests from the controller, applies business rules using domain entities and repository ports, and returns results transformed by the presenter.

#### `certificate_interactor.go`

| Method | Description |
|---|---|
| `GetCertificate(id)` | Fetch a single cert, apply presenter transform |
| `ListCertificates()` | Fetch all certs |
| `SearchCertificates(query)` | Search by subject or fingerprint |
| `GetExpiringCertificates(withinDays)` | Alert window query |
| `ImportCertificate(pem, keyPEM)` | Parse PEM, check dedup, persist |
| `RevokeCertificate(id, reason)` | Soft-revoke inside a transaction |
| `RenewCertificate(id, newCSR)` | Replace cert with new CSR-backed cert |
| `DeleteCertificate(id)` | Soft-delete |
| `ValidateChain(chainID)` | Verify trust chain structure |
| `ValidateCertificate(id)` | Domain-level validity check |

#### `csr_interactor.go`

| Method | Description |
|---|---|
| `SubmitCSR(pem, email)` | Parse CSR, persist in `pending` state |
| `ApproveCSR(id, approverID)` | Transition → `approved` in a transaction |
| `RejectCSR(id, notes)` | Transition → `rejected` in a transaction |
| `GetCSR(id)` | Fetch single CSR |
| `ListPendingCSRs()` | Fetch all pending CSRs |
| `ListAllCSRs()` | Fetch all CSRs |

---

## 6. Interface Layer

**Purpose:** Implement the ports defined in the use case layer. This layer adapts external concerns (HTTP, MySQL, JSON) to the interfaces the use case layer declares.

### `interface/controller/` — Driving Adapters (HTTP Handlers)

Controllers receive HTTP requests, parse input, call the appropriate interactor method, and return HTTP responses. They depend only on interactor **interfaces**.

#### `context.go` — Context Abstraction

Echo's `echo.Context` is wrapped in an interface to make controllers unit-testable:

```go
type Context interface {
    JSON(code int, i interface{}) error
    Bind(i interface{}) error
    Param(name string) string
    QueryParam(name string) string
    Status(code int) error
    Get(key string) any
    Set(key string, val any)
}
```

A real `echoContext` adapter and a mock implementation for tests.

#### `app_controller.go` — Top-Level Container

Holds all domain controllers in a single struct injected into the router:

```go
type AppController struct {
    Certificate interface {
        GetCertificates(c Context) error
        GetCertificate(c Context) error
        ImportCertificate(c Context) error
        // ...
    }
    CSR   interface { /* CSR methods */ }
    Chain interface { /* Chain methods */ }
}
```

Uses **structural typing** — no named interface type is declared; the concrete struct implicitly satisfies it. This is the exact pattern from manakuro's `golang-clean-architecture`.

### `interface/repository/` — Driven Adapters (GORM)

Each repository here implements the corresponding `usecase/repository` interface using GORM.

```go
// interface/repository/certificate_repository.go
type certificateRepository struct { db *gorm.DB }

func NewCertificateRepository(db *gorm.DB) repository.CertificateRepository {
    return &certificateRepository{db: db}
}
```

#### `db_repository.go` — Transaction Adapter

Wraps GORM's `Begin()`/`Commit()`/`Rollback()` in a closure-based API matching `DBRepository`:

```go
func (r *dbRepository) Transaction(
    fn func(interface{}) (interface{}, error),
) (data interface{}, err error) {
    tx := r.db.Begin()
    defer func() {
        if p := recover(); p != nil { tx.Rollback(); panic(p) }
        else if err != nil        { tx.Rollback() }
        else                      { tx.Commit() }
    }()
    return fn(tx)
}
```

### `interface/presenter/` — Driven Adapters (Response Transformers)

```go
// interface/presenter/certificate_presenter.go
func (cp *certificatePresenter) ResponseCert(cert *model.Certificate) *model.Certificate {
    cert.KeyPEM = "" // mask private key in API responses
    return cert
}
```

---

## 7. Infrastructure Layer

**Purpose:** House all external-framework code — database connections, HTTP routing, x509 crypto primitives, and format converters. This layer may import domain entities but must never be imported by the domain or use case layers.

### `infrastructure/crypto/`

Pure cryptographic utilities using only `crypto/x509`, `crypto/rand`, `encoding/pem`, and `golang.org/x/crypto`.

| File | Responsibility |
|---|---|
| `pem.go` | `DecodePEMBlock`, `ParseCertificatePEM`, `ParsePrivateKeyPEM`, `ParseCSRPEM` |
| `fingerprint.go` | `SHA256Fingerprint(cert)` → colon-separated hex string |
| `signer.go` | `SignCertificate(csr, caCert, caKey, template)` — core signing engine; `GenerateSelfSignedCA` |
| `validator.go` | `ValidateCertificate`, `ValidateChain` using `x509.Verify`, `ParseChainPEM` |

### `infrastructure/certificate/`

| File | Responsibility |
|---|---|
| `converter.go` | `ToDomain(x509.Certificate)` → `model.Certificate`; `ToX509` → `*x509.Certificate` |

### `infrastructure/datastore/`

```go
// db.go
func NewDB() (*gorm.DB, error) {
    // Builds DSN from config.C.Database, opens MySQL via GORM
    // Sets MaxIdleConns(10), MaxOpenConns(100)
}
```

### `infrastructure/router/`

```go
// router.go
func NewRouter(e *echo.Echo, ac controller.AppController) *echo.Echo {
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())
    e.Use(middleware.RequestID())

    cert := e.Group("/certificates")
    cert.GET("",  func(c) { return ac.Certificate.GetCertificates(c) })
    cert.GET("/:id", ...)
    cert.POST("", ...)
    cert.DELETE("/:id", ...)
    cert.POST("/:id/revoke", ...)
    cert.GET("/expiring", ...)
    cert.POST("/validate", ...)

    csr := e.Group("/csrs")
    csr.GET("", ...)
    csr.POST("", ...)        // SubmitCSR
    csr.POST("/:id/approve", ...)
    csr.POST("/:id/reject", ...)

    chain := e.Group("/chains")
    // ... chain CRUD + validate
    return e
}
```

---

## 8. Registry — Composition Root

**Purpose:** Instantiate all concrete types in the correct order. This is the **only** place in the codebase where a GORM struct, an Echo handler, and a crypto utility are all in the same file.

```
registry/
├── registry.go     ← Registry interface + NewAppController()
└── wire.go        ← private newCertificateInteractor(), newCSRController(), etc.
```

The registry follows manakuro's pattern exactly:

```go
// registry/registry.go
type Registry interface {
    NewAppController() controller.AppController
}

// registry/wire.go
func (r *registry) newCertificateInteractor() interactor.CertificateInteractor {
    return interactor.NewCertificateInteractor(
        ir.NewCertificateRepository(r.db),    // driven adapter: GORM
        ir.NewChainRepository(r.db),         // driven adapter: GORM
        ip.NewCertificatePresenter(),        // driven adapter: presenter
        ir.NewDBRepository(r.db),             // driven adapter: transaction
    )
}

func (r *registry) NewCertificateController() controller.CertificateController {
    return controller.NewCertificateController(r.newCertificateInteractor())
}
```

Startup sequence in `cmd/app/main.go`:

```
config.ReadConfig()                              ← load config
    ↓
datastore.NewDB()                                ← open MySQL
    ↓
registry.NewRegistry(db)                        ← build DI container
    ↓
router.NewRouter(e, r.NewAppController())       ← wire routes
    ↓
e.Start(":" + config.C.Server.Address)           ← start server
```

---

## 9. Config

Configuration is loaded via [Viper](https://github.com/spf13/viper) from `config/config.yml`.

```yaml
# config/config.yml

database:
  user: root
  password: root
  addr: 127.0.0.1:3306
  dbname: x509_clean_architecture

server:
  address: "8080"
  mode: debug                    # debug | release | test

pki:
  default_validity_days: 365
  default_key_bits: 2048
  default_profile: tls-server
  allowed_profiles:
    - tls-server
    - tls-client
    - code-signing
    - smime
    - ca
  ca_common_name: "My Internal CA"
  crl_url: "http://localhost:8080/crl"
  ocsp_url: "http://localhost:8080/ocsp"

app:
  env: development
  notify_days: 30               # days before expiry to send alerts
  storage_dir: ./certs          # where to store PEM files on disk
  private_key_pass: ""          # passphrase for encrypting stored private keys
```

Config is accessed globally via `config.C`:

```go
addr := ":" + config.C.Server.Address
```

---

## 10. Entry Points

| File | Purpose |
|---|---|
| `main.go` | Root-level entry point (mirrors manakuro's legacy style) |
| `cmd/app/main.go` | Primary HTTP server — loads config, opens DB, starts Echo |
| `cmd/migration/main.go` | Runs `db.AutoMigrate()` for all domain models |
| `cmd/cli/main.go` | CLI tool: `x509-cli -cmd info -cert file.pem` |

---

## 11. Database Migrations

Three goose-style SQL migration files:

| File | Table | Key Columns |
|---|---|---|
| `001_setup.sql` | `certificates` | `serial` (UNIQUE), `fingerprint` (UNIQUE), `parent_id` FK, `profile` |
| `002_csrs.sql` | `csrs` | `status` enum, `approver_id` FK, `pem` LONGTEXT |
| `003_chains.sql` | `certificate_chains` | `name` (UNIQUE), `cert_ids` JSON array, `pem_chain` |

Run with:
```bash
make migrate     # via GORM AutoMigrate (dev)
# or
goose mysql "user:pass@tcp(127.0.0.1:3306)/x509_clean_architecture" up  # production
```

---

## 12. CLI Tool

```bash
# Info — print certificate details
x509-cli -cmd info -cert ./cert.pem

# Validate — check certificate structure
x509-cli -cmd validate -cert ./cert.pem

# Fingerprint — print SHA-256 fingerprint
x509-cli -cmd fingerprint -cert ./cert.pem

# Sign — sign a CSR with a CA (TODO: wire signer.go)
x509-cli -cmd sign -csr ./req.csr -ca-cert ./ca.pem -ca-key ./ca-key.pem -days 365
```

---

## 13. Domain Entities

### Certificate Lifecycle

```
import (pending) → active → expired
                          ↘ revoked
                          ↘ renewed → (new cert, old marked expired)
```

### CSR State Machine

```
pending → approved → issued (signed certificate created)
       ↘ rejected
```

### Certificate Profiles

| Profile | Key Usage | Extended Key Usage |
|---|---|---|
| `tls-server` | digitalSignature, keyEncipherment | serverAuth |
| `tls-client` | digitalSignature | clientAuth |
| `code-signing` | digitalSignature | codeSigning |
| `smime` | digitalSignature, keyEncipherment | emailProtection |
| `ca` | keyCertSign, cRLSign, digitalSignature | — |

---

## 14. HTTP API Reference

### Certificates

| Method | Path | Description |
|---|---|---|
| `GET` | `/certificates` | List all certificates |
| `GET` | `/certificates/:id` | Get a single certificate |
| `POST` | `/certificates` | Import a PEM certificate |
| `DELETE` | `/certificates/:id` | Soft-delete a certificate |
| `POST` | `/certificates/:id/revoke` | Revoke a certificate |
| `GET` | `/certificates/expiring?days=30` | List expiring certificates |
| `POST` | `/certificates/validate` | Validate a certificate |

### CSRs

| Method | Path | Description |
|---|---|---|
| `GET` | `/csrs` | List all CSRs |
| `GET` | `/csrs?status=pending` | List pending CSRs |
| `GET` | `/csrs/:id` | Get a single CSR |
| `POST` | `/csrs` | Submit a new CSR |
| `POST` | `/csrs/:id/approve` | Approve a CSR |
| `POST` | `/csrs/:id/reject` | Reject a CSR |

### Chains

| Method | Path | Description |
|---|---|---|
| `GET` | `/chains` | List all chains |
| `GET` | `/chains/:id` | Get a single chain |
| `POST` | `/chains` | Build a chain from cert IDs |
| `DELETE` | `/chains/:id` | Delete a chain |
| `POST` | `/chains/:id/validate` | Validate a trust chain |

### Health Check

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | `{"status": "ok"}` |

---

## 15. Development Workflow

```bash
# 1. Install dependencies
make deps

# 2. Start MySQL
make db-up

# 3. Run migrations
make migrate

# 4. Start the server with hot-reload
make start          # uses air

# Or run directly
make run

# 5. Run tests (with sqlmock — no DB needed)
make test

# 6. Build binaries
make build

# 7. Stop MySQL
make db-down
```

---

## 16. Next Steps — What to Implement

The scaffold is complete but the following stubs need real implementations:

| File | Stub | What to implement |
|---|---|---|
| `infrastructure/crypto/fingerprint.go` | `SHA1Fingerprint`, `MD5Fingerprint` | Use `crypto/sha1` and `crypto/md5` packages |
| `infrastructure/certificate/converter.go` | `parseCertPEM` | Delegate to `crypto.ParseCertificatePEM` |
| `usecase/interactor/certificate_interactor.go` | `RenewCertificate` | Wire `signer.SignCertificate` for CSR → new cert |
| `usecase/interactor/csr_interactor.go` | `decodeCSRPEM` | Delegate to `crypto.ParseCSRPEM` |
| `interface/controller/csr_controller.go` | `GetCSRs` body | Replace `nil, nil` stubs with real interactor calls |
| `interface/controller/chain_controller.go` | All methods | Wire chain interactor once `usecase/interactor/chain_interactor.go` is created |
| `cmd/cli/main.go` | `signCmd` | Wire `infrastructure/crypto/signer` for CSR signing |
| `config/config.yml` | `private_key_pass` | Set a real passphrase; encrypt/decrypt `KeyPEM` in the repository layer |
| `infrastructure/crypto/validator.go` | `KeyUsageFromString` | Fill in all `pkix.KeyUsage` cases |
| `infrastructure/certificate/converter.go` | `ipStrings` | Use `net.IP.String()` correctly (currently panics) |
| `testutil/` | DB mocks | Add mock implementations of all repository interfaces |

---

## Naming Conventions

| Element | Convention | Example |
|---|---|---|
| Domain entities | PascalCase, singular noun | `Certificate`, `CSR` |
| Repository interfaces (use case) | Noun + `Repository` | `CertificateRepository` |
| Presenter interfaces (use case) | Noun + `Presenter` | `CertificatePresenter` |
| Interactor interfaces | Noun + `Interactor` | `CertificateInteractor` |
| Concrete structs (interface layer) | same as interface or camelCase | `certificateController` |
| Constructor functions | `New` + type name | `NewCertificateRepository()` |
| DB tables | snake_case plural | `certificates`, `csrs` |
| HTTP handler methods | Verb + Resource | `GetCertificate`, `ImportCertificate` |
| Domain errors (errdomain) | `Err` + Entity + Condition | `ErrCertExpired`, `ErrCSRRejected` |
| Config fields | camelCase (matches YAML) | `server.address`, `pki.default_validity_days` |
