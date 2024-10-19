-- Create syntax for TABLE 'organizations'
CREATE TABLE `organizations` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `uuid` varchar(255) NOT NULL,
  `data_location` varchar(255) NOT NULL DEFAULT 'EU',
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `is_data_source_created` tinyint(1) NOT NULL DEFAULT 0,
  `is_destination_created` tinyint(1) NOT NULL DEFAULT 0,
  `is_pipeline_created` tinyint(1) NOT NULL DEFAULT 0,
  `data_key` blob DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `name` (`name`,`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'users'
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `first_name` varchar(255) NOT NULL,
  `last_name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL DEFAULT '',
  `source` varchar(255) DEFAULT 'app',
  `password` varchar(255) NOT NULL DEFAULT '',
  `phone` varchar(255) DEFAULT '',
  `uuid` varchar(255) NOT NULL,
  `external_system_id` varchar(255) DEFAULT NULL,
  `roles` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '',
  `organization_id` int(11) unsigned DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `is_newsletter_allowed` tinyint(1) NOT NULL DEFAULT 0,
  `is_email_confirmed` tinyint(1) NOT NULL DEFAULT 0,
  `email_confirmation_token` char(24) NOT NULL,
  `timezone` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `users_external_system_id` (`source`,`external_system_id`),
  KEY `organization_id` (`organization_id`),
  CONSTRAINT `users_ibfk_1` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE `organizations` ADD `creator_id` int(11) unsigned DEFAULT NULL AFTER `data_location`;
ALTER TABLE `organizations` ADD KEY `creator_id` (`creator_id`);
ALTER TABLE `organizations` ADD CONSTRAINT `organizations_ibfk_1` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`);

-- Create syntax for TABLE 'invitations'
CREATE TABLE `invitations` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `invitation_code` varchar(255) NOT NULL,
  `organization_id` int(11) unsigned NOT NULL,
  `sender_id` int(11) unsigned NOT NULL,
  `accepter_id` int(11) unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `email` (`email`,`organization_id`),
  UNIQUE KEY `invitation_code` (`invitation_code`,`organization_id`),
  UNIQUE KEY `accepter_id` (`accepter_id`),
  KEY `organization_id` (`organization_id`),
  KEY `sender_id` (`sender_id`),
  CONSTRAINT `invitations_ibfk_1` FOREIGN KEY (`organization_id`) REFERENCES `organizations` (`id`),
  CONSTRAINT `invitations_ibfk_2` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`),
  CONSTRAINT `invitations_ibfk_3` FOREIGN KEY (`accepter_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
