-- 创建数据库
CREATE DATABASE IF NOT EXISTS file_share
    DEFAULT CHARACTER SET utf8mb4
    DEFAULT COLLATE utf8mb4_unicode_ci;

USE file_share;

-- 用户表
CREATE TABLE IF NOT EXISTS `user` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(64) NOT NULL COMMENT '用户名',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（BCrypt加密）',
    `email` VARCHAR(128) DEFAULT NULL COMMENT '邮箱',
    `role` ENUM('USER', 'ADMIN') NOT NULL DEFAULT 'USER' COMMENT '角色',
    `status` ENUM('ACTIVE', 'DISABLED') NOT NULL DEFAULT 'ACTIVE' COMMENT '状态',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    KEY `idx_email` (`email`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 文件夹表
CREATE TABLE IF NOT EXISTS `folder` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '文件夹ID',
    `folder_uuid` VARCHAR(64) NOT NULL COMMENT '文件夹UUID',
    `folder_name` VARCHAR(255) NOT NULL COMMENT '文件夹名称',
    `parent_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '上级目录ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '所属用户ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_folder_uuid` (`folder_uuid`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_folder_name` (`folder_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件夹表';

-- 文件表
CREATE TABLE IF NOT EXISTS `file` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '文件ID',
    `file_uuid` VARCHAR(64) NOT NULL COMMENT '文件UUID',
    `file_name` VARCHAR(255) NOT NULL COMMENT '文件名',
    `file_size` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '文件大小(字节)',
    `file_type` VARCHAR(64) DEFAULT NULL COMMENT '文件类型(MIME)',
    `md5` VARCHAR(64) DEFAULT NULL COMMENT '文件MD5校验值',
    `cos_key` VARCHAR(512) DEFAULT NULL COMMENT 'COS存储路径',
    `bucket_name` VARCHAR(128) DEFAULT NULL COMMENT 'COS Bucket名称',
    `storage_type` ENUM('COS', 'LOCAL') NOT NULL DEFAULT 'COS' COMMENT '存储类型',
    `folder_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '所属文件夹ID',
    `user_id` BIGINT UNSIGNED NOT NULL COMMENT '上传用户ID',
    `download_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '下载次数',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_file_uuid` (`file_uuid`),
    KEY `idx_file_name` (`file_name`),
    KEY `idx_folder_id` (`folder_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_md5` (`md5`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件表';

-- 分享表
CREATE TABLE IF NOT EXISTS `share` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '分享ID',
    `share_uuid` VARCHAR(64) NOT NULL COMMENT '分享UUID',
    `file_id` BIGINT UNSIGNED NOT NULL COMMENT '分享文件ID',
    `share_code` VARCHAR(16) DEFAULT NULL COMMENT '分享短码',
    `password` VARCHAR(255) DEFAULT NULL COMMENT '访问密码(BCrypt加密)',
    `expire_time` DATETIME DEFAULT NULL COMMENT '过期时间',
    `view_count` INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '访问次数',
    `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_share_uuid` (`share_uuid`),
    UNIQUE KEY `uk_share_code` (`share_code`),
    KEY `idx_file_id` (`file_id`),
    KEY `idx_created_by` (`created_by`),
    KEY `idx_expire_time` (`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分享表';

-- 分片上传记录表
CREATE TABLE IF NOT EXISTS `file_chunk` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '记录ID',
    `file_uuid` VARCHAR(64) NOT NULL COMMENT '文件UUID',
    `chunk_index` INT UNSIGNED NOT NULL COMMENT '分片序号',
    `chunk_size` BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '分片大小',
    `upload_id` VARCHAR(255) NOT NULL COMMENT 'COS分片上传ID',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '上传状态(0:进行中,1:完成)',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_file_chunk` (`file_uuid`, `chunk_index`),
    KEY `idx_upload_id` (`upload_id`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分片上传记录表';

-- 插入管理员账号 (密码: admin123)
INSERT INTO `user` (`username`, `password`, `email`, `role`, `status`)
VALUES ('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'admin@example.com', 'ADMIN', 'ACTIVE')
ON DUPLICATE KEY UPDATE username = username;
