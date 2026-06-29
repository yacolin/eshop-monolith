-- ============================================================================
-- RBAC 用户数据脚本（第 5 步）
-- 密码均为 "123456"，bcrypt hash（cost=10）
--
-- 执行: mysql -u root -p eshop_db < scripts/seed_rbac_05_users.sql
-- 顺序: 01 → 02 → 03 → 04 → 05
-- ============================================================================

USE eshop_db;

-- 管理员 admin（密码: 123456）
INSERT INTO usr_users (username, password_hash, nickname, email, phone, status, register_source) VALUES
('admin', '$2a$10$HFzEUNEVKJQCZ4aPYVb/YONrhix2jwj8iiJWM5TUZdXM4wPdkEllC', '管理员', 'admin@eshop.dev', '13800000001', 1, 'admin');
INSERT INTO usr_infos (user_id) VALUES (1);
INSERT INTO usr_user_roles (user_id, role_id) VALUES (1, (SELECT id FROM usr_roles WHERE name = 'admin'));

-- 普通用户 colin（密码: 123456）
INSERT INTO usr_users (username, password_hash, nickname, email, phone, status, register_source) VALUES
('colin', '$2a$10$HFzEUNEVKJQCZ4aPYVb/YONrhix2jwj8iiJWM5TUZdXM4wPdkEllC', 'Colin', 'colin@eshop.dev', '13800000002', 1, 'web');
INSERT INTO usr_infos (user_id) VALUES (2);
INSERT INTO usr_user_roles (user_id, role_id) VALUES (2, (SELECT id FROM usr_roles WHERE name = 'user'));


-- ==================== 收货地址 ====================
INSERT INTO usr_addresses (user_id, consignee, phone, country, province, city, district, detail, zip_code, tag, is_default) VALUES
(1, '张管理', '13800138001', '中国', '北京市', '北京市', '朝阳区', '建国路88号SOHO现代城A座1508', '100022', 'office', TRUE),
(1, '张管理', '13800138002', '中国', '北京市', '北京市', '海淀区', '中关村大街1号银谷大厦2005', '100080', 'office', FALSE);

INSERT INTO usr_addresses (user_id, consignee, phone, country, province, city, district, detail, zip_code, tag, is_default) VALUES
(2, '陈科林', '13900139001', '中国', '广东省', '深圳市', '南山区', '科技园南区高新南一道2号飞亚达科技大厦12F', '518057', 'company', TRUE),
(2, '陈科林', '13900139002', '中国', '广东省', '广州市', '天河区', '珠江新城华夏路16号富力盈凯广场3001', '510623', 'company', FALSE);