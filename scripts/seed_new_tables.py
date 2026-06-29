"""
为新表（sp_ / tx_ / mkt_）批量生成测试数据。

用法：
    python scripts/seed_new_tables.py                  # 生成全部
    python scripts/seed_new_tables.py --clean           # 先清空再生成
    python scripts/seed_new_tables.py --module product  # 只生成商品域
"""
import os
import random
import sys
import argparse
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
    ("Huawei", "华为", "H"), ("OPPO", "OPPO", "O"), ("Vivo", "vivo", "V"),
    ("Sony", "索尼", "S"), ("Lenovo", "联想", "L"), ("Dell", "戴尔", "D"),
    ("HP", "惠普", "H"), ("Asus", "华硕", "A"), ("Nike", "耐克", "N"),
    ("Adidas", "阿迪达斯", "A"), ("Uniqlo", "优衣库", "U"), ("Zara", "飒拉", "Z"),
    ("H&M", "海恩斯莫里斯", "H"), ("Lululemon", "露露乐蒙", "L"), ("New Balance", "新百伦", "N"),
    ("Converse", "匡威", "C"), ("Vans", "范斯", "V"), ("Supreme", "至尊", "S"),
    ("Starbucks", "星巴克", "S"), ("Coca-Cola", "可口可乐", "C"), ("Pepsi", "百事", "P"),
    ("Nestle", "雀巢", "N"), ("L'Oreal", "欧莱雅", "L"), ("Estee Lauder", "雅诗兰黛", "E"),
    ("Dior", "迪奥", "D"), ("Chanel", "香奈儿", "C"), ("Gucci", "古驰", "G"),
    ("Prada", "普拉达", "P"), ("Louis Vuitton", "路易威登", "L"),
]

CATEGORIES = [
    # 一级类目
    (0, "电子产品", 1), (0, "服装鞋帽", 1), (0, "食品饮料", 1),
    (0, "美妆护肤", 1), (0, "家居生活", 1), (0, "运动户外", 1),
    (0, "图书文具", 1), (0, "汽车用品", 1), (0, "母婴玩具", 1),
    (0, "宠物用品", 1),
    # 二级类目：电子产品
    (1, "手机通讯", 2), (1, "电脑办公", 2), (1, "数码配件", 2),
    (1, "智能设备", 2), (1, "影音娱乐", 2),
    # 二级类目：服装鞋帽
    (2, "男装", 2), (2, "女装", 2), (2, "童装", 2),
    (2, "运动鞋", 2), (2, "休闲鞋", 2), (2, "箱包皮具", 2),
    # 二级类目：食品饮料
    (3, "休闲零食", 2), (3, "粮油调味", 2), (3, "饮料冲调", 2),
    (3, "进口食品", 2), (3, "生鲜果蔬", 2),
    # 二级类目：美妆护肤
    (4, "面部护肤", 2), (4, "彩妆", 2), (4, "香水", 2),
    (4, "美发造型", 2), (4, "身体护理", 2),
    # 二级类目：运动户外
    (6, "运动器材", 2), (6, "户外装备", 2), (6, "骑行用品", 2),
    # 三级类目：手机通讯
    (11, "智能手机", 3), (11, "功能机", 3), (11, "对讲机", 3),
    # 三级类目：电脑办公
    (12, "笔记本", 3), (12, "平板电脑", 3), (12, "台式机", 3),
    (12, "显示器", 3), (12, "打印机", 3),
    # 三级类目：数码配件
    (13, "手机壳", 3), (13, "充电器", 3), (13, "数据线", 3),
    (13, "耳机", 3), (13, "移动电源", 3),
    # 三级类目：男装
    (16, "T恤", 3), (16, "衬衫", 3), (16, "外套", 3),
    (16, "牛仔裤", 3), (16, "西服", 3),
    # 三级类目：女装
    (17, "连衣裙", 3), (17, "上衣", 3), (17, "半身裙", 3),
    (17, "外套", 3), (17, "毛衣", 3),
    # 三级类目：运动鞋
    (19, "跑步鞋", 3), (19, "篮球鞋", 3), (19, "休闲鞋", 3),
    # 三级类目：休闲零食
    (23, "膨化食品", 3), (23, "糖果巧克力", 3), (23, "饼干糕点", 3),
    # 三级类目：面部护肤
    (28, "洁面", 3), (28, "爽肤水", 3), (28, "精华", 3),
    (28, "面霜", 3), (28, "防晒", 3),
    # 三级类目：彩妆
    (29, "粉底", 3), (29, "口红", 3), (29, "眼影", 3),
    (29, "腮红", 3), (29, "睫毛膏", 3),
]

