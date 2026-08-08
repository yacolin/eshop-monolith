-- ------------------------------------------------------------
-- Table: base_notification_reads
-- ------------------------------------------------------------
CREATE TABLE `base_notification_reads` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `notification_id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `read_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_notification_user` (`notification_id`,`user_id`),
  KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: base_notification_templates
-- ------------------------------------------------------------
CREATE TABLE `base_notification_templates` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `template_code` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `channel` tinyint NOT NULL,
  `title_template` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content_template` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `category` tinyint DEFAULT NULL,
  `priority` tinyint NOT NULL DEFAULT '1',
  `status` tinyint NOT NULL DEFAULT '1',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_base_notification_templates_template_code` (`template_code`),
  KEY `idx_code_channel` (`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: base_notifications
-- ------------------------------------------------------------
CREATE TABLE `base_notifications` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `merchant_id` bigint NOT NULL DEFAULT '0',
  `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `content_template` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `template_params` json DEFAULT NULL,
  `channel` tinyint NOT NULL,
  `category` tinyint NOT NULL,
  `target_type` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `target_id` bigint DEFAULT NULL,
  `redirect_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `icon_url` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `priority` tinyint NOT NULL DEFAULT '1',
  `created_by` bigint NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  `is_read` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_channel` (`channel`),
  KEY `idx_category` (`category`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: mkt_promotion_products
-- ------------------------------------------------------------
CREATE TABLE `mkt_promotion_products` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `promotion_id` bigint NOT NULL,
  `product_type` tinyint NOT NULL,
  `product_id` bigint DEFAULT NULL,
  `category_id` bigint DEFAULT NULL,
  `created_by` bigint DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_promo_product` (`promotion_id`,`product_type`,`product_id`,`category_id`),
  KEY `idx_promotion` (`promotion_id`),
  KEY `idx_mkt_promotion_products_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: mkt_promotion_rules
-- ------------------------------------------------------------
CREATE TABLE `mkt_promotion_rules` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `promotion_id` bigint NOT NULL,
  `rule_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `condition_type` tinyint NOT NULL,
  `condition_value` decimal(10,2) DEFAULT NULL,
  `benefit_type` tinyint NOT NULL,
  `benefit_value` decimal(10,2) DEFAULT NULL,
  `is_stackable` tinyint NOT NULL DEFAULT '0',
  `stack_priority` bigint NOT NULL DEFAULT '0',
  `created_by` bigint DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_promotion` (`promotion_id`),
  KEY `idx_stack` (`is_stackable`,`stack_priority`),
  KEY `idx_mkt_promotion_rules_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: mkt_promotion_usage_logs
