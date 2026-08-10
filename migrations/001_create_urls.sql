CREATE TABLE urls (
                      id BIGSERIAL PRIMARY KEY,

                      code VARCHAR(10) NOT NULL UNIQUE,

                      original_url TEXT NOT NULL,

                      clicks BIGINT NOT NULL DEFAULT 0,

                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                      expires_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_urls_code
    ON urls(code);