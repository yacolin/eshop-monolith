"""
Generate sample data for categories, products, and inventory for the eshop-microservices project.
This script connects directly to the database and inserts sample data.
"""

import json
import uuid
from datetime import datetime
import random
import os
import pymysql  # Requires: pip install PyMySQL
import sys


def mysql_utc_now():
    return datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S")


def generate_categories():
    """Generate sample categories"""
    categories = []

    # Root categories
    root_categories = [
        "Electronics",
        "Clothing",
        "Books",
        "Home & Garden",
        "Sports & Outdoors",
        "Beauty & Health",
        "Toys & Games",
        "Automotive",
        "Food & Grocery",
        "Office Supplies"
    ]

    # Sub-categories
    sub_categories_map = {
        "Electronics": ["Smartphones", "Laptops", "Tablets", "TVs", "Cameras", "Audio Equipment"],
        "Clothing": ["Men's Clothing", "Women's Clothing", "Kids' Clothing", "Shoes", "Accessories"],
        "Books": ["Fiction", "Non-Fiction", "Science", "Technology", "Biography", "Cookbooks"],
        "Home & Garden": ["Furniture", "Decor", "Kitchen", "Bedding", "Tools", "Outdoor"]
    }

    # Generate root categories
    for name in root_categories:
        category = {
            "id": None,
            "name": name,
            "description": f"Category for {name.lower()}",
            "parent_id": None,
            "created_at": mysql_utc_now(),
            "updated_at": mysql_utc_now(),
        }
        categories.append(category)

        # Generate sub-categories if they exist
        if name in sub_categories_map:
            for sub_name in sub_categories_map[name]:
                sub_category = {
                    "id": None,
                    "name": sub_name,
                    "description": f"Sub-category for {sub_name.lower()} under {name}",
                    "parent_name": name,
                    "parent_id": None,
                    "created_at": mysql_utc_now(),
                    "updated_at": mysql_utc_now(),
                }
                categories.append(sub_category)

    return categories


def generate_products(categories):
    """Generate sample products"""
    products = []

    # Sample product names by category
    product_samples = {
        "Smartphones": [
            {"name": "iPhone 15 Pro",
                "description": "Latest iPhone with advanced camera system"},
            {"name": "iPhone 14 Pro",
                "description": "Previous-generation Pro phone with premium display"},
            {"name": "iPhone 15",
                "description": "Mainstream iPhone with great everyday performance"},
            {"name": "Samsung Galaxy S24",
                "description": "Android flagship with AI capabilities"},
            {"name": "Samsung Galaxy S23 Ultra",
                "description": "Ultra model with zoom camera and S Pen"},
            {"name": "Google Pixel 8",
                "description": "Pure Android experience with smart photography"},
            {"name": "Google Pixel 7 Pro",
                "description": "Flagship Pixel with advanced camera processing"},
            {"name": "OnePlus 12",
                "description": "Fast Android flagship with smooth performance"},
            {"name": "Xiaomi 14 Ultra",
                "description": "Camera-focused phone with flagship specs"},
            {"name": "Sony Xperia 1 V",
                "description": "Creator-oriented phone with pro-grade camera controls"},
            {"name": "ASUS ROG Phone 8",
                "description": "Gaming phone optimized for high FPS and cooling"},
            {"name": "Motorola Edge 50",
                "description": "Stylish mid-to-high range phone with clean software"},
            {"name": "Nothing Phone (2)",
                "description": "Minimal design with unique transparent interface"},
            {"name": "Huawei P60 Pro",
                "description": "High-end Huawei device with advanced imaging"},
            {"name": "OPPO Find X7 Ultra",
                "description": "Ultra imaging flagship with premium materials"},
        ],
        "Laptops": [
            {"name": "MacBook Pro 16-inch",
                "description": "Professional laptop with M3 chip"},
            {"name": "MacBook Air 13-inch (M2)",
                "description": "Lightweight daily driver with Apple Silicon"},
            {"name": "MacBook Pro 14-inch",
                "description": "Compact Pro laptop with powerful performance"},
            {"name": "Dell XPS 15",
                "description": "Premium Windows laptop with 4K display"},
            {"name": "Dell XPS 13",
                "description": "Compact XPS with strong build quality"},
            {"name": "ThinkPad X1 Carbon",
                "description": "Business laptop with excellent keyboard"},
            {"name": "ThinkPad T14",
                "description": "Practical ThinkPad with great productivity features"},
            {"name": "HP Spectre x360 14",
                "description": "Convertible laptop with premium design"},
            {"name": "HP Envy 15",
                "description": "Everyday premium laptop for work and entertainment"},
            {"name": "ASUS Zenbook 14",
                "description": "Thin-and-light laptop with sleek aluminum body"},
            {"name": "ASUS ROG Zephyrus G14",
                "description": "Gaming laptop balanced for portability"},
            {"name": "Microsoft Surface Laptop 6",
                "description": "Sleek Surface laptop optimized for productivity"},
            {"name": "Microsoft Surface Pro 11",
                "description": "2-in-1 tablet PC for flexible workflows"},
            {"name": "Acer Swift 3",
                "description": "Affordable lightweight laptop with good everyday specs"},
            {"name": "Lenovo Yoga 9i",
                "description": "Premium 2-in-1 with strong performance and build"},
        ],
        "Fiction": [
            {"name": "The Silent Patient",
                "description": "Psychological thriller novel"},
            {"name": "Where the Crawdads Sing",
                "description": "Coming-of-age story and mystery"},
            {"name": "The Midnight Library",
                "description": "Novel about life's infinite possibilities"}
            ,
            {"name": "1984",
                "description": "Dystopian novel about surveillance and control"},
            {"name": "Brave New World",
                "description": "Classic novel exploring a futuristic engineered society"},
            {"name": "The Great Gatsby",
                "description": "American classic about wealth, ambition, and tragedy"},
            {"name": "To Kill a Mockingbird",
                "description": "Pulitzer-winning novel about justice and morality"},
            {"name": "The Hobbit",
                "description": "Fantasy adventure with a reluctant hero"},
            {"name": "The Alchemist",
                "description": "Inspirational tale about pursuing personal legend"},
            {"name": "Catch-22",
                "description": "Satirical novel about war's absurd logic"},
            {"name": "The Da Vinci Code",
                "description": "Mystery thriller involving symbols and secrets"},
            {"name": "Dune",
                "description": "Epic science fiction about politics and ecology"},
            {"name": "Harry Potter and the Sorcerer's Stone",
                "description": "Magical beginning of a legendary wizard journey"},
            {"name": "The Girl with the Dragon Tattoo",
                "description": "Dark mystery thriller featuring a fearless investigator"},
            {"name": "Sapiens",
                "description": "Bestselling nonfiction exploring humanity's history"},
            {"name": "The Book Thief",
                "description": "Historical novel about books, memory, and resilience"}
        ],
        "Furniture": [
            {"name": "Modern Sofa",
                "description": "Comfortable 3-seater sofa in gray fabric"},
            {"name": "Oak Dining Table",
                "description": "Solid oak dining table for 6 people"},
            {"name": "Ergonomic Office Chair",
                "description": "Adjustable chair with lumbar support"},
            {"name": "Recliner Chair",
                "description": "Comfort recliner for relaxing evenings"},
            {"name": "Coffee Table",
                "description": "Simple coffee table for living room setup"},
            {"name": "Bookshelf",
                "description": "Multi-tier bookshelf for books and decor"},
            {"name": "Queen Bed Frame",
                "description": "Sturdy queen bed frame for comfortable sleep"},
            {"name": "Nightstand",
                "description": "Bedside table for lamps and essentials"},
            {"name": "Wardrobe Cabinet",
                "description": "Storage cabinet for clothes and household items"},
            {"name": "Standing Desk",
                "description": "Height-adjustable desk for healthier working posture"},
            {"name": "Bar Stools",
                "description": "Set of bar stools for kitchen counters or bars"},
            {"name": "Side Table",
                "description": "Compact side table for snacks and remote controls"},
            {"name": "Armchair",
                "description": "Supportive armchair for reading and lounging"},
            {"name": "Outdoor Patio Set",
                "description": "Patio seating set for outdoor gatherings"},
            {"name": "Chest of Drawers",
                "description": "Drawer chest for organized clothing and linens"}
        ]
    }

    # Products table in current monolith model has no category_id relation.
    # But we keep a non-persisted `category_hint` to generate product_categories links.

    # Fixed target number of sample products.
    target_count = 80

    # Build base (cat_name, prod_data) combinations from product_samples.
    base_entries = []
    for cat in categories:
        cat_name = cat.get("name")
        if cat_name in product_samples:
            for prod_data in product_samples[cat_name]:
                base_entries.append((cat_name, prod_data))

    # Insert base entries first.
    for cat_name, prod_data in base_entries:
        products.append(
            {
                "id": None,
                "name": prod_data["name"],
                "description": f"{prod_data['description']} (Category hint: {cat_name})",
                # Not persisted to DB; only used for generating product_categories links.
                "category_hint": cat_name,
                # Price in cents (e.g., $10.00 to $500.00)
                "price": random.randint(1000, 50000),
                "sku": f"PS-{str(uuid.uuid4()).split('-')[0].upper()}",
                "created_at": mysql_utc_now(),
                "updated_at": mysql_utc_now(),
            }
        )

    # If we still need more, generate variants based on base entries.
    while len(products) < target_count and base_entries:
        cat_name, prod_data = random.choice(base_entries)
        variant_idx = len(products) + 1
        products.append(
            {
                "id": None,
                "name": f"{prod_data['name']} (Variant {variant_idx})",
                "description": f"{prod_data['description']} - variant {variant_idx} (Category hint: {cat_name})",
                # Not persisted to DB; only used for generating product_categories links.
                "category_hint": cat_name,
                "price": random.randint(500, 80000),
                "sku": f"PSV-{str(uuid.uuid4()).split('-')[0].upper()}",
                "created_at": mysql_utc_now(),
                "updated_at": mysql_utc_now(),
            }
        )

    return products


