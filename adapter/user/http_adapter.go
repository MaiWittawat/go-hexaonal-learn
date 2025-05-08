// adapter/user/http_adapter.go
package adapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wittawat/go-hex/core/entities"
	port "github.com/wittawat/go-hex/core/port/user"
)

type HttpUserHandler struct {
	ib port.UserInbound
}

func NewHttpUserHandler(ib port.UserInbound) *HttpUserHandler {
	return &HttpUserHandler{ib: ib}
}

func (h *HttpUserHandler) Register(c *gin.Context) {

	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.ib.Save(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created user successfully"})
}

func (h *HttpUserHandler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	user, err := h.ib.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get user successfully", "user": user})
}

func (h *HttpUserHandler) GetAllUser(c *gin.Context) {
	users, err := h.ib.Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid input"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get all user successfully", "users": users})
}

func (h *HttpUserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	var user entities.User
	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON input"})
		return
	}

	existUser, err := h.ib.FindById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if user.Username == "" {
		user.Username = existUser.Username
	}
	if user.Email == "" {
		user.Email = existUser.Email
	}
	if user.Password == "" {
		user.Password = existUser.Password
	}

	if err = h.ib.UpdateOne(&user, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated user successfully"})
}

func (h *HttpUserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}
	if err = h.ib.DeleteOne(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted user successfully"})
}

func (h *HttpUserHandler) Login(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.ib.Login(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Login successfully", "token": token})
}
