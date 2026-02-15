---
name: attendify-go-backend-best-practices
description: Backend engineering guidelines for building the Attendify real-time attendance system using Go, Gin, PostgreSQL, and WebSockets. This skill should be used when writing, reviewing, or refactoring backend code to ensure clean architecture, concurrency safety, data integrity, and production-grade design patterns.
license: MIT
metadata:
  author: tahir-iqbal
  version: "1.0.0"
---

# Attendify Go Backend Best Practices

Comprehensive backend engineering guide for the Attendify real-time attendance system.  
Contains structured rules across 8 categories, prioritized by architectural impact and production relevance.

These rules ensure correctness, maintainability, concurrency safety, and long-term scalability.

---

## When to Apply

Reference these guidelines when:

- Writing new API handlers
- Implementing authentication or authorization
- Designing database schema or writing SQL queries
- Building WebSocket functionality
- Refactoring service or repository layers
- Reviewing concurrency logic
- Improving error handling and logging
- Preparing the system for production deployment

---

## Rule Categories by Priority

| Priority | Category | Impact | Prefix |
|----------|----------|--------|--------|
| 1 | Data Integrity & Transactions | CRITICAL | `db-` |
| 2 | Authentication & Security | CRITICAL | `auth-` |
| 3 | Concurrency & WebSockets | HIGH | `ws-` |
| 4 | Architecture & Layering | HIGH | `arch-` |
| 5 | API Design & Contracts | MEDIUM-HIGH | `api-` |
| 6 | Error Handling & Logging | MEDIUM | `err-` |
| 7 | Performance & Optimization | LOW-MEDIUM | `perf-` |
| 8 | Production Readiness | LOW | `prod-` |

---

## Quick Reference

---

### 1. Data Integrity & Transactions (CRITICAL)

- `db-use-transactions` — Use transactions for multi-step database operations.
- `db-enforce-foreign-keys` — Always use foreign key constraints.
- `db-unique-constraints` — Prevent duplicate attendance records.
- `db-index-frequently-queried` — Index foreign keys and lookup fields.
- `db-avoid-n-plus-one` — Avoid repeated queries inside loops.
- `db-context-timeouts` — Always use context with timeouts for DB calls.

---

### 2. Authentication & Security (CRITICAL)

- `auth-hash-passwords` — Use bcrypt for password hashing.
- `auth-validate-jwt` — Validate JWT signature and expiration.
- `auth-role-enforcement` — Enforce role checks in middleware, not handlers.
- `auth-no-plain-secrets` — Never hardcode secrets.
- `auth-validate-input` — Validate all request payloads.
- `auth-minimal-claims` — Store minimal required data inside JWT.

---

### 3. Concurrency & WebSockets (HIGH)

- `ws-single-hub-owner` — Hub must own client map to avoid race conditions.
- `ws-no-shared-mutable-state` — Avoid unsynchronized shared memory.
- `ws-cleanup-connections` — Properly unregister clients on disconnect.
- `ws-validate-before-broadcast` — Persist data before broadcasting.
- `ws-handle-backpressure` — Avoid blocking broadcasts.
- `ws-auth-on-connect` — Validate JWT during WebSocket upgrade.

---

### 4. Architecture & Layering (HIGH)

- `arch-separate-layers` — Keep handler, service, and repository separated.
- `arch-no-db-in-handler` — Handlers should not call DB directly.
- `arch-dependency-injection` — Pass dependencies explicitly.
- `arch-single-responsibility` — One responsibility per component.
- `arch-no-global-state` — Avoid global mutable variables.
- `arch-context-propagation` — Pass context through all layers.

---

### 5. API Design & Contracts (MEDIUM-HIGH)

- `api-standard-response-format` — Use consistent success/error JSON format.
- `api-http-status-correct` — Use correct HTTP status codes.
- `api-validate-request-structs` — Validate all input structs.
- `api-no-internal-error-leak` — Do not expose internal DB errors.
- `api-versioning-ready` — Design routes with future versioning in mind.

---

### 6. Error Handling & Logging (MEDIUM)

- `err-wrap-errors` — Wrap errors with context.
- `err-structured-logging` — Use structured logs (zerolog).
- `err-no-log-flooding` — Avoid excessive logging inside loops.
- `err-log-at-boundaries` — Log errors at system boundaries only.
- `err-no-panic-in-prod` — Avoid panic outside startup.

---

### 7. Performance & Optimization (LOW-MEDIUM)

- `perf-use-connection-pool` — Use pgx pool, not single connections.
- `perf-index-joins` — Index foreign key joins.
- `perf-broadcast-efficiently` — Broadcast without re-marshaling repeatedly.
- `perf-avoid-unnecessary-json` — Avoid repeated encoding/decoding.
- `perf-timeouts-everywhere` — Apply timeouts to HTTP and DB.

---

### 8. Production Readiness (LOW)

- `prod-graceful-shutdown` — Implement graceful shutdown.
- `prod-env-config-only` — Load configuration from environment.
- `prod-health-endpoint` — Provide health check route.
- `prod-no-debug-mode` — Disable Gin debug mode in production.
- `prod-clean-folder-structure` — Maintain clean project layout.
- `prod-swagger-docs` — Add OpenAPI documentation after stabilization.

---

## Architectural Philosophy

The Attendify backend must:

- Prefer correctness over cleverness.
- Favor explicitness over abstraction.
- Avoid premature optimization.
- Protect data integrity first.
- Enforce authorization centrally.
- Treat concurrency carefully.
- Separate concerns strictly.

---

## How to Use

Read individual rule files for database and concurrency-related detailed explanations and code examples:

