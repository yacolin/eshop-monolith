"""
为新表（sp_ / tx_ / mkt_）批量生成测试数据。

用法：
    python scripts/seed_new_tables.py                  # 生成全部
    python scripts/seed_new_tables.py --clean           # 先清空再生成
    python scripts/seed_new_tables.py --module product  # 只生成商品域
"""
import os, random, sys, argparse
from datetime import datetime, timedelta

import pymysql

random.seed(42)
FMT = "%Y-%m-%d %H:%M:%S"

MYSQL_CFG = {
    "host": os.getenv("DB_HOST", "localhost"),
    "port": int(os.getenv("DB_PORT", 3306)),
    "user": os.getenv("DB_USER", "root"),
    "password": os.getenv("DB_PASSWORD", "123456"),
    "database": os.getenv("DB_NAME", "eshop_db"),
    "charset": "utf8mb4",
}

# ── 测试数据 ──────────────────────────────────────

BRANDS = [
    ("Apple", "苹果", "A"), ("Samsung", "三星", "S"), ("Xiaomi", "小米", "X"),
    ("Huawei", "华为", "H"), ("OPPO", "OPPO", "O"), ("Sony", "索尼", "S"),
    ("Nike", "耐克", "N"), ("Adidas", "阿迪达斯", "A"), ("Uniqlo", "优衣库", "U"),
]

CATEGORIES = [
    (0, "电子产品", 1), (0, "服装", 1), (0, "食品", 1),
    (1, "手机", 2), (1, "电脑", 2), (1, "耳机", 2),
    (2, "男装", 2), (2, "女装", 2),
    (4, "智能手机", 3), (4, "功能机", 3), (5, "笔记本", 3), (5, "平板", 3),
]

ATTRS = [
    (4, "颜色", 2, '["黑色","白色","红色","蓝色"]', 1, 1),
    (4, "内存", 4, '["128G","256G","512G"]', 1, 1),
    (4, "屏幕尺寸", 4, '["6.1","6.7"]', 0, 0),
    (5, "处理器", 1, None, 0, 0),
    (5, "内存", 4, '["8G","16G","32G"]', 1, 1),
    (5, "硬盘", 4, '["256G","512G","1T"]', 1, 1),
]

PHONES = [
    ("iPhone 16 Pro", "钛金属旗舰", 4, 1, 799900, 899900),
    ("Galaxy S25", "AI 智能旗舰", 4, 2, 699900, 799900),
    ("小米 15 Pro", "徕卡影像", 4, 3, 499900, 599900),
    ("Mate 70 Pro", "鸿蒙旗舰", 4, 4, 699900, 799900),
]

LAPTOPS = [
    ("MacBook Pro 14", "M4 芯片", 11, 1, 1299900, 1499900),
    ("ThinkPad X1", "商务旗舰", 11, 3, 999900, 1199900),
]

SKU_SPECS = [
    ('{"颜色":"黑色","内存":"256G"}', 1),
    ('{"颜色":"白色","内存":"256G"}', 1),
    ('{"颜色":"黑色","内存":"512G"}', 2),
    ('{"颜色":"白色","内存":"512G"}', 2),
]


def connect():
    try:
        conn = pymysql.connect(**MYSQL_CFG)
        print("MySQL connected")
        return conn
    except Exception as e:
        print(f"MySQL 连接失败: {e}")
        sys.exit(1)


def clean(conn):
    tables = [
        "mkt_promotion_usage_logs", "mkt_user_promotions", "mkt_promotion_products",
        "mkt_promotion_rules", "mkt_promotions",
        "tx_refunds", "tx_payment_logs", "tx_payments",
        "tx_order_logs", "tx_order_items", "tx_orders",
        "tx_cart_items", "tx_carts",
        "sp_inventory_logs", "sp_inventories",
        "sp_product_attributes", "sp_product_descriptions", "sp_skus",
        "sp_products", "sp_attributes", "sp_category_brands", "sp_categories", "sp_brands",
    ]
    with conn.cursor() as cur:
        cur.execute("SET FOREIGN_KEY_CHECKS = 0")
        for t in tables:
            cur.execute(f"TRUNCATE TABLE {t}")
        cur.execute("SET FOREIGN_KEY_CHECKS = 1")
    conn.commit()
    print("已清空所有新表\n")


# ── 商品中心 ──────────────────────────────────────

