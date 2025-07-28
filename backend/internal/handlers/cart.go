package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"hotel-backend/internal/database"
	"hotel-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	db *database.DB
}

func NewCartHandler(db *database.DB) *CartHandler {
	return &CartHandler{db: db}
}

// GetCartItems handles GET /api/cart
func (h *CartHandler) GetCartItems(c *gin.Context) {
	// For demo purposes, we'll use a user_id from header or default to 1
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
		SELECT ci.id, ci.item_type, ci.item_id, ci.quantity, ci.check_in_date, ci.check_out_date,
		       CASE 
		           WHEN ci.item_type = 'room' THEN r.title
		           WHEN ci.item_type = 'product' THEN p.name
		       END as item_name,
		       CASE 
		           WHEN ci.item_type = 'room' THEN r.price_per_night
		           WHEN ci.item_type = 'product' THEN p.price
		       END as unit_price,
		       CASE 
		           WHEN ci.item_type = 'room' THEN r.images
		           WHEN ci.item_type = 'product' THEN p.images
		       END as images
		FROM cart_items ci
		LEFT JOIN rooms r ON ci.item_type = 'room' AND ci.item_id = r.id
		LEFT JOIN products p ON ci.item_type = 'product' AND ci.item_id = p.id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC`

	rows, err := h.db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
		return
	}
	defer rows.Close()

	var cartItems []map[string]interface{}
	totalAmount := 0.0

	for rows.Next() {
		var item struct {
			ID           int                `json:"id"`
			ItemType     string             `json:"item_type"`
			ItemID       int                `json:"item_id"`
			Quantity     int                `json:"quantity"`
			CheckInDate  sql.NullString     `json:"check_in_date"`
			CheckOutDate sql.NullString     `json:"check_out_date"`
			ItemName     string             `json:"item_name"`
			UnitPrice    float64            `json:"unit_price"`
			Images       models.StringArray `json:"images"`
		}

		err := rows.Scan(
			&item.ID, &item.ItemType, &item.ItemID, &item.Quantity,
			&item.CheckInDate, &item.CheckOutDate, &item.ItemName,
			&item.UnitPrice, &item.Images,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan cart item"})
			return
		}

		// Calculate total price for this item
		itemTotal := item.UnitPrice * float64(item.Quantity)

		// For rooms, calculate nights if dates are provided
		nights := 1
		if item.ItemType == "room" && item.CheckInDate.Valid && item.CheckOutDate.Valid {
			checkIn, _ := time.Parse("2006-01-02", item.CheckInDate.String)
			checkOut, _ := time.Parse("2006-01-02", item.CheckOutDate.String)
			nights = int(checkOut.Sub(checkIn).Hours() / 24)
			if nights < 1 {
				nights = 1
			}
			itemTotal = item.UnitPrice * float64(nights) * float64(item.Quantity)
		}

		cartItem := map[string]interface{}{
			"id":          item.ID,
			"item_type":   item.ItemType,
			"item_id":     item.ItemID,
			"item_name":   item.ItemName,
			"quantity":    item.Quantity,
			"unit_price":  item.UnitPrice,
			"total_price": itemTotal,
			"images":      item.Images,
		}

		if item.CheckInDate.Valid {
			cartItem["check_in_date"] = item.CheckInDate.String
		}
		if item.CheckOutDate.Valid {
			cartItem["check_out_date"] = item.CheckOutDate.String
		}
		if item.ItemType == "room" {
			cartItem["nights"] = nights
		}

		cartItems = append(cartItems, cartItem)
		totalAmount += itemTotal
	}

	c.JSON(http.StatusOK, gin.H{
		"cart_items":   cartItems,
		"total_amount": totalAmount,
		"item_count":   len(cartItems),
		"user_id":      userID,
	})
}

// AddToCart handles POST /api/cart/add
func (h *CartHandler) AddToCart(c *gin.Context) {
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
		ItemType     string `json:"item_type" binding:"required"` // 'room' or 'product'
		ItemID       int    `json:"item_id" binding:"required"`
		Quantity     int    `json:"quantity" binding:"required"`
		CheckInDate  string `json:"check_in_date,omitempty"`
		CheckOutDate string `json:"check_out_date,omitempty"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate item exists
	var exists bool
	if request.ItemType == "room" {
		err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1 AND is_available = true)", request.ItemID).Scan(&exists)
	} else if request.ItemType == "product" {
		err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1 AND is_active = true)", request.ItemID).Scan(&exists)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item type. Must be 'room' or 'product'"})
		return
	}

	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Check if item already exists in cart
	var existingID int
	checkQuery := `
		SELECT id FROM cart_items 
		WHERE user_id = $1 AND item_type = $2 AND item_id = $3 
		AND ($4::date IS NULL OR check_in_date = $4)
		AND ($5::date IS NULL OR check_out_date = $5)`

	var checkInDate, checkOutDate interface{}
	if request.CheckInDate != "" {
		checkInDate = request.CheckInDate
	}
	if request.CheckOutDate != "" {
		checkOutDate = request.CheckOutDate
	}

	err = h.db.QueryRow(checkQuery, userID, request.ItemType, request.ItemID, checkInDate, checkOutDate).Scan(&existingID)

	if err == sql.ErrNoRows {
		// Insert new cart item
		insertQuery := `
			INSERT INTO cart_items (user_id, item_type, item_id, quantity, check_in_date, check_out_date)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`

		var cartItemID int
		err = h.db.QueryRow(insertQuery, userID, request.ItemType, request.ItemID, request.Quantity, checkInDate, checkOutDate).Scan(&cartItemID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":      "Item added to cart successfully",
			"cart_item_id": cartItemID,
		})
	} else if err == nil {
		// Update existing cart item quantity
		updateQuery := `UPDATE cart_items SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
		_, err = h.db.Exec(updateQuery, request.Quantity, existingID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":      "Cart item quantity updated successfully",
			"cart_item_id": existingID,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing cart items"})
		return
	}
}

// RemoveFromCart handles DELETE /api/cart/:id
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		userIDStr = "1" // Default user for demo
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	cartItemIDStr := c.Param("id")
	cartItemID, err := strconv.Atoi(cartItemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	query := `DELETE FROM cart_items WHERE id = $1 AND user_id = $2`
	result, err := h.db.Exec(query, cartItemID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}

// ClearCart handles DELETE /api/cart/clear
func (h *CartHandler) ClearCart(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		userIDStr = "1" // Default user for demo
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	query := `DELETE FROM cart_items WHERE user_id = $1`
	_, err = h.db.Exec(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart cleared successfully"})
}
