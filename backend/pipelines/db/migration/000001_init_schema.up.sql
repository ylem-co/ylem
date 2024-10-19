-- Create syntax for TABLE 'aggregators'
CREATE TABLE `aggregators` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `expression` blob NOT NULL,
  `variable_name` varchar(255) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'api_calls'
CREATE TABLE `api_calls` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `type` varchar(255) NOT NULL,
  `payload` blob NOT NULL DEFAULT '',
  `query_string` blob NOT NULL DEFAULT '',
  `headers` blob NOT NULL DEFAULT '',
  `attached_file_name` varchar(255) NOT NULL DEFAULT '',
  `destination_uuid` varchar(36) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'codes'
CREATE TABLE `codes` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `type` varchar(255) NOT NULL DEFAULT 'python',
  `uuid` varchar(36) NOT NULL,
  `code` blob NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'conditions'
CREATE TABLE `conditions` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `expression` blob NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'env_variables'
CREATE TABLE `env_variables` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `organization_uuid` varchar(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `value` varchar(255) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `name_organization` (`organization_uuid`,`name`,`is_active`),
  KEY `org_uuid_index` (`organization_uuid`,`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'external_triggers'
CREATE TABLE `external_triggers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `test_data` blob DEFAULT '',
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'filters'
CREATE TABLE `filters` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `expression` blob NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'folders'
CREATE TABLE `folders` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `type` varchar(255) NOT NULL DEFAULT 'generic',
  `name` varchar(255) NOT NULL,
  `organization_uuid` varchar(36) NOT NULL,
  `parent_id` int(11) unsigned DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `name_in_parent_folder` (`parent_id`,`name`),
  KEY `name` (`name`),
  KEY `main_folder_selection_index` (`organization_uuid`,`parent_id`,`is_active`,`name`,`type`),
  CONSTRAINT `folders_ibfk_1` FOREIGN KEY (`parent_id`) REFERENCES `folders` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'for_eaches'
CREATE TABLE `for_eaches` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'gpts'
CREATE TABLE `gpts` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `prompt` blob NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'merges'
CREATE TABLE `merges` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `field_names` varchar(255) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'metrics_run_counts_monthly'
CREATE TABLE `metrics_run_counts_monthly` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `organization_uuid` varchar(36) NOT NULL,
  `year_month` varchar(255) NOT NULL,
  `run_count` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_organization_uuid_year_month` (`organization_uuid`,`year_month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'notifications'
CREATE TABLE `notifications` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `type` varchar(255) NOT NULL,
  `body` blob NOT NULL DEFAULT '',
  `attached_file_name` varchar(255) NOT NULL DEFAULT '',
  `destination_uuid` varchar(36) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'processors'
CREATE TABLE `processors` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `strategy` varchar(255) NOT NULL DEFAULT 'inclusive',
  `expression` blob NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'queries'
CREATE TABLE `queries` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `source_uuid` varchar(36) NOT NULL,
  `sql_query` blob NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'run_pipelines'
CREATE TABLE `run_pipelines` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `pipeline_uuid` varchar(36) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'scheduled_runs'
CREATE TABLE `scheduled_runs` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `pipeline_id` int(11) unsigned NOT NULL,
  `pipeline_run_uuid` varchar(36) DEFAULT NULL,
  `input` blob DEFAULT NULL,
  `env_vars` text DEFAULT NULL,
  `config` text DEFAULT NULL,
  `execute_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `execution_index` (`execute_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'shared_pipelines'
CREATE TABLE `shared_pipelines` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `pipeline_uuid` varchar(36) NOT NULL,
  `organization_uuid` varchar(36) NOT NULL,
  `creator_uuid` varchar(36) NOT NULL,
  `share_link` varchar(255) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `is_link_published` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_share_link` (`share_link`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'pipelines'
CREATE TABLE `pipelines` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `type` varchar(255) NOT NULL DEFAULT 'generic',
  `name` varchar(255) NOT NULL,
  `organization_uuid` varchar(36) NOT NULL,
  `creator_uuid` varchar(36) NOT NULL,
  `folder_id` int(11) unsigned DEFAULT NULL,
  `elements_layout` blob NOT NULL DEFAULT '',
  `preview` longblob NOT NULL,
  `schedule` varchar(255) NOT NULL DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `is_paused` tinyint(1) NOT NULL DEFAULT 0,
  `is_template` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `pipelines_ibfk_1` (`folder_id`),
  KEY `name_search_active_index` (`organization_uuid`,`is_active`,`name`),
  KEY `folder_active_index` (`organization_uuid`,`is_active`,`folder_id`),
  KEY `schedule_index` (`is_paused`,`is_active`,`schedule`),
  KEY `pipeline_active_type_index` (`organization_uuid`,`is_active`,`type`,`is_template`),
  CONSTRAINT `pipelines_ibfk_1` FOREIGN KEY (`folder_id`) REFERENCES `folders` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'tasks'
CREATE TABLE `tasks` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `severity` varchar(36) NOT NULL DEFAULT 'medium',
  `pipeline_id` int(11) unsigned NOT NULL,
  `type` varchar(255) NOT NULL,
  `implementation_id` int(10) unsigned DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  UNIQUE KEY `implementation_and_type` (`implementation_id`,`type`),
  KEY `pipeline_id` (`pipeline_id`),
  KEY `uuid_active_index` (`uuid`,`is_active`),
  KEY `pipeline_id_active_index` (`pipeline_id`,`is_active`),
  KEY `name_active_index` (`name`,`is_active`),
  KEY `type_active_index` (`type`,`is_active`),
  CONSTRAINT `tasks_ibfk_1` FOREIGN KEY (`pipeline_id`) REFERENCES `pipelines` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'task_run_results'
CREATE TABLE `task_run_results` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `state` varchar(255) NOT NULL,
  `task_id` int(11) unsigned NOT NULL,
  `task_run_uuid` varchar(36) NOT NULL,
  `pipeline_run_uuid` varchar(36) NOT NULL,
  `is_successful` tinyint(1) NOT NULL,
  `output` blob DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `task_id` (`task_id`),
  CONSTRAINT `task_run_results_fk_task_id` FOREIGN KEY (`task_id`) REFERENCES `tasks` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'task_run_errors'
CREATE TABLE `task_run_errors` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `task_run_result_id` int(11) unsigned NOT NULL,
  `code` int(11) NOT NULL,
  `severity` varchar(255) DEFAULT NULL,
  `message` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `task_run_errors_fk_task_run_result_id` (`task_run_result_id`),
  CONSTRAINT `task_run_errors_fk_task_run_result_id` FOREIGN KEY (`task_run_result_id`) REFERENCES `task_run_results` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'task_triggers'
CREATE TABLE `task_triggers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `pipeline_id` int(11) unsigned NOT NULL,
  `trigger_task_id` int(11) unsigned DEFAULT NULL,
  `triggered_task_id` int(11) unsigned DEFAULT NULL,
  `trigger_type` varchar(255) NOT NULL,
  `schedule` varchar(255) NOT NULL DEFAULT '',
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `schedule` (`schedule`,`trigger_type`,`is_active`),
  KEY `output_task_id` (`trigger_task_id`),
  KEY `input_task_id` (`triggered_task_id`),
  KEY `pipeline_id` (`pipeline_id`),
  CONSTRAINT `task_triggers_ibfk_1` FOREIGN KEY (`trigger_task_id`) REFERENCES `tasks` (`id`),
  CONSTRAINT `task_triggers_ibfk_2` FOREIGN KEY (`triggered_task_id`) REFERENCES `tasks` (`id`),
  CONSTRAINT `task_triggers_ibfk_3` FOREIGN KEY (`pipeline_id`) REFERENCES `pipelines` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'transformers'
CREATE TABLE `transformers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `type` varchar(255) NOT NULL,
  `json_query_expression` blob NOT NULL DEFAULT '',
  `delimiter` varchar(255) NOT NULL DEFAULT '',
  `decode_format` varchar(255) NOT NULL DEFAULT '',
  `encode_format` varchar(255) NOT NULL DEFAULT '',
  `cast_to_type` varchar(255) NOT NULL DEFAULT '',
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create syntax for TABLE 'pipeline_run_counts_monthly'
CREATE TABLE `pipeline_run_counts_monthly` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `organization_uuid` varchar(36) NOT NULL,
  `year_month` varchar(255) NOT NULL,
  `run_count` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_organization_uuid_year_month` (`organization_uuid`,`year_month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