# 属性定义：(类目索引, 属性名, 输入类型(1文本2单选3多选4数值), 可选值JSON, 是否SKU规格, 是否可搜索)
ATTRS = [
    # 手机类（智能手机）
    (35, "颜色", 2, '["黑色","白色","银色","金色","红色","蓝色","紫色","绿色"]', 1, 1),
    (35, "存储容量", 4, '["64G","128G","256G","512G","1T"]', 1, 1),
    (35, "运行内存", 4, '["4G","6G","8G","12G","16G","24G"]', 1, 1),
    (35, "屏幕尺寸", 4, '["5.5","6.1","6.5","6.7","6.9","7.0"]', 0, 0),
    (35, "处理器", 1, None, 0, 0),
    (35, "后置摄像头", 1, None, 0, 0),
    (35, "电池容量", 4, '["3000","4000","5000","6000","7000"]', 0, 0),
    # 笔记本类
    (38, "颜色", 2, '["银色","深空灰","金色","黑色","白色"]', 1, 1),
    (38, "内存", 4, '["8G","16G","32G","64G"]', 1, 1),
    (38, "硬盘", 4, '["256G SSD","512G SSD","1T SSD","2T SSD"]', 1, 1),
    (38, "屏幕尺寸", 4, '["13.3","14","15.6","16","17.3"]', 0, 0),
    (38, "处理器型号", 1, None, 0, 0),
    (38, "显卡型号", 1, None, 0, 0),
    # 耳机类
    (46, "颜色", 2, '["黑色","白色","银色","金色","蓝色","红色"]', 1, 1),
    (46, "连接方式", 2, '["有线","蓝牙","双模"]', 1, 1),
    (46, "降噪", 2, '["主动降噪","被动降噪","无降噪"]', 0, 0),
    # T恤（男装 T恤）
    (48, "颜色", 2, '["黑色","白色","红色","蓝色","绿色","黄色","粉色","灰色","卡其"]', 1, 1),
    (48, "尺码", 2, '["S","M","L","XL","XXL","XXXL"]', 1, 1),
    (48, "面料", 1, None, 0, 0),
    # 连衣裙（女装 连衣裙）
    (53, "颜色", 2, '["黑色","白色","红色","蓝色","绿色","黄色","粉色","灰色","卡其"]', 1, 1),
    (53, "尺码", 2, '["S","M","L","XL","XXL"]', 1, 1),
    (53, "面料", 1, None, 0, 0),
    # 运动鞋（跑步鞋）
    (58, "颜色", 2, '["黑色","白色","红色","蓝色","绿色","灰色","荧光黄"]', 1, 1),
    (58, "尺码", 2, '["36","37","38","39","40","41","42","43","44","45"]', 1, 1),
    # 口红（彩妆 口红）
    (70, "色号", 2, '["#001 象牙白","#002 自然色","#003 小麦色","#004 粉调白"]', 1, 1),
    (70, "质地", 2, '["哑光","水润","雾面","奶油肌"]', 0, 0),
]

