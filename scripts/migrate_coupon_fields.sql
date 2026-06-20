-- ============================================================================
-- 优惠券与促销系统 — 数据迁移脚本
--
-- 在运行 AutoMigrate 后执行，用于维护订单表 discount_amount 新字段
-- 运行: mysql -u root -p eshop_db < scripts/migrate_coupon_fields.sql
-- ============================================================================

USE eshop_db;

-- ==================== 1. 补充默认值 ====================
-- AutoMigrate 已自动添加 discount_amount 列（默认 0），
-- 历史数据如有 NULL 或负数则纠正
UPDATE orders
SET discount_amount = 0
WHERE discount_amount IS NULL OR discount_amount < 0;

-- ==================== 2. 验证迁移结果 ====================
SELECT 'orders discount_amount validation' AS check_name,
       COUNT(*) AS total_rows,
       SUM(CASE WHEN discount_amount IS NULL THEN 1 ELSE 0 END) AS null_discounts,
       SUM(CASE WHEN discount_amount < 0 THEN 1 ELSE 0 END) AS negative_discounts
FROM orders;

-- ==================== 3. discount_amount 索引 ====================
-- coupon_id 已有 GORM 自动创建的索引，只需添加 discount_amount
-- MySQL 8.0+ 支持 CREATE INDEX IF NOT EXISTS
SET @db = DATABASE();
SET @idx_exists = (
    SELECT COUNT(*) FROM information_schema.STATISTICS
    WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'orders' AND INDEX_NAME = 'idx_orders_discount_amount'
);
SET @sql = IF(@idx_exists = 0,
    'CREATE INDEX idx_orders_discount_amount ON orders(discount_amount)',
    'SELECT ''index already exists'' AS msg'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
