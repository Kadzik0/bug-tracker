-- +goose Up
CREATE TABLE incident_events (
    id UUID NOT NULL PRIMARY KEY,
    incident_id UUID NOT NULL REFERENCES incidents(id),
    actor_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(255) NOT NULL CHECK (type IN ('CREATED', 'STATUS_CHANGED', 'ASSIGNED', 'PRIORITY_CHANGED', 'COMMENT_ADDED')),
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE incident_events;