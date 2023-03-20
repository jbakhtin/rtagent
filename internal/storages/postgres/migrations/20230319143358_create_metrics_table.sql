-- +goose Up
CREATE TABLE metrics(
    id varchar(256),
    type varchar(256),
    delta bigint null,
    value double precision null,
    PRIMARY KEY (id, type)
);

-- +goose Down
DROP TABLE metrics;
