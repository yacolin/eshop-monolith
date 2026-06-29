#!/usr/bin/env python3
"""
RBAC 数据初始化脚本
合并所有 SQL 文件为一次执行，避免多次输入密码

Usage: python scripts/seed_rbac.py [--user root] [--password] [--database eshop_db]
"""
import subprocess
import sys
import os
import argparse
import tempfile

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))

FILES = [
    "seed_rbac_01_clean.sql",
    "seed_rbac_02_permissions.sql",
    "seed_rbac_03_roles.sql",
    "seed_rbac_04_assignments.sql",
    "seed_rbac_05_users.sql",
]


def combine_sql() -> str:
    lines: list[str] = []
    for fname in FILES:
        path = os.path.join(SCRIPT_DIR, fname)
        with open(path, "r", encoding="utf-8") as f:
            content = f.read().strip()
        lines.append(content)
        if not content.endswith(";"):
            lines.append("")
    return "\n\n".join(lines)


def main():
    parser = argparse.ArgumentParser(description="Initialize RBAC data")
    parser.add_argument("--user", default="root", help="MySQL user")
    parser.add_argument("--password", "-p", action="store_true", help="Prompt for password")
    parser.add_argument("--database", default="eshop_db", help="Database name")
    parser.add_argument("--host", default="127.0.0.1", help="MySQL host")
    parser.add_argument("--port", default=3306, type=int, help="MySQL port")
    args = parser.parse_args()

    print("=" * 40)
    print(f" RBAC seed script")
    print(f" Database: {args.database}@{args.host}:{args.port}")
    print("=" * 40)
    print()

    sql = combine_sql()

    with tempfile.NamedTemporaryFile(mode="wb", suffix=".sql", delete=False) as tmp:
        tmp.write(sql.encode("utf-8"))
        tmp_path = tmp.name

    try:
        cmd = ["mysql", f"--user={args.user}"]
        if args.password:
            cmd.append("-p")
        cmd += [f"--host={args.host}", f"--port={args.port}",
                "--default-character-set=utf8mb4", args.database]

        print(f"Running: mysql -u {args.user} {'-p ' if args.password else ''}--host={args.host} ... < (combined)")
        print()

        with open(tmp_path, "rb") as f:
            result = subprocess.run(cmd, stdin=f)

        if result.returncode == 0:
            print("[SUCCESS] RBAC data initialized.")
        else:
            print("[FAILED] Check your MySQL connection and try again.")
            sys.exit(1)
    finally:
        os.unlink(tmp_path)


if __name__ == "__main__":
    main()
