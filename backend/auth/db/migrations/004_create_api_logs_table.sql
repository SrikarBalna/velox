-- +goose Up
CREATE TABLE api_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    api_key_id UUID NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
    submission_id UUID, -- Link to the job
    endpoint VARCHAR(255) NOT NULL,
    method VARCHAR(10) NOT NULL,
    status_code INTEGER NOT NULL,
    duration_ms INTEGER NOT NULL,
    overall_state VARCHAR(50), -- From SubmissionResponse (Accepted, TLE, etc)
    language VARCHAR(50),
    error_msg TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_api_logs_api_key_id_created_at ON api_logs(api_key_id, created_at);

-- +goose Down
DROP TABLE api_logs;
