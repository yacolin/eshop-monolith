-- ============================================================================
-- RBAC 数据初始化脚本
-- 覆盖所有业务模块的权限与角色分配
--
-- 使用前请确保数据库已通过 AutoMigrate 创建表结构
-- 运行: mysql -u root -p eshop_db < scripts/seed_rbac.sql
--
-- 【增量新增权限】参见 add_permission_template.sql
--   - 往已有模块中插入新操作
--   - 在模块末尾追加操作
--   - 在模块间插入全新模块（有空位/无空位）
--   - 在大类末尾追加新模块
-- ============================================================================
-- 【排序编号方案】
--   层    级  | 间隔 | 示例
--   ──────────┼──────┼──────────
--   大类间隔  │ 10000│ 10000 20000 30000 ...
--   模块基准  │   500│ 11000 11500 12000 ...
--   操作间隔  │    50│ 11000 11050 11100 ...
--
-- 扩展示例：
--   在 coupon:delete(23650) 和 coupon:claim(23700) 之间加 coupon:export
--   → 用 23675（50 和 100 中间）
--
--   在订单模块(21000)和购物车模块(21500)之间加售后模块
--   → 用 21250（500 中间）
--
--   模块整体移动（例如优惠券管理后移 1000 到 24500）：
--   → UPDATE permissions SET sort = sort + 1000 WHERE resource = 'coupon';
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

-- ========================================================================
-- 权限数据
-- 格式: resource:action
-- 四大类分组，组间 10000 间隔，模块间 500 间隔，操作间 50 间隔
-- ========================================================================

-- ═══════════════════════════════════════════════════════════════════════════
-- 大类一：商品库存 (10000-19999)
-- ═══════════════════════════════════════════════════════════════════════════

-- ==================== 商品管理 (11000) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('product:read',   '查看产品',   '查看产品列表和详情',          'product',  'read',   '商品管理', 11000, 1),
('product:create', '创建产品',   '创建新产品',                  'product',  'create', '商品管理', 11050, 1),
('product:update', '编辑产品',   '编辑产品信息',                'product',  'update', '商品管理', 11100, 1),
('product:delete', '删除产品',   '删除产品',                    'product',  'delete', '商品管理', 11150, 1);

-- ==================== 分类管理 (11500) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('category:read',   '查看分类',   '查看分类列表和详情',          'category', 'read',   '分类管理', 11500, 1),
('category:create', '创建分类',   '创建分类',                    'category', 'create', '分类管理', 11550, 1),
('category:update', '编辑分类',   '编辑分类信息',                'category', 'update', '分类管理', 11600, 1),
('category:delete', '删除分类',   '删除分类',                    'category', 'delete', '分类管理', 11650, 1);

