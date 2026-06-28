-- /products 列表查询覆盖索引
-- 使用场景: ORDER BY created_at DESC / ASC 时消除 filesort
-- 主键索引已覆盖 ORDER BY id ASC 的默认排序场景
CREATE INDEX idx_products_deleted_created ON products (deleted_at, created_at DESC);
