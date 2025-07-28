-- Migration: 001_initial_schema.sql
-- Description: Create initial database schema for hotel booking and e-commerce system

-- Users table for authentication
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100),
    email VARCHAR(255),
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- OTP table for phone authentication
CREATE TABLE otps (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(20) NOT NULL,
    otp_code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    is_used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Hotel rooms table
CREATE TABLE rooms (
    id SERIAL PRIMARY KEY,
    room_number VARCHAR(10) UNIQUE NOT NULL,
    room_type VARCHAR(50) NOT NULL, -- single, double, suite, etc.
    title VARCHAR(255) NOT NULL,
    description TEXT,
    price_per_night DECIMAL(10,2) NOT NULL,
    max_occupancy INTEGER NOT NULL,
    amenities JSONB, -- WiFi, AC, TV, etc.
    images JSONB, -- Array of image URLs
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products table for e-commerce
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category VARCHAR(100), -- souvenirs, food, amenities, etc.
    stock_quantity INTEGER DEFAULT 0,
    images JSONB, -- Array of image URLs
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Orders table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- pending, confirmed, cancelled, completed
    customer_name VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(20) NOT NULL,
    customer_email VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order items table (for both rooms and products)
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    item_type VARCHAR(20) NOT NULL, -- 'room' or 'product'
    item_id INTEGER NOT NULL, -- references rooms.id or products.id
    item_name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    
    -- Room-specific fields
    check_in_date DATE,
    check_out_date DATE,
    nights INTEGER,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Shopping cart table
CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    item_type VARCHAR(20) NOT NULL, -- 'room' or 'product'
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1,
    
    -- Room-specific fields
    check_in_date DATE,
    check_out_date DATE,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, item_type, item_id, check_in_date, check_out_date)
);

-- Room bookings table (for tracking room reservations)
CREATE TABLE room_bookings (
    id SERIAL PRIMARY KEY,
    room_id INTEGER REFERENCES rooms(id),
    order_id INTEGER REFERENCES orders(id),
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    guest_count INTEGER NOT NULL,
    status VARCHAR(50) DEFAULT 'confirmed', -- confirmed, checked_in, checked_out, cancelled
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX idx_users_phone ON users(phone_number);
CREATE INDEX idx_otps_phone_expires ON otps(phone_number, expires_at);
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_cart_items_user_id ON cart_items(user_id);
CREATE INDEX idx_room_bookings_room_id_dates ON room_bookings(room_id, check_in_date, check_out_date);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_active ON products(is_active);

-- Insert sample data

-- Sample admin user
INSERT INTO users (phone_number, name, email, is_admin) VALUES 
('+1234567890', 'Admin User', 'admin@hotel.com', true);

-- Sample rooms
INSERT INTO rooms (room_number, room_type, title, description, price_per_night, max_occupancy, amenities, images) VALUES 
('101', 'single', 'Cozy Single Room', 'Perfect for solo travelers with modern amenities', 89.99, 1, 
 '["WiFi", "AC", "TV", "Mini Fridge"]', '["room101_1.jpg", "room101_2.jpg"]'),
('201', 'double', 'Deluxe Double Room', 'Spacious room with city view and premium facilities', 149.99, 2, 
 '["WiFi", "AC", "TV", "Mini Fridge", "City View", "Work Desk"]', '["room201_1.jpg", "room201_2.jpg"]'),
('301', 'suite', 'Presidential Suite', 'Luxury suite with separate living area and panoramic view', 299.99, 4, 
 '["WiFi", "AC", "TV", "Mini Fridge", "City View", "Work Desk", "Jacuzzi", "Living Area"]', '["suite301_1.jpg", "suite301_2.jpg"]');

-- Sample products
INSERT INTO products (name, description, price, category, stock_quantity, images) VALUES 
('Hotel T-Shirt', 'Comfortable cotton t-shirt with hotel logo', 24.99, 'souvenirs', 50, '["tshirt.jpg"]'),
('Local Coffee Beans', 'Premium locally sourced coffee beans', 18.99, 'food', 30, '["coffee.jpg"]'),
('Spa Voucher', 'Relaxing spa treatment voucher (2 hours)', 89.99, 'services', 20, '["spa.jpg"]'),
('City Guide Book', 'Complete guide to local attractions and restaurants', 15.99, 'souvenirs', 25, '["guidebook.jpg"]'),
('Breakfast Package', 'Continental breakfast for two people', 29.99, 'food', 100, '["breakfast.jpg"]');
