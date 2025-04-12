-- 1. Enum Types
CREATE TYPE order_status AS ENUM ('active', 'closed');
CREATE TYPE inventory_transaction_type AS ENUM ('increment', 'decrement', 'adjustment');

-- 2. Create Tables

-- Customers table with id as SERIAL
CREATE TABLE customers(
    customer_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    number VARCHAR(20)
);


-- Menu Items table with id as SERIAL
CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    current_price DECIMAL(10, 2) NOT NULL,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT positive_price CHECK (current_price > 0)
);

-- Inventory table with id as SERIAL (renamed column to id
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    quantity INT,
    unit VARCHAR(255),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Join Table: Menu Item Ingredients
CREATE TABLE menu_item_ingredients (
    menu_item_id INT REFERENCES menu_items(id) ON DELETE CASCADE,
    inventory_id INT REFERENCES inventory(id) ON DELETE CASCADE,
    quantity DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (menu_item_id, inventory_id),
    CONSTRAINT positive_quantity CHECK (quantity > 0)
);

-- Orders table with id as SERIAL
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(customer_id) ON DELETE CASCADE,
    status order_status NOT NULL DEFAULT 'active',
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(), 
    CONSTRAINT positive_total_amount CHECK (total_amount >= 0)
);

-- Order Items table with id as SERIAL
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id INT REFERENCES menu_items(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    CONSTRAINT positive_quantity CHECK (quantity > 0),
    CONSTRAINT positive_unit_price CHECK (unit_price > 0)
);

-- Order Status History table with id as SERIAL
CREATE TABLE order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id) ON DELETE CASCADE,
    status order_status NOT NULL,
    notes TEXT,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Price History table with id as SERIAL
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_items(id) ON DELETE CASCADE,
    price DECIMAL(10, 2) NOT NULL,
    effective_from TIMESTAMP WITH TIME ZONE NOT NULL,
    effective_to TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT positive_price CHECK (price > 0)
);

-- Inventory Transactions table with id as SERIAL
CREATE TABLE inventory_transactions (
    id SERIAL PRIMARY KEY,
    inventory_id INT REFERENCES inventory(id) ON DELETE CASCADE,
    transaction_type inventory_transaction_type NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    order_id INT REFERENCES orders(id) ON DELETE SET NULL,
    notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    CONSTRAINT non_zero_quantity CHECK (quantity != 0)
);

-- 3. Create Indexes

CREATE INDEX idx_customer_email ON customers(email);
CREATE INDEX idx_menu_item_name ON menu_items(name);
CREATE INDEX idx_inventory_name ON inventory(name);
CREATE INDEX idx_order_customer ON orders(customer_id);
CREATE INDEX idx_order_status ON orders(status);
CREATE INDEX idx_order_item_order ON order_items(order_id);
CREATE INDEX idx_order_item_menu ON order_items(menu_item_id);
CREATE INDEX idx_price_history_menu ON price_history(menu_item_id);
CREATE INDEX idx_inventory_transaction_inventory ON inventory_transactions(inventory_id);
CREATE INDEX idx_inventory_transaction_order ON inventory_transactions(order_id);

-- 4. Insert Sample Data

-- Insert Customers (IDs will be auto-generated)
INSERT INTO customers (name, email, number)
VALUES
    ('Alice Johnson', 'alice@example.com', '123-456-7890'),
    ('Bob Smith', 'bob@example.com', '234-567-8901'),
    ('Charlie Brown', 'charlie@example.com', '345-678-9012'),
    ('Diana Green', 'diana@example.com', '456-789-0123'),
    ('Edward White', 'edward@example.com', '567-890-1234');

-- Insert Menu Items
INSERT INTO menu_items (name, description, category, current_price, is_available)
VALUES
    ('Espresso coffee', 'Strong black coffee', 'Beverage', 3.50, true),
    ('Cappuccino', 'Espresso with steamed milk and foam', 'Beverage', 4.00, true),
    ('Latte', 'Espresso with steamed milk', 'Beverage', 4.50, true),
    ('Mocha', 'Chocolate-flavored coffee', 'Beverage', 5.00, true),
    ('Americano coffee', 'Diluted espresso', 'Beverage', 3.00, true);

