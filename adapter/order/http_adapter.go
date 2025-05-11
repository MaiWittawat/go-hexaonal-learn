package orderAdapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wittawat/go-hex/core/entities"
	"github.com/wittawat/go-hex/core/entities/request"
	orderPort "github.com/wittawat/go-hex/core/port/order"
)

type HttpOrderHandler struct {
	service orderPort.OrderService
}

func NewHttpOrderHandler(service orderPort.OrderService) *HttpOrderHandler {
	return &HttpOrderHandler{service: service}
}

func newOrderFromRequest(req *request.OrderRequest) entities.Order {
	return entities.Order{
		UserId:    req.UserId,
		ProductId: req.ProductId,
	}
}

func (h *HttpOrderHandler) CreateOrder(c *gin.Context) {
	var orderReq request.OrderRequest
	if err := c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := newOrderFromRequest(&orderReq)
	if err := h.service.Create(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Created order successfully"})
}

func (h *HttpOrderHandler) FindOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orders, err := h.service.GetByUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get orders successfully", "orders": orders})
}

func (h *HttpOrderHandler) UpdateOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var orderReq request.OrderRequest
	if err = c.ShouldBindJSON(&orderReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order := newOrderFromRequest(&orderReq)
	if err = h.service.Update(&order, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Updated order successfully"})
}

func (h *HttpOrderHandler) DeleteOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err = h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted order successfully"})
}
