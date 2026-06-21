-- ============================================================================
-- RBAC 清理脚本（第 1 步）
-- 清理旧数据并重置自增 ID（按外键顺序）
--
-- 执行: mysql -u root -p eshop_db < scripts/seed_rbac_01_clean.sql
-- 顺序: 01 → 02 → 03 → 04 → 05
-- ============================================================================

USE eshop_db;

-- ==================== 清理旧数据（按外键顺序） ====================
DELETE FROM login_histories;
DELETE FROM auth_tokens;
DELETE FROM user_roles;
DELETE FROM role_permissions;
DELETE FROM user_identities;
DELETE FROM user_infos;
DELETE FROM users;
DELETE FROM permissions;
DELETE FROM roles;

-- 重置自增 ID
ALTER TABLE roles AUTO_INCREMENT = 1;
ALTER TABLE permissions AUTO_INCREMENT = 1;
ALTER TABLE users AUTO_INCREMENT = 1;
ALTER TABLE user_infos AUTO_INCREMENT = 1;
ALTER TABLE user_identities AUTO_INCREMENT = 1;
ALTER TABLE user_roles AUTO_INCREMENT = 1;
ALTER TABLE role_permissions AUTO_INCREMENT = 1;
