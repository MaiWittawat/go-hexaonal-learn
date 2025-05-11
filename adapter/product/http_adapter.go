package productAdapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wittawat/go-hex/core/entities"
	"github.com/wittawat/go-hex/core/entities/request"
	productPort "github.com/wittawat/go-hex/core/port/product"
)

type HttpProductHandler struct {
	service productPort.ProductService
}

func NewHttpProductHandler(service productPort.ProductService) *HttpProductHandler {
	return &HttpProductHandler{service: service}
}

func newProductFromRequest(req *request.ProductRequest) entities.Product {
	return entities.Product{
		Title:  req.Title,
		Price:  req.Price,
		Detail: req.Detail,
	}
}

func (h *HttpProductHandler) CreateProduct(c *gin.Context) {
	var productReq request.ProductRequest
	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}
	product := newProductFromRequest(&productReq)
	if err := h.service.Save(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created product successfully"})
}

func (h *HttpProductHandler) GetAllProduct(c *gin.Context) {
	products, err := h.service.Find()
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
	product, err := h.service.FindById(id)
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

	var productReq request.ProductRequest
	if err = c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	product := newProductFromRequest(&productReq)
	if err = h.service.UpdateOne(&product, id); err != nil {
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
	if err = h.service.DeleteOne(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted product successfully"})
}
