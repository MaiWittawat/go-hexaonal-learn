// adapter/user/http_adapter.go
package userAdapter

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wittawat/go-hex/adapter/helper"
	"github.com/wittawat/go-hex/core/entities"
	userPort "github.com/wittawat/go-hex/core/port/user"
)

// ------------------------ Entities ------------------------
type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type HttpUserHandler struct {
	service userPort.UserService
}

// ------------------------ Constructor ------------------------
func newUserFromRequest(req *UserRequest) *entities.User {
	return &entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
}

func NewHttpUserHandler(service userPort.UserService) *HttpUserHandler {
	return &HttpUserHandler{service: service}
}

// ------------------------ Method ------------------------
func (h *HttpUserHandler) Register(c *gin.Context) {
	var userReq UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := newUserFromRequest(&userReq)
	if err := h.service.Save(c.Request.Context(), user, "user"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created user successfully"})
}

func (h *HttpUserHandler) SellerRegister(c *gin.Context) {
	var userReq UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := newUserFromRequest(&userReq)
	if err := h.service.Save(c.Request.Context(), user, "seller"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created user successfully"})
}

func (h *HttpUserHandler) Login(c *gin.Context) {
	var userReq UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), newUserFromRequest(&userReq))
	if token == "" {
		log.Println("error: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successfully", "token": token})
}

func (h *HttpUserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.FindById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get user successfully", "user": user})
}

func (h *HttpUserHandler) GetAllUser(c *gin.Context) {
	users, err := h.service.Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get all user successfully", "users": users})
}

func (h *HttpUserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	email, ok := helper.GetUserEmail(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userReq UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := newUserFromRequest(&userReq)
	token, claims, err := h.service.UpdateOne(c.Request.Context(), user, id, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	email, ok = claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.Set("claims", claims)
	c.Set("userEmail", email)
	c.JSON(http.StatusOK, gin.H{"message": "Updated user successfully", "token": token})
}

func (h *HttpUserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	email, ok := helper.GetUserEmail(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if err := h.service.DeleteOne(c.Request.Context(), id, email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted user successfully"})
}
