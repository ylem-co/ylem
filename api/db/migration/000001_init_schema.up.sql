CREATE TABLE `oauth_clients` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `uuid` VARCHAR(255) NOT NULL,
    `user_uuid` VARCHAR(255) NOT NULL,
    `organization_uuid` VARCHAR(255) NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `secret` VARCHAR(255) NOT NULL,
    `allowed_scopes` TEXT NOT NULL,
    `created_at` TIMESTAMP DEFAULT NOW() NOT NULL,
    `updated_at` TIMESTAMP DEFAULT NULL ON UPDATE NOW(),
    `deleted_at` TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_uuid` (`uuid`)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `oauth_tokens` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `uuid` VARCHAR(255) NOT NULL,
    `oauth_client_uuid` VARCHAR(255) NOT NULL,
    `access_token` TEXT,
    `refresh_token` TEXT,
    `internal_token` TEXT,
    `scope` VARCHAR(255) DEFAULT NULL,
    `access_token_expires_in` INT(11) DEFAULT NULL,
    `refresh_token_expires_in` INT(11) DEFAULT NULL,
    `created_at` TIMESTAMP DEFAULT NOW() NOT NULL,
    `updated_at` TIMESTAMP DEFAULT NULL ON UPDATE NOW(),
    `deleted_at` TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_uuid` (`uuid`),
    FOREIGN KEY (oauth_client_uuid) REFERENCES oauth_clients(uuid)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8mb4
COLLATE=utf8mb4_unicode_ci;
