"""秒杀活动批量创建 100 个 + Redis 缓存验证"""
import requests, time, random
from datetime import datetime, timedelta

BASE = "http://localhost:8080"

try:
    requests.get(f"{BASE}/api/v1/products", timeout=3)
except:
    print("服务未启动，请先运行 go run ./cmd/server")
    exit(1)

fmt = "%Y-%m-%d %H:%M:%S"
now = datetime.now()
random.seed(42)

activities = []
prices = [500, 999, 1990, 2990, 4990, 9990, 19900, 29900, 49900, 99900]
stocks = [10, 20, 30, 50, 100]

# ── 100 个活动，每个不同 product_id ──
for pid in range(1, 101):
    price = random.choice(prices)
    stock = random.choice(stocks)
    # 前 10 个为进行中，其余为未来活动
    if pid <= 10:
        start_offset = random.randint(-120, -5)          # 2h前 ~ 5min前
        end_offset   = random.randint(30, 240)            # 30min后 ~ 4h后
    else:
        start_offset = random.randint(60, 1440)           # 1h后 ~ 24h后
        end_offset   = start_offset + random.randint(60, 720)

    r = requests.post(f"{BASE}/api/v1/flash/activities", json={
        "product_id": pid,
        "flash_price": price,
        "total_stock": stock,
        "start_time": (now + timedelta(minutes=start_offset)).strftime(fmt),
        "end_time":   (now + timedelta(minutes=end_offset)).strftime(fmt),
    })
    a = r.json()["data"]
    activities.append(a)
    if pid % 10 == 0:
        print(f"  已创建 {pid}/100 (id={a['id']})")

# ── 加载库存到 Redis ──
print(f"\n=== 加载 Redis 闪存库存（{len(activities)} 个）===")
for a in activities:
    r = requests.post(f"{BASE}/api/v1/flash/activities/{a['id']}/load-stock")
    if not r.ok:
        print(f"  活动 {a['id']} FAIL")

# ── 状态汇总 ──
active = [a for a in activities if a['status'] == 'active']
pending = [a for a in activities if a['status'] == 'pending']
print(f"\n=== 活动状态汇总 ===")
print(f"  总计 {len(activities)} 个 | 进行中 {len(active)} | 未来 {len(pending)}")
print(f"\n=== 游标分页测试 ===")
r = requests.get(f"{BASE}/api/v1/flash/activities-cursor?cursor=0&size=10")
data = r.json()["data"]
print(f"  第1页: {len(data['list'])} 条, next_cursor={data['next_cursor']}, has_more={data['has_more']}")
if data['has_more']:
    nc = data['next_cursor']
    r2 = requests.get(f"{BASE}/api/v1/flash/activities-cursor?cursor={nc}&size=10")
    d2 = r2.json()["data"]
    print(f"  第2页: {len(d2['list'])} 条, next_cursor={d2['next_cursor']}")

# ── 按状态筛选测试 ──
r = requests.get(f"{BASE}/api/v1/flash/activities-cursor?cursor=0&size=100&status=active")
data = r.json()["data"]
print(f"\n  筛选 status=active: {len(data['list'])} 条")
