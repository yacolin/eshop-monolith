"""
秒杀活动批量生成 —— 直连 MySQL 写入数据，Redis 预热走后端 API。

用法：
    python scripts/gen_flash_data.py
    python scripts/gen_flash_data.py --count 500     # 自定义数量
    python scripts/gen_flash_data.py --clean          # 先清空再生成
"""
import os, random, sys, argparse
from datetime import datetime, timedelta

import pymysql

random.seed(42)

# ── 默认值（与 configs/config.yaml 一致）──
MYSQL_CFG = {
    "host": os.getenv("DB_HOST", "localhost"),
    "port": int(os.getenv("DB_PORT", 3306)),
    "user": os.getenv("DB_USER", "root"),
    "password": os.getenv("DB_PASSWORD", "123456"),
    "database": os.getenv("DB_NAME", "eshop_db"),
    "charset": "utf8mb4",
    "cursorclass": pymysql.cursors.DictCursor,
}
API_BASE = os.getenv("API_BASE", "http://localhost:8080/api/v1")

PRICES = [500, 999, 1990, 2990, 4990, 9990, 19900, 29900, 49900, 99900]
STOCKS = [10, 20, 30, 50, 100]
FMT = "%Y-%m-%d %H:%M:%S"


def connect_db():
    try:
        conn = pymysql.connect(**MYSQL_CFG)
        print("MySQL connected")
        return conn
    except Exception as e:
        print(f"MySQL 连接失败: {e}")
        sys.exit(1)


def clean(conn):
    """清空秒杀相关数据"""
    with conn.cursor() as cur:
        cur.execute("SET FOREIGN_KEY_CHECKS = 0")
        cur.execute("TRUNCATE TABLE flash_orders")
        cur.execute("TRUNCATE TABLE flash_activities")
        cur.execute("SET FOREIGN_KEY_CHECKS = 1")
    conn.commit()
    print("已清空 flash_activities / flash_orders\n")


def generate(conn, count=100):
    now = datetime.now()
    now_str = now.strftime(FMT)

    active_count = 0
    with conn.cursor() as cur:
        for pid in range(1, count + 1):
            price = random.choice(PRICES)
            stock = random.choice(STOCKS)

            if pid <= max(10, count // 10):
                start_time = now + timedelta(minutes=random.randint(-120, -5))
                end_time = now + timedelta(minutes=random.randint(30, 240))
            else:
                start_time = now + timedelta(minutes=random.randint(60, 1440))
                end_time = start_time + timedelta(minutes=random.randint(60, 720))

            start_str = start_time.strftime(FMT)
            end_str = end_time.strftime(FMT)
            status = "active" if start_time <= now < end_time else "pending"

            cur.execute(
                "INSERT INTO flash_activities (product_id, flash_price, total_stock, sold_stock, "
                "start_time, end_time, status, created_at, updated_at) "
                "VALUES (%s, %s, %s, 0, %s, %s, %s, %s, %s)",
                (pid, price, stock, start_str, end_str, status, now_str, now_str),
            )
            activity_id = cur.lastrowid

            if status == "active":
                active_count += 1

            if pid % max(1, count // 10) == 0:
                print(f"  已创建 {pid}/{count} (id={activity_id}, status={status})")

    conn.commit()
    print(f"\n=== 完成 === 总计 {count} | 进行中 {active_count} | 未来 {count - active_count}")


def main():
    parser = argparse.ArgumentParser(description="秒杀活动批量生成")
    parser.add_argument("--count", type=int, default=100, help="生成数量 (默认 100)")
    parser.add_argument("--clean", action="store_true", help="先清空再生成")
    args = parser.parse_args()

    conn = connect_db()

    if args.clean:
        clean(conn)

    generate(conn, args.count)

    with conn.cursor() as cur:
        cur.execute("SELECT COUNT(*) AS cnt FROM flash_activities")
        total = cur.fetchone()["cnt"]
    print(f"\nflash_activities 表总数: {total}")

    conn.close()


if __name__ == "__main__":
    main()
