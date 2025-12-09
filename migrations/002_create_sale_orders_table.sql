-- Create sale_orders table
CREATE TABLE IF NOT EXISTS sale_orders (
    id UUID PRIMARY KEY,
    order_number VARCHAR(255) NOT NULL UNIQUE,
    customer_name VARCHAR(255) NOT NULL,
    total_amount DECIMAL(15, 2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_sale_orders_order_number ON sale_orders(order_number);

CREATE INDEX IF NOT EXISTS idx_sale_orders_created_by ON sale_orders(created_by);

CREATE INDEX IF NOT EXISTS idx_sale_orders_deleted_at ON sale_orders(deleted_at);

CREATE INDEX IF NOT EXISTS idx_sale_orders_status ON sale_orders(status);
