package handlers

import (
	"net/http"
	"strconv"

	"hotel-backend/internal/database"
	"hotel-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	db *database.DB
}

func NewProductHandler(db *database.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

// GetProducts handles GET /api/products
func (h *ProductHandler) GetProducts(c *gin.Context) {
	category := c.Query("category")

	var query string
	var args []interface{}

	if category != "" {
		query = `
			SELECT id, name, description, price, category, stock_quantity, 
			       images, is_active, created_at, updated_at
			FROM products 
			WHERE is_active = true AND category = $1
			ORDER BY name`
		args = append(args, category)
	} else {
		query = `
			SELECT id, name, description, price, category, stock_quantity, 
			       images, is_active, created_at, updated_at
			FROM products 
			WHERE is_active = true
			ORDER BY category, name`
	}

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price,
			&product.Category, &product.StockQuantity, &product.Images,
			&product.IsActive, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product data"})
			return
		}
		products = append(products, product)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
		"category": category,
	})
}

// GetProductByID handles GET /api/products/:id
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	query := `
		SELECT id, name, description, price, category, stock_quantity, 
		       images, is_active, created_at, updated_at
		FROM products 
		WHERE id = $1 AND is_active = true`

	var product models.Product
	err = h.db.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price,
		&product.Category, &product.StockQuantity, &product.Images,
		&product.IsActive, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// GetProductCategories handles GET /api/products/categories
func (h *ProductHandler) GetProductCategories(c *gin.Context) {
	query := `
		SELECT DISTINCT category, COUNT(*) as product_count
		FROM products 
		WHERE is_active = true 
		GROUP BY category 
		ORDER BY category`

	rows, err := h.db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer rows.Close()

	type CategoryInfo struct {
		Category     string `json:"category"`
		ProductCount int    `json:"product_count"`
	}

	var categories []CategoryInfo
	for rows.Next() {
		var cat CategoryInfo
		err := rows.Scan(&cat.Category, &cat.ProductCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan category data"})
			return
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"count":      len(categories),
	})
}
