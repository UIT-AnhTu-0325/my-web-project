package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// User represents a user in the system
type User struct {
	ID          int       `json:"id" db:"id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	IsAdmin     bool      `json:"is_admin" db:"is_admin"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// OTP represents an OTP record
type OTP struct {
	ID          int       `json:"id" db:"id"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	OTPCode     string    `json:"otp_code" db:"otp_code"`
	ExpiresAt   time.Time `json:"expires_at" db:"expires_at"`
	IsUsed      bool      `json:"is_used" db:"is_used"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// StringArray represents a JSON array of strings in PostgreSQL
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = StringArray{}
		return nil
	}
	return json.Unmarshal(value.([]byte), a)
}

// Room represents a hotel room
type Room struct {
	ID            int         `json:"id" db:"id"`
	RoomNumber    string      `json:"room_number" db:"room_number"`
	RoomType      string      `json:"room_type" db:"room_type"`
	Title         string      `json:"title" db:"title"`
	Description   string      `json:"description" db:"description"`
	PricePerNight float64     `json:"price_per_night" db:"price_per_night"`
	MaxOccupancy  int         `json:"max_occupancy" db:"max_occupancy"`
	Amenities     StringArray `json:"amenities" db:"amenities"`
	Images        StringArray `json:"images" db:"images"`
	IsAvailable   bool        `json:"is_available" db:"is_available"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}

// Product represents a product for sale
type Product struct {
	ID            int         `json:"id" db:"id"`
	Name          string      `json:"name" db:"name"`
	Description   string      `json:"description" db:"description"`
	Price         float64     `json:"price" db:"price"`
	Category      string      `json:"category" db:"category"`
	StockQuantity int         `json:"stock_quantity" db:"stock_quantity"`
	Images        StringArray `json:"images" db:"images"`
	IsActive      bool        `json:"is_active" db:"is_active"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
}

// Order represents a customer order
type Order struct {
	ID            int         `json:"id" db:"id"`
	UserID        int         `json:"user_id" db:"user_id"`
	OrderNumber   string      `json:"order_number" db:"order_number"`
	TotalAmount   float64     `json:"total_amount" db:"total_amount"`
	Status        string      `json:"status" db:"status"`
	CustomerName  string      `json:"customer_name" db:"customer_name"`
	CustomerPhone string      `json:"customer_phone" db:"customer_phone"`
	CustomerEmail string      `json:"customer_email" db:"customer_email"`
	Notes         string      `json:"notes" db:"notes"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	Items         []OrderItem `json:"items,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID           int       `json:"id" db:"id"`
	OrderID      int       `json:"order_id" db:"order_id"`
	ItemType     string    `json:"item_type" db:"item_type"` // 'room' or 'product'
	ItemID       int       `json:"item_id" db:"item_id"`
	ItemName     string    `json:"item_name" db:"item_name"`
	Quantity     int       `json:"quantity" db:"quantity"`
	UnitPrice    float64   `json:"unit_price" db:"unit_price"`
	TotalPrice   float64   `json:"total_price" db:"total_price"`
	CheckInDate  *string   `json:"check_in_date,omitempty" db:"check_in_date"`
	CheckOutDate *string   `json:"check_out_date,omitempty" db:"check_out_date"`
	Nights       *int      `json:"nights,omitempty" db:"nights"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// CartItem represents an item in shopping cart
type CartItem struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	ItemType     string    `json:"item_type" db:"item_type"`
	ItemID       int       `json:"item_id" db:"item_id"`
	Quantity     int       `json:"quantity" db:"quantity"`
	CheckInDate  *string   `json:"check_in_date,omitempty" db:"check_in_date"`
	CheckOutDate *string   `json:"check_out_date,omitempty" db:"check_out_date"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// RoomBooking represents a room reservation
type RoomBooking struct {
	ID           int       `json:"id" db:"id"`
	RoomID       int       `json:"room_id" db:"room_id"`
	OrderID      int       `json:"order_id" db:"order_id"`
	CheckInDate  string    `json:"check_in_date" db:"check_in_date"`
	CheckOutDate string    `json:"check_out_date" db:"check_out_date"`
	GuestCount   int       `json:"guest_count" db:"guest_count"`
	Status       string    `json:"status" db:"status"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
