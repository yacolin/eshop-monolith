"""
Generate sample data for categories, products, and inventory for the eshop-microservices project.
This script connects directly to the database and inserts sample data.
"""

import json
import uuid
from datetime import datetime, timezone
import random
import os
import time
import pymysql  # Requires: pip install PyMySQL
import sys


def mysql_utc_now():
    return datetime.now(timezone.utc).strftime("%Y-%m-%d %H:%M:%S")


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
        "Home & Garden": ["Furniture", "Decor", "Kitchen", "Bedding", "Tools", "Outdoor"],
        "Sports & Outdoors": ["Fitness Equipment", "Camping Gear", "Team Sports", "Cycling", "Water Sports"],
        "Beauty & Health": ["Skincare", "Makeup", "Hair Care", "Fragrance", "Vitamins & Supplements"],
        "Toys & Games": ["Board Games", "Action Figures", "Building Sets", "Puzzles", "Video Games"],
        "Automotive": ["Car Care", "Interior Accessories", "Exterior Accessories", "Automotive Tools", "Car Electronics"],
        "Food & Grocery": ["Snacks", "Beverages", "Pantry Staples", "Organic Foods", "Gourmet Foods"],
        "Office Supplies": ["Writing Instruments", "Paper Products", "Desk Accessories", "Office Storage", "Office Technology"],
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
            {"name": "iPhone 15 Pro", "description": "Latest iPhone with advanced camera system"},
            {"name": "iPhone 14 Pro", "description": "Previous-generation Pro phone with premium display"},
            {"name": "iPhone 15", "description": "Mainstream iPhone with great everyday performance"},
            {"name": "Samsung Galaxy S24", "description": "Android flagship with AI capabilities"},
            {"name": "Samsung Galaxy S23 Ultra", "description": "Ultra model with zoom camera and S Pen"},
            {"name": "Google Pixel 8", "description": "Pure Android experience with smart photography"},
            {"name": "OnePlus 12", "description": "Fast Android flagship with smooth performance"},
            {"name": "Xiaomi 14 Ultra", "description": "Camera-focused phone with flagship specs"},
            {"name": "ASUS ROG Phone 8", "description": "Gaming phone optimized for high FPS and cooling"},
            {"name": "Nothing Phone (2)", "description": "Minimal design with unique transparent interface"},
        ],
        "Laptops": [
            {"name": "MacBook Pro 16-inch", "description": "Professional laptop with M3 chip"},
            {"name": "MacBook Air 13-inch", "description": "Lightweight daily driver with Apple Silicon"},
            {"name": "MacBook Pro 14-inch", "description": "Compact Pro laptop with powerful performance"},
            {"name": "Dell XPS 15", "description": "Premium Windows laptop with 4K display"},
            {"name": "Dell XPS 13", "description": "Compact XPS with strong build quality"},
            {"name": "ThinkPad X1 Carbon", "description": "Business laptop with excellent keyboard"},
            {"name": "HP Spectre x360 14", "description": "Convertible laptop with premium design"},
            {"name": "ASUS Zenbook 14", "description": "Thin-and-light laptop with sleek aluminum body"},
            {"name": "ASUS ROG Zephyrus G14", "description": "Gaming laptop balanced for portability"},
            {"name": "Microsoft Surface Laptop 6", "description": "Sleek Surface laptop optimized for productivity"},
        ],
        "Tablets": [
            {"name": "iPad Pro 12.9", "description": "Apple flagship tablet with M2 chip"},
            {"name": "iPad Air 11", "description": "Mid-range Apple tablet with great performance"},
            {"name": "iPad 10th Gen", "description": "Entry-level iPad for everyday use"},
            {"name": "iPad mini 7", "description": "Compact tablet for reading and note-taking"},
            {"name": "Samsung Galaxy Tab S9", "description": "Android flagship tablet with AMOLED display"},
            {"name": "Samsung Galaxy Tab A9", "description": "Budget-friendly tablet for entertainment"},
            {"name": "Microsoft Surface Pro 11", "description": "2-in-1 tablet PC for productivity"},
            {"name": "Amazon Fire HD 10", "description": "Affordable tablet for media consumption"},
            {"name": "Lenovo Tab P12", "description": "Large-screen Android tablet for work"},
            {"name": "Xiaomi Pad 6", "description": "Value tablet with high-resolution display"},
        ],
        "TVs": [
            {"name": "Samsung Neo QLED 65", "description": "65-inch QLED smart TV with quantum HDR"},
            {"name": "LG OLED C3 55", "description": "55-inch OLED TV with perfect blacks"},
            {"name": "Sony Bravia XR 65", "description": "65-inch LED TV with cognitive processor"},
            {"name": "TCL QLED 75", "description": "75-inch QLED TV with great value"},
            {"name": "Hisense U8H 65", "description": "65-inch ULED TV with mini-LED backlight"},
            {"name": "Samsung Frame 55", "description": "55-inch lifestyle TV that displays art"},
            {"name": "LG NanoCell 86", "description": "86-inch large-screen smart TV"},
            {"name": "Vizio OLED 65", "description": "65-inch OLED with Dolby Vision"},
            {"name": "Sony X90L 75", "description": "75-inch full-array LED 4K TV"},
            {"name": "Panasonic MX950 55", "description": "55-inch mid-range LED TV"},
        ],
        "Cameras": [
            {"name": "Sony A7 IV", "description": "Full-frame mirrorless camera with 33MP sensor"},
            {"name": "Canon EOS R6 Mark II", "description": "Full-frame mirrorless with fast AF"},
            {"name": "Nikon Z8", "description": "High-resolution full-frame mirrorless camera"},
            {"name": "Fujifilm X-T5", "description": "APS-C mirrorless with retro design"},
            {"name": "Sony ZV-E1", "description": "Compact vlogging camera with full-frame sensor"},
            {"name": "Canon EOS R100", "description": "Entry-level mirrorless camera"},
            {"name": "DJI Osmo Pocket 3", "description": "Compact gimbal camera for vlogging"},
            {"name": "GoPro Hero 12", "description": "Action camera with stabilization"},
            {"name": "Fujifilm Instax Mini 12", "description": "Instant film camera for prints"},
            {"name": "Panasonic Lumix S5 II", "description": "Full-frame hybrid camera"},
        ],
        "Audio Equipment": [
            {"name": "Sony WH-1000XM5", "description": "Wireless noise-cancelling headphones"},
            {"name": "Apple AirPods Pro 2", "description": "Premium wireless earbuds with ANC"},
            {"name": "Bose QuietComfort Ultra", "description": "Over-ear headphones with spatial audio"},
            {"name": "Sennheiser Momentum 4", "description": "Wireless headphones with audiophile sound"},
            {"name": "Marshall Stanmore III", "description": "Bluetooth speaker with iconic design"},
            {"name": "Sonos Era 300", "description": "Multi-room smart speaker with spatial audio"},
            {"name": "JBL Flip 6", "description": "Portable Bluetooth speaker with deep bass"},
            {"name": "Audio-Technica ATH-M50x", "description": "Professional studio monitor headphones"},
            {"name": "Samsung Galaxy Buds 2 Pro", "description": "Wireless earbuds with 24-bit audio"},
            {"name": "Google Pixel Buds Pro", "description": "Wireless earbuds with real-time translation"},
        ],
        "Men's Clothing": [
            {"name": "Classic Fit Oxford Shirt", "description": "Button-down Oxford shirt in premium cotton"},
            {"name": "Slim Fit Chinos", "description": "Stretch cotton chinos for smart casual wear"},
            {"name": "Merino Wool Sweater", "description": "Lightweight sweater for layering in cold months"},
            {"name": "Linen Blazer", "description": "Unstructured blazer for summer events"},
            {"name": "Denim Jacket", "description": "Classic blue denim jacket with button closure"},
            {"name": "Tailored Wool Trousers", "description": "Wool dress pants for office or formal occasions"},
            {"name": "Polo Shirt", "description": "Cotton pique polo with embroidered logo"},
            {"name": "Leather Chelsea Boots", "description": "Sleek leather boots with elastic side panels"},
            {"name": "Cashmere Scarf", "description": "Luxurious cashmere scarf in heather grey"},
            {"name": "Performance Joggers", "description": "Moisture-wicking joggers for active lifestyle"},
        ],
        "Women's Clothing": [
            {"name": "Wrap Dress", "description": "Flattering wrap dress in printed viscose"},
            {"name": "High-Waist Trousers", "description": "Wide-leg trousers with crepe fabric"},
            {"name": "Cashmere Cardigan", "description": "Soft cashmere cardigan with pearl buttons"},
            {"name": "Silk Blouse", "description": "Elegant silk blouse with French cuffs"},
            {"name": "Denim Skirt", "description": "A-line denim skirt with button-front closure"},
            {"name": "Trench Coat", "description": "Classic double-breasted trench coat"},
            {"name": "Leather Handbag", "description": "Genuine leather tote with magnetic closure"},
            {"name": "Yoga Leggings", "description": "High-waist leggings with compression fabric"},
            {"name": "Linen Jumpsuit", "description": "Casual linen jumpsuit with tie waist"},
            {"name": "Wool Blend Coat", "description": "Long-length coat for cold winter days"},
        ],
        "Kids' Clothing": [
            {"name": "Rainbow T-Shirt", "description": "Colorful t-shirt for everyday adventures"},
            {"name": "Denim Dungarees", "description": "Adjustable dungarees with fun patches"},
            {"name": "Fleece Hoodie", "description": "Soft fleece hoodie with animal ears"},
            {"name": "School Uniform Polo", "description": "Stain-resistant polo shirt for school"},
            {"name": "Printed Leggings", "description": "Fun pattern leggings with stretch fabric"},
            {"name": "Puffer Jacket", "description": "Lightweight puffer jacket with detachable hood"},
            {"name": "Cotton Pajama Set", "description": "Soft organic cotton pajama set"},
            {"name": "Character Socks Pack", "description": "Pack of 5 colorful character socks"},
            {"name": "Waterproof Raincoat", "description": "Bright raincoat with reflective strips"},
            {"name": "Kid's Sneakers", "description": "Lightweight sneakers with easy lace system"},
        ],
        "Shoes": [
            {"name": "Nike Air Max 90", "description": "Classic sneaker with visible air cushioning"},
            {"name": "Adidas Ultraboost Light", "description": "Running shoe with responsive boost foam"},
            {"name": "Clarks Desert Boot", "description": "Timeless desert boot in beeswax leather"},
            {"name": "Dr. Martens 1460", "description": "Iconic 8-eye leather boot with air sole"},
            {"name": "New Balance 990v6", "description": "Premium stability running shoe"},
            {"name": "Vans Old Skool", "description": "Classic skate shoe with side stripe"},
            {"name": "Timberland 6-Inch Boot", "description": "Waterproof leather boot with padded collar"},
            {"name": "Birkenstock Arizona", "description": "Two-strap sandal with contoured footbed"},
            {"name": "Converse Chuck Taylor", "description": "Canvas high-top sneaker in black"},
            {"name": "Salomon Speedcross 5", "description": "Trail running shoe with aggressive grip"},
        ],
        "Accessories": [
            {"name": "Ray-Ban Aviator", "description": "Classic aviator sunglasses with gold frame"},
            {"name": "Daniel Wellington Watch", "description": "Minimalist leather strap watch"},
            {"name": "Silk Twilly Scarf", "description": "Printed silk scarf for elegant styling"},
            {"name": "Leather Belt", "description": "Genuine leather belt with brushed buckle"},
            {"name": "Baseball Cap", "description": "Cotton twill cap with embroidered logo"},
            {"name": "Leather Wallet", "description": "Slim bifold wallet with RFID protection"},
            {"name": "Silver Chain Necklace", "description": "Stainless steel chain necklace"},
            {"name": "Wool Beanie", "description": "Ribbed knit beanie for cold weather"},
            {"name": "Canvas Backpack", "description": "Vintage-style canvas backpack with laptop sleeve"},
            {"name": "Leather Gloves", "description": "Cashmere-lined leather gloves"},
        ],
        "Fiction": [
            {"name": "The Silent Patient", "description": "Psychological thriller novel"},
            {"name": "Where the Crawdads Sing", "description": "Coming-of-age story and mystery"},
            {"name": "The Midnight Library", "description": "Novel about life's infinite possibilities"},
            {"name": "1984", "description": "Dystopian novel about surveillance and control"},
            {"name": "Brave New World", "description": "Classic novel exploring a futuristic engineered society"},
            {"name": "The Great Gatsby", "description": "American classic about wealth, ambition, and tragedy"},
            {"name": "To Kill a Mockingbird", "description": "Pulitzer-winning novel about justice and morality"},
            {"name": "The Hobbit", "description": "Fantasy adventure with a reluctant hero"},
            {"name": "Dune", "description": "Epic science fiction about politics and ecology"},
            {"name": "The Book Thief", "description": "Historical novel about books, memory, and resilience"},
        ],
        "Non-Fiction": [
            {"name": "Sapiens", "description": "Brief history of humankind through the ages"},
            {"name": "Atomic Habits", "description": "Practical guide to building good habits"},
            {"name": "Thinking Fast and Slow", "description": "Exploration of human decision-making"},
            {"name": "Educated", "description": "Memoir of a woman who leaves her survivalist family"},
            {"name": "The Power of Habit", "description": "Science behind habit formation in daily life"},
            {"name": "Outliers", "description": "Study of what makes high achievers different"},
            {"name": "Freakonomics", "description": "Economics applied to unconventional topics"},
            {"name": "Born a Crime", "description": "Trevor Noah's memoir of growing up in South Africa"},
            {"name": "In Cold Blood", "description": "True crime classic about a Kansas murder case"},
            {"name": "The Art of War", "description": "Ancient Chinese military treatise on strategy"},
        ],
        "Science": [
            {"name": "A Brief History of Time", "description": "Hawking's exploration of time and the universe"},
            {"name": "The Selfish Gene", "description": "Gene-centered view of evolution by Dawkins"},
            {"name": "Cosmos", "description": "Journey through space and time by Carl Sagan"},
            {"name": "The Gene", "description": "Intimate history of genetics and its future"},
            {"name": "Silent Spring", "description": "Groundbreaking work on environmental science"},
            {"name": "The Double Helix", "description": "Personal account of DNA structure discovery"},
            {"name": "Guns Germs and Steel", "description": "Geography's role in shaping human civilization"},
            {"name": "The Elegant Universe", "description": "Superstring theory explained for non-scientists"},
            {"name": "The Immortal Life of Henrietta Lacks", "description": "Story of HeLa cells and medical ethics"},
            {"name": "The Origin of Species", "description": "Darwin's foundational work on natural selection"},
        ],
        "Technology": [
            {"name": "Clean Code", "description": "Principles for writing maintainable code"},
            {"name": "Designing Data-Intensive Applications", "description": "Fundamentals of distributed systems"},
            {"name": "The Pragmatic Programmer", "description": "Timeless tips for software craftsmanship"},
            {"name": "Introduction to Algorithms", "description": "Comprehensive guide to algorithm design"},
            {"name": "Structure and Interpretation", "description": "Foundational CS textbook using Scheme"},
            {"name": "The Mythical Man-Month", "description": "Classic essays on software engineering"},
            {"name": "Code Complete", "description": "Practical guide to software construction"},
            {"name": "Artificial Intelligence", "description": "Modern approach to AI by Russell and Norvig"},
            {"name": "Operating Systems Three Easy Pieces", "description": "Introduction to operating system concepts"},
            {"name": "Computer Networking", "description": "Top-down approach to networking concepts"},
        ],
        "Biography": [
            {"name": "Steve Jobs", "description": "Biography of Apple's visionary co-founder"},
            {"name": "Becoming", "description": "Michelle Obama's memoir of her life and career"},
            {"name": "Elon Musk", "description": "Story of Tesla and SpaceX's controversial CEO"},
            {"name": "The Diary of a Young Girl", "description": "Anne Frank's wartime diary"},
            {"name": "Long Walk to Freedom", "description": "Nelson Mandela's autobiography of struggle"},
            {"name": "Leonardo da Vinci", "description": "Biography of the ultimate Renaissance genius"},
            {"name": "Einstein", "description": "Life and work of the iconic physicist"},
            {"name": "The Story of My Experiments with Truth", "description": "Gandhi's autobiography of non-violent struggle"},
            {"name": "Churchill", "description": "Biography of Britain's wartime prime minister"},
            {"name": "I Am Malala", "description": "Malala's fight for girls' education"},
        ],
        "Cookbooks": [
            {"name": "Salt Fat Acid Heat", "description": "Mastering the four elements of good cooking"},
            {"name": "The Joy of Cooking", "description": "Comprehensive home cooking reference"},
            {"name": "Plenty", "description": "Ottolenghi's vibrant vegetable cookbook"},
            {"name": "Mastering the Art of French Cooking", "description": "Classic French recipes by Julia Child"},
            {"name": "The Food Lab", "description": "Science-based home cooking techniques"},
            {"name": "Korean American", "description": "Fusion Korean-American recipes by Eric Kim"},
            {"name": "Ottolenghi Simple", "description": "Middle Eastern inspired easy recipes"},
            {"name": "Half Baked Harvest", "description": "Super tasty everyday comfort foods"},
            {"name": "The Wok", "description": "Essential wok cooking techniques and recipes"},
            {"name": "Flour Water Salt Yeast", "description": "Fundamentals of artisan bread baking"},
        ],
        "Furniture": [
            {"name": "Modern Sofa", "description": "Comfortable 3-seater sofa in gray fabric"},
            {"name": "Oak Dining Table", "description": "Solid oak dining table for 6 people"},
            {"name": "Ergonomic Office Chair", "description": "Adjustable chair with lumbar support"},
            {"name": "Coffee Table", "description": "Minimalist coffee table with storage shelf"},
            {"name": "Bookshelf", "description": "Multi-tier bookshelf for books and decor"},
            {"name": "Queen Bed Frame", "description": "Sturdy queen bed frame with upholstered headboard"},
            {"name": "Nightstand", "description": "Bedside table with drawer and open shelf"},
            {"name": "Standing Desk", "description": "Height-adjustable desk for healthier working posture"},
            {"name": "Wardrobe Cabinet", "description": "Spacious wardrobe with sliding doors"},
            {"name": "Recliner Chair", "description": "Power recliner with massage and heat functions"},
        ],
        "Decor": [
            {"name": "Ceramic Vase Set", "description": "Set of 3 matte ceramic vases in natural tones"},
            {"name": "Framed Wall Art", "description": "Abstract canvas print with wooden frame"},
            {"name": "Scented Candle Collection", "description": "Set of 4 soy wax candles with varied scents"},
            {"name": "Decorative Throw Pillow", "description": "Linen pillow with geometric embroidery"},
            {"name": "Indoor Potted Plant", "description": "Faux fiddle leaf fig in woven basket"},
            {"name": "Table Lamp", "description": "Ceramic table lamp with linen shade"},
            {"name": "Wall Mirror", "description": "Round wall mirror with gold metal frame"},
            {"name": "Rug 5x7", "description": "Handwoven wool rug with traditional pattern"},
            {"name": "Curtain Set", "description": "Blackout curtains in neutral beige, pair"},
            {"name": "Photo Frame Collage", "description": "Set of 6 magnetic photo frames"},
        ],
        "Kitchen": [
            {"name": "Stainless Steel Cookware Set", "description": "10-piece set with tri-ply construction"},
            {"name": "Chef's Knife 8-Inch", "description": "Forged high-carbon stainless steel knife"},
            {"name": "Cast Iron Dutch Oven", "description": "Enameled cast iron pot for slow cooking"},
            {"name": "Non-Stick Frying Pan", "description": "12-inch pan with ceramic non-stick coating"},
            {"name": "Instant Pot Duo", "description": "7-in-1 electric pressure cooker"},
            {"name": "Stand Mixer", "description": "5-quart tilt-head mixer with attachments"},
            {"name": "Cutting Board Set", "description": "Bamboo cutting board set with juice groove"},
            {"name": "KitchenAid Blender", "description": "High-performance blender with BPA-free jar"},
            {"name": "Measuring Cups Set", "description": "Stainless steel measuring cups with magnetic handle"},
            {"name": "Food Storage Container Set", "description": "Glass container set with airtight lids"},
        ],
        "Bedding": [
            {"name": "Egyptian Cotton Sheet Set", "description": "1000-thread count queen sheet set"},
            {"name": "Down Comforter", "description": "All-season down comforter with baffle box"},
            {"name": "Memory Foam Pillow", "description": "Contour pillow with cooling gel layer"},
            {"name": "Weighted Blanket", "description": "15-pound weighted blanket with glass beads"},
            {"name": "Bamboo Mattress Protector", "description": "Breathable waterproof mattress cover"},
            {"name": "Flannel Sheet Set", "description": "Brushed flannel sheets for winter warmth"},
            {"name": "European Pillowcase Set", "description": "Luxury sateen pillowcases set of 2"},
            {"name": "Wool Throw Blanket", "description": "Merino wool throw in herringbone pattern"},
            {"name": "Silk Sleep Mask", "description": "Mulberry silk eye mask with adjustable strap"},
            {"name": "Duvet Cover Set", "description": "Cotton duvet cover with hidden button closure"},
        ],
        "Tools": [
            {"name": "Cordless Drill Kit", "description": "20V brushless drill with battery pack"},
            {"name": "Tool Set 150-Piece", "description": "General household tool kit in carrying case"},
            {"name": "Hammer 16oz", "description": "Curved claw hammer with shock reduction grip"},
            {"name": "Screwdriver Set", "description": "Precision screwdriver set with magnetic tips"},
            {"name": "Tape Measure 25ft", "description": "Self-locking tape measure with belt clip"},
            {"name": "Adjustable Wrench", "description": "Chrome vanadium wrench with cushioned grip"},
            {"name": "Level 48-Inch", "description": "Box beam level with rare earth magnets"},
            {"name": "Utility Knife", "description": "Retractable blade utility knife with storage"},
            {"name": "Work Gloves", "description": "Leather palm work gloves with knuckle protection"},
            {"name": "Flashlight LED Rechargeable", "description": "1000-lumen rechargeable flashlight"},
        ],
        "Outdoor": [
            {"name": "Gas Grill 4-Burner", "description": "Stainless steel propane grill with side burner"},
            {"name": "Patio Umbrella 10ft", "description": "Offset cantilever umbrella with UV coating"},
            {"name": "Adirondack Chair", "description": "Weather-resistant Adirondack chair in teak"},
            {"name": "Outdoor String Lights", "description": "48ft LED string lights for patio"},
            {"name": "Fire Pit Table", "description": "Propane fire pit with lava rock filler"},
            {"name": "Hammock with Stand", "description": "Cotton hammock with folding steel stand"},
            {"name": "Garden Hose 100ft", "description": "Expandable garden hose with spray nozzle"},
            {"name": "Planter Box Set", "description": "Set of 3 cedar planter boxes for deck"},
            {"name": "Bird Feeder", "description": "Clear tube bird feeder with metal ports"},
            {"name": "Outdoor Cushion Set", "description": "Patio chair cushion set in fade-resistant fabric"},
        ],
        "Fitness Equipment": [
            {"name": "Adjustable Dumbbell Set", "description": "Space-saving dumbbells from 5-52.5 lbs"},
            {"name": "Yoga Mat Premium", "description": "Non-slip exercise mat with alignment lines"},
            {"name": "Treadmill Folding", "description": "Compact treadmill with incline adjustment"},
            {"name": "Resistance Bands Set", "description": "Set of 5 bands with varying tension levels"},
            {"name": "Kettlebell 35lbs", "description": "Cast iron kettlebell with flat base"},
            {"name": "Jump Rope Speed", "description": "Ball bearing jump rope for cardio workouts"},
            {"name": "Foam Roller", "description": "High-density foam roller for muscle recovery"},
            {"name": "Pull Up Bar", "description": "Doorway pull-up bar with foam grips"},
            {"name": "Exercise Bike Stationary", "description": "Indoor cycling bike with magnetic resistance"},
            {"name": "Ab Roller Wheel", "description": "Ab wheel with knee pad for core training"},
        ],
        "Camping Gear": [
            {"name": "Camping Tent 4-Person", "description": "Waterproof tent with easy setup design"},
            {"name": "Sleeping Bag 20F", "description": "Mummy sleeping bag with compression sack"},
            {"name": "Camping Stove", "description": "Portable butane stove with wind shield"},
            {"name": "Cooler 45-Quart", "description": "Rotomolded cooler with bear-proof lock"},
            {"name": "Camping Lantern", "description": "LED lantern with 360-degree lighting"},
            {"name": "Hiking Backpack 50L", "description": "Internal frame pack with rain cover"},
            {"name": "Camping Chair", "description": "Padded camp chair with cup holder"},
            {"name": "Portable Water Filter", "description": "Backpacking filter for safe drinking water"},
            {"name": "Hammock Camping", "description": "Double parachute hammock with tree straps"},
            {"name": "Camping Cookware Set", "description": "Mess kit with pot pan and utensils"},
        ],
        "Team Sports": [
            {"name": "Soccer Ball", "description": "FIFA-approved match ball size 5"},
            {"name": "Basketball Indoor Outdoor", "description": "Composite leather basketball full size"},
            {"name": "Volleyball Set", "description": "Complete set with net and ball"},
            {"name": "Baseball Glove", "description": "Leather baseball glove 12-inch"},
            {"name": "Football Official Size", "description": "Regulation football with laces"},
            {"name": "Tennis Racket Pro", "description": "Graphite tennis racket with dampener"},
            {"name": "Badminton Set", "description": "Portable badminton net with 4 rackets"},
            {"name": "Hockey Stick", "description": "Composite hockey stick flex 85"},
            {"name": "Rugby Ball", "description": "Match rugby ball size 5 with grip"},
            {"name": "Pickleball Paddle Set", "description": "Set of 4 paddles with balls and bag"},
        ],
        "Cycling": [
            {"name": "Road Bike Carbon", "description": "Carbon fiber road bike with Shimano 105"},
            {"name": "Mountain Bike 29er", "description": "Full suspension MTB with 29-inch wheels"},
            {"name": "Helmet Road", "description": "Aero road helmet with ventilation"},
            {"name": "Bike Lock U-Shape", "description": "Hardened steel U-lock with mounting bracket"},
            {"name": "Cycling Jersey", "description": "Breathable cycling jersey with rear pockets"},
            {"name": "Bike Pump Floor", "description": "Presta Schrader compatible floor pump"},
            {"name": "Bike Light Set", "description": "USB rechargeable front and rear light set"},
            {"name": "Bike Computer GPS", "description": "Wireless cycling computer with GPS tracking"},
            {"name": "Cycling Shorts", "description": "Padded cycling shorts with chamois"},
            {"name": "Bike Repair Stand", "description": "Portable work stand for bike maintenance"},
        ],
        "Water Sports": [
            {"name": "Inflatable Kayak", "description": "2-person inflatable kayak with paddles"},
            {"name": "Wetsuit 3/2mm", "description": "Full wetsuit for surfing and diving"},
            {"name": "Snorkel Set", "description": "Mask snorkel and fins set for adults"},
            {"name": "Stand Up Paddleboard", "description": "Inflatable SUP with pump and paddle"},
            {"name": "Swimming Goggles", "description": "Anti-fog swim goggles with UV protection"},
            {"name": "Life Jacket", "description": "USCG-approved life vest for adults"},
            {"name": "Waterproof Dry Bag", "description": "Roll-top dry bag 20L for gear storage"},
            {"name": "Diving Fins", "description": "Adjustable open-heel diving fins"},
            {"name": "Pool Float Lounge", "description": "Inflatable pool float with cup holder"},
            {"name": "Surfboard 7ft", "description": "Foam surfboard for beginners"},
        ],
        "Skincare": [
            {"name": "Vitamin C Serum", "description": "Brightening serum with hyaluronic acid"},
            {"name": "Retinol Cream", "description": "Anti-aging night cream with retinol"},
            {"name": "Sunscreen SPF 50", "description": "Mineral sunscreen with zinc oxide"},
            {"name": "Cleansing Balm", "description": "Oil-based cleanser for makeup removal"},
            {"name": "Moisturizer Face", "description": "Daily face cream with ceramides"},
            {"name": "Eye Cream", "description": "Anti-puffiness eye cream with caffeine"},
            {"name": "Niacinamide Serum", "description": "Pore-minimizing serum with vitamin B3"},
            {"name": "Face Mask Sheet Set", "description": "Hydrating sheet mask set of 10"},
            {"name": "Toner Alcohol-Free", "description": "Soothing toner with rose water"},
            {"name": "Lip Balm SPF 30", "description": "Moisturizing lip balm with sun protection"},
        ],
        "Makeup": [
            {"name": "Foundation Liquid", "description": "Full coverage foundation with SPF 20"},
            {"name": "Eyeshadow Palette", "description": "Neutral eyeshadow palette with 18 shades"},
            {"name": "Mascara Waterproof", "description": "Volumizing mascara with curved brush"},
            {"name": "Lipstick Matte", "description": "Long-wear matte lipstick in classic red"},
            {"name": "Concealer Full Coverage", "description": "Cream concealer for blemishes and dark circles"},
            {"name": "Blush Powder", "description": "Silky blush powder with buildable color"},
            {"name": "Eyeliner Pen", "description": "Waterproof liquid eyeliner with fine tip"},
            {"name": "Setting Spray", "description": "Makeup setting spray with 16hr hold"},
            {"name": "Highlighter Stick", "description": "Cream highlighter for natural glow"},
            {"name": "Makeup Brush Set", "description": "Professional brush set of 12 pieces"},
        ],
        "Hair Care": [
            {"name": "Shampoo Repair", "description": "Keratin-infused shampoo for damaged hair"},
            {"name": "Conditioner Moisture", "description": "Deep conditioning treatment with argan oil"},
            {"name": "Hair Dryer Ionic", "description": "Professional blow dryer with diffuser"},
            {"name": "Hair Straightener", "description": "Ceramic flat iron with adjustable temp"},
            {"name": "Hair Oil Treatment", "description": "Moroccan argan oil for shine and smoothness"},
            {"name": "Dry Shampoo", "description": "Volumizing dry shampoo for fresh hair"},
            {"name": "Hair Mask Deep", "description": "Intensive hair mask with biotin"},
            {"name": "Curl Cream", "description": "Defining cream for curly and wavy hair"},
            {"name": "Hairbrush Detangling", "description": "Cushion hairbrush with nylon bristles"},
            {"name": "Heat Protectant Spray", "description": "Thermal spray for protection up to 450F"},
        ],
        "Fragrance": [
            {"name": "Eau de Parfum Floral", "description": "Feminine floral scent with jasmine notes"},
            {"name": "Cologne Citrus", "description": "Fresh citrus cologne for daily wear"},
            {"name": "Perfume Gift Set", "description": "Set of 3 mini perfumes in varied scents"},
            {"name": "Rollerball Fragrance", "description": "Portable rollerball in vanilla musk"},
            {"name": "Body Spray Fresh", "description": "Light body spray in ocean breeze"},
            {"name": "Oud Perfume Oil", "description": "Concentrated perfume oil with woody notes"},
            {"name": "Rose Eau de Toilette", "description": "Romantic rose scent with soft musk"},
            {"name": "Men's Aftershave Balm", "description": "Soothing aftershave with sandalwood"},
            {"name": "Home Fragrance Diffuser", "description": "Reed diffuser set with lavender essential oil"},
            {"name": "Scented Sachet Set", "description": "Linen sachets in lavender rose and vanilla"},
        ],
        "Vitamins & Supplements": [
            {"name": "Multivitamin Daily", "description": "Complete multivitamin for daily wellness"},
            {"name": "Vitamin D3 2000IU", "description": "Vitamin D supplement for immune support"},
            {"name": "Omega 3 Fish Oil", "description": "High-potency fish oil with EPA and DHA"},
            {"name": "Probiotic 10 Strains", "description": "Digestive health probiotic capsules"},
            {"name": "Magnesium Glycinate", "description": "Highly absorbable magnesium for sleep"},
            {"name": "Whey Protein Powder", "description": "Chocolate whey protein 2lbs"},
            {"name": "Vitamin C 1000mg", "description": "Time-release vitamin C for immunity"},
            {"name": "Collagen Peptides", "description": "Hydrolyzed collagen for skin and joints"},
            {"name": "Melatonin 5mg", "description": "Sleep aid with quick-release tablets"},
            {"name": "Greens Powder Superfood", "description": "Organic greens powder with probiotics"},
        ],
        "Board Games": [
            {"name": "Catan", "description": "Strategy board game of trade and settlement"},
            {"name": "Ticket to Ride", "description": "Cross-country train adventure board game"},
            {"name": "Codenames", "description": "Word-based party game for teams"},
            {"name": "Carcassonne", "description": "Tile-placement game of medieval landscape"},
            {"name": "Pandemic", "description": "Cooperative game to save humanity from diseases"},
            {"name": "Wingspan", "description": "Engine-building game about bird species"},
            {"name": "Azul", "description": "Abstract tile-placement game with beautiful design"},
            {"name": "Splendor", "description": "Gem-collecting strategy game"},
            {"name": "7 Wonders", "description": "Card drafting game of ancient civilizations"},
            {"name": "Terraforming Mars", "description": "Complex strategy game about colonizing Mars"},
        ],
        "Action Figures": [
            {"name": "Marvel Spider-Man Figure", "description": "6-inch articulated Spider-Man action figure"},
            {"name": "Star Wars Darth Vader", "description": "Collectible 6-inch Darth Vader figure"},
            {"name": "Transformers Optimus Prime", "description": "Robot-to-truck transforming action figure"},
            {"name": "LEGO Batman Minifigure", "description": "Buildable Batman figure with cape"},
            {"name": "Gundam RX-78-2", "description": "Model kit of the original Gundam"},
            {"name": "Pokemon Pikachu Plush", "description": "Soft Pikachu plush 10-inches tall"},
            {"name": "Harry Potter Wizard Set", "description": "Harry Hermione and Ron figure set"},
            {"name": "Jurassic World T-Rex", "description": "Large dinosaur figure with action features"},
            {"name": "Minecraft Steve Figure", "description": "Buildable Steve figure with accessories"},
            {"name": "Naruto Shippuden Set", "description": "Naruto and Sasuke 2-pack figures"},
        ],
        "Building Sets": [
            {"name": "LEGO Classic Brick Box", "description": "Large box of 1500 assorted LEGO bricks"},
            {"name": "LEGO City Fire Station", "description": "Detailed fire station with vehicle set"},
            {"name": "Magnetic Building Tiles", "description": "100-piece magnetic tile set"},
            {"name": "Wooden Block Set", "description": "Natural wood blocks 100-piece set"},
            {"name": "Marble Run Track", "description": "Construction set with marbles and tracks"},
            {"name": "Technic Car Set", "description": "Buildable sports car with working engine"},
            {"name": "Architecture Skyline Set", "description": "Buildable city skyline model"},
            {"name": "Building Brick Baseplate", "description": "Green baseplate 32x32 studs"},
            {"name": "Straw Constructor Set", "description": "Connecting straw building kit 400 pieces"},
            {"name": "Foam Building Blocks", "description": "Soft foam blocks for toddler construction"},
        ],
        "Puzzles": [
            {"name": "Jigsaw 1000-Piece Landscape", "description": "Scenic mountain landscape puzzle"},
            {"name": "Wooden Puzzle 3D", "description": "3D model puzzle of a sailing ship"},
            {"name": "Rubic's Cube Speed", "description": "Speed cube with smooth rotation"},
            {"name": "Jigsaw 500-Piece Wildlife", "description": "Animal collage puzzle for families"},
            {"name": "Brain Teaser Set", "description": "Set of 6 metal wire puzzle challenges"},
            {"name": "Sudoku Book Hard", "description": "200 hard sudoku puzzles in travel size"},
            {"name": "Crossword Puzzle Book", "description": "Daily crossword puzzle collection"},
            {"name": "Floor Puzzle 48-Piece", "description": "Large kid's floor puzzle dinosaur themed"},
            {"name": "Logic Puzzle Grid", "description": "Grid logic puzzle book with 100 puzzles"},
            {"name": "Hanayama Cast Puzzle", "description": "Metal disentanglement puzzle level 4/6"},
        ],
        "Video Games": [
            {"name": "Elden Ring", "description": "Action RPG from the creators of Dark Souls"},
            {"name": "The Legend of Zelda TOTK", "description": "Open-world adventure in Hyrule"},
            {"name": "God of War Ragnarok", "description": "Epic Norse mythology action game"},
            {"name": "Baldur's Gate 3", "description": "Award-winning D&D-based RPG"},
            {"name": "Spider-Man 2 PS5", "description": "Open-world superhero action game"},
            {"name": "Cyberpunk 2077", "description": "Open-world sci-fi RPG with Phantom Liberty"},
            {"name": "Final Fantasy VII Rebirth", "description": "Expansive JRPG in the FF7 universe"},
            {"name": "Starfield", "description": "Space exploration RPG from Bethesda"},
            {"name": "Hogwarts Legacy", "description": "Open-world action RPG in the Wizarding World"},
            {"name": "Resident Evil 4 Remake", "description": "Reimagined survival horror classic"},
        ],
        "Car Care": [
            {"name": "Car Wash Kit", "description": "Complete car wash bucket and mitt set"},
            {"name": "Microfiber Towel Pack", "description": "Pack of 12 premium microfiber detailing towels"},
            {"name": "Car Wax Premium", "description": "Carnauba wax for deep shine paint protection"},
            {"name": "Tire Shine Spray", "description": "Long-lasting tire dressing with UV protectant"},
            {"name": "Interior Cleaner", "description": "All-purpose automotive interior cleaner spray"},
            {"name": "Glass Cleaner Automotive", "description": "Streak-free glass cleaner for windows"},
            {"name": "Car Vacuum Portable", "description": "High-suction handheld car vacuum"},
            {"name": "Clay Bar Kit", "description": "Clay bar and lubricant for smooth paint"},
            {"name": "Leather Cleaner and Conditioner", "description": "Clean and protect leather car seats"},
            {"name": "Pressure Washer", "description": "Electric pressure washer for car washing"},
        ],
        "Interior Accessories": [
            {"name": "Car Floor Mats All-Weather", "description": "Laser-measured floor liner set front and rear"},
            {"name": "Seat Cover Set", "description": "Waterproof neoprene seat covers for front seats"},
            {"name": "Steering Wheel Cover", "description": "Leather steering wheel cover with stitching"},
            {"name": "Phone Mount Dashboard", "description": "Magnetic phone mount for car dashboard"},
            {"name": "Car Trash Can", "description": "Compact car trash can with lid"},
            {"name": "Sun Shade Windshield", "description": "Foldable sun shade blocks UV and heat"},
            {"name": "Air Freshener Long Lasting", "description": "Vent clip air freshener 60-day scent"},
            {"name": "Back Seat Organizer", "description": "Multi-pocket seat back storage for kids"},
            {"name": "Car Charger USB-C", "description": "Fast-charge 45W USB-C car charger"},
            {"name": "LED Interior Light Strip", "description": "RGB LED strip with app control"},
        ],
        "Exterior Accessories": [
            {"name": "Roof Rack Crossbars", "description": "Adjustable roof rack for cargo and gear"},
            {"name": "License Plate Frame", "description": "Stainless steel license plate frame"},
            {"name": "Hitch Receiver", "description": "Class III trailer hitch 2-inch receiver"},
            {"name": "Mud Flaps Set", "description": "Heavy-duty mud flaps for truck or SUV"},
            {"name": "Side Window Deflectors", "description": "In-channel rain guards for windows"},
            {"name": "Car Cover All-Weather", "description": "Weatherproof car cover for outdoor storage"},
            {"name": "Bike Rack Trunk", "description": "Trunk-mounted bike carrier for 2 bikes"},
            {"name": "Spoiler Lip", "description": "ABS plastic rear spoiler lip spoiler"},
            {"name": "Hood Protector", "description": "Smoke hood shield bug deflector"},
            {"name": "Tailgate Step", "description": "Retractable truck tailgate step"},
        ],
        "Automotive Tools": [
            {"name": "Car Jack 3-Ton", "description": "Hydraulic floor jack with lifting range"},
            {"name": "Jump Starter Battery", "description": "Portable jump starter with USB output"},
            {"name": "OBD2 Scanner", "description": "Bluetooth OBD2 car diagnostic tool"},
            {"name": "Tool Set Auto 100-Piece", "description": "Complete auto repair tool kit"},
            {"name": "Tire Inflator Portable", "description": "Digital tire inflator with auto shut-off"},
            {"name": "Oil Filter Wrench Set", "description": "Set of 4 oil filter wrenches"},
            {"name": "Mechanic's Creeper", "description": "Low-profile mechanic creeper with headrest"},
            {"name": "Funnel Set Automotive", "description": "Set of 3 funnels for fluids"},
            {"name": "Zip Tie Assortment", "description": "Assorted automotive zip ties in case"},
            {"name": "Work Light Rechargeable", "description": "Magnetic LED work light for car repairs"},
        ],
        "Car Electronics": [
            {"name": "Dash Cam 4K", "description": "4K front and rear dash camera with GPS"},
            {"name": "GPS Navigation System", "description": "Portable GPS with real-time traffic"},
            {"name": "Car Stereo Touchscreen", "description": "Double DIN touchscreen with Apple CarPlay"},
            {"name": "Backup Camera Wireless", "description": "Wireless backup camera with night vision"},
            {"name": "Radar Detector", "description": "Long-range radar detector with laser detection"},
            {"name": "Bluetooth FM Transmitter", "description": "Wireless FM transmitter hands-free calling"},
            {"name": "Subwoofer Car Audio", "description": "12-inch powered subwoofer"},
            {"name": "Parking Sensor Kit", "description": "4-sensor parking assist with alarm"},
            {"name": "Tire Pressure Monitor", "description": "Solar TPMS with real-time tire readings"},
            {"name": "Car Battery Charger", "description": "Smart battery maintainer 12V"},
        ],
        "Snacks": [
            {"name": "Mixed Nuts Premium", "description": "Roasted salted mixed nuts 32oz"},
            {"name": "Dark Chocolate Bar 72%", "description": "Belgian dark chocolate with cocoa nibs"},
            {"name": "Potato Chips Sea Salt", "description": "Kettle-cooked potato chips 8oz bag"},
            {"name": "Beef Jerky Original", "description": "Smoked beef jerky high protein"},
            {"name": "Trail Mix Energy", "description": "Nuts seeds and dried fruit mix"},
            {"name": "Rice Crackers Wasabi", "description": "Crispy wasabi rice crackers 6-pack"},
            {"name": "Granola Bar Variety Pack", "description": "Mixed flavor granola bars 24-pack"},
            {"name": "Popcorn Gourmet", "description": "Butter-free gourmet popcorn tin"},
            {"name": "Dried Mango Slices", "description": "Natural dried mango no sugar added"},
            {"name": "Edamame Roasted", "description": "Protein-packed roasted edamame snack"},
        ],
        "Beverages": [
            {"name": "Coffee Beans Single Origin", "description": "Ethiopian Yirgacheffe medium roast 12oz"},
            {"name": "Green Tea Premium", "description": "Japanese matcha green tea ceremonial grade"},
            {"name": "Sparkling Water Variety", "description": "Assorted sparkling water 12-pack"},
            {"name": "Orange Juice Pure", "description": "Not-from-concentrate orange juice 64oz"},
            {"name": "Chai Tea Concentrate", "description": "Spiced chai tea concentrate 32oz"},
            {"name": "Coconut Water Hydrate", "description": "Pure coconut water with pulp 12-pack"},
            {"name": "Kombucha Ginger", "description": "Raw kombucha probiotic drink 4-pack"},
            {"name": "Protein Shake Chocolate", "description": "Ready-to-drink protein shake 4-pack"},
            {"name": "Herbal Tea Assortment", "description": "Caffeine-free herbal tea 20-bag variety"},
            {"name": "Cold Brew Concentrate", "description": "Cold brew coffee concentrate 32oz"},
        ],
        "Pantry Staples": [
            {"name": "Organic Olive Oil", "description": "Extra virgin olive oil cold-pressed"},
            {"name": "Basmati Rice 5lb", "description": "Aged basmati rice long grain"},
            {"name": "Pasta Spaghetti Organic", "description": "Organic durum wheat spaghetti 1lb"},
            {"name": "Peanut Butter Natural", "description": "Creamy peanut butter no added sugar"},
            {"name": "Honey Pure Wildflower", "description": "Raw wildflower honey 12oz"},
            {"name": "Coconut Milk Canned", "description": "Organic full-fat coconut milk 13.5oz"},
            {"name": "Soy Sauce Traditional", "description": "Brewed soy sauce naturally aged"},
            {"name": "Maple Syrup Pure", "description": "Grade A amber maple syrup 8oz"},
            {"name": "Sea Salt Fine Ground", "description": "All-natural sea salt 1lb bag"},
            {"name": "Black Peppercorns Whole", "description": "Tellicherry black peppercorns 4oz"},
        ],
        "Organic Foods": [
            {"name": "Organic Quinoa", "description": "Certified organic white quinoa 2lb"},
            {"name": "Organic Chicken Breast", "description": "Free-range organic chicken breast 1lb"},
            {"name": "Organic Baby Spinach", "description": "Pre-washed organic spinach 5oz"},
            {"name": "Organic Eggs 12-Count", "description": "Pasture-raised organic brown eggs"},
            {"name": "Organic Avocados 4-Pack", "description": "Ripe organic Hass avocados"},
            {"name": "Organic Almond Milk", "description": "Unsweetened vanilla almond milk 32oz"},
            {"name": "Organic Sweet Potatoes", "description": "Organic sweet potato 3lb bag"},
            {"name": "Organic Strawberries", "description": "Fresh organic strawberries 1lb"},
            {"name": "Organic Oatmeal", "description": "Rolled organic oats 42oz canister"},
            {"name": "Organic Tomato Sauce", "description": "Organic crushed tomatoes 24oz jar"},
        ],
        "Gourmet Foods": [
            {"name": "Truffle Oil White", "description": "White truffle infused olive oil"},
            {"name": "Balsamic Vinegar Aged", "description": "12-year aged balsamic vinegar from Modena"},
            {"name": "Smoked Salmon", "description": "Scottish smoked salmon sliced 8oz"},
            {"name": "Italian Prosciutto", "description": "Prosciutto di Parma aged 18 months"},
            {"name": "French Cheese Selection", "description": "Assorted French cheese gift box"},
            {"name": "Caviar Sturgeon", "description": "Osetra caviar 1oz tin"},
            {"name": "Saffron Threads", "description": "Premium Spanish saffron 1g tin"},
            {"name": "Artisan Chocolate Box", "description": "Handcrafted chocolate truffles 24-piece"},
            {"name": "Matcha Green Tea Powder", "description": "Ceremonial grade Japanese matcha 4oz"},
            {"name": "Wagyu Beef Striploin", "description": "A5 Wagyu striploin frozen 8oz"},
        ],
        "Writing Instruments": [
            {"name": "Ballpoint Pen Premium", "description": "Rollerball pen with smooth ink flow"},
            {"name": "Mechanical Pencil 0.5mm", "description": "Drafting pencil with ergonomic grip"},
            {"name": "Fountain Pen Set", "description": "Luxury fountain pen with ink cartridges"},
            {"name": "Highlighter Set Pastel", "description": "Set of 6 pastel highlighters"},
            {"name": "Gel Pen Set Colorful", "description": "Pack of 12 vibrant gel pens"},
            {"name": "Fine Liner Pen Set", "description": "Archival ink fine liners set of 8"},
            {"name": "Whiteboard Markers", "description": "Dry erase markers assorted colors 8-pack"},
            {"name": "Calligraphy Pen Set", "description": "Calligraphy set with nibs and ink"},
            {"name": "Pencil Case Large", "description": "Canvas pencil case with 50-slot capacity"},
            {"name": "Permanent Markers", "description": "Oil-based permanent markers 5-pack"},
        ],
        "Paper Products": [
            {"name": "Notebook A5 Lined", "description": "Hardcover notebook with lay-flat binding"},
            {"name": "Sticky Notes Assorted", "description": "Super sticky neon notes 12-pad pack"},
            {"name": "Printer Paper 5000-Sheet", "description": "Multipurpose copy paper 8.5x11"},
            {"name": "Planner Weekly Monthly", "description": "Undated daily planner with goal tracker"},
            {"name": "Index Cards 3x5", "description": "Blank ruled index cards 100-pack"},
            {"name": "Composition Notebook", "description": "Marble composition book wide ruled"},
            {"name": "Binder 3-Ring 2-Inch", "description": "Durable binder with clear cover"},
            {"name": "Graph Paper Pad", "description": "Quad ruled graph paper 8.5x11 100 sheets"},
            {"name": "Legal Pad Yellow", "description": "Professional legal pad 50 sheets"},
            {"name": "Sketchbook A4", "description": "Heavyweight paper sketchbook 80 pages"},
        ],
        "Desk Accessories": [
            {"name": "Desk Organizer Mesh", "description": "Multi-compartment mesh desk organizer"},
            {"name": "Monitor Stand Riser", "description": "Adjustable monitor stand with storage"},
            {"name": "Desk Lamp LED", "description": "Adjustable LED lamp with color temperature"},
            {"name": "Wrist Rest Keyboard", "description": "Memory foam wrist rest for ergonomic typing"},
            {"name": "Cable Management Sleeve", "description": "Cable organizer sleeve 20-inch"},
            {"name": "Paper Tray Stackable", "description": "Stackable letter tray for documents"},
            {"name": "Mouse Pad Large", "description": "Extended desk mouse pad with stitched edge"},
            {"name": "Stapler Heavy Duty", "description": "50-sheet capacity stapler with staple remover"},
            {"name": "Tape Dispenser Desktop", "description": "Weighted tape dispenser for desk use"},
            {"name": "Business Card Holder", "description": "Acrylic business card holder stand"},
        ],
        "Office Storage": [
            {"name": "File Box Plastic", "description": "Stackable file storage box with lid"},
            {"name": "Hanging File Folders", "description": "Letter size hanging folders 25-pack"},
            {"name": "Storage Cabinet 3-Drawer", "description": "Mobile file cabinet with lock"},
            {"name": "Magazine Holder", "description": "Acrylic magazine holder set of 4"},
            {"name": "Bookend Metal Set", "description": "Non-slip metal bookends black pair"},
            {"name": "Drawer Organizer Set", "description": "Divided drawer organizer for small items"},
            {"name": "Label Maker", "description": "Thermal label maker with Bluetooth"},
            {"name": "Shipping Supplies Kit", "description": "Bubble mailers and tape assortment"},
            {"name": "Copy Paper Box", "description": "Bankers box for document storage"},
            {"name": "Sealable File Bags", "description": "Document storage bags waterproof 10-pack"},
        ],
        "Office Technology": [
            {"name": "USB-C Hub 7-in-1", "description": "Multi-port adapter with HDMI and SD"},
            {"name": "Wireless Mouse", "description": "Ergonomic wireless mouse silent clicks"},
            {"name": "External SSD 1TB", "description": "Portable solid state drive USB-C"},
            {"name": "USB Extension Cable 10ft", "description": "High-speed USB 3.0 extension cable"},
            {"name": "Desk Phone Headset", "description": "Wireless headset with noise cancellation"},
            {"name": "Surge Protector 6-Outlet", "description": "Power strip surge protector with USB ports"},
            {"name": "Webcam 1080p", "description": "Full HD webcam with privacy shutter"},
            {"name": "Microphone USB", "description": "Condenser microphone for conference calls"},
            {"name": "Ethernet Cable 25ft", "description": "Cat6 Ethernet cable high-speed"},
            {"name": "Adapter VGA to HDMI", "description": "VGA to HDMI converter with audio"},
        ],
    }

    # Products table is now SPU (no price/sku, only min_price).
    # We generate SPU entries, then create SKUs separately.

    # Build base (cat_name, prod_data) combinations from product_samples.
    base_entries = []
    for cat in categories:
        cat_name = cat.get("name")
        if cat_name in product_samples:
            for prod_data in product_samples[cat_name]:
                base_entries.append((cat_name, prod_data))

    # Generate SPUs — each base entry becomes one Product (SPU) with min_price
    for cat_name, prod_data in base_entries:
        products.append({
            "id": None,
            "name": prod_data["name"],
            "description": prod_data["description"],
            "category_hint": cat_name,
            "min_price": random.randint(1000, 50000),
            "created_at": mysql_utc_now(),
            "updated_at": mysql_utc_now(),
        })

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


