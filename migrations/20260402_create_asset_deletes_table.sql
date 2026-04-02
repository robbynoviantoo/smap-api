-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS asset_deletes (
    id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    asset_id         BIGINT UNSIGNED NOT NULL,
    name             VARCHAR(255) NOT NULL,
    image            VARCHAR(255) NULL,
    asset            VARCHAR(100) NULL,
    no_asset         VARCHAR(100) NULL,
    location         VARCHAR(255) NULL,
    building         VARCHAR(255) NULL,
    category         VARCHAR(100) NULL,
    sub_category     VARCHAR(100) NULL,
    merk             VARCHAR(100) NULL,
    size             VARCHAR(50)  NULL,
    unit             VARCHAR(50)  NULL,
    status           VARCHAR(50)  NULL,
    available_status VARCHAR(50)  NULL,
    last_maintenance DATETIME     NULL,
    next_maintenance DATETIME     NULL,
    remarks          TEXT         NULL,
    asset_created_at DATETIME     NULL,
    asset_updated_at DATETIME     NULL,
    deleted_by       BIGINT UNSIGNED NULL,
    created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS asset_deletes;
-- +goose StatementEnd