def generate_product_categories_links(products, categories, min_categories=1, max_categories=3):
    """Generate sample product_categories links (many-to-many)"""
    product_categories = []

    # Name -> id (ids are expected to be filled after insertion)
    name_to_id = {c["name"]: c["id"] for c in categories if c.get("id") is not None}
    all_categories = [c for c in categories if c.get("id") is not None]

    for product in products:
        if product.get("id") is None:
            # Should not happen if called after insert_products()
            continue

        if all_categories:
            base_category_id = None
            hint = product.get("category_hint")
            if hint and hint in name_to_id:
                base_category_id = name_to_id[hint]

            if base_category_id is None:
                base_category_id = random.choice(all_categories)["id"]

            target_count = random.randint(min_categories, max_categories)
            assigned_ids = {base_category_id}

            # Fill remaining categories randomly (avoid duplicates)
            if len(assigned_ids) < target_count:
                remaining = [c for c in all_categories if c["id"] not in assigned_ids]
                if target_count - len(assigned_ids) > len(remaining):
                    target_count = len(assigned_ids) + len(remaining)
                for c in random.sample(remaining, target_count - len(assigned_ids)):
                    assigned_ids.add(c["id"])

            for category_id in assigned_ids:
                product_categories.append(
                    {
                        "product_id": product["id"],
                        "category_id": category_id,
                    }
                )

    return product_categories


def generate_inventory(products):
    """Generate sample inventory records"""
    inventory = []

    for product in products:
        inv_record = {
            "product_id": product["id"],
            "quantity": random.randint(0, 100),
            # Reserved items
            "reserved": random.randint(0, min(10, random.randint(0, 100))),
            "threshold": 10,  # Low stock threshold
            "created_at": mysql_utc_now(),
            "updated_at": mysql_utc_now(),
        }
        # Calculate status based on quantity and threshold
        if inv_record["quantity"] <= 0:
            inv_record["status"] = "outofstock"
        elif inv_record["quantity"] <= inv_record["threshold"]:
            inv_record["status"] = "lowstock"
        else:
            inv_record["status"] = "instock"

        inventory.append(inv_record)

    return inventory


def connect_to_database():
    """Connect to the MySQL database"""
    try:
        connection = pymysql.connect(
            host=os.getenv('DB_HOST', 'localhost'),
            port=int(os.getenv('DB_PORT', 3306)),
            user=os.getenv('DB_USER', 'root'),
            password=os.getenv('DB_PASSWORD', '123456'),
            database=os.getenv('DB_NAME', 'eshop_db'),
            charset='utf8mb4',
            cursorclass=pymysql.cursors.DictCursor
        )
        print("Successfully connected to the database!")
        return connection
    except Exception as e:
        print(f"Error connecting to database: {e}")
        sys.exit(1)