-- Insert Inventory Items
INSERT INTO inventory (name, quantity, unit)
VALUES
    ('Espresso Shot', 500, 'shots'),
    ('Milk', 1000, 'ml'),
    ('Chocolate Syrup', 200, 'ml'),
    ('Sugar', 1000, 'grams'),
    ('Coffee Beans', 2000, 'grams');

-- Insert Menu Item Ingredients (Assuming menu_item_id and inventory_id values correspond to the auto-generated IDs)
-- For example, suppose:
-- Espresso (menu_item id=1) uses Coffee Beans (inventory id=5)
-- Cappuccino (menu_item id=2) uses Coffee Beans (5) and Milk (2)
-- Latte (menu_item id=3) uses Coffee Beans (5) and Milk (2)
INSERT INTO menu_item_ingredients (menu_item_id, inventory_id, quantity)
VALUES
    (1, 5, 20),      -- Espresso needs 20 units of Coffee Beans
    (2, 5, 20),      -- Cappuccino needs 20 units of Coffee Beans
    (2, 2, 150),     -- Cappuccino needs 150 ml of Milk
    (3, 5, 20),      -- Latte needs 20 units of Coffee Beans
    (3, 2, 200);     -- Latte needs 200 ml of Milk

-- Insert Orders (Assuming customer IDs 1 through 5)
INSERT INTO orders (customer_id, status, total_amount, created_at)
VALUES
    (1, 'closed', 7.50, '2025-03-01'),
    (2, 'closed', 8.50, '2025-03-03'),
    (3, 'active', 10.00, '2025-03-03'),
    (4, 'closed', 9.00, '2025-03-04'),
    (5, 'active', 6.50, '2025-03-05');

-- Insert Order Items (Assuming orders and menu_items have the following IDs)
INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price)
VALUES
    (1, 1, 1, 3.50),
    (1, 2, 1, 4.00),
    (2, 3, 1, 4.50),
    (3, 4, 1, 5.00),
    (4, 5, 1, 3.00);

-- Insert Order Status History
INSERT INTO order_status_history (order_id, status, notes, updated_at)
VALUES
    (1, 'closed', 'Order placed', '2025-03-01'),
    (2, 'closed', 'Being prepared', '2025-03-03'),
    (3, 'active', 'Order delivered', '2025-03-03'),
    (4, 'closed', 'Waiting for payment', '2025-03-04'),
    (5, 'active', 'Customer cancelled the order', '2025-03-05');

-- Insert Price History (Using NOW() for effective_from)
INSERT INTO price_history (menu_item_id, price, effective_from)
VALUES
    (1, 3.50, NOW()),
    (2, 4.00, NOW()),
    (3, 4.50, NOW()),
    (4, 5.00, NOW()),
    (5, 3.00, NOW());

-- Insert Inventory Transactions
INSERT INTO inventory_transactions (inventory_id, transaction_type, quantity, order_id, notes, price)
VALUES
    (5, 'decrement', 20, 1, 'Used for Espresso', 200),
    (5, 'decrement', 20, 2, 'Used for Cappuccino', 300),
    (2, 'decrement', 150, 2, 'Used for Cappuccino', 400),
    (5, 'decrement', 20, 3, 'Used for Latte', 100),
    (2, 'decrement', 200, 3, 'Used for Latte', 350);


CREATE OR REPLACE FUNCTION update_menu_availability_on_inventory_change()
RETURNS TRIGGER AS $$
BEGIN
    -- Loop through all menu items using this inventory item
    UPDATE menu_items SET is_available = (
        -- Check if ALL ingredients for the menu item have enough inventory
        NOT EXISTS (
            SELECT 1
            FROM menu_item_ingredients mii
            JOIN inventory i ON mii.inventory_id = i.id
            WHERE mii.menu_item_id = menu_items.id
              AND i.quantity < mii.quantity
        )
    )
    WHERE id IN (
        SELECT menu_item_id FROM menu_item_ingredients WHERE inventory_id = NEW.id
    );

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER trg_update_menu_availability
AFTER UPDATE ON inventory
FOR EACH ROW
EXECUTE FUNCTION update_menu_availability_on_inventory_change();
