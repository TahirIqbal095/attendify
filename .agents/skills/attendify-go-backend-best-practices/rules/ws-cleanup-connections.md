# ws-cleanup-connections

## Why This Rule Exists

If connections are not cleaned up:

- Memory leaks occur
- Goroutines remain active
- Hub state grows indefinitely

Disconnected clients must be unregistered.

---

## ❌ Incorrect

```go
for {
    _, _, err := conn.ReadMessage()
    if err != nil {
        break
    }
}