def insert_categories(connection, categories):
    """Insert categories into the database"""
    with connection.cursor() as cursor:
        name_to_id = {}

        root_categories = [c for c in categories if c.get("parent_id") is None and not c.get("parent_name")]
        child_categories = [c for c in categories if c.get("parent_name")]

        for category in root_categories:
            sql = """
            INSERT INTO categories (name, description, parent_id, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(
                    sql,
                    (
                        category["name"],
                        category["description"],
                        category["parent_id"],
                        category["created_at"],
                        category["updated_at"],
                    ),
                )
                category["id"] = cursor.lastrowid
                name_to_id[category["name"]] = category["id"]
            except Exception as e:
                print(f"Error inserting category {category['name']}: {e}")

        for category in child_categories:
            category["parent_id"] = name_to_id.get(category["parent_name"])
            sql = """
            INSERT INTO categories (name, description, parent_id, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(
                    sql,
                    (
                        category["name"],
                        category["description"],
                        category["parent_id"],
                        category["created_at"],
                        category["updated_at"],
                    ),
                )
                category["id"] = cursor.lastrowid
                name_to_id[category["name"]] = category["id"]
            except Exception as e:
                print(f"Error inserting category {category['name']}: {e}")
        connection.commit()
    print(f"Inserted {len(categories)} categories into the database")


def insert_products(connection, products):
    """Insert products into the database"""
    with connection.cursor() as cursor:
        for product in products:
            sql = """
            INSERT INTO products (name, description, price, sku, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(
                    sql,
                    (
                        product["name"],
                        product["description"],
                        product["price"],
                        product["sku"],
                        product["created_at"],
                        product["updated_at"],
                    ),
                )
                product["id"] = cursor.lastrowid
            except Exception as e:
                print(f"Error inserting product {product['name']}: {e}")
        connection.commit()
    print(f"Inserted {len(products)} products into the database")


def insert_inventory(connection, inventory):
    """Insert inventory records into the database"""
    with connection.cursor() as cursor:
        for inv in inventory:
            sql = """
            INSERT INTO inventories (product_id, quantity, status, reserved, threshold, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(sql, (
                    inv["product_id"],
                    inv["quantity"],
                    inv["status"],
                    inv["reserved"],
                    inv["threshold"],
                    inv["created_at"],
                    inv["updated_at"]
                ))
            except Exception as e:
                print(
                    f"Error inserting inventory for product {inv['product_id']}: {e}")
        connection.commit()
    print(f"Inserted {len(inventory)} inventory records into the database")


def insert_product_categories(connection, product_categories):
    """Insert sample product_categories associations"""
    with connection.cursor() as cursor:
        for pc in product_categories:
            sql = """
            INSERT IGNORE INTO product_categories (product_id, category_id)
            VALUES (%s, %s)
            """
            try:
                cursor.execute(
                    sql,
                    (
                        pc["product_id"],
                        pc["category_id"],
                    ),
                )
            except Exception as e:
                # Likely duplicate from reruns; skip noisy failures.
                print(f"Error inserting product_categories link: {e}")
        connection.commit()


def save_to_json(data, filename):
    """Save data to a JSON file"""
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=2, ensure_ascii=False)
    print(f"Saved {len(data)} items to {filename}")


def main():
    print("Generating and inserting sample data for eshop-microservices...")

    # Connect to database
    connection = connect_to_database()

    # Generate data
    categories = generate_categories()
    products = generate_products(categories)

    # Insert data into database
    insert_categories(connection, categories)
    insert_products(connection, products)

    # Insert many-to-many product-category relationships
    product_categories = generate_product_categories_links(products, categories)
    insert_product_categories(connection, product_categories)

    inventory = generate_inventory(products)
    insert_inventory(connection, inventory)

    # Close the database connection
    connection.close()

    # Optionally save to JSON files as well
    output_dir = os.path.join(os.path.dirname(__file__), "..", "sample_data")
    os.makedirs(output_dir, exist_ok=True)

    save_to_json(categories, os.path.join(output_dir, "categories.json"))
    save_to_json(products, os.path.join(output_dir, "products.json"))
    save_to_json(inventory, os.path.join(output_dir, "inventory.json"))

    # Also create a combined file
    combined_data = {
        "categories": categories,
        "products": products,
        "product_categories": product_categories,
        "inventory": inventory,
        "generated_at": mysql_utc_now(),
    }
    save_to_json(combined_data, os.path.join(output_dir, "sample_data.json"))

    print("\nSample data generation and insertion complete!")
    print(f"Inserted:")
    print(f"- {len(categories)} categories")
    print(f"- {len(products)} products")
    print(f"- {len(product_categories)} product-category associations")
    print(f"- {len(inventory)} inventory records")
    print(f"\nFiles saved to: {output_dir}")


if __name__ == "__main__":
    main()