-- ==================== 库存管理 (12000) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('inventory:read',    '查看库存',   '查看库存信息',              'inventory', 'read',    '库存管理', 12000, 1),
('inventory:create',  '创建库存',   '创建库存记录',              'inventory', 'create',  '库存管理', 12050, 1),
('inventory:update',  '编辑库存',   '编辑库存信息',              'inventory', 'update',  '库存管理', 12100, 1),
('inventory:reserve', '库存操作',   '库存预订与释放操作',        'inventory', 'reserve', '库存管理', 12150, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- 大类二：交易订单 (20000-29999)
-- ═══════════════════════════════════════════════════════════════════════════

-- ==================== 订单管理 (21000) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('order:read',         '查看订单',   '查看订单列表和详情',        'order', 'read',         '订单管理', 21000, 1),
('order:create',       '创建订单',   '创建订单',                  'order', 'create',       '订单管理', 21050, 1),
('order:update',       '编辑订单',   '编辑订单信息',              'order', 'update',       '订单管理', 21100, 1),
('order:cancel',       '取消订单',   '取消订单',                  'order', 'cancel',       '订单管理', 21150, 1),
('order:delete',       '删除订单',   '删除订单',                  'order', 'delete',       '订单管理', 21200, 1);

-- ==================== 购物车管理 (21500) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('cart:read',   '查看购物车',   '查看购物车内容',               'cart', 'read',   '购物车管理', 21500, 1),
('cart:add',    '添加商品',     '添加商品到购物车',             'cart', 'add',    '购物车管理', 21550, 1),
('cart:update', '编辑购物车',   '更新购物车商品信息',           'cart', 'update', '购物车管理', 21600, 1),
('cart:delete', '删除商品',     '删除购物车中的商品',           'cart', 'delete', '购物车管理', 21650, 1);

-- ==================== 支付管理 (22000) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('payment:read',   '查看支付',   '查看支付记录和详情',          'payment', 'read',   '支付管理', 22000, 1),
('payment:create', '发起支付',   '发起支付请求',                'payment', 'create', '支付管理', 22050, 1),
('payment:update', '更新支付',   '更新支付状态',                'payment', 'update', '支付管理', 22100, 1);

-- ==================== 退款管理 (22500) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('refund:read',   '查看退款',   '查看退款记录和详情',           'refund', 'read',   '退款管理', 22500, 1),
('refund:create', '申请退款',   '创建退款申请',                 'refund', 'create', '退款管理', 22550, 1),
('refund:update', '处理退款',   '更新退款状态',                 'refund', 'update', '退款管理', 22600, 1);

-- ==================== 秒杀管理 (23000) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('flash:read',   '查看秒杀',   '查看秒杀活动列表和详情',        'flash', 'read',   '秒杀管理', 23000, 1),
('flash:create', '创建秒杀',   '创建秒杀活动',                  'flash', 'create', '秒杀管理', 23050, 1),
('flash:update', '编辑秒杀',   '编辑秒杀订单信息',              'flash', 'update', '秒杀管理', 23100, 1),
('flash:cancel', '取消秒杀',   '取消秒杀订单',                  'flash', 'cancel', '秒杀管理', 23150, 1),
('flash:buy',    '秒杀购买',   '参与秒杀购买',                  'flash', 'buy',    '秒杀管理', 23200, 1),
('flash:manage', '管理秒杀',   '秒杀加载库存/上架等操作',       'flash', 'manage', '秒杀管理', 23250, 1);

-- ==================== 优惠券管理 (23500) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('coupon:read',   '查看优惠券',   '查看优惠券模板列表和详情',     'coupon', 'read',   '优惠券管理', 23500, 1),
('coupon:create', '创建优惠券',   '创建优惠券模板',               'coupon', 'create', '优惠券管理', 23550, 1),
('coupon:update', '编辑优惠券',   '编辑优惠券模板信息',           'coupon', 'update', '优惠券管理', 23600, 1),
('coupon:delete', '删除优惠券',   '删除优惠券模板',               'coupon', 'delete', '优惠券管理', 23650, 1),
('coupon:claim',  '领取优惠券',   '领取优惠券到个人账户',         'coupon', 'claim',  '优惠券管理', 23700, 1);

-- ==================== 促销管理 (24000) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('promotion:read',   '查看促销',   '查看促销活动列表和详情',      'promotion', 'read',   '促销管理', 24000, 1),
('promotion:create', '创建促销',   '创建促销活动',                'promotion', 'create', '促销管理', 24050, 1),
('promotion:update', '编辑促销',   '编辑促销活动信息',            'promotion', 'update', '促销管理', 24100, 1),
('promotion:delete', '删除促销',   '删除促销活动',                'promotion', 'delete', '促销管理', 24150, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- 大类三：评价反馈 (30000-39999)
-- ═══════════════════════════════════════════════════════════════════════════

-- ==================== 评论管理 (31000) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('review:read',     '查看评论',   '查看评论列表和详情',         'review', 'read',     '评论管理', 31000, 1),
('review:create',   '发表评论',   '发表商品评论',               'review', 'create',   '评论管理', 31050, 1),
('review:delete',   '删除评论',   '删除评论',                   'review', 'delete',   '评论管理', 31100, 1),
('review:moderate', '审核评论',   '审核/回复评论',              'review', 'moderate', '评论管理', 31150, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- 大类四：用户与系统 (40000-49999)
-- ═══════════════════════════════════════════════════════════════════════════

-- ==================== 通知管理 (41000) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('notification:read',   '查看通知',   '查看通知列表和详情',     'notification', 'read',   '通知管理', 41000, 1),
('notification:update', '标记已读',   '标记通知为已读',         'notification', 'update', '通知管理', 41050, 1),
('notification:delete', '删除通知',   '删除通知',               'notification', 'delete', '通知管理', 41100, 1),
('notification:send',   '发送通知',   '发送系统通知',           'notification', 'send',   '通知管理', 41150, 1);

-- ==================== 用户管理 (41500) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('user:read',   '查看用户',   '查看用户列表和详情',            'user', 'read',   '用户管理', 41500, 1),
('user:create', '创建用户',   '创建用户',                      'user', 'create', '用户管理', 41550, 1),
('user:update', '编辑用户',   '编辑用户信息',                  'user', 'update', '用户管理', 41600, 1),
('user:delete', '删除用户',   '删除用户',                      'user', 'delete', '用户管理', 41650, 1);

-- ==================== 权限管理 (42000) ====================
INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status) VALUES
('role:read',   '查看角色',   '查看角色列表和详情',            'role', 'read',   '权限管理(角色)', 42000, 1),
('role:create', '创建角色',   '创建角色',                      'role', 'create', '权限管理(角色)', 42050, 1),
('role:update', '编辑角色',   '编辑角色信息',                  'role', 'update', '权限管理(角色)', 42100, 1),
('role:delete', '删除角色',   '删除角色',                      'role', 'delete', '权限管理(角色)', 42150, 1);

-- ========================================================================
-- 角色数据
-- ========================================================================
INSERT INTO roles (name, display_name, description, status, sort, is_system) VALUES
('admin',     '管理员',     '系统管理员，拥有所有权限',            1, 1,  1),
('operator',  '运营人员',   '订单处理、退款审核、评论管理、通知发送', 1, 2,  1),
('editor',    '内容编辑',   '商品/分类/秒杀活动内容维护',           1, 3,  1),
('warehouse', '仓库管理员', '库存管理、订单发货处理',              1, 4,  1),
('finance',   '财务人员',   '支付对账、退款审核处理',              1, 5,  1),
('user',      '普通用户',   '普通用户，拥有基本操作权限',          1, 6,  1);

-- ========================================================================
-- 角色-权限关联
-- 使用子查询以名称关联，避免对 ID 的硬编码依赖，便于扩展维护
-- ========================================================================

-- admin 角色：拥有所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'admin'), id FROM permissions;

-- user 角色：拥有基础操作权限（浏览、下单、购物车、评论等）
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'user'), id FROM permissions WHERE name IN (
    -- 商品：可浏览
    'product:read',
    -- 分类：可浏览
    'category:read',
    -- 库存：可查看
    'inventory:read',
    -- 订单：基础操作（查看/创建/取消，无权编辑删除他人订单）
    'order:read',
    'order:create',
    'order:cancel',
    -- 购物车：完全管理（归属当前用户）
    'cart:read',
    'cart:add',
    'cart:update',
    'cart:delete',
    -- 支付：查看和发起（归属当前用户）
    'payment:read',
    'payment:create',
    -- 退款：查看和申请（归属当前用户）
    'refund:read',
    'refund:create',
    -- 秒杀：浏览和购买
    'flash:read',
    'flash:buy',
    -- 评论：查看/发表/删除自己评论
    'review:read',
    'review:create',
    'review:delete',
    -- 通知：查看和标记已读
    'notification:read',
    'notification:update',
    -- 优惠券：查看和领取（参与营销活动）
    'coupon:read',
    'coupon:claim',
    -- 促销：查看
    'promotion:read',
    -- 用户：查看自己和编辑自己信息
    'user:read',
    'user:update'
);

