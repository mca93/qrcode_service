CREATE TABLE templates (
    id UUID PRIMARY KEY,
    client_app_id UUID REFERENCES client_apps(id),
    name TEXT,
    description TEXT,
    active BOOLEAN DEFAULT TRUE
);
