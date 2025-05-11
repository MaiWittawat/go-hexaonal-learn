// adapter/user/http_adapter.go
package userAdapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wittawat/go-hex/core/entities"
	"github.com/wittawat/go-hex/core/entities/request"
	userPort "github.com/wittawat/go-hex/core/port/user"
)

type HttpUserHandler struct {
	service userPort.UserService
}

func NewHttpUserHandler(service userPort.UserService) *HttpUserHandler {
	return &HttpUserHandler{service: service}
}

func newUserFromRequet(req *request.UserRequest) entities.User {
	return entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
}

func (h *HttpUserHandler) Register(c *gin.Context) {
	var userReq request.UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user := newUserFromRequet(&userReq)
	if err := h.service.Save(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created user successfully"})
}

func (h *HttpUserHandler) Login(c *gin.Context) {
	var userReq request.UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(&userReq)
	if token == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Login fail invalid username, email or password"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successfully", "token": token})
}

func (h *HttpUserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user, err := h.service.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get user successfully", "user": user})
}

func (h *HttpUserHandler) GetAllUser(c *gin.Context) {
	users, err := h.service.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid input"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get all user successfully", "users": users})
}

func (h *HttpUserHandler) UpdateUser(c *gin.Context) {
	emailVal, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	email, ok := emailVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userEmail in context"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	var userReq request.UserRequest
	if err = c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	user := newUserFromRequet(&userReq)
	if err = h.service.UpdateOne(&user, id, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated user successfully"})
}

func (h *HttpUserHandler) DeleteUser(c *gin.Context) {
	emailVal, exists := c.Get("userEmail")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	email, ok := emailVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userEmail"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	if err = h.service.DeleteOne(id, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted user successfully"})
}