-- operator 角色：运营类操作（处理订单、退款审核、评论管理、通知推送）
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'operator'), id FROM permissions WHERE name IN (
    -- 商品：可查看
    'product:read',
    -- 分类：可查看
    'category:read',
    -- 库存：可查看
    'inventory:read',
    -- 订单：查看/编辑/取消（处理订单流转）
    'order:read',
    'order:update',
    'order:cancel',
    -- 支付：可查看
    'payment:read',
    -- 退款：查看和处理退款
    'refund:read',
    'refund:update',
    -- 秒杀：可浏览
    'flash:read',
    -- 评论：审核/回复评论
    'review:read',
    'review:moderate',
    -- 通知：管理（含发送系统通知）
    'notification:read',
    'notification:update',
    'notification:send',
    -- 优惠券：查看
    'coupon:read',
    -- 促销：查看
    'promotion:read',
    -- 用户：查看用户信息
    'user:read'
);

-- editor 角色：内容维护（商品/分类/评论/秒杀活动/优惠券/促销）
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'editor'), id FROM permissions WHERE name IN (
    -- 商品：查看/创建/编辑（无权删除）
    'product:read',
    'product:create',
    'product:update',
    -- 分类：查看/创建/编辑（无权删除）
    'category:read',
    'category:create',
    'category:update',
    -- 库存：可查看
    'inventory:read',
    -- 订单：可查看
    'order:read',
    -- 秒杀：查看/创建/编辑/管理
    'flash:read',
    'flash:create',
    'flash:update',
    'flash:manage',
    -- 优惠券：查看/创建/编辑
    'coupon:read',
    'coupon:create',
    'coupon:update',
    -- 促销：查看/创建/编辑
    'promotion:read',
    'promotion:create',
    'promotion:update',
    -- 评论：查看和审核
    'review:read',
    'review:moderate',
    -- 通知：查看
    'notification:read',
    -- 用户：查看用户信息
    'user:read'
);

