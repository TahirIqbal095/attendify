# ws-no-shared-mutable-state

## Why This Rule Exists

Sharing mutable state across goroutines without synchronization
causes data races.

WebSocket systems are concurrency-heavy.

Shared state must be:

- Owned by one goroutine
- Or protected with synchronization

---

## ❌ Incorrect

```go
var connectedUsers []User

go func() {
    connectedUsers = append(connectedUsers, user)
}()