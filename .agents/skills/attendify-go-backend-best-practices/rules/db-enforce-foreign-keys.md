# db-enforce-foreign-keys

## Why This Rule Exists

Application-level validation is not enough.

If relationships are not enforced at the database level,
invalid data can enter the system.

The database must guarantee relational integrity.

---

## ❌ Incorrect Schema

```sql
CREATE TABLE attendance (
    id UUID PRIMARY KEY,
    class_id UUID,
    student_id UUID
);