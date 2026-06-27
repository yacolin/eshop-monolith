-- ============================================================================
-- RBAC 角色-权限关联脚本（第 4 步）
-- 使用子查询以名称关联，避免对 ID 的硬编码依赖
--
-- 执行: mysql -u root -p eshop_db < scripts/seed_rbac_04_assignments.sql
-- 顺序: 01 → 02 → 03 → 04 → 05
-- ============================================================================

USE eshop_db;

-- ========================================================================
-- 角色-权限关联
-- ========================================================================

-- ==================== admin 角色：拥有所有权限 ====================
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'admin'), id FROM permissions;

-- ==================== user 角色：基础购物操作 ====================
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'user'), id FROM permissions WHERE name IN (
    -- 商品：可浏览
    'product:read',
    -- 分类：可浏览
    'category:read',
    -- 库存：可查看
    'inventory:read',
    -- SKU：可查看
    'sku:read',
    -- 规格属性：可查看
    'attr:read',
    'attr_val:read',
    -- 地址：管理自身地址
    'address:read',
    'address:create',
    'address:update',
    'address:delete',
    -- 订单：基础操作（查看/创建/取消）
    'order:read',
    'order:create',
    'order:cancel',
    -- 购物车：完全管理（归属当前用户）
    'cart:read',
    'cart:create',
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
    -- 优惠券：查看和领取
    'coupon:read',
    'coupon:claim',
    -- 促销：查看
    'promotion:read',
    -- 用户：查看自己和编辑自己信息
    'user:read',
    'user:update'
);

-- ==================== operator 角色：运营操作 ====================
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'operator'), id FROM permissions WHERE name IN (
    -- 商品：可查看
    'product:read',
    -- 分类：可查看
    'category:read',
    -- 库存：可查看
    'inventory:read',
    -- SKU：可查看
    'sku:read',
    -- 规格属性：可查看
    'attr:read',
    'attr_val:read',
    -- 地址：可查看
    'address:read',
    -- 订单：查看/编辑/取消
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

-- ==================== editor 角色：内容维护 ====================
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
    -- SKU：可查看
    'sku:read',
    -- 规格属性：可查看
    'attr:read',
    'attr_val:read',
    -- 地址：可查看
    'address:read',
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

-- ==================== warehouse 角色：库存与发货 ====================
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'warehouse'), id FROM permissions WHERE name IN (
    -- 商品：可查看
    'product:read',
    -- 分类：可查看
    'category:read',
    -- 库存：全部权限
    'inventory:read',
    'inventory:create',
    'inventory:update',
    'inventory:reserve',
    -- SKU：可查看
    'sku:read',
    -- 地址：可查看（发货需要）
    'address:read',
    -- 订单：查看和更新状态（发货）
    'order:read',
    'order:update',
    -- 通知：查看
    'notification:read'
);

-- ==================== finance 角色：财务与退款审核 ====================
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
    -- SKU：可查看
    'sku:read',
    -- 地址：可查看（退款上下文）
    'address:read',
    -- 通知：查看
    'notification:read',
    -- 用户：查看用户信息
    'user:read'
);

-- ==================== merchant 角色：商户管理 ====================
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'merchant'), id FROM permissions WHERE name IN (
    -- 商品：管理自身商品（不含删除）
    'product:read',
    'product:create',
    'product:update',
    -- 分类：可查看
    'category:read',
    -- 库存：可查看
    'inventory:read',
    -- SKU：管理自身商品 SKU
    'sku:read',
    'sku:create',
    'sku:update',
    -- 规格属性：可查看
    'attr:read',
    'attr_val:read',
    -- 地址：可查看
    'address:read',
    -- 订单：查看和处理（管理自身订单）
    'order:read',
    'order:update',
    'order:cancel',
    -- 支付：可查看
    'payment:read',
    -- 退款：可查看
    'refund:read',
    -- 秒杀：可查看
    'flash:read',
    -- 优惠券：可查看
    'coupon:read',
    -- 促销：可查看
    'promotion:read',
    -- 评论：可查看
    'review:read',
    -- 通知：查看和标记已读
    'notification:read',
    'notification:update',
    -- 仪表盘：查看数据
    'dashboard:read',
    -- 用户：查看用户信息
    'user:read'
);

