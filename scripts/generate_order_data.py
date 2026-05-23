"""
Generate sample order data for the eshop-monolith project.
This script connects directly to the database and inserts sample order data.
"""

import json
import uuid
from datetime import datetime, timedelta
import random
import os
import pymysql  # Requires: pip install PyMySQL
import sys


def mysql_utc_now():
    return datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S")


def generate_orders(products, min_orders=50, max_orders=100):
    """Generate sample orders"""
    orders = []
    order_items = []
    
    # Order statuses
    order_statuses = ["pending", "paid", "shipped", "delivered", "cancelled"]
    
    # Generate random customer IDs (using numeric IDs directly)
    customer_ids = []
    for i in range(1, 21):  # 20 unique customers
        customer_ids.append(i)
    
    # Generate orders
    order_count = random.randint(min_orders, max_orders)
    order_id = 1
    
    for _ in range(order_count):
        # Random customer
        customer_id = random.choice(customer_ids)
        
        # Random order status
        status = random.choice(order_statuses)
        
        # Random order items (1-3 items per order)
        item_count = random.randint(1, 3)
        selected_products = random.sample(products, item_count)
        
        total_amount = 0
        order_items_batch = []
        
        for product in selected_products:
            # Random quantity (1-3)
            quantity = random.randint(1, 3)
            unit_price = product["price"]
            amount = unit_price * quantity
            total_amount += amount
            
            # Create order item
            order_item = {
                "id": None,
                "order_id": order_id,
                "product_id": product["id"],
                "quantity": quantity,
                "unit_price": unit_price,
                "amount": amount
            }
            order_items_batch.append(order_item)
        
        # Random created time (within last 30 days)
        created_at = (datetime.utcnow() - timedelta(days=random.randint(0, 30))).strftime("%Y-%m-%d %H:%M:%S")
        updated_at = created_at
        
        # Create order
        order = {
            "id": order_id,
            "customer_id": customer_id,
            "total_amount": total_amount,
            "currency": "CNY",
            "status": status,
            "created_at": created_at,
            "updated_at": updated_at
        }
        
        orders.append(order)
        order_items.extend(order_items_batch)
        order_id += 1
    
    return orders, order_items


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


def load_products_from_json():
    """Load products from JSON file"""
    products_file = os.path.join(os.path.dirname(__file__), "..", "sample_data", "products.json")
    try:
        with open(products_file, 'r', encoding='utf-8') as f:
            products = json.load(f)
        print(f"Loaded {len(products)} products from JSON file")
        return products
    except Exception as e:
        print(f"Error loading products from JSON: {e}")
        sys.exit(1)


def insert_orders(connection, orders):
    """Insert orders into the database"""
    with connection.cursor() as cursor:
        for order in orders:
            sql = """
            INSERT INTO orders (id, customer_id, total_amount, currency, status, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(
                    sql,
                    (
                        order["id"],
                        order["customer_id"],
                        order["total_amount"],
                        order["currency"],
                        order["status"],
                        order["created_at"],
                        order["updated_at"]
                    ),
                )
            except Exception as e:
                print(f"Error inserting order {order['id']}: {e}")
        connection.commit()
    print(f"Inserted {len(orders)} orders into the database")


def insert_order_items(connection, order_items):
    """Insert order items into the database"""
    with connection.cursor() as cursor:
        for item in order_items:
            sql = """
            INSERT INTO order_items (order_id, product_id, quantity, unit_price, amount)
            VALUES (%s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(
                    sql,
                    (
                        item["order_id"],
                        item["product_id"],
                        item["quantity"],
                        item["unit_price"],
                        item["amount"]
                    ),
                )
            except Exception as e:
                print(f"Error inserting order item: {e}")
        connection.commit()
    print(f"Inserted {len(order_items)} order items into the database")


def save_to_json(data, filename):
    """Save data to a JSON file"""
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=2, ensure_ascii=False)
    print(f"Saved {len(data)} items to {filename}")


def main():
    print("Generating and inserting sample order data for eshop-monolith...")

    # Connect to database
    connection = connect_to_database()

    # Load products from JSON
    products = load_products_from_json()

    # Generate order data
    orders, order_items = generate_orders(products)

    # Insert data into database
    insert_orders(connection, orders)
    insert_order_items(connection, order_items)

    # Close the database connection
    connection.close()

    # Save to JSON files
    output_dir = os.path.join(os.path.dirname(__file__), "..", "sample_data")
    os.makedirs(output_dir, exist_ok=True)

    save_to_json(orders, os.path.join(output_dir, "orders.json"))
    save_to_json(order_items, os.path.join(output_dir, "order_items.json"))

    # Also create a combined order data file
    combined_order_data = {
        "orders": orders,
        "order_items": order_items,
        "generated_at": mysql_utc_now(),
    }
    save_to_json(combined_order_data, os.path.join(output_dir, "order_data.json"))

    print("\nSample order data generation and insertion complete!")
    print(f"Inserted:")
    print(f"- {len(orders)} orders")
    print(f"- {len(order_items)} order items")
    print(f"\nFiles saved to: {output_dir}")


if __name__ == "__main__":
    main()
