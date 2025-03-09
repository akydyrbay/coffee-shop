-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "pg_trgm";  -- For full-text search capabilities

-- Enum types
CREATE TYPE order_status AS ENUM ('pending', 'processing', 'completed', 'cancelled', 'closed');
CREATE TYPE inventory_transaction_type AS ENUM ('increment', 'decrement', 'adjustment');

-- Create tables with proper constraints and indexes
CREATE TABLE customers (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE menu_items (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    current_price DECIMAL(10, 2) NOT NULL,
    is_available BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_price CHECK (current_price > 0)
);

CREATE TABLE inventory (
    ingredient_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    quantity INT,
    unit VARCHAR(255),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    -- CONSTRAINT positive_quantity CHECK (quantity >= 0),
    -- CONSTRAINT positive_minimum_quantity CHECK (minimum_quantity >= 0),
    -- CONSTRAINT positive_unit_price CHECK (unit_price > 0)
);

CREATE TABLE menu_item_ingredients (
    menu_item_id VARCHAR(255) REFERENCES menu_items(id) ON DELETE CASCADE,
    inventory_id VARCHAR(255) REFERENCES inventory(ingredient_id) ON DELETE RESTRICT,
    quantity DECIMAL(10, 2) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    PRIMARY KEY (menu_item_id, inventory_id),
    CONSTRAINT positive_quantity CHECK (quantity > 0)
);

CREATE TABLE orders (
    id VARCHAR(255) PRIMARY KEY,
    customer_id VARCHAR(255) REFERENCES customers(id) ON DELETE RESTRICT,
    status order_status NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT positive_total_amount CHECK (total_amount >= 0)
);

CREATE TABLE order_items (
    id VARCHAR(255) PRIMARY KEY,
    order_id VARCHAR(255) REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id VARCHAR(255) REFERENCES menu_items(id) ON DELETE RESTRICT,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_quantity CHECK (quantity > 0),
    CONSTRAINT positive_unit_price CHECK (unit_price > 0)
);

CREATE TABLE order_status_history (
    id VARCHAR(255) PRIMARY KEY,
    order_id VARCHAR(255) REFERENCES orders(id) ON DELETE CASCADE,
    status order_status NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE price_history (
    id VARCHAR(255) PRIMARY KEY,
    menu_item_id VARCHAR(255) REFERENCES menu_items(id) ON DELETE CASCADE,
    price DECIMAL(10, 2) NOT NULL,
    effective_from TIMESTAMP WITH TIME ZONE NOT NULL,
    effective_to TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT positive_price CHECK (price > 0)
);

CREATE TABLE inventory_transactions (
    id VARCHAR(255) PRIMARY KEY,
    inventory_id VARCHAR(255) REFERENCES inventory(ingredient_id) ON DELETE RESTRICT,
    transaction_type inventory_transaction_type NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL,
    order_id VARCHAR(255) REFERENCES orders(id) ON DELETE SET NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT non_zero_quantity CHECK (quantity != 0)
);

-- Create indexes for better query performance
CREATE INDEX idx_menu_items_category ON menu_items(category);
CREATE INDEX idx_menu_items_availability ON menu_items(is_available);
-- CREATE INDEX idx_inventory_minimum_quantity ON inventory(minimum_quantity) WHERE quantity <= minimum_quantity;
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX idx_price_history_menu_item_id ON price_history(menu_item_id);
CREATE INDEX idx_inventory_transactions_inventory_id ON inventory_transactions(inventory_id);

-- Full-text search indexes
CREATE INDEX idx_menu_items_name_description ON menu_items 
    USING gin((setweight(to_tsvector('english', name), 'A') || 
               setweight(to_tsvector('english', COALESCE(description, '')), 'B')));

CREATE INDEX idx_customers_name ON customers 
    USING gin(to_tsvector('english', name));

-- Triggers for maintaining updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_menu_items_updated_at
    BEFORE UPDATE ON menu_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_inventory_updated_at
    BEFORE UPDATE ON inventory
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Trigger for maintaining order status history
CREATE OR REPLACE FUNCTION track_order_status_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') OR (OLD.status IS DISTINCT FROM NEW.status) THEN
        INSERT INTO order_status_history (id, order_id, status)
        VALUES (NEW.id || '_' || EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::TEXT, NEW.id, NEW.status);
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER track_order_status_changes
    AFTER INSERT OR UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION track_order_status_changes();

-- Trigger for maintaining price history
CREATE OR REPLACE FUNCTION track_price_changes()
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') OR (OLD.current_price IS DISTINCT FROM NEW.current_price) THEN
        UPDATE price_history
        SET effective_to = CURRENT_TIMESTAMP
        WHERE menu_item_id = NEW.id AND effective_to IS NULL;
        
        INSERT INTO price_history (
            id,
            menu_item_id,
            price,
            effective_from
        )
        VALUES (
            NEW.id || '_' || EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::TEXT,
            NEW.id,
            NEW.current_price,
            CURRENT_TIMESTAMP
        );
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER track_price_changes
    AFTER INSERT OR UPDATE ON menu_items
    FOR EACH ROW
    EXECUTE FUNCTION track_price_changes();

-- Function for updating inventory
CREATE OR REPLACE FUNCTION update_inventory_quantity(
    p_inventory_id VARCHAR(255),
    p_quantity DECIMAL,
    p_transaction_type inventory_transaction_type,
    p_order_id VARCHAR(255) DEFAULT NULL,
    p_notes TEXT DEFAULT NULL
)
RETURNS VOID AS $$
DECLARE
    v_transaction_id VARCHAR(255);
BEGIN
    -- Generate transaction ID
    v_transaction_id := p_inventory_id || '_' || EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)::TEXT;

    -- Update inventory quantity
    UPDATE inventory
    SET quantity = CASE
        WHEN p_transaction_type = 'increment' THEN quantity + p_quantity
        WHEN p_transaction_type = 'decrement' THEN quantity - p_quantity
        WHEN p_transaction_type = 'adjustment' THEN p_quantity
    END
    WHERE id = p_inventory_id;

    -- Record transaction
    INSERT INTO inventory_transactions (
        id,
        inventory_id, 
        transaction_type, 
        quantity, 
        order_id, 
        notes
    )
    VALUES (
        v_transaction_id,
        p_inventory_id, 
        p_transaction_type, 
        p_quantity, 
        p_order_id, 
        p_notes
    );
END;
$$ LANGUAGE plpgsql;
