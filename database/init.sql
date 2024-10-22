CREATE DATABASE IF NOT EXISTS users CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON users.* TO 'dtmnuser'@'%';

CREATE DATABASE IF NOT EXISTS api CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON api.* TO 'dtmnuser'@'%';

CREATE DATABASE IF NOT EXISTS integrations CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON integrations.* TO 'dtmnuser'@'%';

CREATE DATABASE IF NOT EXISTS pipelines CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON pipelines.* TO 'dtmnuser'@'%';

FLUSH PRIVILEGES;

SET GLOBAL max_allowed_packet=1073741824;

CREATE USER 'dtmntestuser'@'%' IDENTIFIED BY 'dtmntestpassword';

CREATE DATABASE IF NOT EXISTS trial_database CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON trial_database.* TO 'dtmntestuser'@'%';

CREATE DATABASE IF NOT EXISTS compliance CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON compliance.* TO 'dtmntestuser'@'%';

CREATE DATABASE IF NOT EXISTS customer_success CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON customer_success.* TO 'dtmntestuser'@'%';

CREATE DATABASE IF NOT EXISTS data_monitoring CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON data_monitoring.* TO 'dtmntestuser'@'%';

CREATE DATABASE IF NOT EXISTS ecommerce CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON ecommerce.* TO 'dtmntestuser'@'%';

CREATE DATABASE IF NOT EXISTS finance CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON finance.* TO 'dtmntestuser'@'%';

CREATE DATABASE IF NOT EXISTS hr CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
GRANT ALL PRIVILEGES ON hr.* TO 'dtmntestuser'@'%';

FLUSH PRIVILEGES;

USE `trial_database`;

CREATE TABLE `customers` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(255) DEFAULT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `IP` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `customers` (`id`, `email`, `first_name`, `last_name`, `status`, `IP`, `created_at`)
VALUES
  (1,'test@example.com','Emily','James','pending','178.5.92.123','2021-02-05 23:16:07'),
  (2,'test2@example.com','Jack','Lonely :)','waiting_for_approval','183.4.42.352','2021-02-05 23:17:06'),
  (3,'test3@example.com','Marc ;-)','Watson','phone_confirmed','3.45.22.11','2021-02-05 23:17:27'),
  (4,'test4@example.com','Liza','Conolly','profile_created','47.12.654.4','2021-02-05 23:17:56'),
  (5,'fraudster1@example.com','Bb','gg','complete','12.135.67.45','2022-02-07 21:47:51'),
  (6,'fraudster2@example.com','Bb','gg','complete','12.135.67.45','2022-02-07 21:47:53'),
  (7,'fraudster3@example.com','Bb','gg','complete','12.135.67.45','2022-02-07 21:47:55'),
  (8,'complete@example.com','Bill','Wilson','complete','3.45.22.12','2022-02-06 09:11:46');

CREATE TABLE `orders` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `amount` decimal(10,2) DEFAULT NULL,
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `customer_id` int unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `orders` (`id`, `amount`, `status`, `customer_id`, `created_at`)
VALUES
  (1,380.25,'shipped',8,'2021-02-05 22:39:32'),
  (2,1380.25,'shipped',8,'2021-02-12 22:39:32'),
  (3,5425.89,'shipped',8,'2021-03-05 22:39:32'),
  (4,10.44,'In_delivery',8,'2021-12-05 22:39:32'),
  (5,110.44,'In_delivery',8,'2021-12-05 22:39:32');

CREATE TABLE `payments` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `amount` decimal(10,2) DEFAULT NULL,
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `order_id` int unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `payments` (`id`, `amount`, `status`, `order_id`, `created_at`)
VALUES
  (1,180.25,'paid',1,'2021-02-05 22:43:42'),
  (2,1380.25,'pending',2,'2022-02-05 22:44:41'),
  (3,200.05,'paid',1,'2021-02-05 22:43:42'),
  (4,5425.89,'paid',3,'2022-02-05 22:45:58'),
  (5,10.44,'failed',4,'2022-02-05 22:46:08');

USE `compliance`;

