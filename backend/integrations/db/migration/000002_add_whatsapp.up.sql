CREATE TABLE `whatsapps`
(
    `id`             int(10) unsigned NOT NULL AUTO_INCREMENT,
    `integration_id` int(10) unsigned NOT NULL,
    `content_sid`    BLOB NOT NULL,
    `updated_at`     timestamp        NULL     DEFAULT NULL ON UPDATE current_timestamp(),
    `created_at`     timestamp        NOT NULL DEFAULT current_timestamp(),
    PRIMARY KEY (`id`),
    KEY `integration_id` (`integration_id`),
    CONSTRAINT `whatsapps_integration_id` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci;
