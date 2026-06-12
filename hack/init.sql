CREATE DATABASE IF NOT EXISTS `bit303_shop` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE `bit303_shop`;

CREATE TABLE IF NOT EXISTS `employee_info` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Employee ID',
  `username` varchar(64) NOT NULL COMMENT 'Login username',
  `password_hash` varchar(100) NOT NULL COMMENT 'BCrypt password hash',
  `real_name` varchar(64) NOT NULL DEFAULT '' COMMENT 'Employee name',
  `phone` varchar(20) NOT NULL DEFAULT '' COMMENT 'Phone number',
  `email` varchar(128) NOT NULL DEFAULT '' COMMENT 'Email address',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT 'Status: 1 active, 0 disabled',
  `last_login_at` datetime DEFAULT NULL COMMENT 'Last login time',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_employee_username` (`username`),
  KEY `idx_employee_status` (`status`),
  KEY `idx_employee_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Employee accounts';

CREATE TABLE IF NOT EXISTS `employee_points_account` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Credit account ID',
  `employee_id` int unsigned NOT NULL COMMENT 'Employee ID',
  `balance` int unsigned NOT NULL DEFAULT 0 COMMENT 'Available credit balance',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT 'Status: 1 active, 0 disabled',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_points_account_employee` (`employee_id`),
  KEY `idx_points_account_status` (`status`),
  KEY `idx_points_account_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Employee credit accounts';

CREATE TABLE IF NOT EXISTS `employee_points_record` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Credit record ID',
  `employee_id` int unsigned NOT NULL COMMENT 'Employee ID',
  `change_type` tinyint NOT NULL COMMENT 'Change type: 1 add, 2 deduct',
  `points` int unsigned NOT NULL COMMENT 'Changed credits',
  `before_balance` int unsigned NOT NULL COMMENT 'Balance before change',
  `after_balance` int unsigned NOT NULL COMMENT 'Balance after change',
  `operator_employee_id` int unsigned NOT NULL DEFAULT 0 COMMENT 'Operator employee ID',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT 'Remark',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  PRIMARY KEY (`id`),
  KEY `idx_points_record_employee_created` (`employee_id`, `created_at`),
  KEY `idx_points_record_operator` (`operator_employee_id`),
  KEY `idx_points_record_change_type` (`change_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Employee credit records';

CREATE TABLE IF NOT EXISTS `goods_category` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Goods category ID',
  `name` varchar(64) NOT NULL COMMENT 'Category name',
  `sort` int unsigned NOT NULL DEFAULT 0 COMMENT 'Sort order',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT 'Status: 1 enabled, 0 disabled',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_goods_category_name` (`name`),
  KEY `idx_goods_category_status` (`status`),
  KEY `idx_goods_category_sort` (`sort`),
  KEY `idx_goods_category_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Goods categories';

INSERT INTO `employee_info` (`username`, `password_hash`, `real_name`, `phone`, `email`, `status`) VALUES
('root', '$2a$10$wkJo.7jih/0EbEehrNG.seMN5Rm3VZP90xxlK6bebLZDoq5K77W8C', 'System Administrator', '', '', 1)
ON DUPLICATE KEY UPDATE
  `password_hash` = VALUES(`password_hash`),
  `real_name` = VALUES(`real_name`),
  `status` = VALUES(`status`),
  `deleted_at` = NULL;

INSERT INTO `goods_category` (`name`, `sort`, `status`) VALUES
('Office Snacks', 1, 1),
('Employee Benefits', 2, 1),
('Office Supplies', 3, 1)
ON DUPLICATE KEY UPDATE
  `sort` = VALUES(`sort`),
  `status` = VALUES(`status`),
  `deleted_at` = NULL;
