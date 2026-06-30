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

# ── 自动生成 SPU 配置 ──────────────────────────────
# 每片叶子类目生成 PRODUCTS_PER_CATEGORY 个 SPU
PRODUCTS_PER_CATEGORY = 25

# 叶子类目生成配置: (类目索引, 品牌池, 价格区间(分), 单位, 名称模板, 副标题模板)
CATEGORY_PROD_CFG = {
    # 手机通讯
    35: (range(1, 7), (199900, 999900), "台",
         ["{b}旗舰手机", "{b}智能手机", "{b}Pro 机型", "{b}轻旗舰", "{b}性能版",
          "{b}Ultra 版", "{b}Plus 版", "{b}e 青春版"],
         ["5G旗舰", "AI智能", "高性能", "超清影像", "长续航", "轻薄设计"]),
    36: (range(1, 7), (9900, 29900), "台",
         ["{b}老年机", "{b}功能手机", "{b}按键手机"],
         ["大字体", "长续航", "超长待机", "简易操作"]),
    37: (range(7, 12), (19900, 89900), "台",
         ["{b}对讲机", "{b}专业对讲机", "{b}远距离对讲机"],
         ["远距离", "防水防尘", "长续航", "清晰通话"]),
    # 电脑办公
    38: (range(1, 12), (399900, 1999900), "台",
         ["{b}笔记本", "{b}轻薄本", "{b}商务本", "{b}性能本", "{b}创作本",
          "{b}电竞本", "{b}Ultrabook"],
         ["高性能", "轻薄便携", "商务办公", "游戏电竞", "创意设计"]),
    39: (range(1, 6), (199900, 999900), "台",
         ["{b}平板电脑", "{b}旗舰平板", "{b}轻薄平板", "{b}学习平板"],
         ["影音娱乐", "移动办公", "学习教育", "创作绘画"]),
    40: (range(7, 12), (299900, 999900), "台",
         ["{b}台式机", "{b}一体机", "{b}家用台式机", "{b}办公台式机"],
         ["高性价比", "商务办公", "家用娱乐", "设计制图"]),
    41: (range(7, 12), (99900, 499900), "台",
         ["{b}显示器", "{b}显示器 2K", "{b}显示器 4K", "{b}曲面显示器", "{b}电竞显示器"],
         ["高清显示", "护眼", "高刷", "专业色准"]),
    42: (range(9, 12), (29900, 199900), "台",
         ["{b}打印机", "{b}彩色打印机", "{b}激光打印机", "{b}多功能一体机"],
         ["家用打印", "办公高效", "无线打印", "彩色打印"]),
    # 数码配件
    43: (range(1, 6), (1900, 12900), "个",
         ["{b}手机壳", "{b}保护壳", "{b}透明手机壳", "{b}硅胶手机壳"],
         ["防摔保护", "轻薄", "高颜值", "简约设计"]),
    44: (range(1, 6), (2900, 19900), "个",
         ["{b}充电器", "{b}快充头", "{b}氮化镓充电器", "{b}无线充电器"],
         ["快充", "氮化镓", "多口输出", "安全充电"]),
    45: (range(1, 6), (900, 5900), "条",
         ["{b}数据线", "{b}Type-C数据线", "{b}快充数据线", "{b}编织数据线"],
         ["快充传输", "耐用编织", "加长版", "磁吸收纳"]),
    46: (range(1, 8), (59900, 439900), "副",
         ["{b}无线耳机", "{b}降噪耳机", "{b}头戴式耳机", "{b}运动耳机",
          "{b}TWS耳机", "{b}入耳式耳机"],
         ["主动降噪", "Hi-Fi音质", "长续航", "舒适佩戴"]),
    47: (range(1, 8), (49900, 199900), "个",
         ["{b}移动电源", "{b}充电宝", "{b}快充移动电源", "{b}大容量充电宝"],
         ["大容量", "快充", "轻薄便携", "安全电芯"]),
    # 男装
    48: (range(12, 18), (5900, 39900), "件",
         ["{b}T恤", "{b}短袖T恤", "{b}基础款T恤", "{b}印花T恤",
          "{b}纯棉T恤", "{b}运动T恤"],
         ["舒适纯棉", "透气速干", "基础百搭", "休闲运动"]),
    49: (range(12, 18), (9900, 59900), "件",
         ["{b}衬衫", "{b}长袖衬衫", "{b}短袖衬衫", "{b}商务衬衫", "{b}休闲衬衫"],
         ["免烫", "修身版型", "商务休闲", "舒适面料"]),
    50: (range(12, 18), (19900, 99900), "件",
         ["{b}外套", "{b}夹克", "{b}风衣", "{b}休闲外套", "{b}棉服"],
         ["春秋款", "防风保暖", "休闲百搭", "轻薄便携"]),
    51: (range(12, 18), (14900, 79900), "条",
         ["{b}牛仔裤", "{b}直筒牛仔裤", "{b}修身牛仔裤", "{b}宽松牛仔裤"],
         ["经典款", "弹力舒适", "复古水洗", "百搭"]),
    52: (range(7, 12), (29900, 199900), "套",
         ["{b}西服", "{b}西装", "{b}正装西服", "{b}商务西服"],
         ["修身版型", "商务正装", "羊毛混纺", "婚礼西服"]),
    # 女装
    53: (range(14, 18), (14900, 69900), "件",
         ["{b}连衣裙", "{b}碎花连衣裙", "{b}通勤连衣裙", "{b}针织连衣裙",
          "{b}吊带连衣裙", "{b}衬衫裙"],
         ["浪漫清新", "优雅通勤", "舒适修身", "夏日清爽"]),
    54: (range(14, 18), (9900, 49900), "件",
         ["{b}上衣", "{b}针织衫", "{b}雪纺衫", "{b}卫衣", "{b}打底衫"],
         ["百搭款", "舒适面料", "时尚设计", "春秋穿搭"]),
    55: (range(14, 18), (14900, 39900), "条",
         ["{b}半身裙", "{b}A字裙", "{b}百褶裙", "{b}包臀裙", "{b}牛仔裙"],
         ["显瘦版型", "高腰设计", "优雅气质", "休闲百搭"]),
    56: (range(14, 18), (19900, 89900), "件",
         ["{b}女装外套", "{b}风衣", "{b}毛呢外套", "{b}牛仔外套", "{b}小香风外套"],
         ["春秋外套", "百搭款", "气质通勤", "保暖时尚"]),
    57: (range(14, 18), (14900, 59900), "件",
         ["{b}毛衣", "{b}针织毛衣", "{b}羊毛衫", "{b}开衫毛衣", "{b}高领毛衣"],
         ["柔软保暖", "舒适羊毛", "百搭基础", "时尚宽松"]),
    # 运动鞋
    58: (range(12, 21), (49900, 189900), "双",
         ["{b}跑鞋", "{b}运动鞋", "{b}缓震跑鞋", "{b}训练鞋",
          "{b}复古跑鞋", "{b}越野跑鞋"],
         ["缓震", "轻量", "稳定支撑", "透气"]),
    59: (range(12, 21), (59900, 159900), "双",
         ["{b}篮球鞋", "{b}实战篮球鞋", "{b}篮球文化鞋"],
         ["缓震回弹", "包裹支撑", "耐磨防滑", "实战利器"]),
    60: (range(12, 21), (29900, 99900), "双",
         ["{b}休闲鞋", "{b}板鞋", "{b}帆布鞋", "{b}滑板鞋", "{b}德训鞋"],
         ["经典百搭", "复古风格", "舒适脚感", "日常穿着"]),
    # 食品
    61: (range(22, 24), (990, 2990), "袋",
         ["{b}薯片", "{b}膨化食品", "{b}虾条", "{b}米饼", "{b}爆米花"],
         ["香脆可口", "多味装", "休闲零食", "分享装"]),
    62: (range(22, 24), (990, 3990), "盒",
         ["{b}巧克力", "{b}糖果", "{b}棒棒糖", "{b}软糖", "{b}夹心巧克力"],
         ["丝滑口感", "纯可可脂", "精美包装", "礼盒装"]),
    63: (range(22, 24), (990, 2990), "袋",
         ["{b}饼干", "{b}曲奇", "{b}蛋卷", "{b}威化饼干", "{b}夹心饼干"],
         ["酥脆", "独立包装", "早餐搭档", "休闲时光"]),
    # 美妆护肤
    64: (range(24, 30), (1900, 19900), "支",
         ["{b}洁面乳", "{b}洗面奶", "{b}洁面泡沫", "{b}卸妆洁面"],
         ["温和清洁", "深层洁净", "保湿", "控油祛痘"]),
    65: (range(24, 30), (2900, 29900), "瓶",
         ["{b}爽肤水", "{b}柔肤水", "{b}化妆水", "{b}收敛水"],
         ["补水保湿", "舒缓", "收缩毛孔", "清爽"]),
    66: (range(24, 30), (4900, 59900), "瓶",
         ["{b}精华液", "{b}精华露", "{b}美白精华", "{b}抗皱精华", "{b}保湿精华"],
         ["抗衰老", "美白淡斑", "深层保湿", "修复肌肤"]),
    67: (range(24, 30), (2900, 49900), "瓶",
         ["{b}面霜", "{b}保湿面霜", "{b}抗皱面霜", "{b}修复面霜", "{b}日霜"],
         ["深层滋润", "抗皱紧致", "修护屏障", "清爽不腻"]),
    68: (range(24, 30), (2900, 39900), "支",
         ["{b}防晒霜", "{b}防晒乳", "{b}防晒喷雾", "{b}隔离防晒"],
         ["SPF50+", "PA+++", "清爽不油腻", "防水持久"]),
    # 彩妆
    69: (range(24, 30), (3900, 39900), "盒",
         ["{b}粉底液", "{b}气垫", "{b}粉饼", "{b}BB霜", "{b}CC霜"],
         ["自然遮瑕", "水润服帖", "持妆长久", "轻薄透气"]),
    70: (range(24, 30), (9900, 49900), "支",
         ["{b}口红", "{b}唇膏", "{b}唇釉", "{b}唇泥", "{b}唇彩"],
         ["哑光", "水润", "丝绒", "雾面", "奶油肌"]),
    71: (range(24, 30), (2900, 29900), "盒",
         ["{b}眼影盘", "{b}单色眼影", "{b}眼影笔", "{b}液体眼影"],
         ["大地色系", "粉棕系", "哑光", "珠光闪粉"]),
    72: (range(24, 30), (1900, 19900), "盒",
         ["{b}腮红", "{b}腮红盘", "{b}液体腮红", "{b}腮红膏"],
         ["自然显色", "粉嫩", "修容", "元气妆"]),
    73: (range(24, 30), (1900, 19900), "支",
         ["{b}睫毛膏", "{b}睫毛打底", "{b}纤长睫毛膏", "{b}浓密睫毛膏"],
         ["卷翘纤长", "浓密", "防水不晕", "自然裸感"]),
}