def generate_inventory(skus):
    """Generate sample inventory records (per SKU)"""
    inventory = []

    for sku in skus:
        inv_record = {
            "sku_id": sku["id"],
            "quantity": random.randint(0, 100),
            "reserved": random.randint(0, min(10, random.randint(0, 100))),
            "threshold": 10,
            "created_at": mysql_utc_now(),
            "updated_at": mysql_utc_now(),
        }
        if inv_record["quantity"] <= 0:
            inv_record["status"] = "outofstock"
        elif inv_record["quantity"] <= inv_record["threshold"]:
            inv_record["status"] = "lowstock"
        else:
            inv_record["status"] = "instock"

        inventory.append(inv_record)

    return inventory


# ─── SPU/SKU variant helpers ──────────────────────────────────────────

SKU_SPECS = {
    "Smartphones": {
        "dimensions": [
            {"name": "系列", "options": ["标准版", "Pro", "Pro Max"]},
            {"name": "内存", "options": ["128GB", "256GB", "512GB"]},
            {"name": "颜色", "options": ["深空黑", "银色"]},
        ],
    },
    "Laptops": {
        "dimensions": [
            {"name": "内存", "options": ["16GB", "32GB"]},
            {"name": "硬盘", "options": ["512GB", "1TB"]},
        ],
    },
    "Shoes": {
        "dimensions": [
            {"name": "尺码", "options": ["39", "40", "41", "42", "43"]},
            {"name": "颜色", "options": ["黑色", "白色"]},
        ],
    },
}