-- warehouse 角色：库存管理与订单发货
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'warehouse'), id FROM permissions WHERE name IN (
    -- 商品：可查看
    'product:read',
    -- 分类：可查看
    'category:read',
    -- 库存：全部权限（创建/编辑/预订操作）
    'inventory:read',
    'inventory:create',
    'inventory:update',
    'inventory:reserve',
    -- 订单：查看和更新状态（发货）
    'order:read',
    'order:update',
    -- 通知：查看
    'notification:read'
);

-- finance 角色：财务对账与退款审核
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'finance'), id FROM permissions WHERE name IN (
    -- 订单：可查看
    'order:read',
    -- 支付：查看和更新状态（对账）
    'payment:read',
    'payment:update',
    -- 退款：查看和处理
    'refund:read',
    'refund:update',
    -- 商品：可查看
    'product:read',
    -- 通知：查看
    'notification:read',
    -- 用户：查看用户信息
    'user:read'
);

-- ========================================================================
-- 【运维工具】模块排序调整模板
-- 当需要在模块间插入新模块时，将现有模块整体偏移即可
-- ========================================================================
--
-- -- 按 resource 整体偏移（推荐）：
-- UPDATE permissions SET sort = sort + 1000 WHERE resource = 'coupon';
-- UPDATE permissions SET sort = sort + 1000 WHERE resource = 'promotion';
--
-- -- 按 category 整体偏移：
-- UPDATE permissions SET sort = sort + 500 WHERE category = '优惠券管理';
--
-- -- 查冲突：确认偏移后的值不与现有值重叠
-- SELECT sort FROM permissions ORDER BY sort;
--
-- ========================================================================
-- 用户数据
-- 密码均为 "123456"，bcrypt hash（cost=10）
-- 如需更换密码: go run ./cmd/genhash/ 生成新 hash 后替换
-- ========================================================================

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
