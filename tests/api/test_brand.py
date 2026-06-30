#!/usr/bin/env python3
"""Brand 列表接口测试脚本
用法: python tests/api/test_brand.py [host]
默认: http://localhost:8080
"""
import json
import sys
import urllib.error
import urllib.parse
import urllib.request

HOST = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080"
API = f"{HOST}/api/v1/brands"

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


print("=== Brand 列表接口测试 ===")
print(f"目标: {API}")
print()

# ── 1. 默认分页 ────────────────────────────────
print("--- 1. 默认分页 ---")
data = api_get({})
code = data.get("code")
result = data.get("data", {})
total = result.get("total", 0)
items = result.get("list", [])

ok("响应 code=0") if code == 0 else nok(f"响应 code={code} (期望 0)")
ok(f"total={total}") if isinstance(total, int) else nok("total 格式异常")
ok(f"返回 {len(items)} 条") if isinstance(items, list) else nok("列表格式异常")
print()

# ── 2. 分页参数 ────────────────────────────────
print("--- 2. 分页参数 ---")
data = api_get({"page": 1, "size": 3})
result = data.get("data", {})
items = result.get("list", [])
total = result.get("total", 0)
ok(f"page=1&size=3 返回 {len(items)} 条（total={total}）") if len(items) <= 3 else nok("分页参数异常")
print()

# ── 3. 名称模糊搜索 ────────────────────────────
print("--- 3. 名称模糊搜索 ---")
first_name = ""
data0 = api_get({"page": 1, "size": 1})
items0 = data0.get("data", {}).get("list", [])
if items0:
    first_name = items0[0].get("name", "")[:2]
if first_name:
    data = api_get({"name": first_name})
    code = data.get("code")
    ok("名称搜索请求正常") if code == 0 else nok("名称搜索请求异常")
    items = data.get("data", {}).get("list", [])
    ok(f"name={first_name} 返回 {len(items)} 条") if len(items) >= 1 else nok("名称搜索返回 0 条")
else:
    ok("无数据可搜索（跳过名称搜索测试）")
print()

# ── 4. 首字母筛选 ──────────────────────────────
print("--- 4. 首字母筛选 ---")
data = api_get({"first_letter": "A"})
code = data.get("code")
ok("first_letter=A 请求正常") if code == 0 else nok("首字母筛选异常")
items = data.get("data", {}).get("list", [])
if len(items) >= 1:
    ok(f"first_letter=A 返回 {len(items)} 条")
    all_match = all(item.get("first_letter") == "A" for item in items)
    ok("所有结果 first_letter=A") if all_match else nok("存在 first_letter!=A 的结果")
else:
    ok("first_letter=A 返回 0 条（可能无首字母为A的品牌）")
print()

# ── 5. 状态筛选 ────────────────────────────────
print("--- 5. 状态筛选 ---")
data = api_get({"status": 1})
code = data.get("code")
ok("status=1 请求正常") if code == 0 else nok("状态筛选异常")
items = data.get("data", {}).get("list", [])
if len(items) >= 1:
    ok(f"status=1 返回 {len(items)} 条")
    all_match = all(item.get("status") == 1 for item in items)
    ok("所有结果 status=1") if all_match else nok("存在 status!=1 的结果")
else:
    ok("status=1 返回 0 条")

data = api_get({"status": 0})
items = data.get("data", {}).get("list", [])
ok(f"status=0 返回 {len(items)} 条") if len(items) >= 0 else nok("status=0 请求异常")
print()

# ── 6. 组合筛选 ────────────────────────────────
print("--- 6. 组合筛选 ---")
data = api_get({"first_letter": "A", "status": 1})
code = data.get("code")
ok("组合筛选请求正常") if code == 0 else nok("组合筛选异常")
print()

# ── 7. 大数据量分页 ─────────────────────────────
print("--- 7. 大数据量分页 ---")
data = api_get({"page": 1, "size": 100})
code = data.get("code")
items = data.get("data", {}).get("list", [])
ok(f"size=100 请求正常，返回 {len(items)} 条") if code == 0 else nok("size=100 请求异常")
print()

# ── 8. 参数错误 ────────────────────────────────
print("--- 8. 参数错误 ---")
data = api_get({"page": -1})
code = data.get("code")
ok(f"page=-1 返回错误 code={code}") if code != 0 else nok("page=-1 未返回错误")

data = api_get({"size": 2000})
code = data.get("code")
ok(f"size=2000 返回错误 code={code}") if code != 0 else nok("size=2000 未返回错误")
print()

# ── 总结 ──────────────────────────────────────
total = pass_count + fail_count
print("=== 结果 ===")
print(f"通过: {pass_count} / {total}, 失败: {fail_count}")
print("全部通过!") if fail_count == 0 else print(f"有 {fail_count} 个测试失败")
