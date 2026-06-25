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

CREATE TABLE IF NOT EXISTS `goods_info` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '商品ID',
  `category_id` int unsigned NOT NULL COMMENT '商品分类ID',
  `name` varchar(128) NOT NULL COMMENT '商品名称',
  `image_url` varchar(255) NOT NULL DEFAULT '' COMMENT '商品图片',
  `points_price` int unsigned NOT NULL COMMENT '兑换所需积分',
  `stock` int unsigned NOT NULL DEFAULT 0 COMMENT '库存',
  `description` text COMMENT '商品简介',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态：1上架 0下架',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_goods_category_id` (`category_id`),
  KEY `idx_goods_status` (`status`),
  KEY `idx_goods_name` (`name`),
  KEY `idx_goods_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='商品信息表';

CREATE TABLE IF NOT EXISTS `cart_info` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '购物车ID',
  `employee_id` int unsigned NOT NULL COMMENT '员工ID',
  `goods_id` int unsigned NOT NULL COMMENT '商品ID',
  `count` int unsigned NOT NULL DEFAULT 1 COMMENT '商品数量',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_cart_employee_id` (`employee_id`),
  KEY `idx_cart_goods_id` (`goods_id`),
  KEY `idx_cart_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='购物车表';

CREATE TABLE IF NOT EXISTS `order_info` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Order ID',
  `order_no` varchar(32) NOT NULL COMMENT 'Order number',
  `employee_id` int unsigned NOT NULL COMMENT 'Employee ID',
  `total_points` int unsigned NOT NULL DEFAULT 0 COMMENT 'Total points',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT 'Status: 1 pending 2 completed 3 cancelled',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT 'Remark',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
  `deleted_at` datetime DEFAULT NULL COMMENT 'Deleted time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_order_no` (`order_no`),
  KEY `idx_order_employee_id` (`employee_id`),
  KEY `idx_order_status` (`status`),
  KEY `idx_order_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Points redemption order table';

CREATE TABLE IF NOT EXISTS `order_item` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Order item ID',
  `order_id` int unsigned NOT NULL COMMENT 'Order ID',
  `employee_id` int unsigned NOT NULL COMMENT 'Employee ID',
  `goods_id` int unsigned NOT NULL COMMENT 'Goods ID',
  `goods_name` varchar(128) NOT NULL COMMENT 'Goods name snapshot',
  `goods_image_url` varchar(255) NOT NULL DEFAULT '' COMMENT 'Goods image snapshot',
  `points_price` int unsigned NOT NULL COMMENT 'Points price snapshot',
  `count` int unsigned NOT NULL COMMENT 'Goods count',
  `total_points` int unsigned NOT NULL COMMENT 'Item total points',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  PRIMARY KEY (`id`),
  KEY `idx_order_item_order_id` (`order_id`),
  KEY `idx_order_item_employee_id` (`employee_id`),
  KEY `idx_order_item_goods_id` (`goods_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Points redemption order item table';

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

UPDATE `goods_info` g
JOIN `goods_category` old_c ON old_c.`id` = g.`category_id` AND old_c.`name` = '办公零食'
JOIN `goods_category` new_c ON new_c.`name` = 'Office Snacks'
SET g.`category_id` = new_c.`id`;

UPDATE `goods_info` g
JOIN `goods_category` old_c ON old_c.`id` = g.`category_id` AND old_c.`name` = '福利商品'
JOIN `goods_category` new_c ON new_c.`name` = 'Employee Benefits'
SET g.`category_id` = new_c.`id`;

UPDATE `goods_info` g
JOIN `goods_category` old_c ON old_c.`id` = g.`category_id` AND old_c.`name` = '办公用品'
JOIN `goods_category` new_c ON new_c.`name` = 'Office Supplies'
SET g.`category_id` = new_c.`id`;

DELETE FROM `goods_category`
WHERE `name` IN ('办公零食', '福利商品', '办公用品');

UPDATE `employee_info` SET `real_name` = 'Test Employee' WHERE `real_name` = '测试员工';
UPDATE `employee_info` SET `real_name` = 'Goods Manager' WHERE `real_name` = '商品管理员';
