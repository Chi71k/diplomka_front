CREATE TYPE match_status AS ENUM ('pending', 'accepted', 'declined', 'canceled');

CREATE TABLE IF NOT EXISTS matches (
                                       id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    requester_id  UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id   UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status        match_status NOT NULL DEFAULT 'pending',
    message       TEXT        NOT NULL DEFAULT '',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT no_self_match CHECK (requester_id <> receiver_id)
    );

CREATE INDEX IF NOT EXISTS idx_matches_requester ON matches(requester_id);
CREATE INDEX IF NOT EXISTS idx_matches_receiver  ON matches(receiver_id);
CREATE INDEX IF NOT EXISTS idx_matches_pair ON matches(
    LEAST(requester_id::text, receiver_id::text),
    GREATEST(requester_id::text, receiver_id::text)
    );