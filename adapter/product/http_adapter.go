package adapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wittawat/go-hex/core/entities"
	port "github.com/wittawat/go-hex/core/port/product"
)

type HttpProductHandler struct {
	ib port.ProductInbound
}

func NewHttpProductHandler(ib port.ProductInbound) *HttpProductHandler {
	return &HttpProductHandler{ib: ib}
}

func (h *HttpProductHandler) CreateProduct(c *gin.Context) {
	var product entities.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	if err := h.ib.Save(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created product successfully"})
}

func (h *HttpProductHandler) GetAllProduct(c *gin.Context) {
	products, err := h.ib.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get all product successfully", "products": products})
}

func (h *HttpProductHandler) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}
	product, err := h.ib.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get all product successfully", "product": product})
}

func (h *HttpProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}

	var product entities.Product
	if err = c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	existProduct, err := h.ib.FindById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if product.Title == "" {
		product.Title = existProduct.Title
	}
	if product.Price == 0 {
		product.Price = existProduct.Price
	}
	if product.Detail == "" {
		product.Detail = existProduct.Detail
	}

	if err = h.ib.UpdateOne(&product, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated product successfully"})
}

func (h *HttpProductHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product id"})
		return
	}
	if err = h.ib.DeleteOne(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted product successfully"})
}
