-- ============================================================================
-- RBAC 用户数据脚本（第 5 步）
-- 密码均为 "123456"，bcrypt hash（cost=10）
-- 如需更换密码: go run ./cmd/genhash/ 生成新 hash 后替换
--
-- 执行: mysql -u root -p eshop_db < scripts/seed_rbac_05_users.sql
-- 顺序: 01 → 02 → 03 → 04 → 05
-- ============================================================================

USE eshop_db;

-- 管理员 admin（密码: 123456）
INSERT INTO users (status) VALUES (1);
INSERT INTO user_infos (user_id, nickname) VALUES (1, '管理员');
INSERT INTO user_identities (user_id, provider, identifier, credential, verified, meta) VALUES
(1, 'password', 'admin', '$2a$10$HFzEUNEVKJQCZ4aPYVb/YONrhix2jwj8iiJWM5TUZdXM4wPdkEllC', 1, '{}');
INSERT INTO user_roles (user_id, role_id) VALUES (1, (SELECT id FROM roles WHERE name = 'admin'));

-- 普通用户 colin（密码: 123456）
INSERT INTO users (status) VALUES (1);
INSERT INTO user_infos (user_id, nickname) VALUES (2, 'Colin');
INSERT INTO user_identities (user_id, provider, identifier, credential, verified, meta) VALUES
(2, 'password', 'colin', '$2a$10$HFzEUNEVKJQCZ4aPYVb/YONrhix2jwj8iiJWM5TUZdXM4wPdkEllC', 1, '{}');
INSERT INTO user_roles (user_id, role_id) VALUES (2, (SELECT id FROM roles WHERE name = 'user'));
