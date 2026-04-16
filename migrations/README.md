# 数据库迁移说明

## 执行迁移脚本

如果遇到 `Error 1452: Cannot add or update a child row: a foreign key constraint fails` 错误，说明数据库中存在外键约束。请执行以下步骤：

### 方法1: 直接执行 SQL

连接到数据库后执行：

```sql
-- 查看外键约束名称
SELECT CONSTRAINT_NAME 
FROM information_schema.KEY_COLUMN_USAGE 
WHERE TABLE_NAME = 'bills' 
AND COLUMN_NAME = 'tenant_id' 
AND REFERENCED_TABLE_NAME = 'tenants';

-- 删除外键约束（替换 CONSTRAINT_NAME 为上面查询的结果）
ALTER TABLE bills DROP FOREIGN KEY CONSTRAINT_NAME;

-- 确保 tenant_id 可以为 NULL
ALTER TABLE bills MODIFY COLUMN tenant_id BIGINT UNSIGNED DEFAULT NULL;
```

### 方法2: 使用迁移脚本

```bash
# 进入 MySQL 容器或直接执行
mysql -u root -p rental_management < migrations/001_fix_bill_tenant_fk.sql
```

## 自动计算水电气费用

创建账单时，可以设置 `auto_calculate: true`，系统会：

1. 根据账单月份查找该月的抄表记录
2. 获取用户配置的水电气费率
3. 自动计算费用：
   - 水费 = 用水量 × 水费单价
   - 电费 = 用电量 × 电费单价
   - 气费 = 用气量 × 气费单价

需要先在系统中配置费率（访问"费率设置"页面）。
