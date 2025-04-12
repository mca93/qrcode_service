CREATE TABLE api_keys (
    id UUID PRIMARY KEY,
    client_app_id UUID REFERENCES client_apps(id),
    key_prefix TEXT,
    status TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMP
);
