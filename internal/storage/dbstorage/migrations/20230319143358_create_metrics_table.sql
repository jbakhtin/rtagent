-- +goose Up
CREATE TABLE IF NOT EXISTS metrics(
    id varchar(256),
    type varchar(256),
    delta bigint null,
    value double precision null,
    PRIMARY KEY (id, type)
);

-- +goose Down
DROP TABLE IF EXISTS metrics;
