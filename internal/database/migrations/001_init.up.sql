CREATE TABLE bins (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT NOT NULL DEFAULT '',
    secret        TEXT NOT NULL,
    max_requests  INTEGER NOT NULL DEFAULT 1000,
    request_count INTEGER NOT NULL DEFAULT 0,
    expires_at    TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE requests (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bin_id       UUID NOT NULL REFERENCES bins(id) ON DELETE CASCADE,
    method       TEXT NOT NULL,
    path         TEXT NOT NULL,
    query        JSONB NOT NULL DEFAULT '{}',
    headers      JSONB NOT NULL DEFAULT '{}',
    body         BYTEA,
    content_type TEXT,
    size_bytes   INTEGER NOT NULL DEFAULT 0,
    received_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_requests_bin_received
    ON requests(bin_id, received_at DESC);

CREATE TABLE replays (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id    UUID NOT NULL REFERENCES requests(id) ON DELETE CASCADE,
    target_url    TEXT NOT NULL,
    status_code   INTEGER,
    response_body BYTEA,
    error         TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);