def generate_skus(spus):
    """Generate SKU entries for each SPU (variant categories get multiple SKUs)"""
    import itertools
    skus = []
    for spu in spus:
        hint = spu.get("category_hint")
        spec_template = SKU_SPECS.get(hint)
        base_price = spu.get("min_price", 2999)

        if spec_template:
            dims = spec_template["dimensions"]
            keys = [d["name"] for d in dims]
            values = [d["options"] for d in dims]
            for combo in itertools.product(*values):
                spec = dict(zip(keys, combo))
                sku_name = " ".join([spu["name"]] + [str(v) for v in combo])
                sku_code = "SKU-{}".format(str(uuid.uuid4()).split("-")[0].upper())
                price = base_price + random.randint(0, 20000)
                skus.append({
                    "id": None, "product_id": spu["id"],
                    "name": sku_name, "price": price,
                    "sku_code": sku_code,
                    "spec": json.dumps(spec, ensure_ascii=False),
                    "created_at": spu["created_at"], "updated_at": spu["updated_at"],
                })
        else:
            sku_code = "SKU-{}".format(str(uuid.uuid4()).split("-")[0].upper())
            skus.append({
                "id": None, "product_id": spu["id"],
                "name": spu["name"], "price": base_price,
                "sku_code": sku_code, "spec": None,
                "created_at": spu["created_at"], "updated_at": spu["updated_at"],
            })
    return skus


