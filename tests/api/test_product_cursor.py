#!/usr/bin/env python3
"""SPU 列表游标分页测试脚本
用法: python tests/api/test_product_cursor.py [host]
默认: http://localhost:8080
"""
import json
import sys
import urllib.error
import urllib.parse
import urllib.request

HOST = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080"
API = f"{HOST}/api/v1/products"

pass_count = 0
fail_count = 0


def api_get(params: dict) -> dict:
    qs = urllib.parse.urlencode(params)
    url = f"{API}?{qs}"
    try:
        with urllib.request.urlopen(url, timeout=10) as resp:
            return json.loads(resp.read())
    except urllib.error.HTTPError as e:
        return json.loads(e.read())


def ok(msg):
    global pass_count
    print(f"  ✅ {msg}")
    pass_count += 1


def nok(msg):
    global fail_count
    print(f"  ❌ {msg}")
    fail_count += 1


print("=== SPU Keyset 游标分页测试 ===")
print(f"目标: {API}")
print()

# ── 1. 首次请求（无 cursor） ──────────────────
print("--- 1. 首次请求（无 cursor） ---")
data = api_get({"size": 3})
code = data.get("code")
result = data.get("data", {})
items = result.get("list", [])
has_more = result.get("has_more")
cursor = result.get("cursor")

ok("响应 code=0") if code == 0 else nok(f"响应 code={code} (期望 0)")
ok(f"返回 {len(items)} 条") if len(items) >= 1 else nok("返回 0 条")
ok(f"has_more={has_more}") if has_more in (True, False) else nok(f"has_more 格式异常: {has_more}")
ok(f"cursor 非空: {cursor[:12]}...") if cursor else nok("cursor 为空")
print()

# ── 2. 翻页 ──────────────────────────────────
print("--- 2. 连续翻页 ---")
cur = cursor
page = 0
while has_more and page < 5:
    data = api_get({"size": 3, "cursor": cur})
    result = data.get("data", {})
    items = result.get("list", [])
    has_more = result.get("has_more")
    next_cursor = result.get("cursor")
    page += 1

    if len(items) >= 1:
        ok(f"翻页 #{page}: {len(items)} 条, has_more={has_more}")
        if next_cursor != cur:
            ok(f"翻页 #{page}: cursor 已更新")
        else:
            nok(f"翻页 #{page}: cursor 未变化")
        if has_more is False and not next_cursor:
            ok(f"翻页 #{page}: 最后一页无 cursor")
    else:
        nok(f"翻页 #{page}: 返回 0 条")
    cur = next_cursor

ok(f"共翻页 {page} 次") if page > 0 else ok("全部数据已在一页内")
print()

# ── 3. 带分类过滤 ─────────────────────────────
print("--- 3. 分类过滤 ---")
data = api_get({"size": 3, "category_id": 1})
result = data.get("data", {})
items = result.get("list", [])
cursor_filter = result.get("cursor")
ok(f"category_id=1 返回 {len(items)} 条")
if items:
    all_match = all(item.get("category_id") == 1 for item in items)
    ok("所有结果 category_id=1") if all_match else nok("存在 category_id!=1 的结果")
print()

# ── 4. 状态过滤 ───────────────────────────────
print("--- 4. 状态过滤 ---")
data = api_get({"size": 3, "status": 2})
result = data.get("data", {})
items = result.get("list", [])
ok(f"status=2 返回 {len(items)} 条")
if items:
    all_match = all(item.get("status") == 2 for item in items)
    ok("所有结果 status=2") if all_match else nok("存在 status!=2 的结果")
print()

# ── 5. 名称模糊搜索（走 DB keyset 降级） ──────
print("--- 5. 名称模糊搜索（走 DB 降级） ---")
first_name = ""
data0 = api_get({"size": 1})
items0 = data0.get("data", {}).get("list", [])
if items0:
    first_name = items0[0].get("name", "")[:2]
if first_name:
    data = api_get({"size": 3, "name": first_name})
    items = data.get("data", {}).get("list", [])
    ok(f"name={first_name} 搜索返回 {len(items)} 条") if len(items) >= 1 else nok("名称搜索返回 0 条")
else:
    ok("无数据可搜索（跳过名称搜索测试）")
print()

# ── 6. 价格范围查询（走 DB keyset 降级） ──────
print("--- 6. 价格范围查询（走 DB 降级） ---")
data = api_get({"size": 3, "price_min": 100000, "price_max": 500000})
result = data.get("data", {})
items = result.get("list", [])
if len(items) >= 1:
    ok(f"price_min=100000&price_max=500000 返回 {len(items)} 条")
    all_match = all(
        item.get("min_price", 0) >= 100000 and item.get("max_price", 0) <= 500000
        for item in items
    )
    ok("所有结果在价格范围内") if all_match else nok("存在超范围的结果")
else:
    ok("价格范围查询 返回 0 条（可能无符合数据）")
print()

# ── 7. 组合过滤 ───────────────────────────────
print("--- 7. 组合过滤 ---")
data = api_get({"size": 3, "category_id": 1, "status": 2})
result = data.get("data", {})
items = result.get("list", [])
ok(f"category_id=1&status=2 返回 {len(items)} 条")
if items:
    all_match = all(
        item.get("category_id") == 1 and item.get("status") == 2
        for item in items
    )
    ok("所有结果符合组合过滤") if all_match else nok("存在不符合条件的结果")
print()

# ── 8. 游标稳定性（重复请求同一游标） ──────────
print("--- 8. 游标稳定性 ---")
if cursor_filter:
    d1 = api_get({"size": 3, "cursor": cursor_filter, "category_id": 1})
    d2 = api_get({"size": 3, "cursor": cursor_filter, "category_id": 1})
    ids1 = [i["id"] for i in d1.get("data", {}).get("list", [])]
    ids2 = [i["id"] for i in d2.get("data", {}).get("list", [])]
    ok("同一游标重复请求结果一致") if ids1 == ids2 else nok("同一游标结果不一致")
print()

# ── 9. 参数错误 ───────────────────────────────
print("--- 9. 参数错误 ---")
data = api_get({"size": -1})
code = data.get("code")
ok(f"size=-1 返回错误 code={code}") if code != 0 else nok("size=-1 未返回错误")

data = api_get({"size": 2000})
code = data.get("code")
ok(f"size=2000 返回错误 code={code}") if code != 0 else nok("size=2000 未返回错误")
print()

# ── 总结 ──────────────────────────────────────
total = pass_count + fail_count
print("=== 结果 ===")
print(f"通过: {pass_count} / {total}, 失败: {fail_count}")
print("全部通过!") if fail_count == 0 else print(f"有 {fail_count} 个测试失败")
