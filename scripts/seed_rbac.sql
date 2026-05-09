-- ============================================================================
-- RBAC 数据初始化脚本
-- 使用前请确保数据库已通过 AutoMigrate 创建表结构
-- 运行: mysql -u root -p eshop_db < scripts/seed_rbac.sql
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

-- ==================== 权限数据 ====================
-- resource: 所属模块, action: 操作
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('product:read',    '查看产品',   '查看产品列表和详情',   'product',  'read',   '商品管理', 1, 1),
('product:create',  '创建产品',   '创建新产品',           'product',  'create', '商品管理', 2, 1),
('product:update',  '编辑产品',   '编辑产品信息',         'product',  'update', '商品管理', 3, 1),
('product:delete',  '删除产品',   '删除产品',             'product',  'delete', '商品管理', 4, 1),

('order:read',      '查看订单',   '查看订单列表和详情',   'order',    'read',   '订单管理', 5, 1),
('order:create',    '创建订单',   '创建订单',             'order',    'create', '订单管理', 6, 1),
('order:update',    '编辑订单',   '编辑订单信息',         'order',    'update', '订单管理', 7, 1),
('order:cancel',    '取消订单',   '取消订单',             'order',    'cancel', '订单管理', 8, 1),

('user:read',       '查看用户',   '查看用户列表和详情',   'user',     'read',   '用户管理', 9, 1),
('user:create',     '创建用户',   '创建用户',             'user',     'create', '用户管理',10, 1),
('user:update',     '编辑用户',   '编辑用户信息',         'user',     'update', '用户管理',11, 1),
('user:delete',     '删除用户',   '删除用户',             'user',     'delete', '用户管理',12, 1),

('role:read',       '查看角色',   '查看角色列表和详情',   'role',     'read',   '权限管理',13, 1),
('role:create',     '创建角色',   '创建角色',             'role',     'create', '权限管理',14, 1),
('role:update',     '编辑角色',   '编辑角色信息',         'role',     'update', '权限管理',15, 1),
('role:delete',     '删除角色',   '删除角色',             'role',     'delete', '权限管理',16, 1),

('inventory:read',   '查看库存',  '查看库存信息',        'inventory','read',   '库存管理',17, 1),
('inventory:manage', '管理库存',  '入库/出库/调整库存',   'inventory','manage', '库存管理',18, 1),

('category:read',   '查看分类',   '查看分类列表',        'category', 'read',   '分类管理',19, 1),
('category:manage', '管理分类',   '创建/编辑/删除分类',   'category','manage', '分类管理',20, 1);

-- ==================== 角色数据 ====================
INSERT INTO roles (name, display_name, description, status, sort, is_system) VALUES
('admin',  '管理员',  '系统管理员，拥有所有权限',      1, 1, 1),
('user',   '普通用户','普通用户，拥有基本操作权限',    1, 2, 1);

-- ==================== 角色-权限关联 ====================
-- admin 角色拥有所有权限（1-20）
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- user 角色拥有基础权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE name IN (
    'product:read',
    'order:read',
    'order:create',
    'order:cancel',
    'category:read',
    'inventory:read'
);

-- ==================== 用户数据 ====================
-- 密码均为 "123456"，bcrypt hash（cost=10）
-- 如需更换密码: go run ./cmd/genhash/ 生成新 hash 后替换

-- 管理员 admin（密码: 123456）
INSERT INTO users (status) VALUES (1);
INSERT INTO user_infos (user_id, nickname) VALUES (1, '管理员');
INSERT INTO user_identities (user_id, provider, identifier, credential, verified, meta) VALUES
(1, 'password', 'admin', '$2a$10$HFzEUNEVKJQCZ4aPYVb/YONrhix2jwj8iiJWM5TUZdXM4wPdkEllC', 1, '{}');
INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);  -- admin 角色

-- 普通用户 colin（密码: 123456）
INSERT INTO users (status) VALUES (1);
INSERT INTO user_infos (user_id, nickname) VALUES (2, 'Colin');
INSERT INTO user_identities (user_id, provider, identifier, credential, verified, meta) VALUES
(2, 'password', 'colin', '$2a$10$HFzEUNEVKJQCZ4aPYVb/YONrhix2jwj8iiJWM5TUZdXM4wPdkEllC', 1, '{}');
INSERT INTO user_roles (user_id, role_id) VALUES (2, 2);  -- user 角色
