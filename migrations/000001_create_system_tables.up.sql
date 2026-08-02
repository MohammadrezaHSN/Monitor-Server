CREATE TABLE IF NOT EXISTS agents (
    id BIGSERIAL PRIMARY KEY,
    system_id TEXT UNIQUE,
    interval_seconds INT NOT NULL,
    is_alive BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT,
    last_active TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    next_request_at TIMESTAMPTZ
    );

CREATE TABLE IF NOT EXISTS requests (
    id BIGSERIAL PRIMARY KEY,
    agent_id BIGINT REFERENCES agents(id) ON DELETE CASCADE,
    context TEXT NOT NULL,
    is_done BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    received_to_agent_at TIMESTAMPTZ,
    estimate_response_at TIMESTAMPTZ
    );

CREATE TABLE IF NOT EXISTS responses (
    id BIGSERIAL PRIMARY KEY,
    agent_id BIGINT REFERENCES agents(id) ON DELETE CASCADE,
    request_id BIGINT REFERENCES requests(id) ON DELETE CASCADE,
    result TEXT NOT NULL,
    is_successful BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ,
    received_to_server_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS files (
    id BIGSERIAL PRIMARY KEY,
    agent_id BIGINT REFERENCES agents(id) ON DELETE CASCADE,
    request_id BIGINT REFERENCES requests(id) ON DELETE CASCADE,
    file_path_from_agent_system TEXT NOT NULL,
    file_path_from_server_system TEXT NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    upload_at TIMESTAMPTZ,
    received_to_server_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
