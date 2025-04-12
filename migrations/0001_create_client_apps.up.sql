CREATE TABLE client_apps (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    contact_email TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
