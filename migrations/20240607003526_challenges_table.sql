-- +goose Up
-- +goose StatementBegin
CREATE TABLE challenges (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    challenge CHAR(64) NOT NULL, -- hex encoded
    bento_id UUID NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT fk_challenge_bento FOREIGN KEY (bento_id) REFERENCES bentos(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE challenges;
-- +goose StatementEnd
