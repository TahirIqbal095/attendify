## Attendify — Real-Time Attendance Backend System

## 1. Project Overview

Attendify is a real-time backend system for managing classroom attendance digitally.

**Capabilities:**
- Teachers create and manage classes
- Students mark attendance live
- Real-time broadcasting of attendance events
- Persistent storage of attendance records
- Secure authentication and role-based authorization

Single-node backend service built in Go.

---

## 2. Core Responsibilities

- Authenticate users securely
- Enforce role-based authorization
- Provide REST endpoints for class management
- Provide WebSocket endpoint for live attendance
- Persist attendance records in PostgreSQL
- Broadcast attendance events only after successful persistence
- Maintain strict data integrity guarantees
- Follow layered clean architecture

---

## 3. Technology Stack

| Component | Technology |
|-----------|------------|
| Language | Go |
| HTTP Framework | Gin |
| Database | PostgreSQL |
| Database Driver | pgx (connection pool) |
| Authentication | JWT (HMAC) |
| Password Hashing | bcrypt |
| Validation | go-playground/validator |
| Logging | zerolog |
| Configuration | Environment-based |
| WebSockets | Gorilla WebSocket |

---

## 4. Architectural Principles

**Layered Structure:** Handler → Service → Repository → Database

**Rules:**
- Handlers parse input and return responses
- Services contain business logic
- Repositories handle database operations only
- Handlers must never access database directly
- Authorization must not live inside handlers
- Context must propagate through all layers
- No global mutable state

**Philosophy:** Explicitness over abstraction. Correctness over cleverness.

---

## 5. Critical: Data Integrity Rules

The database is the final authority on correctness.

- Use transactions for multi-step operations
- Define foreign key constraints
- Define UNIQUE constraints for attendance
- Index foreign keys and frequently queried columns
- Avoid N+1 query patterns
- Always use context timeouts for database calls

Never rely solely on application-level validation for data safety.

---

## 6. Critical: Authentication & Security Rules

Security logic must be centralized and enforced consistently.

- Hash passwords using bcrypt
- Validate JWT signature and expiration
- Enforce roles via middleware
- Keep JWT claims minimal (user_id, role, exp)
- Never hardcode secrets
- Never return sensitive fields in API responses
- Never leak internal database errors

**Principle:** Authentication verifies identity. Authorization verifies permission. They must remain separate.

---

## 7. WebSocket Concurrency Rules

- Use a single Hub goroutine to own all client state
- Never modify shared maps outside the hub
- Authenticate users before WebSocket registration
- Persist attendance before broadcasting
- Clean up connections on disconnect
- Prevent broadcast blocking from slow clients
- Avoid shared mutable state across goroutines

**Principle:** Use ownership via channels rather than ad-hoc synchronization.

---

## 8. API Design Standards

**Success Response:**
```json
{
    "success": true,
    "data": { ... }
}
```

**Error Response:**
```json
{
    "success": false,
    "error": "message"
}
```

**Rules:**
- Use correct HTTP status codes
- Validate all request inputs
- Do not expose internal error details
- Maintain consistent response shape across all endpoints

---

## 9. Logging & Error Handling

- Use structured logging (zerolog)
- Wrap errors with context
- Log errors at system boundaries only
- Avoid duplicate logging of the same error
- Do not panic in production flows
- Logs must be structured and meaningful

---

## 10. Production Readiness Requirements

- Graceful shutdown handling
- Health check endpoint (`/health`)
- Environment-driven configuration
- Disabled debug mode in production
- Clean folder structure
- Connection pooling for PostgreSQL
- Secrets must never be committed to version control

---

## 11. Attendance Flow (Authoritative Order)

1. Validate JWT
2. Validate role (student)
3. Validate class existence
4. Begin transaction (if required)
5. Insert attendance record
6. Commit transaction
7. Broadcast event to WebSocket clients

**Critical:** Broadcasting must never occur before successful persistence.

---

## 12. Non-Goals

Intentionally not included:
- Frontend UI
- Distributed scaling
- Multi-node pub/sub
- Redis caching
- Multi-room WebSocket segmentation
- Background job processing

Single-instance real-time backend focused on correctness and structure.

---

## 13. Agent Operational Constraints

**Prohibited Actions:**
- Never bypass service layer
- Never perform DB operations in handlers
- Never broadcast uncommitted state
- Never modify shared WebSocket state outside hub

**Required Actions:**
- Always propagate context
- Always enforce role checks via middleware
- Always return standardized API responses
- Always prefer explicit logic over abstraction

Violations introduce instability.

---

## 14. Engineering Philosophy

Attendify is a backend engineering exercise emphasizing:
- Data integrity
- Security correctness
- Concurrency safety
- Clean architecture
- Production discipline

Agents must prioritize structural correctness over convenience.
