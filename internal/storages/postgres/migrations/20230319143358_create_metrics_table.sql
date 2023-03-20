-- +goose Up
CREATE TABLE metrics(
    id varchar(20),
    type varchar(20),
    delta bigint null,
    value double precision null,
    PRIMARY KEY (id, type)
);

-- +goose Down
DROP TABLE metrics;
