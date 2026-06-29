"""
清理被新表（sp_ / tx_ / mkt_）替代的旧业务表。

用法：
    python scripts/clean_old_tables.py --dry-run    # 只预览，不删除
    python scripts/clean_old_tables.py --execute    # 执行删除
    python scripts/clean_old_tables.py --list       # 分类列出所有表
"""
import os, sys, argparse

import pymysql

MYSQL_CFG = {
    "host": os.getenv("DB_HOST", "localhost"),
    "port": int(os.getenv("DB_PORT", 3306)),
    "user": os.getenv("DB_USER", "root"),
    "password": os.getenv("DB_PASSWORD", "123456"),
    "database": os.getenv("DB_NAME", "eshop_db"),
    "charset": "utf8mb4",
}

# ── 表分类 ────────────────────────────────────────

# 已确认可删除的旧表（已被新表完整替代）
REPLACED = {
    "products": "→ sp_products",
    "categories": "→ sp_categories",
    "category_attributes": "→（类目-属性关联已被 product center 替代）",
    "product_categories": "→（产品-类目关联已被 sp_products.category_id 替代）",
    "skus": "→ sp_skus",
    "sku_attributes": "→ sp_attributes.is_sku_spec",
    "product_attribute_values": "→ sp_product_attributes",
    "attribute_attributes": "→ sp_attributes",
    "attribute_values": "→ sp_attributes.values",
    "inventories": "→ sp_inventories",
    # 交易中心
    "carts": "→ tx_carts",
    "cart_items": "→ tx_cart_items",
    "orders": "→ tx_orders",
    "order_items": "→ tx_order_items",
    "payments": "→ tx_payments",
    "payment_methods": "→（支付方式合并到 tx_payments.payment_method）",
    "payment_transactions": "→ tx_payment_logs",
    "refunds": "→ tx_refunds",
    # 营销中心
    "coupons": "→ mkt_promotions（promo_type=1,2）",
    "user_coupons": "→ mkt_user_promotions",
    "promotions": "→ mkt_promotions（promo_type=4,5,6）",
    "promotion_products": "→ mkt_promotion_products",
    "flash_activities": "→ mkt_promotions（promo_type=3）",
    "flash_orders": "→（秒杀订单合并到 tx_orders）",
}

# 仍被旧模块使用的表（暂不删除）
KEPT = [
    "users", "user_infos", "user_identities", "user_roles",
    "roles", "permissions", "role_permissions", "auth_tokens",
    "login_histories", "addresses", "notifications",
    "reviews", "product_rating_summaries",
]


def connect():
    try:
        conn = pymysql.connect(**MYSQL_CFG)
        return conn
    except Exception as e:
        print(f"MySQL 连接失败: {e}")
        sys.exit(1)


def list_tables(conn):
    with conn.cursor() as cur:
        cur.execute("SHOW TABLES")
        all_tables = {row[0] for row in cur.fetchall()}

    print(f"\n{'='*60}")
    print(f"数据库: {MYSQL_CFG['database']}  共 {len(all_tables)} 张表")
    print(f"{'='*60}\n")

    print(f"■ 已确认可删除（{len(REPLACED)} 张）：")
    for t, note in sorted(REPLACED.items()):
        exists = "✓" if t in all_tables else " "
        print(f"   [{exists}] {t:30s} {note}")
    print()

    print(f"■ 保留（仍被旧模块使用，{len(KEPT)} 张）：")
    for t in KEPT:
        exists = "✓" if t in all_tables else " "
        print(f"   [{exists}] {t}")

    other = [t for t in all_tables if t not in REPLACED and t not in KEPT
             and not t.startswith(("sp_", "tx_", "mkt_"))]
    if other:
        print(f"\n■ 未分类（{len(other)} 张）：")
        for t in sorted(other):
            print(f"   [×] {t}")


def dry_run(conn):
    print(f"\n{'='*60}")
    print("预览模式 --dry-run，以下表将被删除：")
    print(f"{'='*60}\n")
    with conn.cursor() as cur:
        for t in sorted(REPLACED.keys()):
            cur.execute(f"SELECT COUNT(*) FROM `{t}`")
            count = cur.fetchone()[0]
            print(f"   DROP TABLE `{t}`  (当前 {count} 行数据)  {REPLACED[t]}")
    print(f"\n共 {len(REPLACED)} 张表，确认后请使用 --execute 执行\n")


def execute(conn):
    print(f"\n{'='*60}")
    print("开始删除旧业务表...")
    print(f"{'='*60}\n")
    with conn.cursor() as cur:
        cur.execute("SET FOREIGN_KEY_CHECKS = 0")
        for t in sorted(REPLACED.keys()):
            try:
                cur.execute(f"DROP TABLE IF EXISTS `{t}`")
                print(f"   ✓ DROP TABLE `{t}`")
            except Exception as e:
                print(f"   ✗ DROP TABLE `{t}`: {e}")
        cur.execute("SET FOREIGN_KEY_CHECKS = 1")
    conn.commit()
    print(f"\n完成！已删除 {len(REPLACED)} 张旧表\n")


def main():
    parser = argparse.ArgumentParser(description="清理旧业务数据表")
    parser.add_argument("--dry-run", action="store_true", help="只预览不删除")
    parser.add_argument("--execute", action="store_true", help="执行删除")
    parser.add_argument("--list", action="store_true", help="列出所有表分类")
    args = parser.parse_args()

    if not any(vars(args).values()):
        parser.print_help()
        print("\n请指定操作模式：--dry-run / --execute / --list")
        return

    conn = connect()
    if args.list:
        list_tables(conn)
    elif args.dry_run:
        dry_run(conn)
    elif args.execute:
        confirm = input("确认删除 23 张旧表？(yes/no): ")
        if confirm == "yes":
            execute(conn)
        else:
            print("取消操作")
    conn.close()


if __name__ == "__main__":
    main()
