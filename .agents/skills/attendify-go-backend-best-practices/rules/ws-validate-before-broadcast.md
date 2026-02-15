# ws-validate-before-broadcast

## Why This Rule Exists

Broadcasting before persistence creates inconsistency.

If broadcast succeeds but DB insert fails,
clients receive false confirmation.

Order matters.

---

## ❌ Incorrect

```go
hub.broadcast <- message
repo.InsertAttendance(ctx, record)
