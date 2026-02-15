# db-context-timeouts

## Why This Rule Exists

Database calls can hang due to:

- Network issues
- Lock contention
- Slow queries

Without timeouts, requests may hang indefinitely.

---

## ❌ Incorrect

```go
rows, err := pool.Query(context.Background(), query)