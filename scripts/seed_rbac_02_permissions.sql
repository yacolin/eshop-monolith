-- ============================================================================
-- RBAC 权限初始化脚本（第 2 步）
-- 按后端模块划分，每模块独立区间
--
-- 执行: mysql -u root -p eshop_db < scripts/seed_rbac_02_permissions.sql
-- 顺序: 01 → 02 → 03 → 04 → 05
-- ============================================================================

USE eshop_db;

-- ═══════════════════════════════════════════════════════════════════════════
-- product 模块 (10000-19999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, resource, action, category, sort_order, status) VALUES
('product:read',   '查看产品',   'product',  'read',   'product', 11000, 1),
('product:create', '创建产品',   'product',  'create', 'product', 11050, 1),
('product:update', '编辑产品',   'product',  'update', 'product', 11100, 1),
('product:delete', '删除产品',   'product',  'delete', 'product', 11150, 1),

('category:read',   '查看分类',   'category', 'read',   'product', 11500, 1),
('category:create', '创建分类',   'category', 'create', 'product', 11550, 1),
('category:update', '编辑分类',   'category', 'update', 'product', 11600, 1),
('category:delete', '删除分类',   'category', 'delete', 'product', 11650, 1),

('brand:read',   '查看品牌',   'brand', 'read',   'product', 12000, 1),
('brand:create', '创建品牌',   'brand', 'create', 'product', 12050, 1),
('brand:update', '编辑品牌',   'brand', 'update', 'product', 12100, 1),
('brand:delete', '删除品牌',   'brand', 'delete', 'product', 12150, 1),

('sku:read',   '查看 SKU',   'sku', 'read',   'product', 12500, 1),
('sku:create', '创建 SKU',   'sku', 'create', 'product', 12550, 1),
('sku:update', '编辑 SKU',   'sku', 'update', 'product', 12600, 1),
('sku:delete', '删除 SKU',   'sku', 'delete', 'product', 12650, 1),

('attr:read',      '查看属性',   'attr',     'read',   'product', 13000, 1),
('attr:create',    '创建属性',   'attr',     'create', 'product', 13050, 1),
('attr:update',    '编辑属性',   'attr',     'update', 'product', 13100, 1),
('attr:delete',    '删除属性',   'attr',     'delete', 'product', 13150, 1),
('attr_val:read',   '查看属性值', 'attr_val', 'read',   'product', 13250, 1),
('attr_val:create', '创建属性值', 'attr_val', 'create', 'product', 13300, 1),
('attr_val:update', '编辑属性值', 'attr_val', 'update', 'product', 13350, 1),
('attr_val:delete', '删除属性值', 'attr_val', 'delete', 'product', 13400, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- inventory 模块 (20000-29999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, resource, action, category, sort_order, status) VALUES
('inventory:read',    '查看库存',   'inventory', 'read',    'inventory', 21000, 1),
('inventory:create',  '创建库存',   'inventory', 'create',  'inventory', 21050, 1),
('inventory:update',  '编辑库存',   'inventory', 'update',  'inventory', 21100, 1),
('inventory:reserve', '库存操作',   'inventory', 'reserve', 'inventory', 21150, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- trade 模块 (30000-39999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, resource, action, category, sort_order, status) VALUES
('order:read',         '查看订单',   'order', 'read',   'trade', 31000, 1),
('order:create',       '创建订单',   'order', 'create', 'trade', 31050, 1),
('order:update',       '编辑订单',   'order', 'update', 'trade', 31100, 1),
('order:cancel',       '取消订单',   'order', 'cancel', 'trade', 31150, 1),

('cart:read',   '查看购物车',   'cart', 'read',   'trade', 31500, 1),
('cart:create', '添加商品',     'cart', 'add',    'trade', 31550, 1),
('cart:update', '编辑购物车',   'cart', 'update', 'trade', 31600, 1),
('cart:delete', '删除商品',     'cart', 'delete', 'trade', 31650, 1),

('payment:read',   '查看支付',   'payment', 'read',   'trade', 32000, 1),
('payment:create', '发起支付',   'payment', 'create', 'trade', 32050, 1),
('payment:update', '更新支付',   'payment', 'update', 'trade', 32100, 1),

('refund:read',   '查看退款',   'refund', 'read',   'trade', 32500, 1),
('refund:create', '申请退款',   'refund', 'create', 'trade', 32550, 1),
('refund:update', '处理退款',   'refund', 'update', 'trade', 32600, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- marketing 模块 (40000-49999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, resource, action, category, sort_order, status) VALUES
('promotion:read',   '查看促销',   'promotion', 'read',   'marketing', 41000, 1),
('promotion:create', '创建促销',   'promotion', 'create', 'marketing', 41050, 1),
('promotion:update', '编辑促销',   'promotion', 'update', 'marketing', 41100, 1),
('promotion:delete', '删除促销',   'promotion', 'delete', 'marketing', 41150, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- review 模块 (50000-59999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, resource, action, category, sort_order, status) VALUES
('review:read',     '查看评论',   'review', 'read',     'review', 51000, 1),
('review:create',   '发表评论',   'review', 'create',   'review', 51050, 1),
('review:delete',   '删除评论',   'review', 'delete',   'review', 51100, 1),
('review:moderate', '审核评论',   'review', 'moderate', 'review', 51150, 1),
('review:reply',    '回复评论',   'review', 'reply',    'review', 51200, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- user 模块 (60000-69999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, resource, action, category, sort_order, status) VALUES
('user:read',   '查看用户',   'user', 'read',   'user', 61000, 1),
('user:create', '创建用户',   'user', 'create', 'user', 61050, 1),
('user:update', '编辑用户',   'user', 'update', 'user', 61100, 1),
('user:delete', '删除用户',   'user', 'delete', 'user', 61150, 1),

('role:read',   '查看角色',   'role', 'read',   'user', 61500, 1),
('role:create', '创建角色',   'role', 'create', 'user', 61550, 1),
('role:update', '编辑角色',   'role', 'update', 'user', 61600, 1),
('role:delete', '删除角色',   'role', 'delete', 'user', 61650, 1),

('address:read',   '查看地址',   'address', 'read',   'user', 62000, 1),
('address:create', '创建地址',   'address', 'create', 'user', 62050, 1),
('address:update', '编辑地址',   'address', 'update', 'user', 62100, 1),
('address:delete', '删除地址',   'address', 'delete', 'user', 62150, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- base 模块 (70000-79999) — 通知
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, resource, action, category, sort_order, status) VALUES
('notification:read',   '查看通知',   'notification', 'read',   'base', 71000, 1),
('notification:update', '标记已读',   'notification', 'update', 'base', 71050, 1),
('notification:delete', '删除通知',   'notification', 'delete', 'base', 71100, 1),
('notification:send',   '发送通知',   'notification', 'send',   'base', 71150, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- dashboard 模块 (80000-89999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, resource, action, category, sort_order, status) VALUES
('dashboard:read', '查看仪表盘', 'dashboard', 'read', 'dashboard', 81000, 1);
