-- ============================================================================
-- RBAC 权限初始化脚本（第 2 步）
--
-- 执行: mysql -u root -p eshop_db < scripts/seed_rbac_02_permissions.sql
-- 顺序: 01 → 02 → 03 → 04 → 05
-- ============================================================================

USE eshop_db;

-- ═══════════════════════════════════════════════════════════════════════════
-- 大类一：商品库存 (10000-19999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('product:read',   '查看产品',   '查看产品列表和详情',          'product',  'read',   '商品管理', 11000, 1),
('product:create', '创建产品',   '创建新产品',                  'product',  'create', '商品管理', 11050, 1),
('product:update', '编辑产品',   '编辑产品信息',                'product',  'update', '商品管理', 11100, 1),
('product:delete', '删除产品',   '删除产品',                    'product',  'delete', '商品管理', 11150, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('category:read',   '查看分类',   '查看分类列表和详情',          'category', 'read',   '分类管理', 11500, 1),
('category:create', '创建分类',   '创建分类',                    'category', 'create', '分类管理', 11550, 1),
('category:update', '编辑分类',   '编辑分类信息',                'category', 'update', '分类管理', 11600, 1),
('category:delete', '删除分类',   '删除分类',                    'category', 'delete', '分类管理', 11650, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('inventory:read',    '查看库存',   '查看库存信息',              'inventory', 'read',    '库存管理', 12000, 1),
('inventory:create',  '创建库存',   '创建库存记录',              'inventory', 'create',  '库存管理', 12050, 1),
('inventory:update',  '编辑库存',   '编辑库存信息',              'inventory', 'update',  '库存管理', 12100, 1),
('inventory:reserve', '库存操作',   '库存预订与释放操作',        'inventory', 'reserve', '库存管理', 12150, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('sku:read',   '查看 SKU',   '查看 SKU 列表和详情',            'sku', 'read',   'SKU 管理', 12500, 1),
('sku:create', '创建 SKU',   '创建新的 SKU',                    'sku', 'create', 'SKU 管理', 12550, 1),
('sku:update', '编辑 SKU',   '更新 SKU 信息',                   'sku', 'update', 'SKU 管理', 12600, 1),
('sku:delete', '删除 SKU',   '删除 SKU',                       'sku', 'delete', 'SKU 管理', 12650, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('attr:read',      '查看属性维度', '查看规格属性维度列表和详情', 'attr',     'read',   '规格属性管理', 13000, 1),
('attr:create',    '创建属性维度', '创建新的规格属性维度',       'attr',     'create', '规格属性管理', 13050, 1),
('attr:update',    '编辑属性维度', '更新规格属性维度信息',       'attr',     'update', '规格属性管理', 13100, 1),
('attr:delete',    '删除属性维度', '删除规格属性维度',           'attr',     'delete', '规格属性管理', 13150, 1),
('attr_val:read',   '查看属性值', '查看属性可选值列表',         'attr_val', 'read',   '规格属性管理', 13250, 1),
('attr_val:create', '创建属性值', '为属性维度创建可选值',       'attr_val', 'create', '规格属性管理', 13300, 1),
('attr_val:update', '编辑属性值', '更新属性可选值信息',         'attr_val', 'update', '规格属性管理', 13350, 1),
('attr_val:delete', '删除属性值', '删除属性可选值',             'attr_val', 'delete', '规格属性管理', 13400, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- 大类二：交易订单 (20000-29999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('order:read',         '查看订单',   '查看订单列表和详情',        'order', 'read',         '订单管理', 21000, 1),
('order:create',       '创建订单',   '创建订单',                  'order', 'create',       '订单管理', 21050, 1),
('order:update',       '编辑订单',   '编辑订单信息',              'order', 'update',       '订单管理', 21100, 1),
('order:cancel',       '取消订单',   '取消订单',                  'order', 'cancel',       '订单管理', 21150, 1),
('order:delete',       '删除订单',   '删除订单',                  'order', 'delete',       '订单管理', 21200, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('address:read',   '查看地址',   '查看收货地址列表和详情',      'address', 'read',   '地址管理', 21250, 1),
('address:create', '创建地址',   '添加新的收货地址',            'address', 'create', '地址管理', 21300, 1),
('address:update', '编辑地址',   '编辑收货地址信息',            'address', 'update', '地址管理', 21350, 1),
('address:delete', '删除地址',   '删除收货地址',                'address', 'delete', '地址管理', 21400, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('cart:read',   '查看购物车',   '查看购物车内容',               'cart', 'read',   '购物车管理', 21500, 1),
('cart:create',    '添加商品',     '添加商品到购物车',          'cart', 'add',    '购物车管理', 21550, 1),
('cart:update', '编辑购物车',   '更新购物车商品信息',           'cart', 'update', '购物车管理', 21600, 1),
('cart:delete', '删除商品',     '删除购物车中的商品',           'cart', 'delete', '购物车管理', 21650, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('payment:read',   '查看支付',   '查看支付记录和详情',          'payment', 'read',   '支付管理', 22000, 1),
('payment:create', '发起支付',   '发起支付请求',                'payment', 'create', '支付管理', 22050, 1),
('payment:update', '更新支付',   '更新支付状态',                'payment', 'update', '支付管理', 22100, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('refund:read',   '查看退款',   '查看退款记录和详情',           'refund', 'read',   '退款管理', 22500, 1),
('refund:create', '申请退款',   '创建退款申请',                 'refund', 'create', '退款管理', 22550, 1),
('refund:update', '处理退款',   '更新退款状态',                 'refund', 'update', '退款管理', 22600, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('promotion:read',   '查看促销',   '查看促销活动列表和详情',      'promotion', 'read',   '营销管理', 24000, 1),
('promotion:create', '创建促销',   '创建促销活动',                'promotion', 'create', '营销管理', 24050, 1),
('promotion:update', '编辑促销',   '编辑促销活动信息',            'promotion', 'update', '营销管理', 24100, 1),
('promotion:delete', '删除促销',   '删除促销活动',                'promotion', 'delete', '营销管理', 24150, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- 大类三：评价反馈 (30000-39999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('review:read',     '查看评论',   '查看评论列表和详情',         'review', 'read',     '评论管理', 31000, 1),
('review:create',   '发表评论',   '发表商品评论',               'review', 'create',   '评论管理', 31050, 1),
('review:delete',   '删除评论',   '删除评论',                   'review', 'delete',   '评论管理', 31100, 1),
('review:moderate', '审核评论',   '审核/回复评论',              'review', 'moderate', '评论管理', 31150, 1);

-- ═══════════════════════════════════════════════════════════════════════════
-- 大类四：用户与系统 (40000-49999)
-- ═══════════════════════════════════════════════════════════════════════════

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('dashboard:read', '查看仪表盘', '查看仪表盘汇总数据',          'dashboard', 'read', '仪表盘管理', 40500, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('notification:read',   '查看通知',   '查看通知列表和详情',     'notification', 'read',   '通知管理', 41000, 1),
('notification:update', '标记已读',   '标记通知为已读',         'notification', 'update', '通知管理', 41050, 1),
('notification:delete', '删除通知',   '删除通知',               'notification', 'delete', '通知管理', 41100, 1),
('notification:send',   '发送通知',   '发送系统通知',           'notification', 'send',   '通知管理', 41150, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('user:read',   '查看用户',   '查看用户列表和详情',            'user', 'read',   '用户管理', 41500, 1),
('user:create', '创建用户',   '创建用户',                      'user', 'create', '用户管理', 41550, 1),
('user:update', '编辑用户',   '编辑用户信息',                  'user', 'update', '用户管理', 41600, 1),
('user:delete', '删除用户',   '删除用户',                      'user', 'delete', '用户管理', 41650, 1);

INSERT INTO usr_permissions (name, display_name, description, resource, action, category, sort_order, status) VALUES
('role:read',   '查看角色',   '查看角色列表和详情',            'role', 'read',   '权限管理', 42000, 1),
('role:create', '创建角色',   '创建角色',                      'role', 'create', '权限管理', 42050, 1),
('role:update', '编辑角色',   '编辑角色信息',                  'role', 'update', '权限管理', 42100, 1),
('role:delete', '删除角色',   '删除角色',                      'role', 'delete', '权限管理', 42150, 1);
