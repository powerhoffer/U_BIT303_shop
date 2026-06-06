CREATE DATABASE IF NOT EXISTS `bit303_shop` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE `bit303_shop`;

CREATE TABLE IF NOT EXISTS `employee_info` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '员工ID',
  `username` varchar(64) NOT NULL COMMENT '登录账号',
  `password_hash` varchar(100) NOT NULL COMMENT 'bcrypt密码哈希',
  `real_name` varchar(64) NOT NULL DEFAULT '' COMMENT '员工姓名',
  `phone` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
  `email` varchar(128) NOT NULL DEFAULT '' COMMENT '邮箱',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1正常 0禁用',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_employee_username` (`username`),
  KEY `idx_employee_status` (`status`),
  KEY `idx_employee_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='员工账号表';

CREATE TABLE IF NOT EXISTS `employee_points_account` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '积分账户ID',
  `employee_id` int unsigned NOT NULL COMMENT '员工ID',
  `balance` int unsigned NOT NULL DEFAULT 0 COMMENT '当前可用积分',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1正常 0停用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_points_account_employee` (`employee_id`),
  KEY `idx_points_account_status` (`status`),
  KEY `idx_points_account_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='员工积分账户表';

CREATE TABLE IF NOT EXISTS `employee_points_record` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '积分流水ID',
  `employee_id` int unsigned NOT NULL COMMENT '员工ID',
  `change_type` tinyint NOT NULL COMMENT '变动类型：1增加 2扣除',
  `points` int unsigned NOT NULL COMMENT '变动积分',
  `before_balance` int unsigned NOT NULL COMMENT '变动前积分',
  `after_balance` int unsigned NOT NULL COMMENT '变动后积分',
  `operator_employee_id` int unsigned NOT NULL DEFAULT 0 COMMENT '操作员工ID',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_points_record_employee_created` (`employee_id`, `created_at`),
  KEY `idx_points_record_operator` (`operator_employee_id`),
  KEY `idx_points_record_change_type` (`change_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='员工积分流水表';

CREATE TABLE IF NOT EXISTS `goods_category` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '商品分类ID',
  `name` varchar(64) NOT NULL COMMENT '分类名称',
  `sort` int unsigned NOT NULL DEFAULT 0 COMMENT '排序值',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1启用 0停用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_goods_category_name` (`name`),
  KEY `idx_goods_category_status` (`status`),
  KEY `idx_goods_category_sort` (`sort`),
  KEY `idx_goods_category_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='商品分类表';

INSERT INTO `goods_category` (`name`, `sort`, `status`) VALUES
('办公零食', 1, 1),
('福利商品', 2, 1),
('办公用品', 3, 1)
ON DUPLICATE KEY UPDATE
  `sort` = VALUES(`sort`),
  `status` = VALUES(`status`),
  `deleted_at` = NULL;
