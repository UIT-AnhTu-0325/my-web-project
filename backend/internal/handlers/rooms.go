package handlers

import (
	"net/http"
	"strconv"

	"hotel-backend/internal/database"
	"hotel-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	db *database.DB
}

func NewRoomHandler(db *database.DB) *RoomHandler {
	return &RoomHandler{db: db}
}

// GetRooms handles GET /api/rooms
func (h *RoomHandler) GetRooms(c *gin.Context) {
	query := `
		SELECT id, room_number, room_type, title, description, price_per_night, 
		       max_occupancy, amenities, images, is_available, created_at, updated_at
		FROM rooms 
		WHERE is_available = true
		ORDER BY room_number`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rooms"})
		return
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID, &room.RoomNumber, &room.RoomType, &room.Title, &room.Description,
			&room.PricePerNight, &room.MaxOccupancy, &room.Amenities, &room.Images,
			&room.IsAvailable, &room.CreatedAt, &room.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan room data"})
			return
		}
		rooms = append(rooms, room)
	}

	c.JSON(http.StatusOK, gin.H{
		"rooms": rooms,
		"count": len(rooms),
	})
}

// GetRoomByID handles GET /api/rooms/:id
func (h *RoomHandler) GetRoomByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room ID"})
		return
	}

	query := `
		SELECT id, room_number, room_type, title, description, price_per_night, 
		       max_occupancy, amenities, images, is_available, created_at, updated_at
		FROM rooms 
		WHERE id = $1`

	var room models.Room
	err = h.db.QueryRow(query, id).Scan(
		&room.ID, &room.RoomNumber, &room.RoomType, &room.Title, &room.Description,
		&room.PricePerNight, &room.MaxOccupancy, &room.Amenities, &room.Images,
		&room.IsAvailable, &room.CreatedAt, &room.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	c.JSON(http.StatusOK, room)
}

// CheckRoomAvailability checks if a room is available for given dates
func (h *RoomHandler) CheckRoomAvailability(c *gin.Context) {
	var request struct {
		RoomID       int    `json:"room_id" binding:"required"`
		CheckInDate  string `json:"check_in_date" binding:"required"`
		CheckOutDate string `json:"check_out_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for overlapping bookings
	query := `
		SELECT COUNT(*) 
		FROM room_bookings 
		WHERE room_id = $1 
		AND status IN ('confirmed', 'checked_in')
		AND (
			(check_in_date <= $2 AND check_out_date > $2) OR
			(check_in_date < $3 AND check_out_date >= $3) OR
			(check_in_date >= $2 AND check_out_date <= $3)
		)`

	var count int
	err := h.db.QueryRow(query, request.RoomID, request.CheckInDate, request.CheckOutDate).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check availability"})
		return
	}

	available := count == 0

	c.JSON(http.StatusOK, gin.H{
		"available":            available,
		"room_id":              request.RoomID,
		"check_in_date":        request.CheckInDate,
		"check_out_date":       request.CheckOutDate,
		"conflicting_bookings": count,
	})
}