def generate_products():
    """根据 CATEGORY_PROD_CFG 自动生成 SPU 数据"""
    products = []
    used_names = set()
    suffix_idx = 0
    for cat_idx, cfg in CATEGORY_PROD_CFG.items():
        brand_range, price_range, unit, name_tpls, sub_tpls = cfg
        brands = list(brand_range)
        for i in range(PRODUCTS_PER_CATEGORY):
            brand = random.choice(brands)
            name_tpl = random.choice(name_tpls)
            sub_tpl = random.choice(sub_tpls)
            brand_name = BRANDS[brand - 1][1]  # 中文品牌名
            name = name_tpl.format(b=brand_name)
            # 避免重名
            while name in used_names:
                suffix_idx += 1
                name = f"{name_tpl.format(b=brand_name)} {suffix_idx}"
            used_names.add(name)
            price_min = price_range[0]
            price_max = price_range[1]
            price = random.randint(price_min, price_max)
            market_price = int(price * random.uniform(1.1, 1.3))
            products.append((name, sub_tpl, cat_idx, brand, price, market_price, unit))
    return products


# 生成 SPU 数据（惰性加载，仅在 seed_product 中调用一次）
_GENERATED_PRODUCTS = None

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
    # 运动鞋/休闲鞋/篮球鞋
    elif category_id in [58, 59, 60]:
        color = random.choice(COLORS[:8])
        size = random.choice(SHOE_SIZES)
        return f'{{"颜色":"{color}","尺码":"{size}"}}', {"color": color, "size": size}
    # 男装外套/衬衫/牛仔裤/西服 + 女装外套/毛衣
    elif category_id in [49, 50, 51, 52, 56, 57]:
        color = random.choice(COLORS[:8])
        size = random.choice(CLOTHES_SIZES)
        return f'{{"颜色":"{color}","尺码":"{size}"}}', {"color": color, "size": size}
    # 食品（口味）
    elif category_id in [61, 62, 63]:
        flavor = random.choice(["原味", "番茄味", "麻辣味", "烧烤味", "海苔味"])
        return f'{{"口味":"{flavor}"}}', {"flavor": flavor}
    # 美妆护肤（洁面/爽肤水/精华/面霜/防晒）
    elif category_id in [64, 68]:
        spec_type = random.choice(["清爽型", "滋润型", "敏感肌用"])
        return f'{{"类型":"{spec_type}","容量":"{random.choice(["100ml","150ml","200ml"])}"}}', {"type": spec_type}
    elif category_id in [65, 66, 67]:
        spec_type = random.choice(["清爽型", "滋润型", "修护型"])
        capacity = random.choice(["30ml", "50ml", "100ml", "120ml"])
        return f'{{"类型":"{spec_type}","容量":"{capacity}"}}', {"type": spec_type, "capacity": capacity}
    # 彩妆（粉底/眼影/腮红/睫毛膏）
    elif category_id in [69, 71, 72, 73]:
        shade = random.choice(["自然色", "象牙白", "小麦色", "粉调"])
        return f'{{"色号":"{shade}"}}', {"shade": shade}
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
    global _GENERATED_PRODUCTS
    if _GENERATED_PRODUCTS is None:
        _GENERATED_PRODUCTS = generate_products()
        print(f"  自动生成 SPU: {len(_GENERATED_PRODUCTS)} 个")
    products = _GENERATED_PRODUCTS

    with conn.cursor() as cur:
        # 品牌
        brand_id_map = {}
        for name, cname, letter in BRANDS:
            cur.execute(
                "INSERT INTO sp_brands (name, english_name, first_letter, sort_order, status) "
                "VALUES (%s, %s, %s, %s, 1)",
                (cname, name, letter, random.randint(1, 100)),
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
                "INSERT INTO sp_categories (name, parent_id, level, path, sort_order, status) "
                "VALUES (%s, %s, %s, %s, %s, 1)",
                (name, parent_id if parent_id == 0 else cat_ids.get(parent_id, 0),
                 level, path, i * 10),
            )
            cat_ids[i] = cur.lastrowid
        print(f"  类目: {len(CATEGORIES)}")

        # 属性
        attr_map = {}
        for cat_idx, name, input_type, values, is_sku, searchable in ATTRS:
            cur.execute(
                "INSERT INTO sp_attributes (name, category_id, input_type, `values`, is_sku_spec, searchable, status) "
                "VALUES (%s, %s, %s, %s, %s, %s, 1)",
                (name, cat_ids[cat_idx], input_type, values, is_sku, searchable),
            )
            attr_map[(cat_idx, name)] = cur.lastrowid
        print(f"  属性: {len(ATTRS)}")

        # SPU & SKU
        total_skus = 0
        product_count = 0
        for name, subtitle, cat_idx, brand_idx, price, market, unit in products:
            # 获取品牌ID
            brand_id = brand_idx  # 直接用索引
            cur.execute(
                "INSERT INTO sp_products (name, subtitle, category_id, brand_id, unit, main_image, "
                "min_price, max_price, status) "
                "VALUES (%s, %s, %s, %s, %s, '', %s, %s, 2)",
                (name, subtitle, cat_ids[cat_idx], brand_id, unit, price, market),
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
                    "INSERT INTO sp_skus (product_id, sku_code, barcode, spec, price, market_price, cost_price, status) "
                    "VALUES (%s, %s, %s, %s, %s, %s, %s, 1)",
                    (spu_id, sku_code, barcode, spec_json, sku_price,
                     sku_price + random.randint(int(sku_price * 0.1), int(sku_price * 0.3)),
                     int(sku_price * 0.6)),
                )
                total_skus += 1

            # 生成图文描述（部分 SPU 有）
            has_desc = 0
            if random.random() < 0.7:
                desc_text = f"{name}是一款优质的{subtitle}产品，给您带来极致体验。采用高品质材料，精心打造每一个细节。"
                mobile_text = f"<h1>{name}</h1><p>{subtitle}，{desc_text}</p>"
                cur.execute(
                    "INSERT INTO sp_product_descriptions (product_id, description, mobile_description) "
                    "VALUES (%s, %s, %s)",
                    (spu_id, desc_text, mobile_text),
                )
                has_desc = 1

            # 同步 has_description
            if has_desc == 1:
                cur.execute("UPDATE sp_products SET has_description = 1 WHERE id = %s", (spu_id,))

            # 生成属性值（为该类目关联的属性随机赋值）
            for (a_cat_idx, attr_name), attr_id in attr_map.items():
                if a_cat_idx == cat_idx:
                    # 给属性一个随机值
                    val = random.choice(["标准", "优质", "普通", "高级", "入门"])
                    if attr_name in ["颜色", "色号"]:
                        val = random.choice(["黑色", "白色", "红色", "蓝色"])
                    elif attr_name in ["面料", "材质", "处理器型号", "显卡型号"]:
                        val = random.choice(["优质材料", "标准款", "高性能版"])
                    cur.execute(
                        "INSERT IGNORE INTO sp_product_attributes (product_id, attribute_id, value, sort_order) "
                        "VALUES (%s, %s, %s, 0)",
                        (spu_id, attr_id, val),
                    )

        print(f"  SPU: {product_count}, SKU: {total_skus}")

        # 类目-品牌关联（按业务合理性匹配）
        # key=CATEGORIES 中二级类目的索引(parent_id in CATEGORIES tuple), value=品牌1-based索引
        cat_brand_groups = {
            11: range(1, 7),        # 手机通讯 → Apple/三星/小米/华为/OPPO/vivo
            12: range(1, 12),       # 电脑办公 → 电子品牌+联想/戴尔/惠普/华硕/索尼
            13: range(1, 8),        # 数码配件 → 主要电子品牌+索尼
            16: [12, 13, 14, 15, 16, 21, 30, 32],  # 男装 → Nike/Adidas/Uniqlo/Zara/H&M/Supreme/Gucci/LV
            17: [14, 15, 16, 17, 30, 31, 32],       # 女装 → Uniqlo/Zara/H&M/Lululemon/Gucci/Prada/LV
            19: [12, 13, 17, 18, 19, 20],            # 运动鞋 → Nike/Adidas/Lululemon/NB/Converse/Vans
            23: range(22, 26),      # 食品饮料 → Starbucks/可口可乐/百事/雀巢
            28: range(26, 30),      # 面部护肤 → L'Oreal/雅诗兰黛/Dior/Chanel
            29: range(26, 30),      # 彩妆 → L'Oreal/雅诗兰黛/Dior/Chanel
        }
        cb_count = 0
        for i, (cat_parent_id, _, cat_level) in enumerate(CATEGORIES, 1):
            # 二级类目索引对应 brand group key
            key = i if cat_level == 2 else cat_parent_id
            group = cat_brand_groups.get(key, range(1, 33))
            brand_ids = random.sample(list(group), min(len(list(group)), random.randint(3, 8)))
            for bid in brand_ids:
                cur.execute(
                    "INSERT IGNORE INTO sp_category_brands (category_id, brand_id) VALUES (%s, %s)",
                    (cat_ids[i], bid),
                )
                cb_count += 1
        print(f"  类目-品牌: {cb_count}")

    conn.commit()
    print("商品中心 ✅\n")