# SPU 数据：每个元组 (名称, 副标题, 类目索引, 品牌索引, 最低价, 最高价, 单位, 主图占位)
PRODUCTS = [
    # 手机（智能手机）
    ("iPhone 16 Pro Max", "钛金属旗舰，A18芯片", 35, 1, 899900, 999900, "台"),
    ("iPhone 16 Pro", "钛金属专业级", 35, 1, 799900, 899900, "台"),
    ("iPhone 16", "A18芯片，超强续航", 35, 1, 599900, 699900, "台"),
    ("Galaxy S25 Ultra", "AI智能旗舰，2亿像素", 35, 2, 899900, 999900, "台"),
    ("Galaxy S25+", "大屏AI旗舰", 35, 2, 749900, 849900, "台"),
    ("Galaxy S25", "AI智能旗舰", 35, 2, 699900, 799900, "台"),
    ("小米15 Pro", "徕卡影像，骁龙8至尊", 35, 3, 499900, 599900, "台"),
    ("小米15", "徕卡影像旗舰", 35, 3, 399900, 499900, "台"),
    ("Redmi K80 Pro", "性能旗舰", 35, 3, 299900, 369900, "台"),
    ("华为Mate 70 Pro", "鸿蒙旗舰，麒麟芯片", 35, 4, 699900, 799900, "台"),
    ("华为Mate 70", "鸿蒙系统旗舰", 35, 4, 549900, 649900, "台"),
    ("华为Pura 70 Pro", "超聚光影像", 35, 4, 649900, 749900, "台"),
    ("OPPO Find X8 Pro", "双潜望影像旗舰", 35, 5, 529900, 629900, "台"),
    ("OPPO Find X8", "轻薄影像旗舰", 35, 5, 399900, 499900, "台"),
    ("vivo X200 Pro", "蔡司影像，天玑9400", 35, 6, 499900, 599900, "台"),
    ("vivo X200", "蔡司影像旗舰", 35, 6, 399900, 479900, "台"),
    ("一加13", "性能旗舰，哈苏影像", 35, 5, 429900, 529900, "台"),
    ("魅族21 Pro", "AI终端旗舰", 35, 3, 399900, 499900, "台"),
    ("荣耀Magic7 Pro", "AI智慧旗舰", 35, 4, 569900, 669900, "台"),
    ("真我GT7 Pro", "电竞性能旗舰", 35, 5, 299900, 369900, "台"),
    ("Nothing Phone (2a)", "Glyph灯光设计", 35, 2, 249900, 299900, "台"),
    # 笔记本
    ("MacBook Pro 14", "M4芯片，专业级", 38, 1, 1299900, 1499900, "台"),
    ("MacBook Pro 16", "M4 Max，极致性能", 38, 1, 1699900, 1999900, "台"),
    ("MacBook Air 13", "M3芯片，超轻薄", 38, 1, 899900, 999900, "台"),
    ("MacBook Air 15", "M3芯片，大屏轻薄", 38, 1, 999900, 1099900, "台"),
    ("ThinkPad X1 Carbon", "商务旗舰，超轻薄", 38, 7, 999900, 1199900, "台"),
    ("ThinkPad T14", "专业商务本", 38, 7, 799900, 999900, "台"),
    ("Dell XPS 14", "全面屏超轻薄", 38, 9, 999900, 1199900, "台"),
    ("Dell XPS 16", "大屏创作本", 38, 9, 1199900, 1399900, "台"),
    ("HP Spectre x360", "翻转触控旗舰", 38, 10, 899900, 1099900, "台"),
    ("HP Envy 16", "创意设计本", 38, 10, 799900, 999900, "台"),
    ("ASUS ROG 幻16 Air", "轻薄性能本", 38, 11, 999900, 1199900, "台"),
    ("ASUS Zenbook 14", "超轻薄商务本", 38, 11, 699900, 899900, "台"),
    ("联想小新Pro 14", "高性能轻薄本", 38, 8, 499900, 699900, "台"),
    ("联想ThinkBook 14", "商务全能本", 38, 8, 399900, 599900, "台"),
    # 平板（平板电脑）
    ("iPad Pro 11", "M4芯片，轻薄专业", 39, 1, 799900, 999900, "台"),
    ("iPad Pro 13", "M4芯片，大屏专业", 39, 1, 999900, 1299900, "台"),
    ("iPad Air 11", "M2芯片，轻薄全能", 39, 1, 599900, 699900, "台"),
    ("iPad Air 13", "M2芯片，大屏全能", 39, 1, 699900, 799900, "台"),
    ("iPad mini", "A17 Pro，便携旗舰", 39, 1, 399900, 499900, "台"),
    ("华为MatePad Pro 13.2", "鸿蒙专业平板", 39, 4, 699900, 799900, "台"),
    ("华为MatePad Air", "轻薄办公平板", 39, 4, 399900, 499900, "台"),
    ("小米平板7 Pro", "高性能创作平板", 39, 3, 299900, 399900, "台"),
    ("小米平板7", "高性价比平板", 39, 3, 199900, 299900, "台"),
    ("三星Galaxy Tab S10+", "AI旗舰平板", 39, 2, 799900, 899900, "台"),
    # 耳机
    ("AirPods Pro 2", "主动降噪，H2芯片", 46, 1, 189900, 199900, "副"),
    ("AirPods 4", "半入耳，降噪版", 46, 1, 129900, 139900, "副"),
    ("AirPods Max", "头戴式旗舰耳机", 46, 1, 399900, 439900, "副"),
    ("Galaxy Buds3 Pro", "智能降噪耳机", 46, 2, 159900, 179900, "副"),
    ("Galaxy Buds3", "开放式耳机", 46, 2, 99900, 119900, "副"),
    ("小米Buds 4 Pro", "旗舰降噪耳机", 46, 3, 99900, 129900, "副"),
    ("小米Buds 4", "轻降噪耳机", 46, 3, 59900, 79900, "副"),
    ("华为FreeBuds Pro 3", "超感知降噪", 46, 4, 129900, 149900, "副"),
    ("华为FreeBuds 5", "全开放舒适佩戴", 46, 4, 89900, 109900, "副"),
    ("Sony WF-1000XM5", "旗舰降噪耳机", 46, 7, 199900, 219900, "副"),
    ("Sony WH-1000XM5", "头戴式旗舰降噪", 46, 7, 279900, 299900, "副"),
    ("OPPO Enco X3", "丹拿调音旗舰", 46, 5, 99900, 119900, "副"),
    ("vivo TWS 4", "Hi-Fi音质旗舰", 46, 6, 79900, 99900, "副"),
    # T恤
    ("Uniqlo 基础款T恤", "舒适百搭", 48, 14, 7900, 9900, "件"),
    ("Uniqlo AIRism T恤", "凉感速干", 48, 14, 9900, 12900, "件"),
    ("Uniqlo U系列T恤", "设计师合作款", 48, 14, 14900, 19900, "件"),
    ("Nike Dri-FIT T恤", "运动速干", 48, 12, 24900, 29900, "件"),
    ("Nike 经典Logo T恤", "运动休闲", 48, 12, 19900, 25900, "件"),
    ("Adidas 三叶草T恤", "经典复古", 48, 13, 22900, 27900, "件"),
    ("Adidas 运动T恤", "吸湿排汗", 48, 13, 19900, 24900, "件"),
    # 连衣裙
    ("Zara 法式碎花连衣裙", "浪漫清新", 53, 15, 29900, 39900, "件"),
    ("Zara 简约通勤连衣裙", "干练优雅", 53, 15, 34900, 44900, "件"),
    ("H&M 针织连衣裙", "舒适修身", 53, 16, 19900, 29900, "件"),
    ("H&M 吊带连衣裙", "夏日必备", 53, 16, 14900, 24900, "件"),
    ("Uniqlo 连衣裙", "简约百搭", 53, 14, 19900, 29900, "件"),
    # 运动鞋（跑步鞋）
    ("Nike Air Force 1", "经典复古篮球鞋", 58, 12, 74900, 89900, "双"),
    ("Nike Air Max 90", "气垫经典跑鞋", 58, 12, 89900, 99900, "双"),
    ("Nike Dunk Low", "复古滑板鞋", 58, 12, 79900, 89900, "双"),
    ("Nike Vomero 17", "顶级缓震跑鞋", 58, 12, 129900, 139900, "双"),
    ("Adidas Ultraboost Light", "超轻缓震跑鞋", 58, 13, 139900, 149900, "双"),
    ("Adidas Samba", "经典复古板鞋", 58, 13, 89900, 99900, "双"),
    ("Adidas Superstar", "贝壳头经典", 58, 13, 79900, 89900, "双"),
    ("Adidas NMD R1", "经典潮流跑鞋", 58, 13, 99900, 109900, "双"),
    ("New Balance 327", "复古运动鞋", 58, 18, 89900, 99900, "双"),
    ("New Balance 990v6", "美产经典", 58, 18, 179900, 189900, "双"),
    ("New Balance 2002R", "复古机能", 58, 18, 119900, 129900, "双"),
    ("Converse Chuck 70", "经典帆布鞋", 58, 19, 59900, 69900, "双"),
    ("Converse Run Star Hike", "厚底增高帆布鞋", 58, 19, 79900, 89900, "双"),
    ("Vans Old Skool", "经典滑板鞋", 58, 20, 59900, 69900, "双"),
    ("Vans Authentic", "经典帆布鞋", 58, 20, 49900, 59900, "双"),
    # 口红
    ("Dior 烈艳蓝金唇膏", "经典缎面口红", 70, 27, 38000, 42000, "支"),
    ("Dior 魅惑唇膏", "水润光泽", 70, 27, 35000, 39000, "支"),
    ("Chanel 炫亮魅力唇膏", "丝绒哑光", 70, 28, 45000, 49000, "支"),
    ("Chanel 可可小姐唇膏", "水润清透", 70, 28, 42000, 46000, "支"),
    ("YSL 小金条口红", "哑光丝绒", 70, 27, 38000, 42000, "支"),
    ("YSL 黑管唇釉", "镜面水光", 70, 27, 36000, 40000, "支"),
    ("MAC 子弹头口红", "经典哑光", 70, 26, 23000, 27000, "支"),
    ("MAC 水漾润泽口红", "滋润保湿", 70, 26, 23000, 27000, "支"),
    ("完美日记 小红钻", "国货之光", 70, 24, 9990, 12990, "支"),
    ("花西子 雕花口红", "东方美学", 70, 24, 19900, 24900, "支"),
]

