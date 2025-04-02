-- +goose Up
-- +goose StatementBegin
-- enable foreign key support
PRAGMA foreign_keys = ON;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- disable foreign key support
PRAGMA foreign_keys = OFF;
-- +goose StatementEnd
