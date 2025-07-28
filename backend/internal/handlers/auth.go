package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"hotel-backend/internal/database"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	db *database.DB
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// SendOTP handles POST /api/auth/send-otp
func (h *AuthHandler) SendOTP(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate 6-digit OTP
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Set expiry time (5 minutes from now)
	expiresAt := time.Now().Add(5 * time.Minute)

	// Store OTP in database
	query := `
		INSERT INTO otps (phone_number, otp_code, expires_at)
		VALUES ($1, $2, $3)`

	_, err := h.db.Exec(query, request.PhoneNumber, otp, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	// In a real application, you would send the OTP via SMS
	// For demo purposes, we'll return it in the response
	c.JSON(http.StatusOK, gin.H{
		"message":      "OTP sent successfully",
		"phone_number": request.PhoneNumber,
		"otp_code":     otp, // Remove this in production
		"expires_at":   expiresAt,
	})
}

// VerifyOTP handles POST /api/auth/verify-otp
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var request struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		OTPCode     string `json:"otp_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if OTP is valid and not expired
	var otpID int
	var isUsed bool
	var expiresAt time.Time

	query := `
		SELECT id, is_used, expires_at 
		FROM otps 
		WHERE phone_number = $1 AND otp_code = $2 
		ORDER BY created_at DESC 
		LIMIT 1`

	err := h.db.QueryRow(query, request.PhoneNumber, request.OTPCode).Scan(&otpID, &isUsed, &expiresAt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	if isUsed {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP already used"})
		return
	}

	if time.Now().After(expiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired"})
		return
	}

	// Mark OTP as used
	_, err = h.db.Exec("UPDATE otps SET is_used = true WHERE id = $1", otpID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify OTP"})
		return
	}

	// Check if user exists, if not create new user
	var userID int
	var userName string
	var isAdmin bool

	userQuery := `SELECT id, name, is_admin FROM users WHERE phone_number = $1`
	err = h.db.QueryRow(userQuery, request.PhoneNumber).Scan(&userID, &userName, &isAdmin)

	if err != nil {
		// Create new user
		insertQuery := `
			INSERT INTO users (phone_number, name) 
			VALUES ($1, $2) 
			RETURNING id, name, is_admin`

		defaultName := "User" // You might want to ask for name in a separate step
		err = h.db.QueryRow(insertQuery, request.PhoneNumber, defaultName).Scan(&userID, &userName, &isAdmin)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// In a real application, you would generate a JWT token here
	// For demo purposes, we'll return user info
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":           userID,
			"name":         userName,
			"phone_number": request.PhoneNumber,
			"is_admin":     isAdmin,
		},
		"token": fmt.Sprintf("demo-token-%d", userID), // Replace with real JWT
	})
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// In a real application, you would invalidate the JWT token
	// For now, just return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetProfile handles GET /api/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user struct {
		ID          int    `json:"id"`
		PhoneNumber string `json:"phone_number"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		IsAdmin     bool   `json:"is_admin"`
		CreatedAt   string `json:"created_at"`
	}

	query := `
		SELECT id, phone_number, name, email, is_admin, created_at
		FROM users 
		WHERE id = $1`

	err := h.db.QueryRow(query, userIDStr).Scan(
		&user.ID, &user.PhoneNumber, &user.Name, &user.Email, &user.IsAdmin, &user.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
