-- +goose Up
CREATE TABLE incidents (
    id UUID NOT NULL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    environment varchar(255) NOT NULL CHECK (environment IN ('prod', 'stage', 'dev')),
    priority varchar(255) NOT NULL CHECK (priority IN ('P1', 'P2', 'P3', 'P4')),
    status varchar(255) NOT NULL CHECK (status IN ('OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED')),
    reporter_id UUID NOT NULL REFERENCES users(id),
    assignee_id UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE incidents;