# db-avoid-n-plus-one

## Why This Rule Exists

N+1 queries destroy performance.

Looping and querying inside the loop scales poorly.

---

## ❌ Incorrect

```go
for _, class := range classes {
    students, _ := repo.GetStudentsByClass(ctx, class.ID)
}