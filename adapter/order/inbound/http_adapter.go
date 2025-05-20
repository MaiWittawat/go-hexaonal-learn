package orderAdapter

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wittawat/go-hex/core/entities"
	orderPort "github.com/wittawat/go-hex/core/port/order"
	"github.com/wittawat/go-hex/utils/helper"
)

// ------------------------ Entities ------------------------ //
type orderRequest struct {
	UserID    string
	ProductID string `json:"product_id"`
}

type HttpOrderHandler struct {
	service orderPort.OrderService
}

// ------------------------ Constructor ------------------------ //
func NewHttpOrderHandler(service orderPort.OrderService) *HttpOrderHandler {
	return &HttpOrderHandler{service: service}
}

func newOrderFromRequest(req *orderRequest) entities.Order {
	return entities.Order{
		UserID:    req.UserID,
		ProductID: req.ProductID,
	}
}

// ------------------------ Method ------------------------ //
func (h *HttpOrderHandler) CreateOrder(c *gin.Context) {
	email, ok := helper.GetUserEmail(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var orderReq orderRequest
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := newOrderFromRequest(&orderReq)
	fmt.Println("order: ", order)
	if err := h.service.Create(c.Request.Context(), &order, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created order successfully"})
}

func (h *HttpOrderHandler) FindOrder(c *gin.Context) {
	id := c.Param("user_id")
	orders, err := h.service.GetByUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get orders successfully", "orders": orders})
}

func (h *HttpOrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	email, ok := helper.GetUserEmail(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var orderReq orderRequest
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := newOrderFromRequest(&orderReq)
	if err := h.service.EditOne(c.Request.Context(), &order, id, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated order successfully"})
}

func (h *HttpOrderHandler) DeleteOrder(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "Deleted order successfully"})
}
