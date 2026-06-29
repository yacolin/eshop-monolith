-- ============================================================================
-- RBAC 清理脚本（第 1 步）
-- 清理旧数据并重置自增 ID（按外键顺序）
--
-- 执行: mysql -u root -p eshop_db < scripts/seed_rbac_01_clean.sql
-- 顺序: 01 → 02 → 03 → 04 → 05
-- ============================================================================

USE eshop_db;

-- ==================== 清理旧数据（按外键顺序） ====================
DELETE FROM usr_login_histories;
DELETE FROM usr_role_permissions;
DELETE FROM usr_user_roles;
DELETE FROM usr_addresses;
DELETE FROM usr_infos;
DELETE FROM usr_users;
DELETE FROM usr_permissions;
DELETE FROM usr_roles;

-- 重置自增 ID
ALTER TABLE usr_roles AUTO_INCREMENT = 1;
ALTER TABLE usr_permissions AUTO_INCREMENT = 1;
ALTER TABLE usr_users AUTO_INCREMENT = 1;
ALTER TABLE usr_infos AUTO_INCREMENT = 1;
ALTER TABLE usr_addresses AUTO_INCREMENT = 1;
ALTER TABLE usr_user_roles AUTO_INCREMENT = 1;
ALTER TABLE usr_role_permissions AUTO_INCREMENT = 1;
