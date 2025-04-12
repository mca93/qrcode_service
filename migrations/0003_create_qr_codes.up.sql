CREATE TABLE qr_codes (
    id UUID PRIMARY KEY,
    type TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    status TEXT,
    scan_count BIGINT DEFAULT 0,
    image_url TEXT,
    deep_link_url TEXT,
    client_app_id UUID REFERENCES client_apps(id),
    template_id UUID REFERENCES templates(id)
);
