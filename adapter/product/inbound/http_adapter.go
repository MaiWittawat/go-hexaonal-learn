package productAdapter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wittawat/go-hex/adapter/helper"
	"github.com/wittawat/go-hex/core/entities"
	productPort "github.com/wittawat/go-hex/core/port/product"
)

// ------------------------ Entities ------------------------ //
type productRequest struct {
	Title  string `json:"title"`
	Price  int32  `json:"price"`
	Detail string `json:"detail"`
}

type HttpProductHandler struct {
	service productPort.ProductService
}

// ------------------------ Constructor ------------------------ //
func newProductFromRequest(req *productRequest) entities.Product {
	return entities.Product{
		Title:  req.Title,
		Price:  req.Price,
		Detail: req.Detail,
	}
}

func NewHttpProductHandler(service productPort.ProductService) *HttpProductHandler {
	return &HttpProductHandler{service: service}
}

// ------------------------ Method ------------------------ //
func (h *HttpProductHandler) CreateProduct(c *gin.Context) {
	email, ok := helper.GetUserEmail(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var productReq productRequest
	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product := newProductFromRequest(&productReq)
	if err := h.service.Create(c.Request.Context(), &product, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created product successfully"})
}

func (h *HttpProductHandler) GetAllProduct(c *gin.Context) {
	products, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get all product successfully", "products": products})
}

func (h *HttpProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get all product successfully", "product": product})
}

func (h *HttpProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	email, ok := helper.GetUserEmail(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var productReq productRequest
	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := newProductFromRequest(&productReq)
	if err := h.service.EditOne(c.Request.Context(), &product, id, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated product successfully"})
}

func (h *HttpProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	email, ok := helper.GetUserEmail(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.DropOne(c.Request.Context(), id, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted product successfully"})
}
