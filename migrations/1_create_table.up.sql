CREATE TABLE go_table (
    id           SERIAL PRIMARY KEY,
    path         varchar(255) NOT NULL,
    geo          integer NOT NULL,
    header_props json NOT NULL,
    footer_props json NOT NULL,
    created_at   timestamp with time zone NOT NULL,
    modified_at  timestamp with time zone NOT NULL,
    deleted_at   timestamp with time zone,
    UNIQUE(path, geo)
);
