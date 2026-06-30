#!/usr/bin/env python3
"""Category 列表接口测试脚本
用法: python tests/api/test_categories.py [host]
默认: http://localhost:8080
"""
import json
import sys
import urllib.error
import urllib.parse
import urllib.request

HOST = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:8080"
BASE_API = f"{HOST}/api/v1/categories"

pass_count = 0
fail_count = 0


def api_get(endpoint: str = "", params: dict = None) -> dict:
    """发送 GET 请求到分类 API"""
    if params is None:
        params = {}
    qs = urllib.parse.urlencode(params)
    url = f"{BASE_API}{endpoint}?{qs}" if qs else f"{BASE_API}{endpoint}"
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


print("=== Category 列表接口测试 ===")
print(f"目标: {BASE_API}")
print()

# ── 1. 默认列表 ────────────────────────────────
print("--- 1. 默认列表 (/) ---")
data = api_get("")
code = data.get("code")
result = data.get("data", {})
items = result.get("list", []) if isinstance(result, dict) else []
total = result.get("total", 0) if isinstance(result, dict) else 0

ok("响应 code=0") if code == 0 else nok(f"响应 code={code} (期望 0)")
ok(f"total={total}") if isinstance(total, int) else nok("total 格式异常")
ok(f"返回 {len(items)} 条") if isinstance(items, list) else nok("列表格式异常")
print()

# ── 2. 根分类 (root) ──────────────────────────
print("--- 2. 根分类 (/root) ---")
data = api_get("/root")
code = data.get("code")
items = data.get("data", [])
ok("根分类请求正常") if code == 0 else nok(f"根分类请求异常 code={code}")
if isinstance(items, list):
    ok(f"返回 {len(items)} 个根分类")
else:
    nok("data 格式异常")
print()

# ── 3. 全部分类 (all) ─────────────────────────
print("--- 3. 全部分类 (/all) ---")
data = api_get("/all")
code = data.get("code")
items = data.get("data", [])
ok("全部分类请求正常") if code == 0 else nok(f"全部分类请求异常 code={code}")
if isinstance(items, list):
    all_total = len(items)
    ok(f"返回 {all_total} 个分类")
else:
    nok("data 格式异常")
print()

# ── 4. 按层级查询 (level) ──────────────────────
print("--- 4. 按层级查询 (/level/:level) ---")
data = api_get("/level/1")
code = data.get("code")
items = data.get("data", [])
ok("level=1 请求正常") if code == 0 else nok(f"level=1 请求异常 code={code}")
if isinstance(items, list):
    ok(f"level=1 返回 {len(items)} 条")
else:
    nok("data 格式异常")

data = api_get("/level/2")
code = data.get("code")
ok("level=2 请求正常") if code == 0 else nok(f"level=2 请求异常 code={code}")

data = api_get("/level/3")
code = data.get("code")
ok("level=3 请求正常") if code == 0 else nok(f"level=3 请求异常 code={code}")
print()

# ── 5. 获取子分类 (children) ────────────────────
print("--- 5. 获取子分类 (/:id/children) ---")
root_data = api_get("/root")
root_items = root_data.get("data", [])
if root_items and len(root_items) > 0:
    first_id = root_items[0].get("id")
    if first_id:
        data = api_get(f"/{first_id}/children")
        code = data.get("code")
        ok(f"获取 id={first_id} 的子分类正常") if code == 0 else nok(f"获取子分类异常 code={code}")
        items = data.get("data", [])
        if isinstance(items, list):
            ok(f"返回 {len(items)} 个子分类")
        else:
            nok("data 格式异常")
    else:
        nok("根分类缺少 id 字段")
else:
    ok("无根分类，跳过子分类测试")
print()

# ── 6. 按 ID 查询详情 ──────────────────────────
print("--- 6. 按 ID 查询详情 (/:id) ---")
all_data = api_get("/all")
all_items = all_data.get("data", [])
if all_items and len(all_items) > 0:
    first_id = all_items[0].get("id")
    if first_id:
        data = api_get(f"/{first_id}")
        code = data.get("code")
        ok(f"获取 id={first_id} 详情正常") if code == 0 else nok(f"获取详情异常 code={code}")
        result = data.get("data", {})
        if isinstance(result, dict) and result.get("id") == first_id:
            ok("返回数据 ID 匹配")
        else:
            nok("返回数据 ID 不匹配")
    else:
        nok("分类缺少 id 字段")
