-- Courses catalog (owned by Courses service).
CREATE TABLE IF NOT EXISTS courses (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title        VARCHAR(255) NOT NULL,
    description  TEXT NOT NULL,
    subject      VARCHAR(255) NOT NULL,
    level        VARCHAR(50) NOT NULL,
    owner_user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_courses_owner_user_id ON courses (owner_user_id);
CREATE INDEX idx_courses_subject ON courses (subject);

COMMENT ON TABLE courses IS 'Courses: study courses created by users.';

