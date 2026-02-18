CREATE TABLE contact_groups (
    contact_id BIGINT NOT NULL REFERENCES contacts(id) ON DELETE CASCADE,
    group_id   BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (contact_id, group_id)
);