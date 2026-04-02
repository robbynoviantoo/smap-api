-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS events (
    id        BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title     VARCHAR(255) NOT NULL,
    start     DATETIME     NULL,
    `end`     DATETIME     NULL,
    all_day   TINYINT(1)   NOT NULL DEFAULT 1,
    color     VARCHAR(50)  NULL,
    created_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pengadaan_asset_pendings (
    id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id      BIGINT UNSIGNED NOT NULL,
    request_date DATE            NOT NULL,
    category     VARCHAR(100)    NULL,
    name         VARCHAR(255)    NOT NULL,
    merk         VARCHAR(100)    NULL,
    spec         TEXT            NULL,
    qty          INT             NOT NULL DEFAULT 1,
    unit         VARCHAR(50)     NULL,
    image        VARCHAR(500)    NULL,
    remark       TEXT            NULL,
    priority     VARCHAR(50)     NOT NULL DEFAULT 'normal',
    status       VARCHAR(20)     NOT NULL DEFAULT 'pending', -- pending, approved, rejected
    reviewed_by  BIGINT UNSIGNED NULL,
    reviewed_at  DATETIME        NULL,
    reject_reason TEXT           NULL,
    created_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_pap_user     FOREIGN KEY (user_id)    REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_pap_reviewer FOREIGN KEY (reviewed_by) REFERENCES users(id) ON DELETE SET NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pengadaan_asset_pendings;
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