CREATE TABLE `organizations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `organizations` (`id`, `name`, `status`, `created_at`, `updated_at`)
VALUES
  (1,'Delivery scale Ltd','complete','2022-10-01 14:31:43','2022-10-01 14:31:43'),
  (2,'Security scale Ltd','complete','2022-10-01 14:31:56','2022-10-01 14:31:56'),
  (3,'Mountain one GmbH','new','2022-10-01 14:32:41','2022-10-01 14:32:41'),
  (4,'Surfing camp Ltd','complete','2022-10-01 14:32:57','2022-10-01 14:32:57'),
  (5,'Organization UG','verification_pending','2022-10-01 14:33:11','2022-10-01 14:33:11');

CREATE TABLE `documents` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `type` varchar(255) DEFAULT NULL,
  `organization_id` int unsigned DEFAULT NULL,
  `uploaded_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `organization_id` (`organization_id`),

  CONSTRAINT `documents_ibfk_1` FOREIGN KEY (organization_id) REFERENCES organizations(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `documents` (`id`, `type`, `organization_id`, `uploaded_at`)
VALUES
  (1,'incorporation_certificate',1,'2022-10-01 15:43:05'),
  (2,'incorporation_certificate',2,'2022-10-01 15:43:17'),
  (3,'shareholders_list',1,'2022-10-01 15:43:41'),
  (4,'shareholders_list',2,'2022-10-01 15:43:41'),
  (5,'shareholders_list',3,'2022-10-01 15:43:41'),
  (6,'shareholders_list',4,'2022-10-01 15:43:41'),
  (9,'incorporation_certificate',5,'2022-10-01 15:43:17'),
  (12,'incorporation_certificate',1,'2024-03-26 16:46:00');

CREATE TABLE `kycs` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` int unsigned DEFAULT NULL,
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `started_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `completed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `organization_id` (`organization_id`),
  CONSTRAINT `kycs_ibfk_1` FOREIGN KEY (organization_id) REFERENCES organizations(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `kycs` (`id`, `organization_id`, `status`, `started_at`, `completed_at`)
VALUES
  (1,1,'new','2022-10-01 15:40:34',NULL),
  (2,2,'complete','2022-10-01 15:40:51','2022-10-01 15:40:51'),
  (3,3,'complete','2022-10-01 15:40:57','2022-10-01 15:40:57'),
  (4,4,'request_sent','2022-10-01 15:41:05',NULL),
  (5,5,'failed','2022-10-01 15:41:23',NULL);

CREATE TABLE `questionaires` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` int unsigned DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `uodated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `organization_id` (`organization_id`),
  CONSTRAINT `questionaires_ibfk_1` FOREIGN KEY (organization_id) REFERENCES organizations(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `questionaires` (`id`, `organization_id`, `status`, `created_at`, `uodated_at`)
VALUES
  (1,1,'filled','2022-10-01 15:49:45','2022-10-01 15:49:45'),
  (2,2,'pending','2022-10-01 15:49:52','2022-10-01 15:49:52'),
  (3,3,'filled','2022-10-01 15:50:00','2022-10-01 15:50:00'),
  (4,4,'filled','2022-10-01 15:50:06','2022-10-01 15:50:06'),
  (5,5,'pending','2022-10-01 15:50:11','2022-10-01 15:50:11');

USE `customer_success`;

CREATE TABLE `organizations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `creator_id` int unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `creator_id` (`creator_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `organizations` (`id`, `name`, `creator_id`, `created_at`, `updated_at`)
VALUES
  (1,'Datacom Ltd',1,'2022-10-01 14:11:05','2022-10-01 14:11:05'),
  (2,'Delivery moon GmbH',4,'2022-10-01 14:11:20','2022-10-01 14:11:20'),
  (3,'Security Cloud Ltd',5,'2022-10-01 14:11:31','2022-10-01 14:11:31'),
  (4,'Step One AG',8,'2022-10-01 14:11:42','2022-10-01 14:11:42'),
  (5,'Recovery Side GmbH',9,'2022-10-01 14:11:54','2022-10-01 14:11:54'),
  (6,'Veggies and fruits UG',10,'2022-10-01 14:12:17','2022-10-01 14:12:17');

CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` int unsigned DEFAULT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `actions` int DEFAULT NULL,
  `last_time_logged_in` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `organization_id` (`organization_id`),
  CONSTRAINT `users_ibfk_1` FOREIGN KEY (organization_id) REFERENCES organizations(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `users` (`id`, `organization_id`, `first_name`, `last_name`, `email`, `status`, `actions`, `last_time_logged_in`, `created_at`, `updated_at`)
VALUES
  (1,1,'Liza','Jensen','liza@example.com','complete',567,'2022-10-01 14:12:52','2022-10-01 14:12:52','2022-10-01 14:12:52'),
  (2,1,'John','Meyer','john@example.com','verification_pending',1,'2022-10-01 14:13:41','2022-10-01 14:13:41','2022-10-01 14:13:41'),
  (3,1,'Carlo','Miglione','carlo@example.com','complete',11,'2022-10-01 14:13:58','2022-10-01 14:13:58','2022-10-01 14:13:58'),
  (4,2,'Rachel','Rocks','rachel@example.com','complete',110,'2022-10-01 14:14:38','2022-10-01 14:14:38','2022-10-01 14:14:38'),
  (5,3,'David','Lee','david@example.com','complete',24,'2022-10-01 14:14:56','2022-10-01 14:14:56','2022-10-01 14:14:56'),
  (6,3,'Jessica','Lambert','jessica@example.com','new',35,'2022-10-01 14:15:27','2022-10-01 14:15:27','2022-10-01 14:15:27'),
  (7,3,'Killian','Chee','killian@example.com','verification_pending',0,'2022-10-01 14:16:08','2022-10-01 14:16:08','2022-10-01 14:16:08'),
  (8,4,'Bianca','Mele','bianca@example.com','complete',59,'2022-10-01 14:16:34','2022-10-01 14:16:34','2022-10-01 14:16:34'),
  (9,5,'Jorg','Schulze','jorg@example.com','complete',70,'2022-10-01 14:17:35','2022-10-01 14:17:35','2022-10-01 14:17:35'),
  (10,6,'Vera','Magnus','vera@example.com','complete',48,'2022-10-01 14:18:11','2022-10-01 14:18:11','2022-10-01 14:18:11'),
  (11,6,'Jan','Johansson','jan@example.com','verification_pending',34,'2022-10-01 14:18:42','2022-10-01 14:18:42','2022-10-01 14:18:42');

USE `data_monitoring`;

CREATE TABLE `third_party_scoring_data` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `organization_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `type` varchar(255) DEFAULT NULL,
  `provider` varchar(255) DEFAULT NULL,
  `score` decimal(5,2) DEFAULT NULL,
  `component_1` decimal(5,2) DEFAULT NULL,
  `component_2` decimal(5,2) DEFAULT NULL,
  `component_3` decimal(5,2) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `third_party_scoring_data` (`id`, `organization_name`, `type`, `provider`, `score`, `component_1`, `component_2`, `component_3`, `created_at`)
VALUES
  (1,'Organization 1','company_score','north_data',67.00,50.00,35.00,85.00,'2022-10-01 14:25:20'),
  (2,'Organization 2','company_score','schufa',49.00,2.00,NULL,87.00,'2022-10-01 14:25:52'),
  (3,'Organization 3','personal_score','scorer',96.00,98.00,91.00,93.00,'2022-10-01 14:26:13'),
  (4,'Organization 4','credit_score','crefo',46.00,21.00,22.00,76.00,'2022-10-01 14:26:37'),
  (5,'Organization 5','personal_score','scorer',84.00,80.00,71.00,69.00,'2022-10-01 14:26:57'),
  (6,'Organization 6','credit_score','crefo',51.00,55.00,62.00,26.00,'2022-10-01 14:27:15'),
  (7,'Organization 7','company_score','north_data',NULL,NULL,25.00,42.00,'2022-10-01 14:27:32');

CREATE TABLE `units` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `units` (`id`, `name`, `status`, `created_at`, `updated_at`)
VALUES
  (1,'First unit','new','2022-10-01 14:20:54','2022-10-01 14:20:54'),
  (2,'Another unit','done','2022-10-01 14:21:12','2022-10-01 14:21:12'),
  (3,'Unit three','done','2022-10-01 14:21:21','2022-10-01 14:21:21'),
  (4,'And one more ','pending','2022-10-01 14:21:31','2022-10-01 14:21:31'),
  (5,'Created recently','new','2022-10-01 14:21:39','2022-10-01 14:21:39'),
  (6,'Complete unit','done','2022-10-01 14:21:53','2022-10-01 14:21:53');

USE `ecommerce`;

CREATE TABLE `orders` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `status` varchar(255) DEFAULT NULL,
  `organization_id` int unsigned DEFAULT NULL,
  `amount` decimal(10,2) DEFAULT NULL,
  `currency` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `orders` (`id`, `status`, `organization_id`, `amount`, `currency`, `created_at`, `updated_at`)
VALUES
  (1,'dispatched',1,250.00,'USD','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (2,'paid',1,250.00,'USD','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (3,'completed',2,3677.20,'USD','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (4,'assembled',1,21.87,'EUR','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (5,'assembled',4,21.87,'EUR','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (6,'complete',1,341.76,'EUR','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (7,'new',3,2110.76,'USD','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (8,'new',3,127.89,'USD','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (9,'new',4,127.90,'EUR','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (10,'paid',5,344.44,'USD','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (11,'paid',1,1344.44,'USD','2022-10-01 16:01:11','2022-10-01 16:01:11'),
  (12,'assembled',2,67.01,'EUR','2022-10-01 16:01:11','2022-10-01 16:01:11');

CREATE TABLE `organizations` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `creator_id` int unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `creator_id` (`creator_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `organizations` (`id`, `name`, `creator_id`, `created_at`, `updated_at`)
VALUES
  (1,'Datacom Ltd',1,'2022-10-01 14:11:05','2022-10-01 14:11:05'),
  (2,'Delivery moon GmbH',4,'2022-10-01 14:11:20','2022-10-01 14:11:20'),
  (3,'Security Cloud Ltd',5,'2022-10-01 14:11:31','2022-10-01 14:11:31'),
  (4,'Step One AG',8,'2022-10-01 14:11:42','2022-10-01 14:11:42'),
  (5,'Recovery Side GmbH',9,'2022-10-01 14:11:54','2022-10-01 14:11:54'),
  (6,'Veggies and fruits UG',10,'2022-10-01 14:12:17','2022-10-01 14:12:17');

CREATE TABLE `products` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `num_in_stock` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `products` (`id`, `name`, `num_in_stock`)
VALUES
  (1,'Pencil',50),
  (2,'Book',3),
  (3,'Notebook',78),
  (4,'Eraser',4);

CREATE TABLE `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `organization_id` int unsigned DEFAULT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `actions` int DEFAULT NULL,
  `last_time_logged_in` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `organization_id` (`organization_id`),
  CONSTRAINT `users_ibfk_1` FOREIGN KEY (organization_id) REFERENCES organizations(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `users` (`id`, `organization_id`, `first_name`, `last_name`, `email`, `status`, `actions`, `last_time_logged_in`, `created_at`, `updated_at`)
VALUES
  (1,1,'Liza','Jensen','liza@example.com','complete',567,'2022-10-01 14:12:52','2022-10-01 14:12:52','2022-10-01 14:12:52'),
  (2,1,'John','Meyer','john@example.com','verification_pending',1,'2022-10-01 14:13:41','2022-10-01 14:13:41','2022-10-01 14:13:41'),
  (3,1,'Carlo','Miglione','carlo@example.com','complete',11,'2022-10-01 14:13:58','2022-10-01 14:13:58','2022-10-01 14:13:58'),
  (4,2,'Rachel','Rocks','rachel@example.com','complete',110,'2022-10-01 14:14:38','2022-10-01 14:14:38','2022-10-01 14:14:38'),
  (5,3,'David','Lee','david@example.com','complete',24,'2022-10-01 14:14:56','2022-10-01 14:14:56','2022-10-01 14:14:56'),
  (6,3,'Jessica','Lambert','jessica@example.com','new',35,'2022-10-01 14:15:27','2022-10-01 14:15:27','2022-10-01 14:15:27'),
  (7,3,'Killian','Chee','killian@example.com','verification_pending',0,'2022-10-01 14:16:08','2022-10-01 14:16:08','2022-10-01 14:16:08'),
  (8,4,'Bianca','Mele','bianca@example.com','complete',59,'2022-10-01 14:16:34','2022-10-01 14:16:34','2022-10-01 14:16:34'),
  (9,5,'Jorg','Schulze','jorg@example.com','complete',70,'2022-10-01 14:17:35','2022-10-01 14:17:35','2022-10-01 14:17:35'),
  (10,6,'Vera','Magnus','vera@example.com','complete',48,'2022-10-01 14:18:11','2022-10-01 14:18:11','2022-10-01 14:18:11'),
  (11,6,'Jan','Johansson','jan@example.com','verification_pending',34,'2022-10-01 14:18:42','2022-10-01 14:18:42','2022-10-01 14:18:42');

USE `finance`;

CREATE TABLE `customers` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `customers` (`id`, `name`, `status`, `created_at`, `updated_at`)
VALUES
  (1,'Detailed Payments Ltd.','complete','2022-09-30 18:37:25','2022-09-30 18:37:25'),
  (2,'Startup Merchant GmbH','verification_pending','2022-09-30 18:37:57','2022-09-30 18:37:57'),
  (3,'Delivery space Ltd.','complete','2022-09-30 18:37:25','2022-09-30 18:37:25'),
  (4,'Existing stars AG','new','2022-09-30 18:38:37','2022-09-30 18:38:37'),
  (5,'Banking owners UG','complete','2022-09-30 18:38:52','2022-09-30 18:38:52');

CREATE TABLE `invoices` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `customer_from_id` int unsigned DEFAULT NULL,
  `customer_to_id` int unsigned DEFAULT NULL,
  `amount` decimal(8,2) DEFAULT NULL,
  `currency` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `customer_from_id` (`customer_from_id`),
  KEY `customer_to_id` (`customer_to_id`),
  CONSTRAINT `invoices_ibfk_1` FOREIGN KEY (customer_from_id) REFERENCES customers(id),
  CONSTRAINT `invoices_ibfk_2` FOREIGN KEY (customer_to_id) REFERENCES customers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `invoices` (`id`, `customer_from_id`, `customer_to_id`, `amount`, `currency`, `created_at`, `updated_at`)
VALUES
  (1,1,2,299.76,'EUR','2022-09-30 18:53:35','2022-09-30 18:53:35'),
  (2,1,2,11.01,'EUR','2022-09-30 18:57:35','2022-09-30 18:53:35'),
  (3,1,2,11.01,'EUR','2022-09-30 18:53:35','2022-09-30 18:53:35'),
  (4,3,5,1300.00,'USD','2022-09-30 18:54:26','2022-09-30 18:54:26'),
  (5,3,5,1300.00,'USD','2022-09-30 18:54:26','2022-09-30 18:54:26'),
  (6,5,3,1504.36,'USD','2022-09-30 18:54:39','2022-09-30 18:54:39'),
  (7,4,1,34.89,'EUR','2022-09-30 18:55:10','2022-09-30 18:55:10'),
  (8,1,4,564.89,'EUR','2022-09-30 18:55:28','2022-09-30 18:55:28'),
  (9,1,4,1567.51,'EUR','2022-09-30 18:55:28','2022-09-30 18:55:28'),
  (10,2,1,21.11,'EUR','2022-09-30 18:56:06','2022-09-30 18:56:06'),
  (11,3,5,34.35,'USD','2022-09-30 18:56:26','2022-09-30 18:56:26'),
  (12,5,3,112.56,'USD','2022-09-30 18:56:44','2022-09-30 18:56:44');

CREATE TABLE `bank_accounts` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `iban` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `owner_id` int unsigned DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `owner_id` (`owner_id`),
  CONSTRAINT `bank_accounts_ibfk_1` FOREIGN KEY (owner_id) REFERENCES customers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `bank_accounts` (`id`, `iban`, `status`, `owner_id`, `created_at`, `updated_at`)
VALUES
  (1,'NL22ABNA9242086444','verified',1,'2022-09-30 18:39:38','2022-09-30 18:39:38'),
  (2,'DE28500105175143257851','pending',2,'2022-09-30 18:39:57','2022-09-30 18:39:57'),
  (3,'IT10E0300203280154176113663','verified',3,'2022-09-30 18:40:18','2022-09-30 18:40:18'),
  (4,'ES0800758582392111774524','pending',4,'2022-09-30 18:40:36','2022-09-30 18:40:36'),
  (5,'SE1693857923159423561388','verified',5,'2022-09-30 18:41:03','2022-09-30 18:41:03');

CREATE TABLE `bank_transfers` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `account_from_id` int unsigned DEFAULT NULL,
  `account_to_id` int unsigned DEFAULT NULL,
  `invoice_id` int unsigned DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `amount` decimal(8,2) DEFAULT NULL,
  `currency` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `account_from_id` (`account_from_id`),
  KEY `account_to_id` (`account_to_id`),
  KEY `invoice_id` (`invoice_id`),
  CONSTRAINT `bank_transfers_ibfk_1` FOREIGN KEY (account_from_id) REFERENCES customers(id),
  CONSTRAINT `bank_transfers_ibfk_2` FOREIGN KEY (account_to_id) REFERENCES customers(id),
  CONSTRAINT `bank_transfers_ibfk_3` FOREIGN KEY (invoice_id) REFERENCES invoices(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `bank_transfers` (`id`, `account_from_id`, `account_to_id`, `invoice_id`, `status`, `amount`, `currency`, `created_at`, `updated_at`)
VALUES
  (1,1,2,1,'complete',299.76,'EUR','2022-09-30 18:53:35','2022-09-30 18:53:35'),
  (2,1,2,2,'complete',11.01,'EUR','2022-09-30 18:57:35','2022-09-30 18:53:35'),
  (3,1,2,3,'complete',11.01,'EUR','2022-09-30 18:53:35','2022-09-30 18:53:35'),
  (4,3,5,4,'failed',1300.00,'USD','2022-09-30 18:54:26','2022-09-30 18:54:26'),
  (5,3,5,5,'failed',1300.00,'USD','2022-09-30 18:54:26','2022-09-30 18:54:26'),
  (6,5,3,6,'pending',1504.36,'USD','2022-09-30 18:54:39','2022-09-30 18:54:39'),
  (7,4,1,7,'complete',34.89,'EUR','2022-09-30 18:55:10','2022-09-30 18:55:10'),
  (8,1,4,8,'pending',564.89,'EUR','2022-09-30 18:55:28','2022-09-30 18:55:28'),
  (9,1,4,9,'pending',1567.51,'EUR','2022-09-30 18:55:28','2022-09-30 18:55:28'),
  (10,2,1,10,'complete',21.11,'EUR','2022-09-30 18:56:06','2022-09-30 18:56:06'),
  (11,3,5,11,'complete',34.35,'USD','2022-09-30 18:56:26','2022-09-30 18:56:26'),
  (12,5,3,12,'pending',112.56,'USD','2022-09-30 18:56:44','2022-09-30 18:56:44');

USE `hr`;

CREATE TABLE `hires` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `first_name` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `department` varchar(255) DEFAULT NULL,
  `position` varchar(255) DEFAULT NULL,
  `hired_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

INSERT INTO `hires` (`id`, `first_name`, `last_name`, `department`, `position`, `hired_at`)
VALUES
  (1,'Emily','Hughes','Engineering','Senior Software Engineer','2023-01-13 13:43:42'),
  (2,'George','Dawson','Engineering','QA','2023-01-13 13:43:54'),
  (3,'Maria','Rodrigues','Marketing','PR specialist','2023-01-13 13:44:18'),
  (4,'Patrick','Mueller','Sales','Account manager','2023-01-13 13:44:32'),
  (5,'Zoe','Walter','Sales','Account manager','2023-01-13 13:44:47'),
  (6,'Rachel','Queen','Product','Senior Product Manager','2023-01-13 13:45:06');
