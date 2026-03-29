USE eshop_db;

-- 权限数据
INSERT INTO permissions (name, description) VALUES
('manage_users', 'Manage users'),
('manage_products', 'Manage products'),
('manage_orders', 'Manage orders'),
('view_products', 'View products'),
('place_orders', 'Place orders'),
('view_orders', 'View orders'),
('manage_inventory', 'Manage inventory'),
('view_inventory', 'View inventory');

-- 角色数据
INSERT INTO roles (name, description) VALUES
('admin', 'Administrator role'),
('user', 'Regular user role'),
('manager', 'Manager role');

-- 角色权限关联
INSERT INTO role_permissions (role_id, permission_id) VALUES
-- Admin permissions
(1, 1),
(1, 2),
(1, 3),
(1, 4),
(1, 5),
(1, 6),
(1, 7),
(1, 8),
-- User permissions
(2, 4),
(2, 5),
(2, 6),
-- Manager permissions
(3, 2),
(3, 3),
(3, 4),
(3, 6),
(3, 7),
(3, 8);