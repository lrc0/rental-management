-- 租房管理系统数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS rental_management DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE rental_management;

-- 房东用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
  `password_hash` VARCHAR(255) NOT NULL COMMENT '密码哈希',
  `name` VARCHAR(50) DEFAULT '' COMMENT '姓名',
  `avatar` VARCHAR(255) DEFAULT '' COMMENT '头像',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 1正常 2禁用',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_phone` (`phone`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='房东用户表';

-- 费率配置表
CREATE TABLE IF NOT EXISTS `fee_rates` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '房东ID',
  `water_rate` DECIMAL(10,4) DEFAULT 0.0000 COMMENT '水费单价(元/吨)',
  `electricity_rate` DECIMAL(10,4) DEFAULT 0.0000 COMMENT '电费单价(元/度)',
  `gas_rate` DECIMAL(10,4) DEFAULT 0.0000 COMMENT '气费单价(元/立方)',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='费率配置表';

-- 房源表
CREATE TABLE IF NOT EXISTS `properties` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '房东ID',
  `name` VARCHAR(100) NOT NULL COMMENT '房源名称',
  `address` VARCHAR(255) DEFAULT '' COMMENT '地址',
  `property_type` TINYINT NOT NULL COMMENT '类型: 1整栋 2单套 3商铺',
  `total_rooms` INT DEFAULT 0 COMMENT '总房间数',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 1正常 2已下架',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='房源表';

-- 房间表
CREATE TABLE IF NOT EXISTS `rooms` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `property_id` BIGINT UNSIGNED NOT NULL COMMENT '房源ID',
  `room_number` VARCHAR(20) NOT NULL COMMENT '房间号',
  `floor` INT DEFAULT 0 COMMENT '楼层',
  `area` DECIMAL(10,2) DEFAULT 0.00 COMMENT '面积(平米)',
  `monthly_rent` DECIMAL(10,2) DEFAULT 0.00 COMMENT '月租金',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 1空置 2已租 3维修中',
  `facilities` JSON COMMENT '设施配置',
  `remark` VARCHAR(255) DEFAULT '' COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_property_id` (`property_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='房间表';

-- 租客表
CREATE TABLE IF NOT EXISTS `tenants` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '房东ID',
  `name` VARCHAR(50) NOT NULL COMMENT '姓名',
  `phone` VARCHAR(20) DEFAULT '' COMMENT '手机号',
  `id_card` VARCHAR(18) DEFAULT '' COMMENT '身份证号',
  `gender` TINYINT DEFAULT 0 COMMENT '性别: 0未知 1男 2女',
  `emergency_contact` VARCHAR(50) DEFAULT '' COMMENT '紧急联系人',
  `emergency_phone` VARCHAR(20) DEFAULT '' COMMENT '紧急联系电话',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 1正常 2已退租',
  `remark` VARCHAR(255) DEFAULT '' COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租客表';

-- 合同表
CREATE TABLE IF NOT EXISTS `contracts` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `room_id` BIGINT UNSIGNED NOT NULL COMMENT '房间ID',
  `tenant_id` BIGINT UNSIGNED NOT NULL COMMENT '租客ID',
  `start_date` DATE NOT NULL COMMENT '开始日期',
  `end_date` DATE NOT NULL COMMENT '结束日期',
  `monthly_rent` DECIMAL(10,2) DEFAULT 0.00 COMMENT '月租金',
  `deposit` DECIMAL(10,2) DEFAULT 0.00 COMMENT '押金',
  `payment_day` TINYINT DEFAULT 1 COMMENT '每月几号交租',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 1生效 2已到期 3已解约',
  `terminate_reason` VARCHAR(255) DEFAULT '' COMMENT '解约原因',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` DATETIME DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_room_id` (`room_id`),
  KEY `idx_tenant_id` (`tenant_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='合同表';

-- 抄表记录表
CREATE TABLE IF NOT EXISTS `meter_readings` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `room_id` BIGINT UNSIGNED NOT NULL COMMENT '房间ID',
  `reading_date` DATE NOT NULL COMMENT '抄表日期',
  `water_reading` DECIMAL(10,2) DEFAULT 0.00 COMMENT '水表读数',
  `electricity_reading` DECIMAL(10,2) DEFAULT 0.00 COMMENT '电表读数',
  `gas_reading` DECIMAL(10,2) DEFAULT 0.00 COMMENT '气表读数',
  `water_usage` DECIMAL(10,2) DEFAULT 0.00 COMMENT '用水量',
  `electricity_usage` DECIMAL(10,2) DEFAULT 0.00 COMMENT '用电量',
  `gas_usage` DECIMAL(10,2) DEFAULT 0.00 COMMENT '用气量',
  `remark` VARCHAR(255) DEFAULT '' COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_room_date` (`room_id`, `reading_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='抄表记录表';

-- 账单表
CREATE TABLE IF NOT EXISTS `bills` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '房东ID',
  `room_id` BIGINT UNSIGNED NOT NULL COMMENT '房间ID',
  `tenant_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '租客ID',
  `bill_type` TINYINT DEFAULT 5 COMMENT '类型: 1租金 2水费 3电费 4气费 5综合',
  `bill_month` VARCHAR(7) NOT NULL COMMENT '账单月份(YYYY-MM)',
  `amount` DECIMAL(10,2) DEFAULT 0.00 COMMENT '总金额',
  `water_fee` DECIMAL(10,2) DEFAULT 0.00 COMMENT '水费',
  `electricity_fee` DECIMAL(10,2) DEFAULT 0.00 COMMENT '电费',
  `gas_fee` DECIMAL(10,2) DEFAULT 0.00 COMMENT '气费',
  `rent_fee` DECIMAL(10,2) DEFAULT 0.00 COMMENT '租金',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 1待支付 2已支付 3已逾期',
  `due_date` DATE DEFAULT NULL COMMENT '应付日期',
  `paid_at` DATETIME DEFAULT NULL COMMENT '支付时间',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_month` (`user_id`, `bill_month`),
  KEY `idx_room_id` (`room_id`),
  KEY `idx_tenant_id` (`tenant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='账单表';

-- 收款记录表
CREATE TABLE IF NOT EXISTS `payments` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `bill_id` BIGINT UNSIGNED NOT NULL COMMENT '账单ID',
  `amount` DECIMAL(10,2) DEFAULT 0.00 COMMENT '金额',
  `payment_method` TINYINT DEFAULT 1 COMMENT '支付方式: 1现金 2微信 3支付宝 4银行转账',
  `paid_at` DATETIME NOT NULL COMMENT '支付时间',
  `note` VARCHAR(255) DEFAULT '' COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_bill_id` (`bill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='收款记录表';
