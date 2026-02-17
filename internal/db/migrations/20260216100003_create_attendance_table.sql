-- migrate:up
CREATE TABLE attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_date DATE NOT NULL,
    marked_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(class_id, student_id, session_date)
);

CREATE INDEX idx_attendance_class_id ON attendance(class_id);
CREATE INDEX idx_attendance_student_id ON attendance(student_id);
CREATE INDEX idx_attendance_session ON attendance(class_id, session_date);

-- migrate:down
DROP TABLE IF EXISTS attendance;
