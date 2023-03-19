-- +goose Up
-- +goose StatementBegin
CREATE TABLE metrics(
    id varchar(20),
    type varchar(20),
    delta bigint null,
    value double precision null,
    PRIMARY KEY (id, type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE metrics;
-- +goose StatementEnd
