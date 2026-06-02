import requests, time, subprocess

BASE = "http://localhost:8080"

def redis_get(key):
    r = subprocess.run(["redis-cli", "get", key], capture_output=True, text=True)
    return r.stdout.strip()

def db_query(sql):
    r = subprocess.run(["mysql", "-uroot", "-p123456", "eshop_db", "-N", "-e", sql],
                       capture_output=True, text=True, timeout=5)
    return r.stdout.strip()

start = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(time.time()))
end = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime(time.time() + 3600))

# 1. 创建活动
r = requests.post(f"{BASE}/api/v1/flash/activities", json={
    "product_id": 1, "flash_price": 99900, "total_stock": 5,
    "start_time": start, "end_time": end
})
data = r.json()["data"]
aid = data["id"]
print(f"1. 创建活动: id={aid} status={data['status']}")

# 2. 加载库存到Redis
requests.post(f"{BASE}/api/v1/flash/activities/{aid}/load-stock")
print(f"2. 加载库存: Redis={redis_get(f'flash:stock:{aid}')}")

# 3. 查看真实库存（通过库存接口）
r = requests.get(f"{BASE}/api/v1/inventories/product/1")
inv = r.json()["data"]
print(f"3. 真实库存: quantity={inv['quantity']} reserved={inv['reserved']}")

# 4. 6个用户抢5个库存
print("\n4. 抢购（6人抢5份）:")
for uid in range(1, 7):
    r = requests.post(f"{BASE}/api/v1/flash/buy", json={
        "activity_id": aid, "user_id": uid
    })
    d = r.json()["data"]
    m = "✅" if d["success"] else "❌"
    print(f"   user={uid:2d} {m} {d['message']}")

# 5. 抢购后状态
redis_left = redis_get(f"flash:stock:{aid}")
order_cnt = db_query(f"SELECT COUNT(*) FROM flash_orders WHERE activity_id={aid}")
sold = requests.get(f"{BASE}/api/v1/flash/activities/{aid}").json()["data"]["sold_stock"]
r = requests.get(f"{BASE}/api/v1/inventories/product/1")
inv2 = r.json()["data"]
print(f"\n5. 抢购后: Redis闪存={redis_left}, 已售={sold}, 订单={order_cnt}")
print(f"   真实库存: quantity={inv2['quantity']}, reserved={inv2['reserved']}")

# 6. 确认前3单
print("\n6. 确认订单:")
for uid in range(1, 4):
    oid = db_query(f"SELECT id FROM flash_orders WHERE activity_id={aid} AND user_id={uid}")
    if oid:
        r = requests.post(f"{BASE}/api/v1/flash/orders/{oid}/confirm")
        print(f"   user={uid} order={oid} -> {r.json()['data']['message']}")

# 7. 取消第4-5单
print("\n7. 取消订单:")
for uid in range(4, 6):
    oid = db_query(f"SELECT id FROM flash_orders WHERE activity_id={aid} AND user_id={uid}")
    if oid:
        r = requests.post(f"{BASE}/api/v1/flash/orders/{oid}/cancel")
        print(f"   user={uid} order={oid} -> {r.json()['data']['message']}")

# 8. 最终状态
r = requests.get(f"{BASE}/api/v1/inventories/product/1")
inv3 = r.json()["data"]
redis_left = redis_get(f"flash:stock:{aid}")
sold_now = requests.get(f"{BASE}/api/v1/flash/activities/{aid}").json()["data"]["sold_stock"]
print(f"\n8. 最终状态:")
print(f"   Redis闪存={redis_left}, 活动已售={sold_now}")
print(f"   真实库存: quantity={inv3['quantity']}, reserved={inv3['reserved']}, 可用={inv3['quantity']-inv3['reserved']}")
print(f"   说明: 3单确认→reserved-3,qty-3 | 2单取消→reserved-2,闪存+2")