-- ==================== support 角色：客服售后 ====================
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'support'), id FROM permissions WHERE name IN (
    -- 商品：可查看
    'product:read',
    -- 分类：可查看
    'category:read',
    -- SKU：可查看（客服上下文）
    'sku:read',
    -- 规格属性：可查看
    'attr:read',
    'attr_val:read',
    -- 地址：可查看（处理订单问题）
    'address:read',
    -- 订单：查看/编辑/取消（处理订单问题）
    'order:read',
    'order:update',
    'order:cancel',
    -- 支付：可查看
    'payment:read',
    -- 退款：查看和处理退款
    'refund:read',
    'refund:update',
    -- 秒杀：可查看
    'flash:read',
    -- 评论：审核和回复
    'review:read',
    'review:moderate',
    -- 通知：管理（含发送通知给用户）
    'notification:read',
    'notification:update',
    'notification:send',
    -- 用户：查看用户信息
    'user:read',
    -- 仪表盘：查看数据
    'dashboard:read'
);

-- ==================== sales 角色：销售营销 ====================
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'sales'), id FROM permissions WHERE name IN (
    -- 商品：可查看
    'product:read',
    -- 分类：可查看
    'category:read',
    -- 订单：可查看（了解销售情况）
    'order:read',
    -- SKU：可查看
    'sku:read',
    -- 规格属性：可查看
    'attr:read',
    'attr_val:read',
    -- 地址：可查看
    'address:read',
    -- 秒杀：创建和管理秒杀活动
    'flash:read',
    'flash:create',
    'flash:update',
    -- 优惠券：创建和管理优惠券（不含删除）
    'coupon:read',
    'coupon:create',
    'coupon:update',
    -- 促销：创建和管理促销（不含删除）
    'promotion:read',
    'promotion:create',
    'promotion:update',
    -- 通知：查看和标记已读
    'notification:read',
    'notification:update',
    -- 用户：查看用户信息
    'user:read',
    -- 仪表盘：查看数据
    'dashboard:read'
);

-- ==================== analyst 角色：数据分析 ====================
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'analyst'), id FROM permissions WHERE name IN (
    -- 只读权限：查看所有业务数据用于分析
    'product:read',
    'category:read',
    'inventory:read',
    'order:read',
    'cart:read',
    'payment:read',
    'refund:read',
    'flash:read',
    'coupon:read',
    'promotion:read',
    'review:read',
    'notification:read',
    'user:read',
    'dashboard:read',
    -- SKU 和规格属性：只读（数据分析）
    'sku:read',
    'attr:read',
    'attr_val:read',
    -- 地址：只读（数据分析）
    'address:read'
);

-- ==================== guest 角色：访客浏览 ====================
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'guest'), id FROM permissions WHERE name IN (
    -- 最基础的浏览权限
    'product:read',
    'category:read',
    'review:read',
    'flash:read',
    'promotion:read',
    -- SKU：浏览商品时可见
    'sku:read',
    'attr:read',
    'attr_val:read'
);

-- ==================== customer 角色：客户售后 ====================
INSERT INTO role_permissions (role_id, permission_id)
SELECT (SELECT id FROM roles WHERE name = 'customer'), id FROM permissions WHERE name IN (
    -- 商品：可查看
    'product:read',
    -- 分类：可查看
    'category:read',
    -- 库存：可查看
    'inventory:read',
    -- SKU：可查看
    'sku:read',
    -- 规格属性：可查看
    'attr:read',
    'attr_val:read',
    -- 地址：管理自身地址
    'address:read',
    'address:create',
    'address:update',
    'address:delete',
    -- 订单：完整操作（可更新售后状态）
    'order:read',
    'order:create',
    'order:update',
    'order:cancel',
    -- 购物车：完整管理
    'cart:read',
    'cart:create',
    'cart:update',
    'cart:delete',
    -- 支付：查看和发起
    'payment:read',
    'payment:create',
    -- 退款：完整售后（含处理）
    'refund:read',
    'refund:create',
    'refund:update',
    -- 秒杀：浏览和购买
    'flash:read',
    'flash:buy',
    -- 评论：查看和发表
    'review:read',
    'review:create',
    -- 通知：查看和标记已读
    'notification:read',
    'notification:update',
    -- 优惠券：查看和领取
    'coupon:read',
    'coupon:claim',
    -- 促销：查看
    'promotion:read',
    -- 用户：查看自己和编辑自己信息
    'user:read',
    'user:update'
);