-- ------------------------------------------------------------
CREATE TABLE `mkt_promotion_usage_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `promotion_id` bigint NOT NULL,
  `user_id` bigint NOT NULL,
  `order_id` bigint DEFAULT NULL,
  `action_type` tinyint NOT NULL,
  `fail_reason` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `request_params` json DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_promotion` (`promotion_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_order` (`order_id`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: mkt_promotions
-- ------------------------------------------------------------
CREATE TABLE `mkt_promotions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `promo_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `promo_type` tinyint NOT NULL,
  `promo_code` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `start_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `total_quantity` bigint NOT NULL DEFAULT '0',
  `per_user_limit` bigint NOT NULL DEFAULT '0',
  `used_quantity` bigint NOT NULL DEFAULT '0',
  `rule_id` bigint DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT '1',
  `created_by` bigint DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`promo_code`),
  KEY `idx_time` (`start_time`,`end_time`),
  KEY `idx_mkt_promotions_rule_id` (`rule_id`),
  KEY `idx_status` (`status`),
  KEY `idx_mkt_promotions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: mkt_user_promotions
-- ------------------------------------------------------------
CREATE TABLE `mkt_user_promotions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `promotion_id` bigint NOT NULL,
  `acquire_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `expire_time` datetime DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT '1',
  `used_time` datetime DEFAULT NULL,
  `order_id` bigint DEFAULT NULL,
  `queue_token` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `created_by` bigint DEFAULT NULL,
  `updated_by` bigint DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_promo` (`user_id`,`promotion_id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_promotion` (`promotion_id`),
  KEY `idx_status_expire` (`status`),
  KEY `idx_order` (`order_id`),
  KEY `idx_mkt_user_promotions_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: rev_review_audit_logs
-- ------------------------------------------------------------
CREATE TABLE `rev_review_audit_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `review_id` bigint NOT NULL,
  `auditor_id` bigint NOT NULL,
  `action` tinyint NOT NULL,
  `reason` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `sensitive_words` json DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_review` (`review_id`),
  KEY `idx_auditor` (`auditor_id`),
  KEY `idx_created` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: rev_review_media
-- ------------------------------------------------------------
CREATE TABLE `rev_review_media` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `review_id` bigint NOT NULL,
  `media_type` tinyint NOT NULL,
  `media_url` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sort_order` bigint NOT NULL DEFAULT '0',
  `audit_status` tinyint NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_review` (`review_id`),
  KEY `idx_audit` (`audit_status`),
  KEY `idx_rev_review_media_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: rev_review_ratings
-- ------------------------------------------------------------
CREATE TABLE `rev_review_ratings` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `spu_id` bigint NOT NULL,
  `avg_overall_rating` decimal(2,1) NOT NULL DEFAULT '0.0',
  `avg_quality_rating` decimal(2,1) NOT NULL DEFAULT '0.0',
  `avg_logistics_rating` decimal(2,1) NOT NULL DEFAULT '0.0',
  `avg_service_rating` decimal(2,1) NOT NULL DEFAULT '0.0',
  `rating5_count` bigint NOT NULL DEFAULT '0',
  `rating4_count` bigint NOT NULL DEFAULT '0',
  `rating3_count` bigint NOT NULL DEFAULT '0',
  `rating2_count` bigint NOT NULL DEFAULT '0',
  `rating1_count` bigint NOT NULL DEFAULT '0',
  `total_reviews` bigint NOT NULL DEFAULT '0',
  `with_media_count` bigint NOT NULL DEFAULT '0',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_spu` (`spu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: rev_review_replies
-- ------------------------------------------------------------
CREATE TABLE `rev_review_replies` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `review_id` bigint NOT NULL,
  `reply_by` bigint NOT NULL,
  `reply_by_type` tinyint NOT NULL DEFAULT '1',
  `content` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint NOT NULL DEFAULT '1',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_review` (`review_id`),
  KEY `idx_reply_by` (`reply_by`),
  KEY `idx_rev_review_replies_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: rev_reviews
-- ------------------------------------------------------------
CREATE TABLE `rev_reviews` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `order_id` bigint NOT NULL,
  `order_item_id` bigint DEFAULT NULL,
  `spu_id` bigint NOT NULL,
  `sku_id` bigint DEFAULT NULL,
  `overall_rating` tinyint NOT NULL,
  `quality_rating` tinyint DEFAULT NULL,
  `logistics_rating` tinyint DEFAULT NULL,
  `service_rating` tinyint DEFAULT NULL,
  `content` text COLLATE utf8mb4_unicode_ci,
  `is_anonymous` tinyint(1) NOT NULL DEFAULT '0',
  `status` tinyint NOT NULL DEFAULT '0',
  `reject_reason` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `latest_reply_id` bigint DEFAULT NULL,
  `reply_count` bigint NOT NULL DEFAULT '0',
  `like_count` bigint NOT NULL DEFAULT '0',
  `helpful_count` bigint NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user` (`user_id`),
  KEY `idx_order` (`order_id`),
  KEY `idx_rev_reviews_order_item_id` (`order_item_id`),
  KEY `idx_spu` (`spu_id`),
  KEY `idx_sku` (`sku_id`),
  KEY `idx_status_created` (`status`),
  KEY `idx_rev_reviews_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sp_attributes
-- ------------------------------------------------------------
CREATE TABLE `sp_attributes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `category_id` bigint NOT NULL,
  `input_type` tinyint NOT NULL DEFAULT '1',
  `values` json DEFAULT NULL,
  `unit` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `required` tinyint NOT NULL DEFAULT '0',
  `searchable` tinyint NOT NULL DEFAULT '0',
  `is_sku_spec` tinyint NOT NULL DEFAULT '0',
  `sort_order` bigint NOT NULL DEFAULT '0',
  `status` tinyint NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category` (`category_id`),
  KEY `idx_searchable` (`searchable`),
  KEY `idx_is_sku_spec` (`is_sku_spec`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sp_brands
-- ------------------------------------------------------------
CREATE TABLE `sp_brands` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `english_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `logo_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `first_letter` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `sort_order` bigint NOT NULL DEFAULT '0',
  `status` tinyint NOT NULL DEFAULT '1',
  `description` text COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_first_letter` (`first_letter`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sp_categories
-- ------------------------------------------------------------
CREATE TABLE `sp_categories` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `parent_id` bigint NOT NULL DEFAULT '0',
  `level` tinyint NOT NULL DEFAULT '1',
  `path` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `icon_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `sort_order` bigint NOT NULL DEFAULT '0',
  `status` tinyint NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_level_status` (`level`,`status`),
  KEY `idx_path` (`path`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sp_category_brands
-- ------------------------------------------------------------
CREATE TABLE `sp_category_brands` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `category_id` bigint NOT NULL,
  `brand_id` bigint NOT NULL,
  `sort_order` bigint NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_category_brand` (`category_id`,`brand_id`),
  KEY `idx_brand_id` (`brand_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sp_inventories
-- ------------------------------------------------------------
CREATE TABLE `sp_inventories` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `sku_id` bigint NOT NULL,
  `warehouse_id` bigint NOT NULL DEFAULT '0',
  `quantity` bigint NOT NULL DEFAULT '0',
  `reserved` bigint NOT NULL DEFAULT '0',
  `in_transit` bigint NOT NULL DEFAULT '0',
  `threshold` bigint NOT NULL DEFAULT '10',
  `max_threshold` bigint NOT NULL DEFAULT '999999',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'instock',
  `last_counted_at` datetime DEFAULT NULL,
  `last_counted_by` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_warehouse` (`sku_id`,`warehouse_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sp_inventory_logs
-- ------------------------------------------------------------
CREATE TABLE `sp_inventory_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `sku_id` bigint NOT NULL,
  `warehouse_id` bigint NOT NULL DEFAULT '0',
  `change_type` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL,
  `before_quantity` bigint NOT NULL DEFAULT '0',
  `after_quantity` bigint NOT NULL DEFAULT '0',
  `before_reserved` bigint NOT NULL DEFAULT '0',
  `after_reserved` bigint NOT NULL DEFAULT '0',
  `change_amount` bigint NOT NULL DEFAULT '0',
  `reference_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `operator` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `note` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_sku_id` (`sku_id`),
  KEY `idx_change_type` (`change_type`),
  KEY `idx_reference_id` (`reference_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sp_product_attributes
-- ------------------------------------------------------------
CREATE TABLE `sp_product_attributes` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint NOT NULL,
  `attribute_id` bigint NOT NULL,
  `value` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sort_order` bigint NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_attribute` (`product_id`,`attribute_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sp_product_descriptions
-- ------------------------------------------------------------
CREATE TABLE `sp_product_descriptions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint NOT NULL,
  `description` longtext COLLATE utf8mb4_unicode_ci,
  `mobile_description` longtext COLLATE utf8mb4_unicode_ci,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sp_products
-- ------------------------------------------------------------
CREATE TABLE `sp_products` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `subtitle` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `category_id` bigint NOT NULL DEFAULT '0',
  `brand_id` bigint NOT NULL DEFAULT '0',
  `unit` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT '件',
  `main_image` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `images` json DEFAULT NULL,
  `video_url` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `min_price` bigint NOT NULL DEFAULT '0',
  `max_price` bigint NOT NULL DEFAULT '0',
  `total_stock` bigint NOT NULL DEFAULT '0',
  `sales_count` bigint NOT NULL DEFAULT '0',
  `rating_average` decimal(3,2) NOT NULL DEFAULT '0.00',
  `rating_count` bigint NOT NULL DEFAULT '0',
  `status` tinyint NOT NULL DEFAULT '0',
  `sort_order` bigint NOT NULL DEFAULT '0',
  `has_description` tinyint NOT NULL DEFAULT '0',
  `created_by` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `updated_by` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_status_sort` (`category_id`),
  KEY `idx_brand_status` (`brand_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sp_skus
-- ------------------------------------------------------------
CREATE TABLE `sp_skus` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `product_id` bigint NOT NULL,
  `sku_code` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `barcode` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `spec` json NOT NULL,
  `price` bigint NOT NULL,
  `market_price` bigint NOT NULL DEFAULT '0',
  `cost_price` bigint NOT NULL DEFAULT '0',
  `weight` decimal(10,2) NOT NULL DEFAULT '0.00',
  `volume` decimal(10,2) NOT NULL DEFAULT '0.00',
  `length` decimal(10,2) NOT NULL DEFAULT '0.00',
  `width` decimal(10,2) NOT NULL DEFAULT '0.00',
  `height` decimal(10,2) NOT NULL DEFAULT '0.00',
  `min_purchase_qty` bigint NOT NULL DEFAULT '1',
  `max_purchase_qty` bigint NOT NULL DEFAULT '0',
  `image` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `status` tinyint NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_sku_code` (`sku_code`),
  UNIQUE KEY `uk_barcode` (`barcode`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sys_login_histories
-- ------------------------------------------------------------
CREATE TABLE `sys_login_histories` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `staff_id` bigint NOT NULL,
  `login_ip` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `login_device` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `login_location` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `login_method` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `login_status` tinyint(1) NOT NULL DEFAULT '1',
  `failure_reason` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_staff_id` (`staff_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sys_permissions
-- ------------------------------------------------------------
CREATE TABLE `sys_permissions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `display_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `resource` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `action` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `parent_id` bigint NOT NULL DEFAULT '0',
  `category` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `sort_order` bigint NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_resource` (`resource`),
  KEY `idx_action` (`action`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sys_role_permissions
-- ------------------------------------------------------------
CREATE TABLE `sys_role_permissions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `role_id` bigint NOT NULL,
  `permission_id` bigint NOT NULL,
  `scope_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'platform',
  `scope_id` bigint NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_role_permission` (`role_id`,`permission_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sys_roles
-- ------------------------------------------------------------
CREATE TABLE `sys_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `display_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `role_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'custom',
  `sort_order` bigint NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sys_staff
-- ------------------------------------------------------------
CREATE TABLE `sys_staff` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `real_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `avatar` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '1',
  `last_login_ip` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `last_login_at` datetime(3) DEFAULT NULL,
  `created_by` bigint NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: sys_staff_roles
-- ------------------------------------------------------------
CREATE TABLE `sys_staff_roles` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `staff_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_staff_role` (`staff_id`,`role_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: tx_cart_items
-- ------------------------------------------------------------
CREATE TABLE `tx_cart_items` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `cart_id` bigint NOT NULL,
  `sku_id` bigint NOT NULL,
  `product_id` bigint NOT NULL DEFAULT '0',
  `product_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `sku_spec` json DEFAULT NULL,
  `image` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `price` bigint NOT NULL,
  `quantity` bigint NOT NULL DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_cart_id` (`cart_id`),
  KEY `idx_sku_id` (`sku_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: tx_carts
-- ------------------------------------------------------------
CREATE TABLE `tx_carts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL DEFAULT '0',
  `session_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `item_count` bigint NOT NULL DEFAULT '0',
  `total_amount` bigint NOT NULL DEFAULT '0',
  `expired_at` datetime DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_expired_at` (`expired_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: tx_order_items
-- ------------------------------------------------------------
CREATE TABLE `tx_order_items` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` bigint NOT NULL,
  `order_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `sku_id` bigint NOT NULL,
  `product_id` bigint NOT NULL DEFAULT '0',
  `sku_code` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `product_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `sku_spec` json DEFAULT NULL,
  `image` varchar(512) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `price` bigint NOT NULL,
  `quantity` bigint NOT NULL DEFAULT '1',
  `subtotal` bigint NOT NULL DEFAULT '0',
  `refund_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'none',
  `refund_amount` bigint NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_sku_id` (`sku_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: tx_order_logs
-- ------------------------------------------------------------
CREATE TABLE `tx_order_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_id` bigint NOT NULL,
  `order_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `from_status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `to_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `operator` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT 'system',
  `operator_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'system',
  `note` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_to_status` (`to_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: tx_orders
-- ------------------------------------------------------------
CREATE TABLE `tx_orders` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `order_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `user_id` bigint NOT NULL,
  `total_amount` bigint NOT NULL DEFAULT '0',
  `discount_amount` bigint NOT NULL DEFAULT '0',
  `shipping_fee` bigint NOT NULL DEFAULT '0',
  `pay_amount` bigint NOT NULL DEFAULT '0',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
  `payment_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'unpaid',
  `payment_method` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `consignee` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `province` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `city` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `district` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `detail_addr` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `zip_code` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `coupon_id` bigint DEFAULT NULL,
  `coupon_snapshot` json DEFAULT NULL,
  `buyer_remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `seller_remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `source` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pc',
  `paid_at` datetime DEFAULT NULL,
  `shipped_at` datetime DEFAULT NULL,
  `delivered_at` datetime DEFAULT NULL,
  `completed_at` datetime DEFAULT NULL,
  `closed_at` datetime DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`),
  KEY `idx_payment_status` (`payment_status`),
  KEY `idx_tx_orders_coupon_id` (`coupon_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: tx_payment_logs
-- ------------------------------------------------------------
CREATE TABLE `tx_payment_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `payment_id` bigint NOT NULL,
  `payment_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `channel` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `transaction_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `action` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL,
  `request_body` text COLLATE utf8mb4_unicode_ci,
  `response_body` text COLLATE utf8mb4_unicode_ci,
  `status` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_payment_id` (`payment_id`),
  KEY `idx_payment_no` (`payment_no`),
  KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_action` (`action`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: tx_payments
-- ------------------------------------------------------------
CREATE TABLE `tx_payments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `payment_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `order_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `order_id` bigint NOT NULL,
  `order_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'order',
  `amount` bigint NOT NULL,
  `currency` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'CNY',
  `payment_method` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `channel` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `transaction_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
  `failure_reason` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `paid_at` datetime DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_payment_no` (`payment_no`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: tx_refunds
-- ------------------------------------------------------------
CREATE TABLE `tx_refunds` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `refund_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `payment_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `order_no` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL,
  `order_id` bigint NOT NULL,
  `amount` bigint NOT NULL,
  `reason` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'pending',
  `channel_transaction_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `failure_reason` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `applied_at` datetime DEFAULT NULL,
  `success_at` datetime DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_refund_no` (`refund_no`),
  KEY `idx_payment_no` (`payment_no`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: usr_addresses
-- ------------------------------------------------------------
CREATE TABLE `usr_addresses` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `consignee` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `country` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `province` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `city` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `district` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `detail` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `zip_code` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `tag` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `is_default` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: usr_infos
-- ------------------------------------------------------------
CREATE TABLE `usr_infos` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `gender` tinyint NOT NULL DEFAULT '0',
  `birthday` date DEFAULT NULL,
  `bio` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `country` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `province` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `city` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `zip_code` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `language` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'zh-CN',
  `timezone` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'Asia/Shanghai',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: usr_login_histories
-- ------------------------------------------------------------
CREATE TABLE `usr_login_histories` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `provider` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `ip` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `user_agent` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `device_id` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `event` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `fail_reason` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ------------------------------------------------------------
-- Table: usr_users
-- ------------------------------------------------------------
CREATE TABLE `usr_users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `password_hash` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `email` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `email_verified` tinyint(1) NOT NULL DEFAULT '0',
  `phone` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `phone_verified` tinyint(1) NOT NULL DEFAULT '0',
  `avatar` varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `nickname` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT '1',
  `register_ip` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `register_source` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `last_login_ip` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `last_login_at` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  UNIQUE KEY `uk_phone` (`phone`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

