-- ============================================================================
-- 地址种子数据
-- 为 admin (user_id=1) 和 colin (user_id=2) 各生成 3 个地址
--
-- 执行: mysql -u root -p eshop_db < scripts/seed_addresses.sql
-- ============================================================================

USE eshop_db;

-- 清空旧数据（按需启用）
-- DELETE FROM addresses WHERE user_id IN (1, 2);

-- ==================== admin（user_id=1） ====================

INSERT INTO addresses (user_id, consignee, phone, province, city, district, detail, zip_code, is_default, created_at, updated_at) VALUES
(1, '张管理', '13800138001', '北京市', '北京市', '朝阳区', '建国路88号SOHO现代城A座1508', '100022', TRUE, NOW(), NOW());

INSERT INTO addresses (user_id, consignee, phone, province, city, district, detail, zip_code, is_default, created_at, updated_at) VALUES
(1, '张管理', '13800138002', '北京市', '北京市', '海淀区', '中关村大街1号银谷大厦2005', '100080', FALSE, NOW(), NOW());

INSERT INTO addresses (user_id, consignee, phone, province, city, district, detail, zip_code, is_default, created_at, updated_at) VALUES
(1, '张管理', '13800138003', '上海市', '上海市', '浦东新区', '张江高科技园区博云路2号501', '201203', FALSE, NOW(), NOW());

-- ==================== colin（user_id=2） ====================

INSERT INTO addresses (user_id, consignee, phone, province, city, district, detail, zip_code, is_default, created_at, updated_at) VALUES
(2, '陈科林', '13900139001', '广东省', '深圳市', '南山区', '科技园南区高新南一道2号飞亚达科技大厦12F', '518057', TRUE, NOW(), NOW());

INSERT INTO addresses (user_id, consignee, phone, province, city, district, detail, zip_code, is_default, created_at, updated_at) VALUES
(2, '陈科林', '13900139002', '广东省', '广州市', '天河区', '珠江新城华夏路16号富力盈凯广场3001', '510623', FALSE, NOW(), NOW());

INSERT INTO addresses (user_id, consignee, phone, province, city, district, detail, zip_code, is_default, created_at, updated_at) VALUES
(2, '陈科林', '13900139003', '浙江省', '杭州市', '西湖区', '文三路478号华星科技大厦801', '310012', FALSE, NOW(), NOW());
