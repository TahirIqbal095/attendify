# ws-handle-backpressure

## Why This Rule Exists

If a single slow client blocks WriteMessage,
the entire hub can stall.

This creates system-wide performance issues.

---

## ❌ Incorrect

```go
for conn := range h.clients {
    conn.WriteMessage(websocket.TextMessage, message)
}