def seed_product(conn):
    now = datetime.now().strftime(FMT)
    with conn.cursor() as cur:
        # 品牌
        for name, cname, letter in BRANDS:
            cur.execute(
                "INSERT INTO sp_brands (name, english_name, first_letter, sort_order, status, created_at) "
                "VALUES (%s, %s, %s, %s, 1, %s)",
                (cname, name, letter, random.randint(1, 100), now),
            )
        print(f"  品牌: {len(BRANDS)}")

        # 类目
        cat_ids = {}
        for i, (parent_id, name, level) in enumerate(CATEGORIES, 1):
            path = f"{parent_id}/" if parent_id else ""
            cur.execute(
                "INSERT INTO sp_categories (name, parent_id, level, path, sort_order, status, created_at) "
                "VALUES (%s, %s, %s, %s, %s, 1, %s)",
                (name, parent_id, level, path, i * 10, now),
            )
            cat_ids[i] = cur.lastrowid
        print(f"  类目: {len(CATEGORIES)}")

        # 属性
        for cat_idx, name, input_type, values, is_sku, searchable in ATTRS:
            cur.execute(
                "INSERT INTO sp_attributes (name, category_id, input_type, `values`, is_sku_spec, searchable, status, created_at) "
                "VALUES (%s, %s, %s, %s, %s, %s, 1, %s)",
                (name, cat_ids[cat_idx], input_type, values, is_sku, searchable, now),
            )
        print(f"  属性: {len(ATTRS)}")

        # SPU
        all_products = PHONES + LAPTOPS
        for name, subtitle, cat_idx, brand_idx, price, market in all_products:
            cur.execute(
                "INSERT INTO sp_products (name, subtitle, category_id, brand_id, unit, main_image, "
                "min_price, max_price, status, sort_order, created_at, updated_at) "
                "VALUES (%s, %s, %s, %s, '件', '', %s, %s, 2, %s, %s, %s)",
                (name, subtitle, cat_ids[cat_idx], brand_idx, price, market,
                 random.randint(1, 100), now, now),
            )
            spu_id = cur.lastrowid

            # SKU（每个 SPU 2-4 个规格）
            sku_count = random.randint(2, 4)
            for j in range(sku_count):
                colors = ["黑色", "白色", "蓝色", "红色"]
                storages = ["128G", "256G", "512G"]
                spec = f'{{"颜色":"{random.choice(colors)}","内存":"{random.choice(storages)}"}}'
                sku_price = price + random.randint(-5000, 5000)
                cur.execute(
                    "INSERT INTO sp_skus (product_id, sku_code, barcode, spec, price, market_price, cost_price, status, created_at) "
                    "VALUES (%s, %s, %s, %s, %s, %s, %s, 1, %s)",
                    (spu_id, f"SKU{spu_id}-{j+1:03d}", f"BAR{spu_id}-{j+1:03d}", spec, sku_price,
                     sku_price + random.randint(5000, 20000), int(sku_price * 0.6), now),
                )
        print(f"  SPU: {len(all_products)}, SKU: generated")

    conn.commit()
    print("商品中心 ✅\n")


# ── 库存中心 ──────────────────────────────────────

def seed_inventory(conn):
    now = datetime.now().strftime(FMT)
    with conn.cursor() as cur:
        cur.execute("SELECT id FROM sp_skus WHERE deleted_at IS NULL")
        skus = cur.fetchall()
        for sku in skus:
            qty = random.randint(50, 500)
            cur.execute(
                "INSERT INTO sp_inventories (sku_id, quantity, reserved, threshold, status, created_at) "
                "VALUES (%s, %s, 0, 10, 'instock', %s)",
                (sku[0], qty, now),
            )
    conn.commit()
    print(f"  库存: {len(skus)} 条")
    print("库存中心 ✅\n")


# ── 营销中心 ──────────────────────────────────────

def seed_marketing(conn):
    now = datetime.now()
    now_str = now.strftime(FMT)
    with conn.cursor() as cur:
        for i, (name, ptype, condition, benefit) in enumerate([
            ("满200减30", 4, 20000, 3000),
            ("全场8折", 5, 0, 20),
            ("新用户满减券", 1, 0, 5000),
            ("会员9折", 6, 0, 10),
            ("限时秒杀", 3, 0, 50),
        ], 1):
            start = (now - timedelta(days=1)).strftime(FMT)
            end = (now + timedelta(days=7)).strftime(FMT)
            cur.execute(
                "INSERT INTO mkt_promotions (promo_name, promo_type, promo_code, start_time, end_time, "
                "total_quantity, per_user_limit, used_quantity, status, created_at) "
                "VALUES (%s, %s, %s, %s, %s, %s, %s, 0, 2, %s)",
                (name, ptype, f"PROMO{i:04d}", start, end, 1000, 1, now_str),
            )
            promo_id = cur.lastrowid

            cur.execute(
                "INSERT INTO mkt_promotion_rules (promotion_id, rule_name, condition_type, condition_value, "
                "benefit_type, benefit_value, created_at) VALUES (%s, %s, %s, %s, %s, %s, %s)",
                (promo_id, f"{name}规则", 2 if condition > 0 else 1, condition / 100,
                 1, benefit / 100, now_str),
            )

            # 秒杀配置商品
            if ptype == 3:
                cur.execute("SELECT id FROM sp_products LIMIT 5")
                for p in cur.fetchall():
                    cur.execute(
                        "INSERT INTO mkt_promotion_products (promotion_id, product_type, product_id, created_at) "
                        "VALUES (%s, 3, %s, %s)", (promo_id, p[0], now_str),
                    )

    conn.commit()
    print("营销中心 ✅\n")


def main():
    parser = argparse.ArgumentParser(description="为新表生成测试数据")
    parser.add_argument("--clean", action="store_true", help="先清空再生成")
    parser.add_argument("--module", choices=["product", "inventory", "trade", "marketing"],
                        help="只生成指定模块")
    args = parser.parse_args()

    conn = connect()
    if args.clean:
        clean(conn)

    modules = {
        "product": seed_product,
        "inventory": seed_inventory,
        "marketing": seed_marketing,
    }

    if args.module:
        modules[args.module](conn)
    else:
        for name, fn in modules.items():
            print(f"正在生成: {name}")
            fn(conn)

    # 统计
    with conn.cursor() as cur:
        for table in ["sp_brands", "sp_categories", "sp_attributes", "sp_products",
                      "sp_skus", "sp_inventories", "mkt_promotions"]:
            cur.execute(f"SELECT COUNT(*) AS cnt FROM {table}")
            row = cur.fetchone()
            print(f"  {table}: {row[0]}")

    conn.close()
    print("\n完成!")


if __name__ == "__main__":
    main()
