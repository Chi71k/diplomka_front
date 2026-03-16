CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE user_embeddings (
    user_id    UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    embedding  vector(768) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON user_embeddings USING hnsw (embedding vector_cosine_ops);
