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
  `operator_admin_id` int unsigned NOT NULL DEFAULT 0 COMMENT 'Operator admin ID',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_points_record_employee_created` (`employee_id`, `created_at`),
  KEY `idx_points_record_operator` (`operator_employee_id`),
  KEY `idx_points_record_change_type` (`change_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='员工积分流水表';

SET @points_operator_admin_column_exists = (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'employee_points_record'
    AND COLUMN_NAME = 'operator_admin_id'
);
SET @points_operator_admin_column_sql = IF(
  @points_operator_admin_column_exists = 0,
  'ALTER TABLE `employee_points_record` ADD COLUMN `operator_admin_id` int unsigned NOT NULL DEFAULT 0 COMMENT ''Operator admin ID'' AFTER `operator_employee_id`',
  'SELECT 1'
);
PREPARE points_operator_admin_column_stmt FROM @points_operator_admin_column_sql;
EXECUTE points_operator_admin_column_stmt;
DEALLOCATE PREPARE points_operator_admin_column_stmt;

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

CREATE TABLE IF NOT EXISTS `goods_stock_record` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '库存流水ID',
  `goods_id` int unsigned NOT NULL COMMENT '商品ID',
  `goods_name` varchar(128) NOT NULL COMMENT '商品名称快照',
  `change_type` tinyint NOT NULL COMMENT '变动类型：1初始库存 2后台增加 3后台扣减 4订单扣减 5取消订单恢复',
  `change_quantity` int NOT NULL COMMENT '库存变动数量，正数增加，负数扣减',
  `before_stock` int unsigned NOT NULL COMMENT '变动前库存',
  `after_stock` int unsigned NOT NULL COMMENT '变动后库存',
  `biz_type` varchar(64) NOT NULL DEFAULT '' COMMENT '业务来源',
  `biz_id` int unsigned NOT NULL DEFAULT 0 COMMENT '关联业务ID',
  `operator_type` tinyint NOT NULL DEFAULT 0 COMMENT '操作者类型：0系统 1管理员 2员工',
  `operator_id` int unsigned NOT NULL DEFAULT 0 COMMENT '操作者ID',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_stock_record_goods_created` (`goods_id`, `created_at`),
  KEY `idx_stock_record_change_type` (`change_type`),
  KEY `idx_stock_record_biz` (`biz_type`, `biz_id`),
  KEY `idx_stock_record_operator` (`operator_type`, `operator_id`),
  KEY `idx_stock_record_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='商品库存变动流水表';

INSERT INTO `goods_stock_record` (
  `goods_id`, `goods_name`, `change_type`, `change_quantity`, `before_stock`, `after_stock`,
  `biz_type`, `biz_id`, `operator_type`, `operator_id`, `remark`
)
SELECT
  g.`id`, g.`name`, 1, g.`stock`, 0, g.`stock`, 'initial', g.`id`, 0, 0, 'Initial stock snapshot'
FROM `goods_info` g
WHERE g.`deleted_at` IS NULL
  AND NOT EXISTS (
    SELECT 1
    FROM `goods_stock_record` r
    WHERE r.`goods_id` = g.`id`
      AND r.`change_type` = 1
  );

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

