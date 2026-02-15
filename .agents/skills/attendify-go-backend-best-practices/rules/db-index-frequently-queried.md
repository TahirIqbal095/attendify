# db-index-frequently-queried

## Why This Rule Exists

Queries degrade as data grows.

Attendance systems will accumulate records quickly.

Without indexes, queries become slow.

---

## Example Query

```sql
SELECT * FROM attendance
WHERE class_id = $1;