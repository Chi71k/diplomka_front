-- Interests catalog (reference data for user interests selection).
CREATE TABLE IF NOT EXISTS interests (
    id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE interests
    ADD COLUMN IF NOT EXISTS created_at timestamptz NOT NULL DEFAULT now();

-- User-interests many-to-many (Users service).
CREATE TABLE IF NOT EXISTS user_interests (
    user_id    UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    interest_id UUID NOT NULL REFERENCES interests (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, interest_id)
);

CREATE INDEX IF NOT EXISTS idx_user_interests_user_id ON user_interests (user_id);
CREATE INDEX IF NOT EXISTS idx_user_interests_interest_id ON user_interests (interest_id);

-- Seed default interests (idempotent).
INSERT INTO interests (name) VALUES
    ('Game Development'), ('Web Development'), ('Frontend Development'), ('Backend Development'), ('Database Management'), ('AI/ML'), ('Cybersecurity'), ('UI/UX Design'), ('DevOps'), ('Cloud Computing'), ('Data Science'), ('Software Engineering'), ('Computer Science'), ('Mathematics'), ('Languages')
ON CONFLICT (name) DO NOTHING;
