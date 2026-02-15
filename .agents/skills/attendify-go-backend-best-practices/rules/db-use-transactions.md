# db-use-transactions

## Why This Rule Exists

Multi-step database operations must be atomic.

If one step succeeds and another fails, the database can enter an inconsistent state.

For the Attendify system, marking attendance involves:

1. Validating class existence
2. Validating student role
3. Inserting attendance record
4. Possibly updating related metadata

These must succeed or fail together.

---

## ❌ Incorrect

```go
_, err := repo.InsertAttendance(ctx, record)
if err != nil {
    return err
}

_, err = repo.UpdateAttendanceCount(ctx, classID)
if err != nil {
    return err
}
