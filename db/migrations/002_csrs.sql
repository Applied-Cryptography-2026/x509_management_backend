-- Migration: 002_csrs.sql
-- Creates the CSRs table.

-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS csrs (
    id                   INT UNSIGNED NOT NULL AUTO_INCREMENT,
    subject              VARCHAR(1024) NOT NULL COMMENT 'Requested subject DN',
    pem                  LONGTEXT NOT NULL COMMENT 'PEM-encoded CSR',
    key_algorithm       VARCHAR(64) NOT NULL COMMENT 'Key algorithm: RSA, ECDSA, Ed25519',
    signature_algorithm VARCHAR(64) NOT NULL COMMENT 'Signature algorithm',
    dns_names            JSON DEFAULT NULL COMMENT 'SAN DNS names',
    ip_addresses         JSON DEFAULT NULL COMMENT 'SAN IP addresses',
    requester_email      VARCHAR(256) DEFAULT NULL COMMENT 'Email of the requester',
    status               VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT 'pending|approved|rejected|issued',
    approved_at           DATETIME DEFAULT NULL COMMENT 'Approval timestamp',
    rejected_at           DATETIME DEFAULT NULL COMMENT 'Rejection timestamp',
    approver_id           INT UNSIGNED DEFAULT NULL COMMENT 'FK to approver user',
    notes                TEXT DEFAULT NULL COMMENT 'Admin notes',
    created_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at            DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at            DATETIME DEFAULT NULL COMMENT 'Soft-delete timestamp',

    PRIMARY KEY (id),
    INDEX idx_status (status),
    INDEX idx_subject (subject(255)),
    INDEX idx_requester_email (requester_email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_as_ci;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS csrs;
-- +goose StatementEnd
