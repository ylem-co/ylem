-- Create syntax for TABLE 'integrations'
CREATE TABLE `integrations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(36) NOT NULL,
  `creator_uuid` char(36) NOT NULL,
  `organization_uuid` char(36) NOT NULL,
  `type` varchar(255) DEFAULT NULL,
  `io_type` varchar(255) NOT NULL DEFAULT 'write',
  `status` varchar(255) NOT NULL DEFAULT '',
  `name` varchar(512) NOT NULL DEFAULT '',
  `value` blob NOT NULL,
  `user_updated_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `uuid_active_index` (`uuid`,`deleted_at`),
  KEY `organization_uuid_active_index` (`organization_uuid`,`deleted_at`),
  KEY `organization_io_type_index` (`organization_uuid`,`io_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'apis'
CREATE TABLE `apis` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `method` varchar(255) NOT NULL DEFAULT 'post',
  `auth_type` varchar(32) NOT NULL,
  `auth_bearer_token` blob DEFAULT NULL,
  `auth_header_value` blob DEFAULT NULL,
  `auth_header_name` varchar(255) DEFAULT NULL,
  `auth_basic_user_password` varchar(255) DEFAULT NULL,
  `auth_basic_user_name` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `apis_integrations` (`integration_id`),
  CONSTRAINT `apis_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'emails'
CREATE TABLE `emails` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `code` char(6) NOT NULL,
  `is_confirmed` tinyint(1) NOT NULL DEFAULT 0,
  `requested_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `email_integrations` (`integration_id`),
  CONSTRAINT `email_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'google_sheets'
CREATE TABLE `google_sheets` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `spreadsheet_id` varchar(255) NOT NULL DEFAULT '',
  `sheet_id` int(10) NOT NULL DEFAULT 0,
  `mode` varchar(255) NOT NULL DEFAULT 'overwrite',
  `credentials` blob NOT NULL,
  `write_header` tinyint(1) NOT NULL DEFAULT 1,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `gs_integrations` (`integration_id`),
  CONSTRAINT `gs_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'hubspot_authorizations'
CREATE TABLE `hubspot_authorizations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(36) NOT NULL,
  `creator_uuid` char(36) NOT NULL,
  `organization_uuid` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `state` char(64) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `access_token` blob DEFAULT NULL,
  `access_token_expires_at` timestamp NULL DEFAULT NULL,
  `refresh_token` blob DEFAULT NULL,
  `scopes` varchar(32) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `state` (`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'hubspots'
CREATE TABLE `hubspots` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `hubspot_authorization_id` int(10) unsigned NOT NULL,
  `pipeline_stage_code` varchar(255) NOT NULL DEFAULT '',
  `owner_code` varchar(255) NOT NULL DEFAULT '',
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `hubspot_authorization_id` (`hubspot_authorization_id`),
  KEY `hubspot_integrations` (`integration_id`),
  CONSTRAINT `hubspot_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`),
  CONSTRAINT `hubspots_jira_authorization_id` FOREIGN KEY (`hubspot_authorization_id`) REFERENCES `hubspot_authorizations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'incident_ios'
CREATE TABLE `incident_ios` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `api_key` blob NOT NULL,
  `visibility` varchar(255) NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `incident_io_integrations` (`integration_id`),
  CONSTRAINT `incident_io_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'jenkinses'
CREATE TABLE `jenkinses` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `base_url` varchar(2048) NOT NULL,
  `token` blob NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `jenkins_integrations` (`integration_id`),
  CONSTRAINT `jenkins_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'jira_authorizations'
CREATE TABLE `jira_authorizations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(36) NOT NULL,
  `creator_uuid` char(36) NOT NULL,
  `organization_uuid` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `state` char(64) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `access_token` blob DEFAULT NULL,
  `scopes` varchar(32) DEFAULT NULL,
  `cloudid` char(36) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `state` (`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'jiras'
CREATE TABLE `jiras` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `jira_authorization_id` int(10) unsigned NOT NULL,
  `issue_type` varchar(255) NOT NULL DEFAULT '',
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `jira_authorization_id` (`jira_authorization_id`),
  KEY `jira_integrations` (`integration_id`),
  CONSTRAINT `jira_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`),
  CONSTRAINT `jiras_jira_authorization_id` FOREIGN KEY (`jira_authorization_id`) REFERENCES `jira_authorizations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'opsgenies'
CREATE TABLE `opsgenies` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `api_key` blob NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `opsgenie_integrations` (`integration_id`),
  CONSTRAINT `opsgenie_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'salesforce_authorizations'
CREATE TABLE `salesforce_authorizations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(36) NOT NULL,
  `creator_uuid` char(36) NOT NULL,
  `organization_uuid` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `state` char(64) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `access_token` blob DEFAULT NULL,
  `refresh_token` blob DEFAULT NULL,
  `scopes` varchar(32) DEFAULT NULL,
  `domain` varchar(256) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `updated_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `state` (`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'salesforces'
CREATE TABLE `salesforces` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `salesforce_authorization_id` int(10) unsigned NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `salesforce_authorization_id` (`salesforce_authorization_id`),
  KEY `salesforce_integrations` (`integration_id`),
  CONSTRAINT `salesforce_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`),
  CONSTRAINT `salesforces_salesforce_authorization_id` FOREIGN KEY (`salesforce_authorization_id`) REFERENCES `salesforce_authorizations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'slack_authorizations'
CREATE TABLE `slack_authorizations` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` char(36) NOT NULL,
  `creator_uuid` char(36) NOT NULL,
  `organization_uuid` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `state` char(64) NOT NULL,
  `is_active` tinyint(1) NOT NULL,
  `access_token` varchar(255) DEFAULT NULL,
  `scopes` varchar(255) DEFAULT NULL,
  `bot_user_id` varchar(32) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `state` (`state`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'slacks'
CREATE TABLE `slacks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `slack_authorization_id` int(10) unsigned NOT NULL,
  `slack_channel_id` varchar(32) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `slacks_slack_authorizations` (`slack_authorization_id`),
  KEY `slack_integrations` (`integration_id`),
  CONSTRAINT `slack_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`),
  CONSTRAINT `slacks_slack_authorizations` FOREIGN KEY (`slack_authorization_id`) REFERENCES `slack_authorizations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'smses'
CREATE TABLE `smses` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `code` char(6) NOT NULL,
  `is_confirmed` tinyint(1) NOT NULL DEFAULT 0,
  `requested_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `sms_integrations` (`integration_id`),
  CONSTRAINT `sms_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'sqls'
CREATE TABLE `sqls` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `type` varchar(64) NOT NULL,
  `data_key` blob DEFAULT NULL,
  `host` blob NOT NULL,
  `port` int(6) NOT NULL DEFAULT 3306,
  `user` varchar(255) NOT NULL,
  `password` blob DEFAULT NULL,
  `database` varchar(255) NOT NULL DEFAULT '',
  `connection_type` varchar(64) NOT NULL DEFAULT 'direct',
  `ssl_enabled` tinyint(1) NOT NULL DEFAULT 0,
  `ssh_host` blob NOT NULL,
  `ssh_port` int(6) NOT NULL DEFAULT 22,
  `ssh_user` varchar(255) NOT NULL DEFAULT '',
  `project_id` varchar(255) DEFAULT NULL,
  `credentials` blob DEFAULT NULL,
  `es_version` tinyint(4) DEFAULT NULL,
  `is_trial` tinyint(4) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `sql_integrations` (`integration_id`),
  CONSTRAINT `sql_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- Create syntax for TABLE 'tableau'
CREATE TABLE `tableau` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `integration_id` int(10) unsigned NOT NULL,
  `username` blob NOT NULL,
  `password` blob NOT NULL,
  `site_name` varchar(255) NOT NULL,
  `project_name` varchar(255) NOT NULL,
  `datasource_name` varchar(255) NOT NULL,
  `mode` varchar(255) NOT NULL,
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `tableau_integrations` (`integration_id`),
  CONSTRAINT `tableau_integrations` FOREIGN KEY (`integration_id`) REFERENCES `integrations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
