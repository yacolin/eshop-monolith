USE eshop_db;

-- 插入测试用户
INSERT INTO users (username, email, password, role) VALUES
('admin', 'admin@example.com', '$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW', 'admin'),
('user1', 'user1@example.com', '$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW', 'user'),
('user2', 'user2@example.com', '$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW', 'user');

-- 插入测试产品
INSERT INTO products (name, description, price, sku) VALUES
('iPhone 15 Pro', 'Apple iPhone 15 Pro with A17 Pro chip', 999900, 'IPHONE15PRO-001'),
('MacBook Pro 14', 'Apple MacBook Pro 14 inch M3 Pro', 1999900, 'MACBOOKPRO14-001'),
('AirPods Pro 2', 'Apple AirPods Pro 2nd generation', 189900, 'AIRPODSPRO2-001');

-- 插入测试库存
INSERT INTO inventories (product_id, quantity, reserved_quantity) VALUES
(1, 100, 0),
(2, 50, 0),
(3, 200, 0);

-- 插入测试角色
INSERT INTO roles (name, description) VALUES
('admin', 'Administrator role'),
('user', 'Regular user role');

-- 插入测试权限
INSERT INTO permissions (name, description) VALUES
('manage_users', 'Manage users'),
('manage_products', 'Manage products'),
('manage_orders', 'Manage orders'),
('view_products', 'View products'),
('place_orders', 'Place orders');

-- 关联角色权限
INSERT INTO role_permissions (role_id, permission_id) VALUES
(1, 1),
(1, 2),
(1, 3),
(1, 4),
(1, 5),
(2, 4),
(2, 5);