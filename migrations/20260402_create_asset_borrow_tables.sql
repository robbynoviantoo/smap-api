-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS asset_borrow_pendings (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    asset_id     BIGINT UNSIGNED NOT NULL,
    user_id      BIGINT UNSIGNED NOT NULL,
    remark       VARCHAR(255)    NOT NULL,
    status       VARCHAR(20)     NOT NULL DEFAULT 'pending', -- pending, approved, rejected
    reviewed_by  BIGINT UNSIGNED NULL,
    reviewed_at  DATETIME        NULL,
    reject_reason VARCHAR(255)   NULL,
    created_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_abp_asset  FOREIGN KEY (asset_id)    REFERENCES assets(id) ON DELETE CASCADE,
    CONSTRAINT fk_abp_user   FOREIGN KEY (user_id)     REFERENCES users(id)  ON DELETE CASCADE,
    CONSTRAINT fk_abp_reviewer FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS asset_transactions (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    asset_id    BIGINT UNSIGNED NOT NULL,
    borrower_id BIGINT UNSIGNED NOT NULL,
    type        VARCHAR(20)     NOT NULL, -- borrow, return
    remark      VARCHAR(255)    NOT NULL,
    action_at   DATETIME        NOT NULL,
    created_at  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_at_asset    FOREIGN KEY (asset_id)    REFERENCES assets(id) ON DELETE CASCADE,
    CONSTRAINT fk_at_borrower FOREIGN KEY (borrower_id) REFERENCES users(id)  ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS asset_transactions;
DROP TABLE IF EXISTS asset_borrow_pendings;
-- +goose StatementEnd