# SKU 规格组合模板
COLORS = ["黑色", "白色", "银色", "金色", "红色", "蓝色", "紫色", "绿色", "粉色", "灰色", "卡其", "荧光黄"]
STORAGES = ["64G", "128G", "256G", "512G", "1T"]
RAMS = ["4G", "6G", "8G", "12G", "16G", "24G"]
SCREEN_SIZES = ["5.5", "6.1", "6.5", "6.7", "6.9", "7.0"]
CLOTHES_SIZES = ["S", "M", "L", "XL", "XXL", "XXXL"]
SHOE_SIZES = ["36", "37", "38", "39", "40", "41", "42", "43", "44", "45"]
LIPSTICK_SHADES = ["#001 经典红", "#002 豆沙粉", "#003 橘红", "#004 玫红", "#005 奶茶色", "#006 复古红"]


def generate_spec(category_id):
    """根据类目生成合理的SKU规格"""
    # 手机类（智能手机）
    if category_id == 35:
        color = random.choice(COLORS[:8])
        storage = random.choice(STORAGES[:4])
        ram = random.choice(RAMS[:5])
        return f'{{"颜色":"{color}","存储容量":"{storage}","运行内存":"{ram}"}}', {
            "color": color, "storage": storage, "ram": ram
        }
    # 笔记本类
    elif category_id == 38:
        color = random.choice(["银色", "深空灰", "金色", "黑色"])
        ram = random.choice(["8G", "16G", "32G"])
        disk = random.choice(["256G SSD", "512G SSD", "1T SSD"])
        return f'{{"颜色":"{color}","内存":"{ram}","硬盘":"{disk}"}}', {
            "color": color, "ram": ram, "disk": disk
        }
    # 耳机类
    elif category_id == 46:
        color = random.choice(COLORS[:6])
        connect = random.choice(["蓝牙", "有线", "双模"])
        return f'{{"颜色":"{color}","连接方式":"{connect}"}}', {
            "color": color, "connect": connect
        }
    # T恤/连衣裙类
    elif category_id in [48, 53]:
        color = random.choice(COLORS[:10])
        size = random.choice(CLOTHES_SIZES)
        return f'{{"颜色":"{color}","尺码":"{size}"}}', {"color": color, "size": size}
    # 运动鞋类（跑步鞋）
    elif category_id == 58:
        color = random.choice(COLORS[:8])
        size = random.choice(SHOE_SIZES)
        return f'{{"颜色":"{color}","尺码":"{size}"}}', {"color": color, "size": size}
    # 口红
    elif category_id == 70:
        shade = random.choice(LIPSTICK_SHADES[:6])
        texture = random.choice(["哑光", "水润", "雾面", "奶油肌"])
        return f'{{"色号":"{shade}","质地":"{texture}"}}', {"shade": shade, "texture": texture}
    else:
        color = random.choice(COLORS[:6])
        return f'{{"颜色":"{color}"}}', {"color": color}


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
        brand_id_map = {}
        for name, cname, letter in BRANDS:
            cur.execute(
                "INSERT INTO sp_brands (name, english_name, first_letter, sort_order, status, created_at) "
                "VALUES (%s, %s, %s, %s, 1, %s)",
                (cname, name, letter, random.randint(1, 100), now),
            )
            brand_id_map[(cname, name)] = cur.lastrowid
        print(f"  品牌: {len(BRANDS)}")

        # 类目
        cat_ids = {}
        for i, (parent_id, name, level) in enumerate(CATEGORIES, 1):
            # 找父级path
            path = ""
            if parent_id > 0:
                parent_id_actual = cat_ids.get(parent_id)
                if parent_id_actual:
                    # 获取父级path
                    cur.execute("SELECT path FROM sp_categories WHERE id = %s", (parent_id_actual,))
                    row = cur.fetchone()
                    if row:
                        path = row[0] + str(parent_id_actual) + "/"
            cur.execute(
                "INSERT INTO sp_categories (name, parent_id, level, path, sort_order, status, created_at) "
                "VALUES (%s, %s, %s, %s, %s, 1, %s)",
                (name, parent_id if parent_id == 0 else cat_ids.get(parent_id, 0),
                 level, path, i * 10, now),
            )
            cat_ids[i] = cur.lastrowid
        print(f"  类目: {len(CATEGORIES)}")

        # 属性
        attr_map = {}
        for cat_idx, name, input_type, values, is_sku, searchable in ATTRS:
            cur.execute(
                "INSERT INTO sp_attributes (name, category_id, input_type, `values`, is_sku_spec, searchable, status, created_at) "
                "VALUES (%s, %s, %s, %s, %s, %s, 1, %s)",
                (name, cat_ids[cat_idx], input_type, values, is_sku, searchable, now),
            )
            attr_map[(cat_idx, name)] = cur.lastrowid
        print(f"  属性: {len(ATTRS)}")

        # SPU & SKU
        total_skus = 0
        product_count = 0
        for name, subtitle, cat_idx, brand_idx, price, market, unit in PRODUCTS:
            # 获取品牌ID
            brand_id = brand_idx  # 直接用索引
            cur.execute(
                "INSERT INTO sp_products (name, subtitle, category_id, brand_id, unit, main_image, "
                "min_price, max_price, status, sort_order, created_at, updated_at) "
                "VALUES (%s, %s, %s, %s, %s, '', %s, %s, 2, %s, %s, %s)",
                (name, subtitle, cat_ids[cat_idx], brand_id, unit, price, market,
                 random.randint(1, 100), now, now),
            )
            spu_id = cur.lastrowid
            product_count += 1

            # 生成 SKU（每个 SPU 2-6 个规格组合）
            sku_count = random.randint(2, 6)
            generated_specs = set()
            for j in range(sku_count):
                # 生成合理的规格
                spec_json, spec_dict = generate_spec(cat_idx)
                if spec_json in generated_specs:
                    # 如果重复，稍作变化
                    color = random.choice(COLORS[:8])
                    if "颜色" in spec_dict:
                        spec_dict["颜色"] = color
                    spec_json = '{"' + '","'.join([f'{k}":"{v}' for k, v in spec_dict.items()]) + '"}'
                    if spec_json in generated_specs:
                        continue
                generated_specs.add(spec_json)

                sku_price = price + random.randint(-int(price * 0.2), int(price * 0.3))
                sku_price = max(price - int(price * 0.3), sku_price)
                sku_code = f"SKU{spu_id}-{j+1:03d}"
                barcode = f"{random.randint(1000000000000, 9999999999999)}"
                cur.execute(
                    "INSERT INTO sp_skus (product_id, sku_code, barcode, spec, price, market_price, cost_price, status, created_at) "
                    "VALUES (%s, %s, %s, %s, %s, %s, %s, 1, %s)",
                    (spu_id, sku_code, barcode, spec_json, sku_price,
                     sku_price + random.randint(int(sku_price * 0.1), int(sku_price * 0.3)),
                     int(sku_price * 0.6), now),
                )
                total_skus += 1

        print(f"  SPU: {product_count}, SKU: {total_skus}")

    conn.commit()
    print("商品中心 ✅\n")


