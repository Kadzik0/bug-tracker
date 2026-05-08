-- +goose Up
CREATE TABLE incident_comments (
    id UUID NOT NULL PRIMARY KEY,
    incident_id UUID NOT NULL REFERENCES incidents(id),
    author_id UUID NOT NULL REFERENCES users(id),
    body TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE incident_comments;