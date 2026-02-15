# ws-auth-on-connect

## Why This Rule Exists

WebSocket connections are long-lived.

If authentication is not enforced at upgrade time,
unauthorized users can maintain open connections.

Authentication must happen before registration.

---

## ❌ Incorrect

```go
conn, _ := upgrader.Upgrade(w, r, nil)
hub.register <- conn