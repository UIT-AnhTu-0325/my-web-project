package handlers

import (
	"net/http"
	"strconv"

	"hotel-backend/internal/database"
	"hotel-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	db *database.DB
}

func NewAdminHandler(db *database.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// GetAllOrders handles GET /api/admin/orders
func (h *AdminHandler) GetAllOrders(c *gin.Context) {
	status := c.Query("status")
	limit := c.DefaultQuery("limit", "50")
	offset := c.DefaultQuery("offset", "0")

	var query string
	var args []interface{}
	argCount := 0

	query = `
		SELECT id, user_id, order_number, total_amount, status, customer_name, 
		       customer_phone, customer_email, notes, created_at, updated_at
		FROM orders`

	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
		argCount++
	}

	query += " ORDER BY created_at DESC"

	if limit != "0" {
		argCount++
		query += " LIMIT $" + strconv.Itoa(argCount)
		args = append(args, limit)
	}

	if offset != "0" {
		argCount++
		query += " OFFSET $" + strconv.Itoa(argCount)
		args = append(args, offset)
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.OrderNumber, &order.TotalAmount, &order.Status,
			&order.CustomerName, &order.CustomerPhone, &order.CustomerEmail,
			&order.Notes, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan order data"})
			return
		}
		orders = append(orders, order)
	}

	// Get total count
	var totalCount int
	countQuery := "SELECT COUNT(*) FROM orders"
	if status != "" {
		countQuery += " WHERE status = $1"
		h.db.QueryRow(countQuery, status).Scan(&totalCount)
	} else {
		h.db.QueryRow(countQuery).Scan(&totalCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":      orders,
		"count":       len(orders),
		"total_count": totalCount,
		"status":      status,
	})
}

// UpdateOrderStatus handles PUT /api/admin/orders/:id
func (h *AdminHandler) UpdateOrderStatus(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var request struct {
		Status string `json:"status" binding:"required"`
		Notes  string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := []string{"pending", "confirmed", "cancelled", "completed"}
	isValid := false
	for _, status := range validStatuses {
		if request.Status == status {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be one of: pending, confirmed, cancelled, completed"})
		return
	}

	query := `
		UPDATE orders 
		SET status = $1, notes = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3`

	result, err := h.db.Exec(query, request.Status, request.Notes, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// If order is cancelled, update room bookings
	if request.Status == "cancelled" {
		h.db.Exec("UPDATE room_bookings SET status = 'cancelled' WHERE order_id = $1", orderID)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Order status updated successfully",
		"order_id": orderID,
		"status":   request.Status,
	})
}

// AddRoom handles POST /api/admin/rooms
func (h *AdminHandler) AddRoom(c *gin.Context) {
	var request struct {
		RoomNumber    string             `json:"room_number" binding:"required"`
		RoomType      string             `json:"room_type" binding:"required"`
		Title         string             `json:"title" binding:"required"`
		Description   string             `json:"description"`
		PricePerNight float64            `json:"price_per_night" binding:"required"`
		MaxOccupancy  int                `json:"max_occupancy" binding:"required"`
		Amenities     models.StringArray `json:"amenities"`
		Images        models.StringArray `json:"images"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO rooms (room_number, room_type, title, description, price_per_night, max_occupancy, amenities, images)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	var roomID int
	var createdAt string
	err := h.db.QueryRow(query, request.RoomNumber, request.RoomType, request.Title, request.Description, request.PricePerNight, request.MaxOccupancy, request.Amenities, request.Images).Scan(&roomID, &createdAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Room created successfully",
		"room_id":    roomID,
		"created_at": createdAt,
	})
}

// AddProduct handles POST /api/admin/products
func (h *AdminHandler) AddProduct(c *gin.Context) {
	var request struct {
		Name          string             `json:"name" binding:"required"`
		Description   string             `json:"description"`
		Price         float64            `json:"price" binding:"required"`
		Category      string             `json:"category" binding:"required"`
		StockQuantity int                `json:"stock_quantity" binding:"required"`
		Images        models.StringArray `json:"images"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO products (name, description, price, category, stock_quantity, images)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	var productID int
	var createdAt string
	err := h.db.QueryRow(query, request.Name, request.Description, request.Price, request.Category, request.StockQuantity, request.Images).Scan(&productID, &createdAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Product created successfully",
		"product_id": productID,
		"created_at": createdAt,
	})
}

// GetDashboardStats handles GET /api/admin/dashboard
func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	// Get order statistics
	var totalOrders, pendingOrders, confirmedOrders, completedOrders int
	var totalRevenue float64

	h.db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&totalOrders)
	h.db.QueryRow("SELECT COUNT(*) FROM orders WHERE status = 'pending'").Scan(&pendingOrders)
	h.db.QueryRow("SELECT COUNT(*) FROM orders WHERE status = 'confirmed'").Scan(&confirmedOrders)
	h.db.QueryRow("SELECT COUNT(*) FROM orders WHERE status = 'completed'").Scan(&completedOrders)
	h.db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM orders WHERE status IN ('confirmed', 'completed')").Scan(&totalRevenue)

	// Get room statistics
	var totalRooms, availableRooms, bookedRooms int
	h.db.QueryRow("SELECT COUNT(*) FROM rooms").Scan(&totalRooms)
	h.db.QueryRow("SELECT COUNT(*) FROM rooms WHERE is_available = true").Scan(&availableRooms)
	h.db.QueryRow("SELECT COUNT(DISTINCT room_id) FROM room_bookings WHERE status = 'confirmed' AND check_in_date <= CURRENT_DATE AND check_out_date > CURRENT_DATE").Scan(&bookedRooms)

	// Get product statistics
	var totalProducts, activeProducts int
	h.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&totalProducts)
	h.db.QueryRow("SELECT COUNT(*) FROM products WHERE is_active = true").Scan(&activeProducts)

	// Get recent orders
	recentOrdersQuery := `
		SELECT id, order_number, customer_name, total_amount, status, created_at
		FROM orders 
		ORDER BY created_at DESC 
		LIMIT 5`

	rows, err := h.db.Query(recentOrdersQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent orders"})
		return
	}
	defer rows.Close()

	var recentOrders []map[string]interface{}
	for rows.Next() {
		var order struct {
			ID           int     `json:"id"`
			OrderNumber  string  `json:"order_number"`
			CustomerName string  `json:"customer_name"`
			TotalAmount  float64 `json:"total_amount"`
			Status       string  `json:"status"`
			CreatedAt    string  `json:"created_at"`
		}

		err := rows.Scan(&order.ID, &order.OrderNumber, &order.CustomerName, &order.TotalAmount, &order.Status, &order.CreatedAt)
		if err != nil {
			continue
		}

		recentOrders = append(recentOrders, map[string]interface{}{
			"id":            order.ID,
			"order_number":  order.OrderNumber,
			"customer_name": order.CustomerName,
			"total_amount":  order.TotalAmount,
			"status":        order.Status,
			"created_at":    order.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": gin.H{
			"total":     totalOrders,
			"pending":   pendingOrders,
			"confirmed": confirmedOrders,
			"completed": completedOrders,
			"revenue":   totalRevenue,
		},
		"rooms": gin.H{
			"total":     totalRooms,
			"available": availableRooms,
			"booked":    bookedRooms,
		},
		"products": gin.H{
			"total":  totalProducts,
			"active": activeProducts,
		},
		"recent_orders": recentOrders,
	})
}