CREATE TABLE IF NOT EXISTS admin_info (
  id int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Admin ID',
  username varchar(64) NOT NULL COMMENT 'Login username',
  password_hash varchar(100) NOT NULL COMMENT 'bcrypt password hash',
  real_name varchar(64) NOT NULL DEFAULT '' COMMENT 'Admin name',
  phone varchar(20) NOT NULL DEFAULT '' COMMENT 'Phone',
  email varchar(128) NOT NULL DEFAULT '' COMMENT 'Email',
  status tinyint NOT NULL DEFAULT 1 COMMENT 'Status: 1 normal 0 disabled',
  is_super tinyint NOT NULL DEFAULT 0 COMMENT 'Is super admin: 1 yes 0 no',
  last_login_at datetime DEFAULT NULL COMMENT 'Last login time',
  created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
  deleted_at datetime DEFAULT NULL COMMENT 'Deleted time',
  PRIMARY KEY (id),
  UNIQUE KEY uniq_admin_username (username),
  KEY idx_admin_status (status),
  KEY idx_admin_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Admin account table';

CREATE TABLE IF NOT EXISTS admin_role (
  id int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Role ID',
  name varchar(64) NOT NULL COMMENT 'Role name',
  description varchar(255) NOT NULL DEFAULT '' COMMENT 'Role description',
  status tinyint NOT NULL DEFAULT 1 COMMENT 'Status: 1 enabled 0 disabled',
  created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
  deleted_at datetime DEFAULT NULL COMMENT 'Deleted time',
  PRIMARY KEY (id),
  UNIQUE KEY uniq_admin_role_name (name),
  KEY idx_admin_role_status (status),
  KEY idx_admin_role_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Admin role table';

CREATE TABLE IF NOT EXISTS admin_permission (
  id int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Permission ID',
  name varchar(128) NOT NULL COMMENT 'Permission name',
  group_name varchar(64) NOT NULL DEFAULT '' COMMENT 'Permission group',
  method varchar(10) NOT NULL COMMENT 'HTTP method',
  path varchar(255) NOT NULL COMMENT 'API path',
  status tinyint NOT NULL DEFAULT 1 COMMENT 'Status: 1 enabled 0 disabled',
  created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  updated_at datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Updated time',
  deleted_at datetime DEFAULT NULL COMMENT 'Deleted time',
  PRIMARY KEY (id),
  UNIQUE KEY uniq_admin_permission_method_path (method, path),
  KEY idx_admin_permission_group (group_name),
  KEY idx_admin_permission_status (status),
  KEY idx_admin_permission_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Admin API permission table';

CREATE TABLE IF NOT EXISTS admin_role_relation (
  id int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Relation ID',
  admin_id int unsigned NOT NULL COMMENT 'Admin ID',
  role_id int unsigned NOT NULL COMMENT 'Role ID',
  created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  PRIMARY KEY (id),
  UNIQUE KEY uniq_admin_role_relation (admin_id, role_id),
  KEY idx_admin_role_relation_admin (admin_id),
  KEY idx_admin_role_relation_role (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Admin role relation table';

CREATE TABLE IF NOT EXISTS admin_role_permission (
  id int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Relation ID',
  role_id int unsigned NOT NULL COMMENT 'Role ID',
  permission_id int unsigned NOT NULL COMMENT 'Permission ID',
  created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  PRIMARY KEY (id),
  UNIQUE KEY uniq_admin_role_permission (role_id, permission_id),
  KEY idx_admin_role_permission_role (role_id),
  KEY idx_admin_role_permission_permission (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Admin role permission relation table';

INSERT INTO admin_info (username, password_hash, real_name, phone, email, status, is_super) VALUES
('root', '$2a$10$wkJo.7jih/0EbEehrNG.seMN5Rm3VZP90xxlK6bebLZDoq5K77W8C', 'System Administrator', '', '', 1, 1)
ON DUPLICATE KEY UPDATE
  password_hash = VALUES(password_hash),
  real_name = VALUES(real_name),
  status = VALUES(status),
  is_super = VALUES(is_super),
  deleted_at = NULL;

INSERT INTO admin_permission (name, group_name, method, path, status) VALUES
('Create employee', 'Employee Management', 'POST', '/backend/employee/manage/create', 1),
('Employee list', 'Employee Management', 'GET', '/backend/employee/manage/list', 1),
('Employee detail', 'Employee Management', 'GET', '/backend/employee/manage/detail', 1),
('Update employee', 'Employee Management', 'POST', '/backend/employee/manage/update', 1),
('Update employee status', 'Employee Management', 'POST', '/backend/employee/manage/status', 1),
('Reset employee password', 'Employee Management', 'POST', '/backend/employee/manage/reset-password', 1),
('Delete employee', 'Employee Management', 'POST', '/backend/employee/manage/delete', 1),
('Add employee credits', 'Credit Management', 'POST', '/backend/points/manage/add', 1),
('Batch add employee credits', 'Credit Management', 'POST', '/backend/points/manage/batch-add', 1),
('Deduct employee credits', 'Credit Management', 'POST', '/backend/points/manage/deduct', 1),
('Employee credit records', 'Credit Management', 'GET', '/backend/points/manage/records', 1),
('Backend category list', 'Goods Management', 'GET', '/backend/category/list', 1),
('Create goods', 'Goods Management', 'POST', '/backend/goods/create', 1),
('Goods list', 'Goods Management', 'GET', '/backend/goods/list', 1),
('Goods detail', 'Goods Management', 'GET', '/backend/goods/detail', 1),
('Update goods', 'Goods Management', 'POST', '/backend/goods/update', 1),
('Update goods status', 'Goods Management', 'POST', '/backend/goods/status', 1),
('Adjust goods stock', 'Stock Management', 'POST', '/backend/stock/adjust', 1),
('Goods stock records', 'Stock Management', 'GET', '/backend/stock/records', 1),
('Backend order list', 'Order Management', 'GET', '/backend/order/list', 1),
('Backend order detail', 'Order Management', 'GET', '/backend/order/detail', 1),
('Complete order', 'Order Management', 'POST', '/backend/order/complete', 1),
('Cancel order', 'Order Management', 'POST', '/backend/order/cancel', 1),
('Create admin', 'Admin Management', 'POST', '/backend/admin/manage/create', 1),
('Admin list', 'Admin Management', 'GET', '/backend/admin/manage/list', 1),
('Admin detail', 'Admin Management', 'GET', '/backend/admin/manage/detail', 1),
('Update admin', 'Admin Management', 'POST', '/backend/admin/manage/update', 1),
('Update admin status', 'Admin Management', 'POST', '/backend/admin/manage/status', 1),
('Reset admin password', 'Admin Management', 'POST', '/backend/admin/manage/reset-password', 1),
('Assign admin roles', 'Admin Management', 'POST', '/backend/admin/manage/roles', 1),
('Create role', 'Role Management', 'POST', '/backend/role/create', 1),
('Role list', 'Role Management', 'GET', '/backend/role/list', 1),
('Role detail', 'Role Management', 'GET', '/backend/role/detail', 1),
('Update role', 'Role Management', 'POST', '/backend/role/update', 1),
('Update role status', 'Role Management', 'POST', '/backend/role/status', 1),
('Assign role permissions', 'Role Management', 'POST', '/backend/role/permissions', 1),
('Permission list', 'Permission Management', 'GET', '/backend/permission/list', 1),
('Permission detail', 'Permission Management', 'GET', '/backend/permission/detail', 1),
('Create permission', 'Permission Management', 'POST', '/backend/permission/create', 1),
('Update permission', 'Permission Management', 'POST', '/backend/permission/update', 1),
('Update permission status', 'Permission Management', 'POST', '/backend/permission/status', 1)
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  group_name = VALUES(group_name),
  status = VALUES(status),
  deleted_at = NULL;

CREATE TABLE IF NOT EXISTS upload_file (
  id int unsigned NOT NULL AUTO_INCREMENT COMMENT 'Upload file ID',
  admin_id int unsigned NOT NULL DEFAULT 0 COMMENT 'Upload admin ID',
  file_name varchar(128) NOT NULL COMMENT 'Saved file name',
  original_name varchar(255) NOT NULL DEFAULT '' COMMENT 'Original file name',
  file_path varchar(500) NOT NULL COMMENT 'Local file path',
  url varchar(500) NOT NULL COMMENT 'Public access URL',
  file_size bigint unsigned NOT NULL DEFAULT 0 COMMENT 'File size bytes',
  mime_type varchar(128) NOT NULL DEFAULT '' COMMENT 'MIME type',
  file_ext varchar(20) NOT NULL DEFAULT '' COMMENT 'File extension',
  biz_type varchar(64) NOT NULL DEFAULT 'goods' COMMENT 'Business type',
  created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT 'Created time',
  PRIMARY KEY (id),
  KEY idx_upload_file_admin_id (admin_id),
  KEY idx_upload_file_biz_type (biz_type),
  KEY idx_upload_file_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Upload file record table';

INSERT INTO admin_permission (name, group_name, method, path, status) VALUES
('Upload goods image', 'Goods Management', 'POST', '/backend/upload/goods-image', 1)
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  group_name = VALUES(group_name),
  status = VALUES(status),
  deleted_at = NULL;