def insert_skus(connection, skus):
    """Insert SKUs into the database"""
    with connection.cursor() as cursor:
        for sku in skus:
            sql = """
            INSERT INTO skus (product_id, name, price, sku_code, spec, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(sql, (
                    sku["product_id"], sku["name"], sku["price"],
                    sku["sku_code"], sku["spec"],
                    sku["created_at"], sku["updated_at"],
                ))
                sku["id"] = cursor.lastrowid
            except Exception as e:
                print(f"Error inserting SKU {sku['name']}: {e}")
        connection.commit()
    print(f"Inserted {len(skus)} SKUs into the database")


def insert_attributes(connection):
    """Insert attribute definitions and values into the database"""
    attrs = [
        {"name": "系列", "values": ["标准版", "Pro", "Pro Max"]},
        {"name": "内存", "values": ["128GB", "256GB", "512GB", "1TB"]},
        {"name": "颜色", "values": ["深空黑", "银色", "金色", "深蓝色", "黑色", "白色"]},
        {"name": "CPU", "values": ["标准版", "高性能版"]},
        {"name": "硬盘", "values": ["512GB", "1TB"]},
        {"name": "尺码", "values": ["39", "40", "41", "42", "43"]},
    ]
    with connection.cursor() as cursor:
        for attr in attrs:
            cursor.execute(
                "INSERT IGNORE INTO attribute_attributes (name, sort_order) VALUES (%s, 0)",
                (attr["name"],),
            )
            if cursor.lastrowid == 0:
                cursor.execute("SELECT id FROM attribute_attributes WHERE name = %s", (attr["name"],))
                attr_id = cursor.fetchone()["id"]
            else:
                attr_id = cursor.lastrowid
            for idx, val in enumerate(attr["values"]):
                cursor.execute(
                    "INSERT IGNORE INTO attribute_values (attribute_id, value, sort_order) VALUES (%s, %s, %s)",
                    (attr_id, val, idx),
                )
        connection.commit()
    print(f"Inserted {len(attrs)} attributes with their values into the database")


def insert_product_attribute_values(connection, spus):
    """Associate products (SPUs) with attribute dimensions and values based on their category."""
    with connection.cursor() as cursor:
        # Load attribute ID → name mapping
        cursor.execute("SELECT id, name FROM attribute_attributes")
        attr_map = {}
        for row in cursor.fetchall():
            attr_map[row['name']] = {"id": row['id'], "values": {}}

        # Load value ID → value mapping per attribute
        cursor.execute("SELECT id, attribute_id, value FROM attribute_values")
        for row in cursor.fetchall():
            for info in attr_map.values():
                if info["id"] == row['attribute_id']:
                    info["values"][row['value']] = row['id']
                    break

        count = 0
        for spu in spus:
            hint = spu.get("category_hint")
            spec_template = SKU_SPECS.get(hint)
            if not spec_template:
                continue

            for dim in spec_template["dimensions"]:
                attr_name = dim["name"]
                if attr_name not in attr_map:
                    continue
                attr_id = attr_map[attr_name]["id"]
                for opt in dim["options"]:
                    if opt in attr_map[attr_name]["values"]:
                        val_id = attr_map[attr_name]["values"][opt]
                        cursor.execute(
                            "INSERT IGNORE INTO product_attribute_values (product_id, attribute_id, attribute_value_id) VALUES (%s, %s, %s)",
                            (spu["id"], attr_id, val_id),
                        )
                        count += 1
        connection.commit()
    print(f"Inserted {count} product-attribute-value links into the database")


def insert_sku_attributes(connection, skus):
    """Link SKUs to their attribute values based on spec JSON"""
    with connection.cursor() as cursor:
        cursor.execute("SELECT id, name FROM attribute_attributes")
        attr_map = {}
        for row in cursor.fetchall():
            attr_map[row['name']] = {"id": row['id'], "values": {}}
        cursor.execute("SELECT id, attribute_id, value FROM attribute_values")
        for row in cursor.fetchall():
            for info in attr_map.values():
                if info["id"] == row['attribute_id']:
                    info["values"][row['value']] = row['id']
                    break
        count = 0
        for sku in skus:
            if not sku.get("spec"):
                continue
            try:
                spec = json.loads(sku["spec"])
            except (json.JSONDecodeError, TypeError):
                continue
            for key, val in spec.items():
                if key in attr_map and val in attr_map[key]["values"]:
                    vid = attr_map[key]["values"][val]
                    cursor.execute(
                        "INSERT IGNORE INTO sku_attributes (sku_id, attribute_id, attribute_value_id) VALUES (%s, %s, %s)",
                        (sku["id"], attr_map[key]["id"], vid),
                    )
                    count += 1
        connection.commit()
    print(f"Inserted {count} sku-attribute links into the database")


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
    """Insert SPUs (products table) into the database"""
    with connection.cursor() as cursor:
        for product in products:
            sql = """
            INSERT INTO products (name, description, min_price, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(sql, (
                    product["name"],
                    product["description"],
                    product["min_price"],
                    product["created_at"],
                    product["updated_at"],
                ))
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
            INSERT INTO inventories (sku_id, quantity, status, reserved, threshold, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(sql, (
                    inv["sku_id"],
                    inv["quantity"],
                    inv["status"],
                    inv["reserved"],
                    inv["threshold"],
                    inv["created_at"],
                    inv["updated_at"]
                ))
            except Exception as e:
                print(
                    f"Error inserting inventory for SKU {inv['sku_id']}: {e}")
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


def insert_orders(connection, orders):
    """Insert orders into the database (without items)."""
    with connection.cursor() as cursor:
        for order in orders:
            sql = """
            INSERT INTO orders (order_no, customer_id, total_amount, currency, status, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(sql, (
                    order["order_no"],
                    order["customer_id"],
                    order["total_amount"],
                    order["currency"],
                    order["status"],
                    order["created_at"],
                    order["updated_at"],
                ))
                order["id"] = cursor.lastrowid
            except Exception as e:
                print(f"Error inserting order {order['order_no']}: {e}")
        connection.commit()
    print(f"Inserted {len(orders)} orders into the database")


def insert_order_items(connection, order_items):
    """Insert order items into the database."""
    with connection.cursor() as cursor:
        for item in order_items:
            sql = """
            INSERT INTO order_items (order_id, product_id, sku_id, quantity, unit_price, amount, created_at, updated_at)
            VALUES (%s, %s, %s, %s, %s, %s, %s, %s)
            """
            try:
                cursor.execute(sql, (
                    item["order_id"],
                    item["product_id"],
                    item["sku_id"],
                    item["quantity"],
                    item["unit_price"],
                    item["amount"],
                    item["created_at"],
                    item["updated_at"],
                ))
                item["id"] = cursor.lastrowid
            except Exception as e:
                print(f"Error inserting order_item for order {item['order_id']}: {e}")
        connection.commit()
    print(f"Inserted {len(order_items)} order items into the database")


def generate_order_no():
    """Generate order_no in format ORD<milliseconds><4-digit-random> (matching Go generateOrderNo())"""
    now = int(time.time() * 1000)
    r = random.randint(0, 9999)
    return f"ORD{now}{r:04d}"


def generate_orders(skus, num_orders=300):
    """Generate sample orders with order_no (order items reference SKUs)"""
    customer_ids = ["1", "2", "3", "4", "5"]

    order_statuses = ["pending", "paid", "paid", "shipped", "delivered", "delivered", "cancelled"]
    orders = []
    now_ts = mysql_utc_now()

    for i in range(num_orders):
        customer_id = random.choice(customer_ids)
        status = random.choice(order_statuses)
        num_items = random.randint(1, 4)
        selected_skus = random.sample(skus, min(num_items, len(skus)))

        total_amount = 0
        items = []
        for sku in selected_skus:
            qty = random.randint(1, 3)
            amount = sku["price"] * qty
            total_amount += amount
            items.append({
                "id": None,
                "order_id": None,
                "product_id": str(sku["product_id"]),
                "sku_id": sku["id"],
                "quantity": qty,
                "unit_price": sku["price"],
                "amount": amount,
                "created_at": now_ts,
                "updated_at": now_ts,
            })

        orders.append({
            "id": None,
            "order_no": generate_order_no(),
            "customer_id": customer_id,
            "total_amount": total_amount,
            "currency": "CNY",
            "status": status,
            "created_at": now_ts,
            "updated_at": now_ts,
            "_items": items,  # internal, won't be persisted directly
        })

    return orders


def generate_order_items(orders):
    """Extract order_items from orders (after orders have IDs assigned)."""
    order_items = []
    for order in orders:
        for item in order["_items"]:
            item["order_id"] = order["id"]
            order_items.append(item)
    return order_items


def clean_database(connection):
    """Clean all existing data from all tables with FK-safe truncation order."""
    tables = [
        "order_items",
        "cart_items",
        "refunds",
        "payment_transactions",
        "orders",
        "carts",
        "payments",
        "notifications",
        "product_attribute_values",
        "sku_attributes",
        "skus",
        "attribute_values",
        "attribute_attributes",
        "inventories",
        "product_categories",
        "products",
        "categories",
        "flash_orders",
        "flash_activities",
    ]
    with connection.cursor() as cursor:
        cursor.execute("SET FOREIGN_KEY_CHECKS = 0")
        for table in tables:
            cursor.execute(f"TRUNCATE TABLE {table}")
        cursor.execute("SET FOREIGN_KEY_CHECKS = 1")
    connection.commit()
    print(f"Cleaned {len(tables)} tables.")


def seed_rbac_data(connection):
    """Seed RBAC data (roles + permissions) and create initial admin/user accounts."""
    now_ts = mysql_utc_now()

    roles_data = [
        ("admin", "管理员", "系统管理员，拥有所有权限", 1, 1, 1),
        ("user", "普通用户", "普通用户，拥有基本操作权限", 1, 2, 1),
    ]

    permissions_data = [
        ("product:read",    "查看产品", "查看产品列表和详情",    "product",  "read",   "商品管理", 1, 1),
        ("product:create",  "创建产品", "创建新产品",            "product",  "create", "商品管理", 2, 1),
        ("product:update",  "编辑产品", "编辑产品信息",          "product",  "update", "商品管理", 3, 1),
        ("product:delete",  "删除产品", "删除产品",              "product",  "delete", "商品管理", 4, 1),
        ("order:read",      "查看订单", "查看订单列表和详情",    "order",    "read",   "订单管理", 5, 1),
        ("order:create",    "创建订单", "创建订单",              "order",    "create", "订单管理", 6, 1),
        ("order:update",    "编辑订单", "编辑订单信息",          "order",    "update", "订单管理", 7, 1),
        ("order:cancel",    "取消订单", "取消订单",              "order",    "cancel", "订单管理", 8, 1),
        ("user:read",       "查看用户", "查看用户列表和详情",    "user",     "read",   "用户管理", 9, 1),
        ("user:create",     "创建用户", "创建用户",              "user",     "create", "用户管理",10, 1),
        ("user:update",     "编辑用户", "编辑用户信息",          "user",     "update", "用户管理",11, 1),
        ("user:delete",     "删除用户", "删除用户",              "user",     "delete", "用户管理",12, 1),
        ("role:read",       "查看角色", "查看角色列表和详情",    "role",     "read",   "权限管理",13, 1),
        ("role:create",     "创建角色", "创建角色",              "role",     "create", "权限管理",14, 1),
        ("role:update",     "编辑角色", "编辑角色信息",          "role",     "update", "权限管理",15, 1),
        ("role:delete",     "删除角色", "删除角色",              "role",     "delete", "权限管理",16, 1),
        ("inventory:read",   "查看库存", "查看库存信息",         "inventory","read",   "库存管理",17, 1),
        ("inventory:manage", "管理库存", "入库/出库/调整库存",   "inventory","manage", "库存管理",18, 1),
        ("category:read",   "查看分类", "查看分类列表",          "category", "read",   "分类管理",19, 1),
        ("category:manage", "管理分类", "创建/编辑/删除分类",    "category","manage", "分类管理",20, 1),
    ]

    user_permission_names = [
        "product:read", "order:read", "order:create",
        "order:cancel", "category:read", "inventory:read",
    ]

    # bcrypt hash of "123456" (cost=10, same as seed_rbac.sql)
    bcrypt_hash = "$2a$10$HFzEUNEVKJQCZ4aPYVb/YONrhix2jwj8iiJWM5TUZdXM4wPdkEllC"

    with connection.cursor() as cursor:
        # Insert roles
        role_ids = {}
        for name, display_name, description, status, sort, is_system in roles_data:
            sql = """INSERT INTO roles (name, display_name, description, status, sort, is_system, created_at, updated_at)
                     VALUES (%s, %s, %s, %s, %s, %s, %s, %s)"""
            cursor.execute(sql, (name, display_name, description, status, sort, is_system, now_ts, now_ts))
            role_ids[name] = cursor.lastrowid

        # Insert permissions
        perm_name_to_id = {}
        for name, display_name, description, resource, action, category, sort, status in permissions_data:
            sql = """INSERT INTO permissions (name, display_name, description, resource, action, category, sort, status, created_at, updated_at)
                     VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s)"""
            cursor.execute(sql, (name, display_name, description, resource, action, category, sort, status, now_ts, now_ts))
            perm_name_to_id[name] = cursor.lastrowid

        # Admin gets all permissions
        admin_role_id = role_ids["admin"]
        for perm_id in perm_name_to_id.values():
            cursor.execute("INSERT IGNORE INTO role_permissions (role_id, permission_id) VALUES (%s, %s)", (admin_role_id, perm_id))

        # User gets basic permissions
        user_role_id = role_ids["user"]
        for name in user_permission_names:
            cursor.execute("INSERT IGNORE INTO role_permissions (role_id, permission_id) VALUES (%s, %s)", (user_role_id, perm_name_to_id[name]))

        # --- Initial users ---
        # Admin user (password: 123456)
        cursor.execute("INSERT INTO users (status, created_at, updated_at) VALUES (1, %s, %s)", (now_ts, now_ts))
        cursor.execute("INSERT INTO user_infos (user_id, nickname, created_at, updated_at) VALUES (1, '管理员', %s, %s)", (now_ts, now_ts))
        cursor.execute("""INSERT INTO user_identities (user_id, provider, identifier, credential, verified, meta, created_at, updated_at)
                          VALUES (1, 'password', 'admin', %s, 1, '{}', %s, %s)""", (bcrypt_hash, now_ts, now_ts))
        cursor.execute("INSERT INTO user_roles (user_id, role_id, created_at) VALUES (1, %s, %s)", (admin_role_id, now_ts))

        # Regular user colin (password: 123456)
        cursor.execute("INSERT INTO users (status, created_at, updated_at) VALUES (1, %s, %s)", (now_ts, now_ts))
        cursor.execute("INSERT INTO user_infos (user_id, nickname, created_at, updated_at) VALUES (2, 'Colin', %s, %s)", (now_ts, now_ts))
        cursor.execute("""INSERT INTO user_identities (user_id, provider, identifier, credential, verified, meta, created_at, updated_at)
                          VALUES (2, 'password', 'colin', %s, 1, '{}', %s, %s)""", (bcrypt_hash, now_ts, now_ts))
        cursor.execute("INSERT INTO user_roles (user_id, role_id, created_at) VALUES (2, %s, %s)", (user_role_id, now_ts))

    connection.commit()
    print(f"Seeded {len(roles_data)} roles, {len(permissions_data)} permissions, and role-permission associations.")
    print("Created initial admin user (admin / 123456) and regular user (colin / 123456).")


def save_to_json(data, filename):
    """Save data to a JSON file"""
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=2, ensure_ascii=False)
    print(f"Saved {len(data)} items to {filename}")


def main():
    print("Generating and inserting sample data for eshop-monolith...")

    # Connect to database
    connection = connect_to_database()

    # 本地开发环境：每次全量清理业务数据并重置自增（保留 RBAC 相关表）
    clean_database(connection)
    print("Existing data cleaned, auto-increment reset.\n")

    # Seed RBAC data — 首次运行或表为空时写入，已有则跳过
    try:
        seed_rbac_data(connection)
    except Exception:
        print("RBAC data already exists, skipping seed.")

    # ── Generate SPUs (products table) ──────────────────────────────────
    categories = generate_categories()
    spus = generate_products(categories)

    insert_categories(connection, categories)
    insert_products(connection, spus)
    insert_attributes(connection)
    insert_product_attribute_values(connection, spus)

    # ── Product-category associations ───────────────────────────────────
    product_categories = generate_product_categories_links(spus, categories)
    if not product_categories:
        print("WARN: product_categories empty, assigning random categories...")
        valid_cats = [c for c in categories if c.get("id") is not None]
        for spu in spus:
            if spu.get("id") is None:
                continue
            cat = random.choice(valid_cats)
            product_categories.append({
                "product_id": spu["id"], "category_id": cat["id"],
            })
    insert_product_categories(connection, product_categories)

    # Close connection
    connection.close()

    # ── Save to JSON ────────────────────────────────────────────────────
    output_dir = os.path.join(os.path.dirname(__file__), "..", "sample_data")
    os.makedirs(output_dir, exist_ok=True)

    save_to_json(categories, os.path.join(output_dir, "categories.json"))
    save_to_json(spus, os.path.join(output_dir, "spus.json"))

    combined_data = {
        "categories": categories,
        "spus": spus,
        "product_categories": product_categories,
        "generated_at": mysql_utc_now(),
    }
    save_to_json(combined_data, os.path.join(output_dir, "sample_data.json"))

    print("\nSample data generation and insertion complete!")
    print(f"Inserted:")
    print(f"- {len(categories)} categories")
    print(f"- {len(spus)} SPUs")
    print(f"- {len(product_categories)} product-category associations")
    print(f"\nSKUs were not generated. Use the product attribute APIs to configure")
    print(f"attribute combinations and create SKUs via POST /api/v1/products/:id/skus/batch")
    print(f"\nFiles saved to: {output_dir}")


if __name__ == "__main__":
    main()