else:
    ok("无分类数据，跳过详情查询测试")
print()

# ── 7. 分页参数 ────────────────────────────────
print("--- 7. 分页参数 ---")
data = api_get("", {"page": 1, "size": 3})
result = data.get("data", {})
items = result.get("list", []) if isinstance(result, dict) else []
total = result.get("total", 0) if isinstance(result, dict) else 0
ok(f"page=1&size=3 返回 {len(items)} 条（total={total}）") if len(items) <= 3 else nok("分页参数异常")
print()

# ── 8. 名称搜索 ────────────────────────────────
print("--- 8. 名称搜索 ---")
first_name = ""
data0 = api_get("", {"page": 1, "size": 1})
data0_result = data0.get("data", {})
items0 = data0_result.get("list", []) if isinstance(data0_result, dict) else []
if items0 and len(items0) > 0:
    first_name = items0[0].get("name", "")[:2]
if first_name:
    data = api_get("", {"name": first_name})
    code = data.get("code")
    ok("名称搜索请求正常") if code == 0 else nok(f"名称搜索请求异常 code={code}")
    result = data.get("data", {})
    items = result.get("list", []) if isinstance(result, dict) else []
    ok(f"name={first_name} 返回 {len(items)} 条") if len(items) >= 1 else nok("名称搜索返回 0 条")
else:
    ok("无数据可搜索（跳过名称搜索测试）")
print()

# ── 9. 参数错误 ────────────────────────────────
print("--- 9. 参数错误 ---")
data = api_get("", {"page": -1})
code = data.get("code")
ok(f"page=-1 返回错误 code={code}") if code != 0 else nok("page=-1 未返回错误")

data = api_get("", {"size": 2000})
code = data.get("code")
ok(f"size=2000 返回错误 code={code}") if code != 0 else nok("size=2000 未返回错误")
print()

# ── 10. 无效层级（边界值）──────────────────────
print("--- 10. 无效层级（边界值）---")
data = api_get("/level/0")
code = data.get("code")
ok("level=0 返回错误 code=4015") if code == 4015 else nok(f"level=0 期望 code=4015 实际 code={code}")

data = api_get("/level/4")
code = data.get("code")
ok("level=4 返回错误 code=4015") if code == 4015 else nok(f"level=4 期望 code=4015 实际 code={code}")

data = api_get("/level/-1")
code = data.get("code")
ok("level=-1 返回错误 code=4015") if code == 4015 else nok(f"level=-1 期望 code=4015 实际 code={code}")

data = api_get("/level/999")
code = data.get("code")
ok("level=999 返回错误 code=4015") if code == 4015 else nok(f"level=999 期望 code=4015 实际 code={code}")
print()

# ── 11. 无效 ID ────────────────────────────────
print("--- 11. 无效 ID ---")
data = api_get("/999999")
code = data.get("code")
ok("id=999999 返回错误 code=4010") if code == 4010 else nok(f"id=999999 期望 code=4010 实际 code={code}")
print()

# ── 12. status 过滤 ────────────────────────────
print("--- 12. status 过滤 ---")
data = api_get("", {"status": 1})
code = data.get("code")
ok("status=1 过滤请求正常") if code == 0 else nok(f"status=1 过滤请求异常 code={code}")

data = api_get("", {"status": 2})
code = data.get("code")
ok("status=2 正常返回（无匹配）") if code == 0 else nok(f"status=2 期望 code=0 实际 code={code}")
print()

# ── 总结 ──────────────────────────────────────
total = pass_count + fail_count
print("=== 结果 ===")
print(f"通过: {pass_count} / {total}, 失败: {fail_count}")
print("全部通过!") if fail_count == 0 else print(f"有 {fail_count} 个测试失败")
