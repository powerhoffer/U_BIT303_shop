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
