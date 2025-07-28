package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"hotel-backend/internal/database"
	"hotel-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	db *database.DB
}

func NewOrderHandler(db *database.DB) *OrderHandler {
	return &OrderHandler{db: db}
}

// GetOrders handles GET /api/orders (user's orders)
func (h *OrderHandler) GetOrders(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		userIDStr = "1" // Default user for demo
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	query := `
		SELECT id, order_number, total_amount, status, customer_name, 
		       customer_phone, customer_email, notes, created_at, updated_at
		FROM orders 
		WHERE user_id = $1 
		ORDER BY created_at DESC`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.OrderNumber, &order.TotalAmount, &order.Status,
			&order.CustomerName, &order.CustomerPhone, &order.CustomerEmail,
			&order.Notes, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order data"})
			return
		}

		// Get order items
		order.Items, _ = h.getOrderItems(order.ID)
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

// CreateOrder handles POST /api/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		userIDStr = "1" // Default user for demo
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var request struct {
		CustomerName  string `json:"customer_name" binding:"required"`
		CustomerPhone string `json:"customer_phone" binding:"required"`
		CustomerEmail string `json:"customer_email"`
		Notes         string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Get cart items
	cartQuery := `
		SELECT ci.item_type, ci.item_id, ci.quantity, ci.check_in_date, ci.check_out_date,
		       CASE 
		           WHEN ci.item_type = 'room' THEN r.title
		           WHEN ci.item_type = 'product' THEN p.name
		       END as item_name,
		       CASE 
		           WHEN ci.item_type = 'room' THEN r.price_per_night
		           WHEN ci.item_type = 'product' THEN p.price
		       END as unit_price
		FROM cart_items ci
		LEFT JOIN rooms r ON ci.item_type = 'room' AND ci.item_id = r.id
		LEFT JOIN products p ON ci.item_type = 'product' AND ci.item_id = p.id
		WHERE ci.user_id = $1`

	rows, err := tx.Query(cartQuery, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	totalAmount := 0.0

	for rows.Next() {
		var item models.OrderItem
		var checkInDate, checkOutDate sql.NullString

		err := rows.Scan(
			&item.ItemType, &item.ItemID, &item.Quantity,
			&checkInDate, &checkOutDate, &item.ItemName, &item.UnitPrice,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan cart item"})
			return
		}

		// Calculate total price
		if item.ItemType == "room" && checkInDate.Valid && checkOutDate.Valid {
			checkIn, _ := time.Parse("2006-01-02", checkInDate.String)
			checkOut, _ := time.Parse("2006-01-02", checkOutDate.String)
			nights := int(checkOut.Sub(checkIn).Hours() / 24)
			if nights < 1 {
				nights = 1
			}
			item.Nights = &nights
			item.CheckInDate = &checkInDate.String
			item.CheckOutDate = &checkOutDate.String
			item.TotalPrice = item.UnitPrice * float64(nights) * float64(item.Quantity)
		} else {
			item.TotalPrice = item.UnitPrice * float64(item.Quantity)
		}

		orderItems = append(orderItems, item)
		totalAmount += item.TotalPrice
	}

	if len(orderItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%d-%d", time.Now().Unix(), userID)

	// Create order
	orderQuery := `
		INSERT INTO orders (user_id, order_number, total_amount, customer_name, customer_phone, customer_email, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	var orderID int
	var createdAt time.Time
	err = tx.QueryRow(orderQuery, userID, orderNumber, totalAmount, request.CustomerName, request.CustomerPhone, request.CustomerEmail, request.Notes).Scan(&orderID, &createdAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Create order items
	for _, item := range orderItems {
		itemQuery := `
			INSERT INTO order_items (order_id, item_type, item_id, item_name, quantity, unit_price, total_price, check_in_date, check_out_date, nights)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

		_, err = tx.Exec(itemQuery, orderID, item.ItemType, item.ItemID, item.ItemName, item.Quantity, item.UnitPrice, item.TotalPrice, item.CheckInDate, item.CheckOutDate, item.Nights)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order item"})
			return
		}

		// Create room booking if item is a room
		if item.ItemType == "room" && item.CheckInDate != nil && item.CheckOutDate != nil {
			bookingQuery := `
				INSERT INTO room_bookings (room_id, order_id, check_in_date, check_out_date, guest_count)
				VALUES ($1, $2, $3, $4, $5)`

			_, err = tx.Exec(bookingQuery, item.ItemID, orderID, *item.CheckInDate, *item.CheckOutDate, item.Quantity)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room booking"})
				return
			}
		}
	}

	// Clear cart
	_, err = tx.Exec("DELETE FROM cart_items WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Send email notifications
	go h.sendOrderNotifications(orderID, orderNumber, request.CustomerName, request.CustomerPhone, request.CustomerEmail, totalAmount, orderItems)

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Order created successfully",
		"order_id":     orderID,
		"order_number": orderNumber,
		"total_amount": totalAmount,
		"created_at":   createdAt,
	})
}

// GetOrderByID handles GET /api/orders/:id
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	query := `
		SELECT id, user_id, order_number, total_amount, status, customer_name, 
		       customer_phone, customer_email, notes, created_at, updated_at
		FROM orders 
		WHERE id = $1`

	var order models.Order
	err = h.db.QueryRow(query, orderID).Scan(
		&order.ID, &order.UserID, &order.OrderNumber, &order.TotalAmount, &order.Status,
		&order.CustomerName, &order.CustomerPhone, &order.CustomerEmail,
		&order.Notes, &order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Get order items
	order.Items, err = h.getOrderItems(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// Helper function to get order items
func (h *OrderHandler) getOrderItems(orderID int) ([]models.OrderItem, error) {
	query := `
		SELECT id, item_type, item_id, item_name, quantity, unit_price, total_price, 
		       check_in_date, check_out_date, nights, created_at
		FROM order_items 
		WHERE order_id = $1 
		ORDER BY id`

	rows, err := h.db.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID, &item.ItemType, &item.ItemID, &item.ItemName,
			&item.Quantity, &item.UnitPrice, &item.TotalPrice,
			&item.CheckInDate, &item.CheckOutDate, &item.Nights, &item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

// Helper function to send order notifications
func (h *OrderHandler) sendOrderNotifications(orderID int, orderNumber, customerName, customerPhone, customerEmail string, totalAmount float64, items []models.OrderItem) {
	emailServiceURL := os.Getenv("EMAIL_SERVICE_URL")
	if emailServiceURL == "" {
		emailServiceURL = "http://localhost:8001"
	}

	// Prepare email data
	emailData := map[string]interface{}{
		"order_number":   orderNumber,
		"customer_name":  customerName,
		"customer_phone": customerPhone,
		"customer_email": customerEmail,
		"total_amount":   totalAmount,
		"status":         "confirmed",
		"items":          items,
	}

	jsonData, _ := json.Marshal(emailData)

	// Send customer confirmation email
	if customerEmail != "" {
		http.Post(emailServiceURL+"/send-order-confirmation", "application/json", bytes.NewBuffer(jsonData))
	}

	// Send admin notification
	http.Post(emailServiceURL+"/send-admin-notification", "application/json", bytes.NewBuffer(jsonData))
}
