-- ============================================================================
-- RBAC 角色-权限关联脚本（第 4 步）
-- 使用子查询以名称关联，避免对 ID 的硬编码依赖
--
-- 执行: mysql -u root -p eshop_db < scripts/seed_rbac_04_assignments.sql
-- 顺序: 01 → 02 → 03 → 04 → 05
-- ============================================================================

USE eshop_db;

-- ==================== admin 角色：拥有所有权限 ====================
INSERT INTO usr_role_permissions (role_id, permission_id)
SELECT (SELECT id FROM usr_roles WHERE name = 'admin'), id FROM usr_permissions;

-- ==================== user 角色：基础购物操作 ====================
INSERT INTO usr_role_permissions (role_id, permission_id)
SELECT (SELECT id FROM usr_roles WHERE name = 'user'), id FROM usr_permissions WHERE name IN (
    'product:read', 'category:read', 'inventory:read', 'sku:read',
    'attr:read', 'attr_val:read',
    'address:read', 'address:create', 'address:update', 'address:delete',
    'order:read', 'order:create', 'order:cancel',
    'cart:read', 'cart:create', 'cart:update', 'cart:delete',
    'payment:read', 'payment:create',
    'refund:read', 'refund:create',
    'review:read', 'review:create', 'review:delete',
    'notification:read', 'notification:update',
    'promotion:read',
    'user:read', 'user:update'
);

-- ==================== operator 角色：运营操作 ====================
INSERT INTO usr_role_permissions (role_id, permission_id)
SELECT (SELECT id FROM usr_roles WHERE name = 'operator'), id FROM usr_permissions WHERE name IN (
    'product:read', 'category:read', 'inventory:read', 'sku:read',
    'attr:read', 'attr_val:read', 'address:read',
    'order:read', 'order:update', 'order:cancel',
    'payment:read', 'refund:read', 'refund:update',
    'review:read', 'review:moderate',
    'notification:read', 'notification:update', 'notification:send',
    'promotion:read',
    'user:read'
);

-- ==================== editor 角色：内容维护 ====================
INSERT INTO usr_role_permissions (role_id, permission_id)
SELECT (SELECT id FROM usr_roles WHERE name = 'editor'), id FROM usr_permissions WHERE name IN (
    'product:read', 'product:create', 'product:update',
    'category:read', 'category:create', 'category:update',
    'inventory:read', 'sku:read',
    'attr:read', 'attr_val:read', 'address:read',
    'order:read',
    'promotion:read', 'promotion:create', 'promotion:update',
    'review:read', 'review:moderate',
    'notification:read',
    'user:read'
);

-- ==================== warehouse 角色：库存与发货 ====================
INSERT INTO usr_role_permissions (role_id, permission_id)
SELECT (SELECT id FROM usr_roles WHERE name = 'warehouse'), id FROM usr_permissions WHERE name IN (
    'product:read', 'category:read',
    'inventory:read', 'inventory:create', 'inventory:update', 'inventory:reserve',
    'sku:read', 'address:read',
    'order:read', 'order:update',
    'notification:read'
);

-- ==================== finance 角色：财务与退款 ====================
INSERT INTO usr_role_permissions (role_id, permission_id)
SELECT (SELECT id FROM usr_roles WHERE name = 'finance'), id FROM usr_permissions WHERE name IN (
    'order:read', 'payment:read', 'payment:update',
    'refund:read', 'refund:update',
    'product:read', 'sku:read', 'address:read',
    'notification:read', 'user:read'
);

-- ==================== merchant 角色：商户管理 ====================
INSERT INTO usr_role_permissions (role_id, permission_id)
SELECT (SELECT id FROM usr_roles WHERE name = 'merchant'), id FROM usr_permissions WHERE name IN (
    'product:read', 'product:create', 'product:update',
    'category:read', 'inventory:read',
    'sku:read', 'sku:create', 'sku:update',
    'attr:read', 'attr_val:read', 'address:read',
    'order:read', 'order:update', 'order:cancel',
    'payment:read', 'refund:read',
    'promotion:read',
    'review:read',
    'notification:read', 'notification:update',
    'dashboard:read',
    'user:read'
);

-- ==================== support 角色：客服售后 ====================
INSERT INTO usr_role_permissions (role_id, permission_id)
SELECT (SELECT id FROM usr_roles WHERE name = 'support'), id FROM usr_permissions WHERE name IN (
    'product:read', 'category:read', 'sku:read',
    'attr:read', 'attr_val:read', 'address:read',
    'order:read', 'order:update', 'order:cancel',
    'payment:read', 'refund:read', 'refund:update',
    'review:read', 'review:moderate',
    'notification:read', 'notification:update', 'notification:send',
    'user:read', 'dashboard:read'
);

-- ==================== analyst 角色：数据分析 ====================
INSERT INTO usr_role_permissions (role_id, permission_id)
SELECT (SELECT id FROM usr_roles WHERE name = 'analyst'), id FROM usr_permissions WHERE name IN (
    'product:read', 'category:read', 'inventory:read',
    'order:read', 'cart:read', 'payment:read', 'refund:read',
    'promotion:read', 'review:read', 'notification:read', 'user:read',
    'sku:read', 'attr:read', 'attr_val:read', 'address:read',
    'dashboard:read'
);