# ── 库存中心 ──────────────────────────────────────

def seed_inventory(conn):
    now = datetime.now().strftime(FMT)
    with conn.cursor() as cur:
        cur.execute("SELECT id FROM sp_skus WHERE deleted_at IS NULL")
        skus = cur.fetchall()
        for sku in skus:
            qty = random.randint(50, 1000)
            reserved = random.randint(0, int(qty * 0.3))
            threshold = random.randint(5, 50)
            status = "instock" if qty > threshold else "lowstock"
            cur.execute(
                "INSERT INTO sp_inventories (sku_id, quantity, reserved, threshold, status, created_at) "
                "VALUES (%s, %s, %s, %s, %s, %s)",
                (sku[0], qty, reserved, threshold, status, now),
            )

            # 部分库存记录日志
            if random.random() < 0.3:
                delta = random.randint(10, 50)
                cur.execute(
                    "INSERT INTO sp_inventory_logs (sku_id, warehouse_id, before_quantity, after_quantity, "
                    "before_reserved, after_reserved, change_amount, "
                    "change_type, reference_id, operator, note, created_at) "
                    "VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)",
                    (sku[0], 0, qty - delta, qty,
                     0, reserved, delta,
                     "purchase", "", "admin", "初始入库", now),
                )
    conn.commit()
    print(f"  库存: {len(skus)} 条")
    print("库存中心 ✅\n")