# ── 库存中心 ──────────────────────────────────────

def seed_inventory(conn):
    with conn.cursor() as cur:
        cur.execute("SELECT id FROM sp_skus WHERE deleted_at IS NULL")
        skus = cur.fetchall()
        for sku in skus:
            qty = random.randint(50, 1000)
            reserved = random.randint(0, int(qty * 0.3))
            threshold = random.randint(5, 50)
            status = "instock" if qty > threshold else "lowstock"
            cur.execute(
                "INSERT INTO sp_inventories (sku_id, quantity, reserved, threshold, status) "
                "VALUES (%s, %s, %s, %s, %s)",
                (sku[0], qty, reserved, threshold, status),
            )

            # 部分库存记录日志
            if random.random() < 0.3:
                delta = random.randint(10, 50)
                cur.execute(
                    "INSERT INTO sp_inventory_logs (sku_id, warehouse_id, before_quantity, after_quantity, "
                    "before_reserved, after_reserved, change_amount, "
                    "change_type, reference_id, operator, note) "
                    "VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)",
                    (sku[0], 0, qty - delta, qty,
                     0, reserved, delta,
                     "purchase", "", "admin", "初始入库"),
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
                      "sp_skus", "sp_product_descriptions", "sp_product_attributes",
                      "sp_inventories", "mkt_promotions"]:
            cur.execute(f"SELECT COUNT(*) AS cnt FROM {table}")
            row = cur.fetchone()
            print(f"  {table}: {row[0]}")

    conn.close()
    print("\n完成!")


if __name__ == "__main__":
    main()