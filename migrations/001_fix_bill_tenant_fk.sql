-- 修复 bills 表的外键约束问题
-- 如果存在外键约束，先删除

-- 查找并删除 bills 表上的外键约束
-- 注意：需要根据实际的外键名称调整

-- 1. 先查看是否有外键约束
-- SELECT CONSTRAINT_NAME FROM information_schema.KEY_COLUMN_USAGE
-- WHERE TABLE_NAME = 'bills' AND COLUMN_NAME = 'tenant_id' AND REFERENCED_TABLE_NAME = 'tenants';

-- 2. 删除外键约束（如果存在）
-- MySQL 8.0 的外键名称通常是 bills_ibfk_1 或类似的格式
-- 如果有多个外键，需要逐个删除

SET @dbname = DATABASE();
SET @tablename = 'bills';
SET @columnname = 'tenant_id';
SET @constraintname = (
    SELECT CONSTRAINT_NAME
    FROM information_schema.KEY_COLUMN_USAGE
    WHERE TABLE_SCHEMA = @dbname
    AND TABLE_NAME = @tablename
    AND COLUMN_NAME = @columnname
    AND REFERENCED_TABLE_NAME IS NOT NULL
    LIMIT 1
);

SET @sql = IF(@constraintname IS NOT NULL,
    CONCAT('ALTER TABLE ', @tablename, ' DROP FOREIGN KEY ', @constraintname),
    'SELECT 1');

PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 3. 确保 tenant_id 可以为 NULL
ALTER TABLE bills MODIFY COLUMN tenant_id BIGINT UNSIGNED DEFAULT NULL;

-- 4. 确保没有 NOT NULL 约束
-- 这条语句确保字段可以为 NULL
