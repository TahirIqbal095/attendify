# db-unique-constraints

## Why This Rule Exists

Duplicate attendance records must not exist.

Application checks are not enough in concurrent systems.

Only the database can reliably prevent duplicates.

---

## ❌ Incorrect

```go
exists, _ := repo.CheckAttendance(ctx, classID, studentID)
if exists {
    return errors.New("already marked")
}
repo.InsertAttendance(ctx, record)