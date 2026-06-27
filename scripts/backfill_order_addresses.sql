-- ============================================================================
-- 批量更新旧订单地址快照
-- 按用户的默认地址补全 orders 表中为空的地址快照字段
--
-- 执行: mysql -u root -p eshop_db --default-character-set=utf8mb4 < scripts\backfill_order_addresses.sql
-- ============================================================================

USE eshop_db;

-- 更新 customer_id 能匹配到 user_id 且地址快照为空的订单
UPDATE orders o
  JOIN addresses a ON (o.customer_id + 0) = a.user_id AND a.is_default = TRUE
SET
  o.consignee   = a.consignee,
  o.phone       = a.phone,
  o.province    = a.province,
  o.city        = a.city,
  o.district    = a.district,
  o.detail_addr = a.detail,
  o.zip_code    = a.zip_code
WHERE
  (o.consignee IS NULL OR o.consignee = '')
  AND (o.customer_id + 0) > 0;