# ── 营销中心 ──────────────────────────────────────

def seed_marketing(conn):
    now = datetime.now()
    now_str = now.strftime(FMT)

    promos = [
        # (名称, 类型, 条件值(分), 优惠值, 每人限领, 总量)
        ("满200减30", 4, 20000, 3000, 1, 1000),
        ("满500减100", 4, 50000, 10000, 1, 500),
        ("满1000减200", 4, 100000, 20000, 1, 300),
        ("全场8折", 5, 0, 20, 0, 0),
        ("全场85折", 5, 0, 15, 0, 0),
        ("新用户满减券", 1, 0, 5000, 1, 500),
        ("新用户专属8折", 1, 0, 20, 1, 300),
        ("会员9折", 6, 0, 10, 0, 0),
        ("会员85折", 6, 0, 15, 0, 0),
        ("限时秒杀-手机", 3, 0, 50, 1, 50),
        ("限时秒杀-耳机", 3, 0, 30, 2, 100),
        ("限时秒杀-运动鞋", 3, 0, 40, 1, 80),
        ("限时秒杀-化妆品", 3, 0, 25, 2, 120),
        ("满300送赠品", 7, 30000, 0, 1, 200),
        ("买2件9折", 8, 0, 10, 0, 0),
        ("买3件8折", 8, 0, 20, 0, 0),
        ("双11预售", 2, 0, 0, 1, 1000),
        ("618大促", 2, 0, 0, 1, 1000),
        ("品牌日特惠", 2, 0, 0, 1, 500),
        ("圣诞限定折扣", 2, 0, 0, 1, 300),
    ]

    with conn.cursor() as cur:
        for i, (name, ptype, condition, benefit, per_limit, total_qty) in enumerate(promos, 1):
            start = (now - timedelta(days=random.randint(0, 5))).strftime(FMT)
            end = (now + timedelta(days=random.randint(3, 30))).strftime(FMT)
            status = random.choices([1, 2, 3], weights=[1, 8, 1])[0]  # 1未开始 2进行中 3已结束

            cur.execute(
                "INSERT INTO mkt_promotions (promo_name, promo_type, promo_code, start_time, end_time, "
                "total_quantity, per_user_limit, used_quantity, status, created_at) "
                "VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)",
                (name, ptype, f"PROMO{i:04d}", start, end,
                 total_qty, per_limit, random.randint(0, int(total_qty * 0.5)), status, now_str),
            )
            promo_id = cur.lastrowid

            # 规则
            rule_condition = 2 if condition > 0 else 1
            cur.execute(
                "INSERT INTO mkt_promotion_rules (promotion_id, rule_name, condition_type, condition_value, "
                "benefit_type, benefit_value, created_at) VALUES (%s, %s, %s, %s, %s, %s, %s)",
                (promo_id, f"{name}规则", rule_condition, condition / 100,
                 1 if ptype in [4, 5, 6, 7, 8] else 2, benefit / 100, now_str),
            )

            # 秒杀/限时折扣 配置商品
            if ptype == 3:
                cur.execute("SELECT id FROM sp_products ORDER BY RAND() LIMIT %s",
                           (random.randint(3, 8),))
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
    parser.add_argument("--module", choices=["product", "inventory", "marketing"],
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
        if args.module in modules:
            modules[args.module](conn)
        else:
            print(f"未知模块: {args.module}")
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